import 'package:gen/domain/entities/auth_result.dart';
import 'package:gen/domain/entities/auth_tokens.dart';

abstract interface class AuthRepository {
  Future<AuthResult> login(String email, String password);

  Future<AuthTokens> refreshToken(String refreshToken);

  Future<void> logout();
}
