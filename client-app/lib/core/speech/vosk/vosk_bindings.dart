import 'dart:ffi';

import 'package:ffi/ffi.dart';

final class VoskModel extends Opaque {}

final class VoskRecognizer extends Opaque {}

typedef _VoskSetLogLevelNative = Void Function(Int32 level);

typedef _VoskSetLogLevelDart = void Function(int level);

typedef _VoskModelNewNative = Pointer<VoskModel> Function(Pointer<Utf8> modelPath);

typedef _VoskModelNewDart = Pointer<VoskModel> Function(Pointer<Utf8> modelPath);

typedef _VoskModelFreeNative = Void Function(Pointer<VoskModel> model);

typedef _VoskModelFreeDart = void Function(Pointer<VoskModel> model);

typedef _VoskRecognizerNewNative = Pointer<VoskRecognizer> Function(
  Pointer<VoskModel> model,
  Float sampleRate,
);

typedef _VoskRecognizerNewDart = Pointer<VoskRecognizer> Function(
  Pointer<VoskModel> model,
  double sampleRate,
);

typedef _VoskRecognizerFreeNative = Void Function(Pointer<VoskRecognizer> recognizer);

typedef _VoskRecognizerFreeDart = void Function(Pointer<VoskRecognizer> recognizer);

typedef _VoskRecognizerAcceptSNative = Int32 Function(
  Pointer<VoskRecognizer> recognizer,
  Pointer<Int16> data,
  Int32 length,
);

typedef _VoskRecognizerAcceptSDart = int Function(
  Pointer<VoskRecognizer> recognizer,
  Pointer<Int16> data,
  int length,
);

typedef _VoskRecognizerResultNative = Pointer<Utf8> Function(Pointer<VoskRecognizer> recognizer);

typedef _VoskRecognizerResultDart = Pointer<Utf8> Function(Pointer<VoskRecognizer> recognizer);

typedef _VoskRecognizerPartialResultNative = Pointer<Utf8> Function(Pointer<VoskRecognizer> recognizer);

typedef _VoskRecognizerPartialResultDart = Pointer<Utf8> Function(Pointer<VoskRecognizer> recognizer);

typedef _VoskRecognizerFinalResultNative = Pointer<Utf8> Function(Pointer<VoskRecognizer> recognizer);

typedef _VoskRecognizerFinalResultDart = Pointer<Utf8> Function(Pointer<VoskRecognizer> recognizer);

typedef _VoskRecognizerResetNative = Void Function(Pointer<VoskRecognizer> recognizer);

typedef _VoskRecognizerResetDart = void Function(Pointer<VoskRecognizer> recognizer);

class VoskBindings {
  VoskBindings(DynamicLibrary lib)
    : _setLogLevel = lib.lookupFunction<_VoskSetLogLevelNative, _VoskSetLogLevelDart>('vosk_set_log_level'),
      _modelNew = lib.lookupFunction<_VoskModelNewNative, _VoskModelNewDart>('vosk_model_new'),
      _modelFree = lib.lookupFunction<_VoskModelFreeNative, _VoskModelFreeDart>('vosk_model_free'),
      _recognizerNew = lib.lookupFunction<_VoskRecognizerNewNative, _VoskRecognizerNewDart>('vosk_recognizer_new'),
      _recognizerFree = lib.lookupFunction<_VoskRecognizerFreeNative, _VoskRecognizerFreeDart>('vosk_recognizer_free'),
      _acceptWaveformS = lib.lookupFunction<_VoskRecognizerAcceptSNative, _VoskRecognizerAcceptSDart>('vosk_recognizer_accept_waveform_s'),
      _result = lib.lookupFunction<_VoskRecognizerResultNative, _VoskRecognizerResultDart>('vosk_recognizer_result'),
      _partialResult = lib.lookupFunction<_VoskRecognizerPartialResultNative, _VoskRecognizerPartialResultDart>('vosk_recognizer_partial_result'),
      _finalResult = lib.lookupFunction<_VoskRecognizerFinalResultNative, _VoskRecognizerFinalResultDart>('vosk_recognizer_final_result'),
      _reset = lib.lookupFunction<_VoskRecognizerResetNative, _VoskRecognizerResetDart>('vosk_recognizer_reset');

  final _VoskSetLogLevelDart _setLogLevel;
  final _VoskModelNewDart _modelNew;
  final _VoskModelFreeDart _modelFree;
  final _VoskRecognizerNewDart _recognizerNew;
  final _VoskRecognizerFreeDart _recognizerFree;
  final _VoskRecognizerAcceptSDart _acceptWaveformS;
  final _VoskRecognizerResultDart _result;
  final _VoskRecognizerPartialResultDart _partialResult;
  final _VoskRecognizerFinalResultDart _finalResult;
  final _VoskRecognizerResetDart _reset;

  void setLogLevel(int level) => _setLogLevel(level);

  Pointer<VoskModel> modelNew(Pointer<Utf8> modelPath) => _modelNew(modelPath);

  void modelFree(Pointer<VoskModel> model) => _modelFree(model);

  Pointer<VoskRecognizer> recognizerNew(Pointer<VoskModel> model, double sampleRate) => _recognizerNew(model, sampleRate);

  void recognizerFree(Pointer<VoskRecognizer> recognizer) => _recognizerFree(recognizer);

  int acceptWaveformS(Pointer<VoskRecognizer> recognizer, Pointer<Int16> data, int length) => _acceptWaveformS(recognizer, data, length);

  Pointer<Utf8> recognizerResult(Pointer<VoskRecognizer> recognizer) => _result(recognizer);

  Pointer<Utf8> recognizerPartialResult(Pointer<VoskRecognizer> recognizer) => _partialResult(recognizer);

  Pointer<Utf8> recognizerFinalResult(Pointer<VoskRecognizer> recognizer) => _finalResult(recognizer);

  void recognizerReset(Pointer<VoskRecognizer> recognizer) => _reset(recognizer);
}
