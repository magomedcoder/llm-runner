import 'package:flutter/material.dart';
import 'package:gen/core/mcp_connection_config.dart';

class McpConnectionJsonDialogSection extends StatelessWidget {
  const McpConnectionJsonDialogSection({
    super.key,
    required this.jsonCtrl,
    required this.isNew,
  });

  final TextEditingController jsonCtrl;
  final bool isNew;

  void _applyExample(String Function() builder) {
    final text = builder();
    jsonCtrl.text = text;
    jsonCtrl.selection = TextSelection.collapsed(offset: text.length);
  }

  @override
  Widget build(BuildContext context) {
    final theme = Theme.of(context);
    final monoSize = theme.textTheme.bodySmall?.fontSize ?? 12;

    return Column(
      crossAxisAlignment: CrossAxisAlignment.stretch,
      children: [
        Text(
          'Примеры',
          style: theme.textTheme.labelSmall?.copyWith(
            color: theme.colorScheme.onSurfaceVariant,
          ),
        ),
        const SizedBox(height: 2),
        Wrap(
          spacing: 4,
          runSpacing: 0,
          children: [
            TextButton(
              onPressed: () => _applyExample(McpConnectionConfig.exampleJsonStdio),
              child: const Text('stdio'),
            ),
            TextButton(
              onPressed: () => _applyExample(McpConnectionConfig.exampleJsonSse),
              child: const Text('sse'),
            ),
            TextButton(
              onPressed: () => _applyExample(McpConnectionConfig.exampleJsonStreamable),
              child: const Text('streamable'),
            ),
          ],
        ),
        const SizedBox(height: 2),
        Text(
          'Полные примеры',
          style: theme.textTheme.labelSmall?.copyWith(
            color: theme.colorScheme.onSurfaceVariant,
          ),
        ),
        const SizedBox(height: 2),
        Wrap(
          spacing: 4,
          runSpacing: 0,
          children: [
            TextButton(
              onPressed: () => _applyExample(McpConnectionConfig.exampleJsonFullStdio),
              child: const Text('Полный stdio'),
            ),
            TextButton(
              onPressed: () => _applyExample(McpConnectionConfig.exampleJsonFullRemote),
              child: const Text('Полный HTTP'),
            ),
          ],
        ),
        const SizedBox(height: 4),
        ExpansionTile(
          tilePadding: EdgeInsets.zero,
          childrenPadding: const EdgeInsets.only(bottom: 8),
          title: Text('Документация'),
          children: [
            Align(
              alignment: Alignment.centerLeft,
              child: SelectableText(
                McpConnectionConfig.documentation,
                style: theme.textTheme.bodySmall?.copyWith(
                  color: theme.colorScheme.onSurfaceVariant,
                  height: 1.45,
                ),
              ),
            ),
          ],
        ),
        const SizedBox(height: 4),
        TextField(
          controller: jsonCtrl,
          maxLines: 28,
          minLines: 14,
          style: TextStyle(
            fontFamily: 'monospace',
            fontSize: monoSize,
          ),
          decoration: InputDecoration(
            alignLabelWithHint: true,
            labelText: 'transport, command, args, env, url, headers, timeoutSeconds',
            helperText: isNew
              ? 'Короткие примеры stdio/sse/streamable'
              : 'Секреты в env/headers могут отображаться как *** оставьте как есть, чтобы не менять на сервере',
            border: const OutlineInputBorder(),
          ),
        ),
      ],
    );
  }
}
