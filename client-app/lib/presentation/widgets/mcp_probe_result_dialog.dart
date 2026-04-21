import 'package:flutter/material.dart';
import 'package:gen/domain/entities/mcp_probe_result_entity.dart';

Future<void> showMcpProbeResultDialog(
  BuildContext context,
  McpProbeResultEntity r,
) {
  return showDialog<void>(
    context: context,
    builder: (ctx) => AlertDialog(
      title: Text(r.ok ? 'MCP: подключение' : 'MCP: ошибка'),
      content: SizedBox(
        width: 480,
        child: SingleChildScrollView(
          child: SelectableText(
            _formatProbe(r),
            style: const TextStyle(fontSize: 13)
          ),
        ),
      ),
      actions: [
        TextButton(
          onPressed: () => Navigator.pop(ctx),
          child: const Text('Закрыть')
        ),
      ],
    ),
  );
}

String _formatProbe(McpProbeResultEntity r) {
  final b = StringBuffer();
  if (!r.ok) {
    b.writeln(r.errorMessage.isNotEmpty ? r.errorMessage : 'Не удалось подключиться');
    return b.toString();
  }

  if (r.serverName.isNotEmpty || r.serverVersion.isNotEmpty) {
    b.writeln('Сервер: ${r.serverName.isNotEmpty ? r.serverName : "?"} ${r.serverVersion.isNotEmpty ? "v${r.serverVersion}" : ""}');
  }

  b.writeln('Возможности: tools=${r.toolsSupported} resources=${r.resourcesSupported} prompts=${r.promptsSupported}');

  if (r.instructions.trim().isNotEmpty) {
    b.writeln();
    b.writeln('Инструкции сервера:');
    b.writeln(r.instructions.trim());
  }

  return b.toString().trim();
}
