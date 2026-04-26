import 'package:gen/domain/entities/runner_info.dart';
import 'package:gen/domain/repositories/runners_repository.dart';

class GetUserRunnersUseCase {
  final RunnersRepository _repo;

  GetUserRunnersUseCase(this._repo);

  Future<List<RunnerInfo>> call() => _repo.getUserRunners();
}

