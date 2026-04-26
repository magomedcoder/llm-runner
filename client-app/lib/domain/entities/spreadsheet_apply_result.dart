import 'dart:typed_data';

class SpreadsheetApplyResult {
  final Uint8List workbookBytes;
  final String previewTsv;

  final String? exportedCsv;

  const SpreadsheetApplyResult({
    required this.workbookBytes,
    required this.previewTsv,
    this.exportedCsv,
  });
}
