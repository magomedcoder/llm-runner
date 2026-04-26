import 'package:flutter/material.dart';
import 'package:gen/domain/entities/loaded_model_status.dart';

class RunnerLoadedModelSection extends StatelessWidget {
  final LoadedModelStatus status;

  const RunnerLoadedModelSection({
    super.key,
    required this.status,
  });

  @override
  Widget build(BuildContext context) {
    final theme = Theme.of(context);
    return Column(
      crossAxisAlignment: CrossAxisAlignment.stretch,
      children: [
        Text(
          'Модель в памяти',
          style: theme.textTheme.labelMedium?.copyWith(
            color: theme.colorScheme.onSurfaceVariant,
          ),
        ),
        const SizedBox(height: 6),
        if (!status.loaded)
          Text(
            'Не загружена',
            style: theme.textTheme.bodySmall?.copyWith(
              color: theme.colorScheme.onSurfaceVariant,
            ),
          )
        else ...[
          if (status.displayName.isNotEmpty)
            Text(
              status.displayName,
              style: theme.textTheme.bodyMedium?.copyWith(
                fontWeight: FontWeight.w600,
              ),
            ),

          if (status.ggufBasename.isNotEmpty) ...[
            if (status.displayName.isNotEmpty) const SizedBox(height: 4),
            SelectableText(
              status.ggufBasename,
              style: theme.textTheme.bodySmall?.copyWith(
                fontFamily: 'monospace',
                color: theme.colorScheme.onSurfaceVariant,
              ),
            ),
          ],

          if (status.displayName.isEmpty && status.ggufBasename.isEmpty)
            Text(
              'Загружена',
              style: theme.textTheme.bodySmall?.copyWith(
                color: theme.colorScheme.primary,
                fontWeight: FontWeight.w500,
              ),
            ),
        ],
      ],
    );
  }
}
