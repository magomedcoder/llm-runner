import 'dart:convert';
import 'dart:ffi';
import 'dart:io';
import 'dart:typed_data';

import 'package:ffi/ffi.dart';
import 'package:path/path.dart' as p;

import 'package:gen/core/speech/vosk/vosk_bindings.dart';
import 'package:gen/core/speech/vosk/vosk_dynamic_lib.dart';

class VoskStreamingEngine {
  static bool get supportedOnBuild {
    if (Platform.isIOS) {
      return const bool.fromEnvironment('GEN_IOS_VOSK', defaultValue: false);
    }

    return true;
  }

  VoskBindings? _bindings;
  Pointer<VoskModel>? _model;
  Pointer<VoskRecognizer>? _recognizer;
  String? _loadedPath;

  String _committedInSession = '';

  bool get isLoaded => _model != null;

  VoskBindings get _b {
    final v = _bindings;
    if (v == null) {
      throw StateError('Vosk: библиотека не инициализирована');
    }
    return v;
  }

  Future<void> loadModel(String modelPath) async {
    _bindings ??= VoskBindings(openVoskDynamicLibrary());
    _b.setLogLevel(-1);

    final normalized = p.normalize(modelPath);
    if (!Directory(normalized).existsSync()) {
      throw StateError('Каталог модели не найден: $normalized');
    }

    if (_loadedPath == normalized && _model != null) {
      return;
    }

    releaseModel();

    final pathPtr = normalized.toNativeUtf8();
    try {
      final m = _b.modelNew(pathPtr);
      if (m == nullptr) {
        throw StateError(
          'vosk_model_new вернула NULL. Проверьте путь к модели и наличие libvosk в bundle.',
        );
      }
      _model = m;
      _loadedPath = normalized;
    } finally {
      calloc.free(pathPtr);
    }
  }

  void startUtterance() {
    _disposeRecognizer();
    _committedInSession = '';
    final model = _model;
    if (model == null) {
      throw StateError('Сначала загрузите модель Vosk');
    }
    final r = _b.recognizerNew(model, 16000);
    if (r == nullptr) {
      throw StateError('vosk_recognizer_new вернула NULL');
    }
    _recognizer = r;
  }

  String applyAudioChunk(Uint8List pcm16le) {
    final rec = _recognizer;
    if (rec == null) {
      return _committedInSession;
    }
    final sampleCount = pcm16le.length ~/ 2;
    if (sampleCount == 0) {
      return _joinLive(_committedInSession, _readPartial(rec));
    }

    final pSamples = calloc<Int16>(sampleCount);
    try {
      final native = pSamples.asTypedList(sampleCount);
      final bd = ByteData.sublistView(pcm16le);
      for (var i = 0; i < sampleCount; i++) {
        native[i] = bd.getInt16(i * 2, Endian.little);
      }
      final code = _b.acceptWaveformS(rec, pSamples, sampleCount);
      if (code < 0) {
        throw StateError('vosk_recognizer_accept_waveform_s: ошибка декодирования');
      }
      if (code == 1) {
        final segment = _parseResultText(_b.recognizerResult(rec));
        _committedInSession = _appendSegment(_committedInSession, segment);
        _b.recognizerReset(rec);
      }
      return _joinLive(_committedInSession, _readPartial(rec));
    } finally {
      calloc.free(pSamples);
    }
  }

  String _readPartial(Pointer<VoskRecognizer> rec) {
    return _parsePartialText(_b.recognizerPartialResult(rec));
  }

  String finishUtterance() {
    final rec = _recognizer;
    if (rec == null) {
      final out = _committedInSession;
      _committedInSession = '';
      return out;
    }
    final tail = _parseResultText(_b.recognizerFinalResult(rec));
    if (tail.isNotEmpty) {
      _committedInSession = _appendSegment(_committedInSession, tail);
    }
    _disposeRecognizer();
    final out = _committedInSession;
    _committedInSession = '';
    return out;
  }

  void abortUtterance() {
    _disposeRecognizer();
    _committedInSession = '';
  }

  void releaseModel() {
    _disposeRecognizer();
    final m = _model;
    if (m != null) {
      _b.modelFree(m);
      _model = null;
      _loadedPath = null;
    }
  }

  void _disposeRecognizer() {
    final r = _recognizer;
    if (r != null) {
      _b.recognizerFree(r);
      _recognizer = null;
    }
  }

  static String _appendSegment(String a, String b) {
    if (b.isEmpty) {
      return a;
    }
    if (a.isEmpty) {
      return b;
    }
    return '$a $b';
  }

  static String _joinLive(String committed, String partial) {
    if (partial.isEmpty) {
      return committed;
    }
    if (committed.isEmpty) {
      return partial;
    }
    return '$committed $partial';
  }

  static String _parseResultText(Pointer<Utf8> ptr) {
    if (ptr == nullptr) {
      return '';
    }
    try {
      final m = jsonDecode(ptr.toDartString()) as Map<String, dynamic>;
      return (m['text'] as String?) ?? '';
    } catch (_) {
      return '';
    }
  }

  static String _parsePartialText(Pointer<Utf8> ptr) {
    if (ptr == nullptr) {
      return '';
    }
    try {
      final m = jsonDecode(ptr.toDartString()) as Map<String, dynamic>;
      return (m['partial'] as String?) ?? '';
    } catch (_) {
      return '';
    }
  }
}
