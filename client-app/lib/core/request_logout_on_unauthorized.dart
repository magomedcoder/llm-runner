import 'package:gen/core/failures.dart';
import 'package:gen/core/log/logs.dart';
import 'package:gen/presentation/screens/auth/bloc/auth_bloc.dart';
import 'package:gen/presentation/screens/auth/bloc/auth_event.dart';

void requestLogoutIfUnauthorized(Object e, AuthBloc authBloc) {
  if (e is UnauthorizedFailure) {
    Logs().w('Требуется выход: не авторизован');
    authBloc.add(const AuthLogoutRequested());
  }
}
