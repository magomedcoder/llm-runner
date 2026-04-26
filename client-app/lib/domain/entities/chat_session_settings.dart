import 'package:equatable/equatable.dart';

class ChatSessionSettings extends Equatable {
  final int sessionId;
  final String systemPrompt;
  final List<String> stopSequences;
  final int timeoutSeconds;
  final double? temperature;
  final int? topK;
  final double? topP;
  final String profile;
  final bool modelReasoningEnabled;
  final bool webSearchEnabled;
  final String webSearchProvider;
  final bool mcpEnabled;
  final List<int> mcpServerIds;

  const ChatSessionSettings({
    required this.sessionId,
    this.systemPrompt = '',
    this.stopSequences = const [],
    this.timeoutSeconds = 0,
    this.temperature,
    this.topK,
    this.topP,
    this.profile = '',
    this.modelReasoningEnabled = false,
    this.webSearchEnabled = false,
    this.webSearchProvider = '',
    this.mcpEnabled = false,
    this.mcpServerIds = const [],
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
    profile,
    modelReasoningEnabled,
    webSearchEnabled,
    webSearchProvider,
    mcpEnabled,
    mcpServerIds,
  ];
}
