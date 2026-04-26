import 'package:gen/domain/entities/session_messages_page.dart';
import 'package:gen/domain/repositories/chat_repository.dart';

class GetSessionMessagesUseCase {
  final ChatRepository repository;

  GetSessionMessagesUseCase(this.repository);

  Future<SessionMessagesPage> call(
    int sessionId, {
    int beforeMessageId = 0,
    int pageSize = 40,
  }) {
    return repository.getSessionMessagesPage(
      sessionId: sessionId,
      beforeMessageId: beforeMessageId,
      pageSize: pageSize,
    );
  }
}
