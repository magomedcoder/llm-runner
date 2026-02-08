import 'package:gen/domain/entities/session.dart';
import 'package:gen/domain/repositories/chat_repository.dart';

class CreateSessionUseCase {
  final ChatRepository repository;

  CreateSessionUseCase(this.repository);

  Future<ChatSession> call({String? title}) async {
    final sessionTitle = title ?? 'Чат от ${DateTime.now().toString()}';
    return await repository.createSession(sessionTitle);
  }
}
