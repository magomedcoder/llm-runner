import 'dart:io';
import 'dart:typed_data';

import 'package:file_picker/file_picker.dart';

Future<bool> saveUserPickedFileImpl(Uint8List bytes, String fileName) async {
  final path = await FilePicker.platform.saveFile(
    dialogTitle: 'Сохранить файл',
    fileName: fileName,
  );

  if (path == null) {
    return false;
  }

  await File(path).writeAsBytes(bytes);

  return true;
}
