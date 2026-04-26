import 'dart:async';

import 'package:flutter/foundation.dart';
import 'package:flutter/services.dart';
import 'package:record/record.dart';
import 'package:shared_preferences/shared_preferences.dart';

import 'package:gen/core/speech/vosk/vosk_streaming_engine.dart';

class LocalVoskDictationService {
  LocalVoskDictationService();

  static const prefsKeyModelPath = 'gen_vosk_model_path';

  final VoskStreamingEngine _engine = VoskStreamingEngine();
  final AudioRecorder _recorder = AudioRecorder();

  StreamSubscription<Uint8List>? _audioSub;
  bool _listening = false;

  bool get isListening => _listening;

  bool get isPlatformSupported => !kIsWeb && VoskStreamingEngine.supportedOnBuild;

  Future<String?> getSavedModelPath() async {
    final p = await SharedPreferences.getInstance();
    return p.getString(prefsKeyModelPath);
  }

  Future<void> saveModelPath(String directoryPath) async {
    final p = await SharedPreferences.getInstance();
    await p.setString(prefsKeyModelPath, directoryPath);
  }

  Future<void> _ensureModelLoaded({String? modelPath}) async {
    final path = modelPath ?? await getSavedModelPath();
    if (path == null || path.trim().isEmpty) {
      throw StateError('Не указан путь к модели Vosk');
    }

    await _engine.loadModel(path.trim());
  }

  Future<void> start({
    required String prefix,
    required String suffix,
    required void Function(String text, int caretOffset) onLive,
    String? modelPath,
  }) async {
    if (_listening) {
      return;
    }

    if (!isPlatformSupported) {
      throw UnsupportedError('Голосовой ввод на этой платформе недоступен');
    }

    await _ensureModelLoaded(modelPath: modelPath);

    final permitted = await _recorder.hasPermission();
    if (permitted != true) {
      throw StateError('Нет разрешения на запись с микрофона');
    }

    _engine.startUtterance();
    _listening = true;

    final stream = await _recorder.startStream(
      const RecordConfig(
        encoder: AudioEncoder.pcm16bits,
        sampleRate: 16000,
        numChannels: 1,
      ),
    );

    _audioSub = stream.listen(
      (chunk) {
        final live = _engine.applyAudioChunk(chunk);
        onLive(prefix + live + suffix, prefix.length + live.length);
      },
      onError: (_) {},
    );
  }

  Future<TextEditingValue> stop({
    required String prefix,
    required String suffix,
  }) async {
    if (!_listening) {
      return TextEditingValue.empty;
    }

    _listening = false;
    await _audioSub?.cancel();
    _audioSub = null;
    await _recorder.stop();
    final spoken = _engine.finishUtterance();
    final text = prefix + spoken + suffix;
    final caret = prefix.length + spoken.length;
    return TextEditingValue(text: text, selection: TextSelection.collapsed(offset: caret));
  }

  Future<void> cancel() async {
    if (!_listening) {
      return;
    }

    _listening = false;
    await _audioSub?.cancel();
    _audioSub = null;
    await _recorder.stop();
    _engine.abortUtterance();
  }
}
