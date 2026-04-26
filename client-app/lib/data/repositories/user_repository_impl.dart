import 'package:gen/core/failures.dart';
import 'package:gen/core/log/logs.dart';
import 'package:gen/data/data_sources/remote/user_remote_datasource.dart';
import 'package:gen/domain/entities/user.dart';
import 'package:gen/domain/repositories/user_repository.dart';

class UserRepositoryImpl implements UserRepository {
  final IUserRemoteDataSource dataSource;

  UserRepositoryImpl(this.dataSource);

  @override
  Future<List<User>> getUsers({required int page, required int pageSize}) async {
    try {
      return await dataSource.getUsers(page: page, pageSize: pageSize);
    } catch (e) {
      if (e is Failure) rethrow;
      Logs().e('UserRepository: ошибка получения пользователей', exception: e);
      throw ApiFailure('Ошибка получения пользователей');
    }
  }

  @override
  Future<User> createUser({
    required String username,
    required String password,
    required String name,
    required String surname,
    required int role,
  }) async {
    try {
      return await dataSource.createUser(
        username: username,
        password: password,
        name: name,
        surname: surname,
        role: role,
      );
    } catch (e) {
      if (e is Failure) rethrow;
      Logs().e('UserRepository: ошибка создания пользователя', exception: e);
      throw ApiFailure('Ошибка создания пользователя');
    }
  }

  @override
  Future<User> editUser({
    required String id,
    required String username,
    required String password,
    required String name,
    required String surname,
    required int role,
  }) async {
    try {
      return await dataSource.editUser(
        id: id,
        username: username,
        password: password,
        name: name,
        surname: surname,
        role: role,
      );
    } catch (e) {
      if (e is Failure) rethrow;
      Logs().e('UserRepository: ошибка обновления пользователя', exception: e);
      throw ApiFailure('Ошибка обновления пользователя');
    }
  }
}
