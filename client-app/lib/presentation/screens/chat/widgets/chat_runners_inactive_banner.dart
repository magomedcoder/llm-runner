import 'package:flutter/material.dart';

class ChatRunnersInactiveBanner extends StatelessWidget {
  const ChatRunnersInactiveBanner({
    super.key,
    required this.isRefreshing,
    required this.onRefresh,
  });

  final bool isRefreshing;
  final VoidCallback onRefresh;

  @override
  Widget build(BuildContext context) {
    final cs = Theme.of(context).colorScheme;
    final fg = cs.onErrorContainer;
    return Material(
      color: cs.errorContainer.withValues(alpha: 0.5),
      child: Padding(
        padding: const EdgeInsets.symmetric(horizontal: 12, vertical: 8),
        child: Row(
          children: [
            Icon(Icons.warning_amber_rounded, color: fg, size: 22),
            const SizedBox(width: 10),
            Expanded(
              child: Text(
                'Нет активных раннеров. Чат недоступен.',
                style: TextStyle(color: fg),
              ),
            ),
            if (isRefreshing)
              const Padding(
                padding: EdgeInsets.only(left: 8),
                child: SizedBox(
                  width: 22,
                  height: 22,
                  child: CircularProgressIndicator(strokeWidth: 2),
                ),
              )
            else
              TextButton.icon(
                onPressed: onRefresh,
                icon: Icon(Icons.refresh_rounded, size: 18, color: fg),
                label: Text('Обновить', style: TextStyle(color: fg)),
                style: TextButton.styleFrom(
                  foregroundColor: fg,
                  padding: const EdgeInsets.symmetric(horizontal: 12),
                ),
              ),
          ],
        ),
      ),
    );
  }
}
