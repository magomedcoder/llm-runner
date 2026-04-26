import 'package:gen/domain/repositories/chat_repository.dart';

class GetSelectedRunnerUseCase {
  final ChatRepository repository;

  GetSelectedRunnerUseCase(this.repository);

  Future<String?> call() => repository.getSelectedRunner();
}
