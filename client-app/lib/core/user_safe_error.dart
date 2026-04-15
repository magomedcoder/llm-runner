import 'package:grpc/grpc.dart';
import 'package:gen/core/failures.dart';

String userSafeErrorMessage(
  Object? error, {
  String fallback = 'Произошла ошибка',
}) {
  if (error == null) {
    return fallback;
  }
  if (error is GrpcError) {
    final m = error.message?.trim();
    if (m != null && m.isNotEmpty) {
      switch (error.code) {
        case StatusCode.invalidArgument:
        case StatusCode.failedPrecondition:
        case StatusCode.resourceExhausted:
        case StatusCode.deadlineExceeded:
        case StatusCode.unavailable:
        case StatusCode.permissionDenied:
          return m;
        default:
          break;
      }
    }
    return 'Ошибка сервера (код ${error.code})';
  }
  if (error is UnauthorizedFailure) {
    return error.message;
  }
  if (error is Failure) {
    return error.message;
  }
  return fallback;
}
