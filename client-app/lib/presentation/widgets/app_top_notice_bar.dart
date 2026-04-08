import 'package:flutter/material.dart';
import 'package:gen/core/injector.dart';
import 'package:gen/core/ui/app_top_notice_controller.dart';

const double _kNoticeRadius = 12;

Widget _noticeCard({
  required BuildContext context,
  required Color color,
  required Widget child,
}) {
  final scheme = Theme.of(context).colorScheme;
  return Material(
    color: color,
    elevation: 3,
    surfaceTintColor: Colors.transparent,
    shadowColor: scheme.shadow.withValues(alpha: 0.22),
    shape: RoundedRectangleBorder(
      borderRadius: BorderRadius.circular(_kNoticeRadius),
    ),
    clipBehavior: Clip.antiAlias,
    child: child,
  );
}

class AppTopNoticeBar extends StatelessWidget {
  const AppTopNoticeBar({super.key});

  @override
  Widget build(BuildContext context) {
    return AnimatedBuilder(
      animation: sl<AppTopNoticeController>(),
      builder: (context, _) {
        final controller = sl<AppTopNoticeController>();
        final entry = controller.current;
        return AnimatedSwitcher(
          duration: const Duration(milliseconds: 280),
          switchInCurve: Curves.easeOutCubic,
          switchOutCurve: Curves.easeInCubic,
          layoutBuilder: (Widget? currentChild, List<Widget> previousChildren) {
            return Stack(
              alignment: Alignment.topCenter,
              fit: StackFit.passthrough,
              clipBehavior: Clip.hardEdge,
              children: <Widget>[
                ...previousChildren,
                ?currentChild,
              ],
            );
          },
          transitionBuilder: (child, animation) {
            final offset = Tween<Offset>(
              begin: const Offset(0, -0.15),
              end: Offset.zero,
            ).animate(
              CurvedAnimation(
                parent: animation,
                curve: Curves.easeOutCubic,
              ),
            );
            return ClipRect(
              clipBehavior: Clip.hardEdge,
              child: SlideTransition(
                position: offset,
                child: child,
              ),
            );
          },
          child: entry == null
              ? const SizedBox(
                  key: ValueKey<String>('app-top-notice-empty'),
                  width: double.infinity,
                )
              : _AppTopNoticeBanner(
                  key: ValueKey<int>(entry.id),
                  entry: entry,
                  onDismiss: controller.dismissCurrent,
                ),
        );
      },
    );
  }
}

class _AppTopNoticeBanner extends StatelessWidget {
  const _AppTopNoticeBanner({
    super.key,
    required this.entry,
    required this.onDismiss,
  });

  final AppTopNoticeEntry entry;
  final VoidCallback onDismiss;

  @override
  Widget build(BuildContext context) {
    final scheme = Theme.of(context).colorScheme;
    final bg = entry.error ? scheme.errorContainer : scheme.secondaryContainer;
    final fg = entry.error ? scheme.onErrorContainer : scheme.onSecondaryContainer;

    return _noticeCard(
      context: context,
      color: bg,
      child: Padding(
        padding: const EdgeInsets.symmetric(vertical: 10, horizontal: 12),
        child: Row(
          crossAxisAlignment: CrossAxisAlignment.start,
          children: [
            Icon(
              entry.error ? Icons.error_outline : Icons.info_outline,
              size: 20,
              color: fg,
            ),
            const SizedBox(width: 10),
            Expanded(
              child: Text(
                entry.message,
                style: TextStyle(
                  fontSize: 14,
                  fontWeight: FontWeight.w500,
                  color: fg,
                  height: 1.35,
                ),
              ),
            ),
            Semantics(
              label: 'Закрыть',
              button: true,
              child: Material(
                color: Colors.transparent,
                child: InkWell(
                  onTap: onDismiss,
                  customBorder: const CircleBorder(),
                  child: Padding(
                    padding: const EdgeInsets.all(8),
                    child: Icon(Icons.close_rounded, size: 20, color: fg),
                  ),
                ),
              ),
            ),
          ],
        ),
      ),
    );
  }
}
