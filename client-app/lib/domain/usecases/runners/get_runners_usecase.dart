import 'package:gen/domain/entities/runner_info.dart';
import 'package:gen/domain/repositories/runners_repository.dart';

class GetRunnersUseCase {
  final RunnersRepository _repo;

  GetRunnersUseCase(this._repo);

  Future<List<RunnerInfo>> call() => _repo.getRunners();
}
