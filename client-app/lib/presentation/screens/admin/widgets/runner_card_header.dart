import 'package:flutter/material.dart';
import 'package:gen/domain/entities/runner_info.dart';
import 'package:gen/presentation/screens/admin/widgets/runner_status.dart';

class RunnerCardHeader extends StatelessWidget {
  final RunnerInfo runner;
  final RunnerStatus status;
  final VoidCallback onToggleEnabled;

  const RunnerCardHeader({
    super.key,
    required this.runner,
    required this.status,
    required this.onToggleEnabled,
  });

  @override
  Widget build(BuildContext context) {
    final theme = Theme.of(context);
    final statusColor = runnerStatusColor(context, status);

    return Row(
      children: [
        _StatusIndicator(color: statusColor),
        const SizedBox(width: 10),
        Expanded(
          child: Column(
            crossAxisAlignment: CrossAxisAlignment.start,
            children: [
              Text(
                runner.name.isNotEmpty ? runner.name : runner.address,
                style: theme.textTheme.titleSmall?.copyWith(
                  fontWeight: FontWeight.w600,
                  color: runner.enabled
                      ? theme.colorScheme.onSurface
                      : theme.colorScheme.onSurfaceVariant,
                ),
                overflow: TextOverflow.ellipsis,
              ),
              if (runner.name.isNotEmpty)
                Text(
                  runner.address,
                  style: theme.textTheme.bodySmall?.copyWith(
                    fontFamily: 'monospace',
                    color: theme.colorScheme.onSurfaceVariant,
                  ),
                  overflow: TextOverflow.ellipsis,
                ),
            ],
          ),
        ),
        Text(
          status.label,
          style: theme.textTheme.labelMedium?.copyWith(
            color: statusColor,
            fontWeight: FontWeight.w500,
          ),
        ),
        const SizedBox(width: 8),
        Switch(
          value: runner.enabled,
          onChanged: (_) => onToggleEnabled(),
        ),
      ],
    );
  }
}

class _StatusIndicator extends StatelessWidget {
  final Color color;

  const _StatusIndicator({required this.color});

  @override
  Widget build(BuildContext context) {
    return Container(
      width: 12,
      height: 12,
      decoration: BoxDecoration(
        color: color,
        shape: BoxShape.circle,
        boxShadow: [
          BoxShadow(
            color: color.withValues(alpha: 0.4),
            blurRadius: 6,
            spreadRadius: 0,
          ),
        ],
      ),
    );
  }
}
