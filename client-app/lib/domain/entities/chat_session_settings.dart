import 'package:equatable/equatable.dart';

class ChatSessionSettings extends Equatable {
  final int sessionId;
  final String systemPrompt;
  final List<String> stopSequences;
  final int timeoutSeconds;
  final double? temperature;
  final int? topK;
  final double? topP;
  final bool jsonMode;
  final String jsonSchema;
  final String toolsJson;
  final String profile;

  const ChatSessionSettings({
    required this.sessionId,
    this.systemPrompt = '',
    this.stopSequences = const [],
    this.timeoutSeconds = 0,
    this.temperature,
    this.topK,
    this.topP,
    this.jsonMode = false,
    this.jsonSchema = '',
    this.toolsJson = '',
    this.profile = '',
  });

  @override
  List<Object?> get props => [
    sessionId,
    systemPrompt,
    stopSequences,
    timeoutSeconds,
    temperature,
    topK,
    topP,
    jsonMode,
    jsonSchema,
    toolsJson,
    profile,
  ];
}
