List<int> extractSessionFileIdsFromText(String text) {
  if (text.isEmpty) {
    return const [];
  }

  final re = RegExp(
    r'"(?:file_id|workbook_file_id)"\s*:\s*(\d+)',
    multiLine: true,
  );

  final out = <int>[];
  final seen = <int>{};
  for (final m in re.allMatches(text)) {
    final id = int.tryParse(m.group(1) ?? '');
    if (id != null && id > 0 && seen.add(id)) {
      out.add(id);
      if (out.length >= 16) {
        break;
      }
    }
  }

  return out;
}
