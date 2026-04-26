import 'dart:ffi';
import 'dart:io';

import 'package:path/path.dart' as p;

DynamicLibrary openVoskDynamicLibrary() {
  final fromEnv = Platform.environment['VOSK_LIBRARY_PATH']?.trim();
  if (fromEnv != null && fromEnv.isNotEmpty) {
    return DynamicLibrary.open(fromEnv);
  }

  if (Platform.isLinux) {
    final bundled = p.join(
      File(Platform.resolvedExecutable).parent.path,
      'lib',
      'libvosk.so',
    );

    if (File(bundled).existsSync()) {
      return DynamicLibrary.open(bundled);
    }

    return DynamicLibrary.open('libvosk.so');
  }

  if (Platform.isMacOS) {
    final exeDir = File(Platform.resolvedExecutable).parent.path;
    final fw = p.normalize(p.join(exeDir, '..', 'Frameworks', 'libvosk.dylib'));
    if (File(fw).existsSync()) {
      return DynamicLibrary.open(fw);
    }

    final bundledLib = p.join(exeDir, 'lib', 'libvosk.dylib');
    if (File(bundledLib).existsSync()) {
      return DynamicLibrary.open(bundledLib);
    }

    return DynamicLibrary.open('libvosk.dylib');
  }

  if (Platform.isWindows) {
    final exeDir = File(Platform.resolvedExecutable).parent;
    for (final name in ['vosk.dll', r'lib\vosk.dll', 'libvosk.dll']) {
      final file = File(p.join(exeDir.path, name));
      if (file.existsSync()) {
        return DynamicLibrary.open(file.path);
      }
    }

    return DynamicLibrary.open('vosk.dll');
  }

  if (Platform.isAndroid) {
    return DynamicLibrary.open('libvosk.so');
  }

  if (Platform.isIOS) {
    final bundleDir = File(Platform.resolvedExecutable).parent.path;
    for (final rel in [
      p.join('Frameworks', 'Vosk.framework', 'Vosk'),
      p.join('Frameworks', 'libvosk.framework', 'libvosk'),
    ]) {
      final path = p.join(bundleDir, rel);
      if (File(path).existsSync()) {
        return DynamicLibrary.open(path);
      }
    }

    return DynamicLibrary.process();
  }

  return DynamicLibrary.open('libvosk.so');
}
