import 'package:gen/domain/entities/auth_result.dart';
import 'package:gen/domain/repositories/auth_repository.dart';

class LoginUseCase {
  final AuthRepository repository;

  LoginUseCase(this.repository);

  Future<AuthResult> call(String email, String password) async {
    return await repository.login(email, password);
  }
}
