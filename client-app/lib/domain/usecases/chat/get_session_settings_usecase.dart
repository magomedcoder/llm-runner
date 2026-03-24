import 'package:gen/domain/entities/chat_session_settings.dart';
import 'package:gen/domain/repositories/chat_repository.dart';

class GetSessionSettingsUseCase {
  final ChatRepository repository;

  GetSessionSettingsUseCase(this.repository);

  Future<ChatSessionSettings> call(int sessionId) => repository.getSessionSettings(sessionId);
}
