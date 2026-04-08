import 'package:flutter/material.dart';
import 'package:gen/core/injector.dart';
import 'package:gen/core/ui/app_top_notice_controller.dart';
import 'package:gen/presentation/screens/chat/bloc/chat_state.dart';
import 'package:gen/presentation/widgets/app_top_notice_bar.dart';

const double _kAlertOuterHorizontal = 12;
const double _kAlertOuterTop = 4;
const double _kAlertOuterBottom = 8;
const double _kAlertRadius = 12;

enum ChatServerConnectionStatus {
  connected,
  connecting,
  disconnected,
}

ChatServerConnectionStatus resolveChatServerConnectionStatus(ChatState state) {
  if (state.isConnected) {
    return ChatServerConnectionStatus.connected;
  }

  if (state.isLoading) {
    return ChatServerConnectionStatus.connecting;
  }

  if (state.hasCompletedInitialConnection) {
    return ChatServerConnectionStatus.disconnected;
  }

  return ChatServerConnectionStatus.connected;
}

enum _AlertLevel {
  none,
  connectionDisconnected,
  connectionConnecting,
  runnersInactive,
  chatError,
  streamNotice,
  appTopNotice,
}

_AlertLevel _resolveAlertLevel(ChatState state, AppTopNoticeController notices) {
  final conn = resolveChatServerConnectionStatus(state);
  if (conn == ChatServerConnectionStatus.disconnected) {
    return _AlertLevel.connectionDisconnected;
  }
  
  if (conn == ChatServerConnectionStatus.connecting) {
    return _AlertLevel.connectionConnecting;
  }

  if (state.hasActiveRunners == false) {
    return _AlertLevel.runnersInactive;
  }

  final err = state.error?.trim();
  if (err != null && err.isNotEmpty) {
    return _AlertLevel.chatError;
  }

  final sn = state.streamNotice?.trim();
  if (sn != null && sn.isNotEmpty) {
    return _AlertLevel.streamNotice;
  }

  if (notices.current != null) {
    return _AlertLevel.appTopNotice;
  }

  return _AlertLevel.none;
}

ShapeBorder _alertShape() {
  return RoundedRectangleBorder(
    borderRadius: BorderRadius.circular(_kAlertRadius),
  );
}

Widget _alertCard({
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
    shape: _alertShape(),
    clipBehavior: Clip.antiAlias,
    child: child,
  );
}

class ChatStatusAlertsColumn extends StatelessWidget {
  const ChatStatusAlertsColumn({
    super.key,
    required this.state,
    this.onRetryConnection,
    this.onRefreshRunners,
    this.onClearChatError,
    this.onDismissStreamNotice,
  });

  final ChatState state;
  final VoidCallback? onRetryConnection;
  final VoidCallback? onRefreshRunners;
  final VoidCallback? onClearChatError;
  final VoidCallback? onDismissStreamNotice;

  @override
  Widget build(BuildContext context) {
    final controller = sl<AppTopNoticeController>();
    return AnimatedBuilder(
      animation: controller,
      builder: (context, _) {
        final level = _resolveAlertLevel(state, controller);
        final strip = AnimatedSwitcher(
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
          child: _buildStripForLevel(
            context,
            level,
            controller,
            key: ValueKey<_AlertLevel>(level),
          ),
        );

        if (level == _AlertLevel.none) {
          return strip;
        }

        return Padding(
          padding: const EdgeInsets.only(
            left: _kAlertOuterHorizontal,
            right: _kAlertOuterHorizontal,
            top: _kAlertOuterTop,
            bottom: _kAlertOuterBottom,
          ),
          child: strip,
        );
      },
    );
  }

  Widget _buildStripForLevel(
    BuildContext context,
    _AlertLevel level,
    AppTopNoticeController controller, {
    required Key key,
  }) {
    switch (level) {
      case _AlertLevel.none:
        return SizedBox.shrink(key: key);
      case _AlertLevel.connectionDisconnected:
        return _ConnectionStatusStrip(
          key: key,
          state: state,
          onRetry: onRetryConnection,
          fixedStatus: ChatServerConnectionStatus.disconnected,
        );
      case _AlertLevel.connectionConnecting:
        return _ConnectionStatusStrip(
          key: key,
          state: state,
          onRetry: onRetryConnection,
          fixedStatus: ChatServerConnectionStatus.connecting,
        );
      case _AlertLevel.runnersInactive:
        return _RunnersInactiveStrip(
          key: key,
          state: state,
          onRefresh: onRefreshRunners,
        );
      case _AlertLevel.chatError:
        return _ChatErrorStrip(
          key: key,
          message: state.error,
          onDismiss: onClearChatError,
        );
      case _AlertLevel.streamNotice:
        return _StreamNoticeStrip(
          key: key,
          message: state.streamNotice,
          onDismiss: onDismissStreamNotice,
        );
      case _AlertLevel.appTopNotice:
        return AppTopNoticeBar(key: key);
    }
  }
}

class _ConnectionStatusStrip extends StatelessWidget {
  const _ConnectionStatusStrip({
    super.key,
    required this.state,
    this.onRetry,
    this.fixedStatus,
  });

  final ChatState state;
  final VoidCallback? onRetry;
  final ChatServerConnectionStatus? fixedStatus;

  @override
  Widget build(BuildContext context) {
    final status = fixedStatus ?? resolveChatServerConnectionStatus(state);
    if (status == ChatServerConnectionStatus.connected) {
      return const SizedBox.shrink();
    }
    
    return _buildBar(context, status);
  }

  Widget _buildBar(BuildContext context, ChatServerConnectionStatus status) {
    final colors = Theme.of(context).colorScheme;

    late Color backgroundColor;
    late Color foregroundColor;
    late IconData icon;
    late String text;
    var showSpinner = false;

    switch (status) {
      case ChatServerConnectionStatus.connecting:
        backgroundColor = colors.primaryContainer;
        foregroundColor = colors.onPrimaryContainer;
        icon = Icons.sync;
        text = 'Подключение к серверу...';
        showSpinner = true;
        break;
      case ChatServerConnectionStatus.disconnected:
        backgroundColor = colors.errorContainer;
        foregroundColor = colors.onErrorContainer;
        icon = Icons.cloud_off;
        text = 'Нет соединения с сервером';
        break;
      case ChatServerConnectionStatus.connected:
        return const SizedBox.shrink();
    }

    return _alertCard(
      context: context,
      color: backgroundColor,
      child: Padding(
        padding: const EdgeInsets.symmetric(vertical: 10, horizontal: 14),
        child: Row(
          children: [
            Icon(icon, size: 18, color: foregroundColor),
            const SizedBox(width: 10),
            Expanded(
              child: Text(
                text,
                style: TextStyle(
                  fontSize: 14,
                  fontWeight: FontWeight.w500,
                  color: foregroundColor,
                ),
              ),
            ),
            if (showSpinner) ...[
              const SizedBox(width: 10),
              SizedBox(
                width: 16,
                height: 16,
                child: CircularProgressIndicator(
                  strokeWidth: 2,
                  valueColor: AlwaysStoppedAnimation<Color>(foregroundColor),
                ),
              ),
            ],
            if (status == ChatServerConnectionStatus.disconnected && onRetry != null) ...[
              const SizedBox(width: 8),
              TextButton(
                onPressed: onRetry,
                style: TextButton.styleFrom(
                  foregroundColor: foregroundColor,
                  padding: const EdgeInsets.symmetric(horizontal: 12),
                ),
                child: const Text('Повторить'),
              ),
            ],
          ],
        ),
      ),
    );
  }
}

class _RunnersInactiveStrip extends StatelessWidget {
  const _RunnersInactiveStrip({
    super.key,
    required this.state,
    this.onRefresh,
  });

  final ChatState state;
  final VoidCallback? onRefresh;

  @override
  Widget build(BuildContext context) {
    final cs = Theme.of(context).colorScheme;
    final fg = cs.onErrorContainer;
    final isRefreshing = state.runnersStatusRefreshing;

    return _alertCard(
      context: context,
      color: cs.errorContainer.withValues(alpha: 0.55),
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
            else if (onRefresh != null)
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

class _ChatErrorStrip extends StatelessWidget {
  const _ChatErrorStrip({
    super.key,
    required this.message,
    this.onDismiss,
  });

  final String? message;
  final VoidCallback? onDismiss;

  @override
  Widget build(BuildContext context) {
    final text = message?.trim();
    if (text == null || text.isEmpty) {
      return const SizedBox.shrink();
    }

    final scheme = Theme.of(context).colorScheme;
    final bg = scheme.errorContainer;
    final fg = scheme.onErrorContainer;

    return _alertCard(
      context: context,
      color: bg,
      child: Padding(
        padding: const EdgeInsets.symmetric(vertical: 10, horizontal: 12),
        child: Row(
          crossAxisAlignment: CrossAxisAlignment.start,
          children: [
            Icon(Icons.error_outline, size: 20, color: fg),
            const SizedBox(width: 10),
            Expanded(
              child: Text(
                text,
                style: TextStyle(
                  fontSize: 14,
                  fontWeight: FontWeight.w500,
                  color: fg,
                  height: 1.35,
                ),
              ),
            ),
            if (onDismiss != null)
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

class _StreamNoticeStrip extends StatelessWidget {
  const _StreamNoticeStrip({
    super.key,
    required this.message,
    this.onDismiss,
  });

  final String? message;
  final VoidCallback? onDismiss;

  @override
  Widget build(BuildContext context) {
    final text = message?.trim();
    if (text == null || text.isEmpty) {
      return const SizedBox.shrink();
    }

    final scheme = Theme.of(context).colorScheme;
    final bg = scheme.secondaryContainer;
    final fg = scheme.onSecondaryContainer;

    return _alertCard(
      context: context,
      color: bg,
      child: Padding(
        padding: const EdgeInsets.symmetric(vertical: 10, horizontal: 12),
        child: Row(
          crossAxisAlignment: CrossAxisAlignment.start,
          children: [
            Icon(Icons.info_outline, size: 20, color: fg),
            const SizedBox(width: 10),
            Expanded(
              child: Text(
                text,
                style: TextStyle(
                  fontSize: 14,
                  fontWeight: FontWeight.w500,
                  color: fg,
                  height: 1.35,
                ),
              ),
            ),
            if (onDismiss != null)
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
