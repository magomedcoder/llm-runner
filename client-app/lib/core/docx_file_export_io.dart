import 'dart:io';
import 'dart:typed_data';

import 'package:file_picker/file_picker.dart';

Future<bool> saveDocxToFileImpl(Uint8List bytes, String fileName) async {
  final path = await FilePicker.platform.saveFile(
    dialogTitle: 'Сохранить документ',
    fileName: fileName,
    type: FileType.custom,
    allowedExtensions: const ['docx'],
  );

  if (path == null) {
    return false;
  }

  await File(path).writeAsBytes(bytes);
  return true;
}
