import 'package:gen/domain/entities/user.dart';
import 'package:gen/domain/repositories/user_repository.dart';

class GetUsersUseCase {
  final UserRepository repo;

  GetUsersUseCase(this.repo);

  Future<List<User>> call({required int page, required int pageSize}) {
    return repo.getUsers(page: page, pageSize: pageSize);
  }
}
