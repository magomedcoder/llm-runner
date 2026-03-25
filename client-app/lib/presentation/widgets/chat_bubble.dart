import 'package:flutter/material.dart';
import 'package:flutter/services.dart';
import 'package:flutter_markdown_plus/flutter_markdown_plus.dart';
import 'package:gen/core/layout/responsive.dart';
import 'package:gen/domain/entities/message.dart';
import 'package:gen/presentation/widgets/code_block_builder.dart';

Color _messageBodyTextColor(ColorScheme cs) {
  return cs.onSurface.withValues(alpha: 0.94);
}

BorderRadius _bubbleRadius(bool isUser) {
  const r = 20.0;
  const tail = 8.0;
  return BorderRadius.only(
    topLeft: const Radius.circular(r),
    topRight: const Radius.circular(r),
    bottomLeft: Radius.circular(isUser ? r : tail),
    bottomRight: Radius.circular(isUser ? tail : r),
  );
}

class ChatBubble extends StatefulWidget {
  final Message message;
  final bool isStreaming;

  const ChatBubble({
    super.key,
    required this.message,
    this.isStreaming = false,
  });

  @override
  State<ChatBubble> createState() => _ChatBubbleState();
}

class _ChatBubbleState extends State<ChatBubble> {
  bool _justCopied = false;

  MarkdownStyleSheet _assistantMarkdownSheet(ThemeData theme) {
    final cs = theme.colorScheme;
    final onVar = _messageBodyTextColor(cs);
    final base = MarkdownStyleSheet.fromTheme(theme);
    return base.copyWith(
      p: base.p?.copyWith(color: onVar, fontSize: 15, height: 1.5),
      h1: base.h1?.copyWith(color: onVar, height: 1.25),
      h2: base.h2?.copyWith(color: onVar, height: 1.3),
      h3: base.h3?.copyWith(color: onVar, height: 1.35),
      h4: base.h4?.copyWith(color: onVar),
      h5: base.h5?.copyWith(color: onVar),
      h6: base.h6?.copyWith(color: onVar),
      em: base.em?.copyWith(color: onVar, fontStyle: FontStyle.italic),
      strong: base.strong?.copyWith(color: onVar, fontWeight: FontWeight.w600),
      a: base.a?.copyWith(color: onVar, decoration: TextDecoration.underline),
      code: base.code?.copyWith(
        color: onVar,
        fontFamily: 'monospace',
        fontSize: 13,
        backgroundColor: Colors.transparent,
      ),
      listIndent: 24,
      listBullet: base.listBullet?.copyWith(color: onVar),
      blockSpacing: 10,
      blockquote: base.blockquote?.copyWith(
        color: onVar.withValues(alpha: 0.92),
        height: 1.45,
      ),
      blockquoteDecoration: BoxDecoration(
        border: Border(
          left: BorderSide(color: onVar.withValues(alpha: 0.45), width: 4),
        ),
      ),
      blockquotePadding: const EdgeInsets.only(left: 12, top: 2, bottom: 2),
      codeblockPadding: const EdgeInsets.all(10),
      codeblockDecoration: BoxDecoration(
        color: Colors.transparent,
        borderRadius: BorderRadius.circular(8),
      ),
      horizontalRuleDecoration: BoxDecoration(
        border: Border(top: BorderSide(color: cs.outlineVariant, width: 1)),
      ),
      tableBorder: TableBorder.all(color: cs.outlineVariant, width: 1),
      tableCellsPadding: const EdgeInsets.symmetric(horizontal: 8, vertical: 6),
      tableHead: base.tableHead?.copyWith(
        color: onVar,
        fontWeight: FontWeight.w600,
      ),
      tableBody: base.tableBody?.copyWith(color: onVar),
    );
  }

  @override
  Widget build(BuildContext context) {
    final message = widget.message;
    final isStreaming = widget.isStreaming;
    final isUser = message.role == MessageRole.user;
    final theme = Theme.of(context);
    final width = Breakpoints.width(context);
    const minBubbleWidth = 64.0;
    final maxBubbleWidth = Breakpoints.isMobile(context)
      ? width * 0.88
      : (Breakpoints.isTablet(context) ? 420.0 : 560.0);
    final semanticsRole = isUser ? 'Ваше сообщение' : 'Ответ ассистента';
    final hasCopyableText = message.content.trim().isNotEmpty;
    final messageTextColor = _messageBodyTextColor(theme.colorScheme);

    return Semantics(
      container: true,
      label: semanticsRole,
      child: Align(
        alignment: isUser ? Alignment.centerRight : Alignment.centerLeft,
        child: Column(
          crossAxisAlignment: isUser ? CrossAxisAlignment.end : CrossAxisAlignment.start,
          mainAxisSize: MainAxisSize.min,
          children: [
            Container(
              margin: const EdgeInsets.symmetric(vertical: 2),
              padding: EdgeInsets.symmetric(
                horizontal: Breakpoints.isMobile(context) ? 12 : 16,
                vertical: Breakpoints.isMobile(context) ? 10 : 12,
              ),
              constraints: BoxConstraints(
                minWidth: minBubbleWidth,
                maxWidth: maxBubbleWidth,
              ),
              decoration: BoxDecoration(
                color: theme.colorScheme.surfaceContainerHigh,
                borderRadius: _bubbleRadius(isUser),
              ),
              child: Column(
                crossAxisAlignment: CrossAxisAlignment.start,
                mainAxisSize: MainAxisSize.min,
                children: [
                  if (message.attachmentFileName != null)
                    Padding(
                      padding: const EdgeInsets.only(bottom: 8),
                      child: Row(
                        mainAxisSize: MainAxisSize.min,
                        children: [
                          Icon(
                            Icons.insert_drive_file_rounded,
                            size: 18,
                            color: messageTextColor,
                          ),
                          const SizedBox(width: 6),
                          Flexible(
                            child: Text(
                              message.attachmentFileName!,
                              style: TextStyle(
                                fontSize: 13,
                                color: messageTextColor,
                              ),
                              overflow: TextOverflow.ellipsis,
                            ),
                          ),
                        ],
                      ),
                    ),
                  if (message.content.isNotEmpty)
                    isUser
                      ? SelectableText(
                          message.content,
                          style: TextStyle(
                            color: messageTextColor,
                            fontSize: 15,
                            height: 1.5,
                          ),
                        )
                      : MarkdownBody(
                          data: message.content,
                          selectable: true,
                          styleSheet: _assistantMarkdownSheet(theme),
                          builders: {
                            'pre': CodeBlockBuilder(
                              textStyle: TextStyle(
                                fontSize: 13,
                                fontFamily: 'monospace',
                                color: messageTextColor,
                              ),
                            ),
                          },
                        ),
                  if (isStreaming)
                    Padding(
                      padding: const EdgeInsets.only(top: 6),
                      child: Row(
                        mainAxisSize: MainAxisSize.min,
                        children: [
                          SizedBox(
                            width: 12,
                            height: 12,
                            child: CircularProgressIndicator(
                              strokeWidth: 1.5,
                              color: messageTextColor.withValues(alpha: 0.75),
                            ),
                          ),
                          const SizedBox(width: 8),
                          Text(
                            'Обрабатываю...',
                            style: TextStyle(
                              fontSize: 12,
                              height: 1.2,
                              color: messageTextColor.withValues(alpha: 0.75),
                            ),
                          ),
                        ],
                      ),
                    ),
                ],
              ),
            ),
            if (hasCopyableText || isStreaming)
              Padding(
                padding: const EdgeInsets.only(left: 4, right: 4, top: 2, bottom: 4),
                child: IconButton(
                  onPressed: hasCopyableText
                      ? () async {
                          await Clipboard.setData(
                            ClipboardData(text: message.content),
                          );

                          if (!mounted) {
                            return;
                          }

                          setState(() => _justCopied = true);

                          Future.delayed(const Duration(seconds: 2), () {
                            if (mounted) {
                              setState(() => _justCopied = false);
                            }
                          });
                        }
                      : null,
                  icon: Icon(
                    _justCopied ? Icons.check_rounded : Icons.copy_rounded,
                    size: 18,
                    color: theme.colorScheme.onSurfaceVariant.withValues(
                      alpha: hasCopyableText ? 1 : 0.4,
                    ),
                  ),
                  tooltip: _justCopied ? 'Скопировано' : 'Копировать',
                  padding: EdgeInsets.zero,
                  visualDensity: VisualDensity.compact,
                  constraints: const BoxConstraints(minWidth: 32, minHeight: 32),
                  style: IconButton.styleFrom(
                    foregroundColor: theme.colorScheme.onSurfaceVariant.withValues(
                      alpha: hasCopyableText ? 1 : 0.4,
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
