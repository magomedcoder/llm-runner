import 'package:gen/domain/entities/auth_tokens.dart';
import 'package:gen/domain/repositories/auth_repository.dart';

class RefreshTokenUseCase {
  final AuthRepository repository;

  RefreshTokenUseCase(this.repository);

  Future<AuthTokens> call(String refreshToken) async {
    return await repository.refreshToken(refreshToken);
  }
}
