import 'package:equatable/equatable.dart';

enum ChatStreamChunkKind { text, toolStatus, notice, reasoning, ragMeta }

class ChatStreamChunk extends Equatable {
  final ChatStreamChunkKind kind;
  final String text;
  final String? toolName;
  final String? ragMode;
  final String? ragSourcesJson;
  final int messageId;

  const ChatStreamChunk({
    required this.kind,
    required this.text,
    this.toolName,
    this.ragMode,
    this.ragSourcesJson,
    this.messageId = 0,
  });

  @override
  List<Object?> get props => [
    kind,
    text,
    toolName,
    ragMode,
    ragSourcesJson,
    messageId,
  ];
}
