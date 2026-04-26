import 'package:gen/domain/entities/chat_session_settings.dart';
import 'package:gen/domain/repositories/chat_repository.dart';

class UpdateSessionSettingsUseCase {
  final ChatRepository repository;

  UpdateSessionSettingsUseCase(this.repository);

  Future<ChatSessionSettings> call({
    required int sessionId,
    required String systemPrompt,
    required List<String> stopSequences,
    required int timeoutSeconds,
    double? temperature,
    int? topK,
    double? topP,
    required String profile,
    required bool modelReasoningEnabled,
    required bool webSearchEnabled,
    required String webSearchProvider,
    required bool mcpEnabled,
    required List<int> mcpServerIds,
  }) {
    return repository.updateSessionSettings(
      sessionId: sessionId,
      systemPrompt: systemPrompt,
      stopSequences: stopSequences,
      timeoutSeconds: timeoutSeconds,
      temperature: temperature,
      topK: topK,
      topP: topP,
      profile: profile,
      modelReasoningEnabled: modelReasoningEnabled,
      webSearchEnabled: webSearchEnabled,
      webSearchProvider: webSearchProvider,
      mcpEnabled: mcpEnabled,
      mcpServerIds: mcpServerIds,
    );
  }
}
