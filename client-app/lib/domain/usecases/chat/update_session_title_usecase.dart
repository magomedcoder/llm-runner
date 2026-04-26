import 'package:gen/domain/entities/session.dart';
import 'package:gen/domain/repositories/chat_repository.dart';

class UpdateSessionTitleUseCase {
  final ChatRepository repository;

  UpdateSessionTitleUseCase(this.repository);

  Future<ChatSession> call(int sessionId, String title) {
    return repository.updateSessionTitle(sessionId, title);
  }
}
