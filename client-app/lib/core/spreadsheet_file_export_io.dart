import 'dart:convert';
import 'dart:io';
import 'dart:typed_data';

import 'package:file_picker/file_picker.dart';

Future<bool> saveSpreadsheetToFileImpl(Uint8List bytes, String fileName) async {
  final path = await FilePicker.platform.saveFile(
    dialogTitle: 'Сохранить таблицу',
    fileName: fileName,
    type: FileType.custom,
    allowedExtensions: const ['xlsx'],
  );

  if (path == null) {
    return false;
  }
  await File(path).writeAsBytes(bytes);

  return true;
}

Future<bool> saveCsvToFileImpl(String utf8Text, String fileName) async {
  final path = await FilePicker.platform.saveFile(
    dialogTitle: 'Сохранить CSV',
    fileName: fileName,
    type: FileType.custom,
    allowedExtensions: const ['csv'],
  );

  if (path == null) {
    return false;
  }

  await File(path).writeAsString(utf8Text, encoding: utf8);

  return true;
}
