import 'package:gen/domain/repositories/chat_repository.dart';

class DeleteSessionUseCase {
  final ChatRepository repository;

  DeleteSessionUseCase(this.repository);

  Future<void> call(int sessionId) {
    return repository.deleteSession(sessionId);
  }
}
