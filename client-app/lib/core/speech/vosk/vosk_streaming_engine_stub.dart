import 'dart:typed_data';

class VoskStreamingEngine {
  static bool get supportedOnBuild => false;

  bool get isLoaded => false;

  Future<void> loadModel(String modelPath) async {
    throw UnsupportedError('Локальное распознавание речи (Vosk) в этом окружении не поддерживается.');
  }

  void startUtterance() {
    throw UnsupportedError('Локальное распознавание речи (Vosk) в этом окружении не поддерживается.');
  }

  String applyAudioChunk(Uint8List pcm16le) => '';

  String finishUtterance() => '';

  void abortUtterance() {}

  void releaseModel() {}
}
