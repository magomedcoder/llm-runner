import 'package:gen/domain/repositories/chat_repository.dart';

class SetSelectedRunnerUseCase {
  final ChatRepository repository;

  SetSelectedRunnerUseCase(this.repository);

  Future<void> call(String? runner) => repository.setSelectedRunner(runner);
}
