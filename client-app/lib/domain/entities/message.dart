import 'dart:typed_data';

import 'package:equatable/equatable.dart';

enum MessageRole { user, assistant, tool }

class Message extends Equatable {
  final int id;
  final String content;
  final MessageRole role;
  final DateTime createdAt;
  final DateTime? updatedAt;
  final String? attachmentFileName;
  final List<String> attachmentFileNames;
  final String? attachmentMime;
  final Uint8List? attachmentContent;
  final int? attachmentFileId;
  final List<int> attachmentFileIds;
  final String? reasoningContent;
  final String? toolCallId;
  final String? toolName;
  final String? toolCallsJson;
  final bool useFileRag;
  final int fileRagTopK;
  final String fileRagEmbedModel;

  const Message({
    required this.id,
    required this.content,
    required this.role,
    required this.createdAt,
    this.updatedAt,
    this.attachmentFileName,
    this.attachmentFileNames = const [],
    this.attachmentMime,
    this.attachmentContent,
    this.attachmentFileId,
    this.attachmentFileIds = const [],
    this.reasoningContent,
    this.toolCallId,
    this.toolName,
    this.toolCallsJson,
    this.useFileRag = false,
    this.fileRagTopK = 0,
    this.fileRagEmbedModel = '',
  });

  Map<String, dynamic> toJson() => {
    'role': switch (role) {
      MessageRole.user => 'user',
      MessageRole.assistant => 'assistant',
      MessageRole.tool => 'tool',
    },
    'content': content,
  };

  @override
  List<Object?> get props => [
    id,
    content,
    role,
    createdAt,
    updatedAt,
    attachmentFileName,
    attachmentFileNames,
    attachmentMime,
    attachmentFileId,
    attachmentFileIds,
    reasoningContent,
    toolCallId,
    toolName,
    toolCallsJson,
    useFileRag,
    fileRagTopK,
    fileRagEmbedModel,
  ];
}
