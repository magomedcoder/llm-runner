import 'package:gen/domain/repositories/chat_repository.dart';

class SetDefaultRunnerModelUseCase {
  final ChatRepository repository;

  SetDefaultRunnerModelUseCase(this.repository);

  Future<void> call(String runner, String? model) =>
      repository.setDefaultRunnerModel(runner, model);
}
