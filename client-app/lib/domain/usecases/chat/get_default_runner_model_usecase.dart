import 'package:gen/domain/repositories/chat_repository.dart';

class GetDefaultRunnerModelUseCase {
  final ChatRepository repository;

  GetDefaultRunnerModelUseCase(this.repository);

  Future<String?> call(String runner) => repository.getDefaultRunnerModel(runner);
}
