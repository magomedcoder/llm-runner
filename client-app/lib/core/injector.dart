import 'package:gen/core/auth_interceptor.dart';
import 'package:gen/core/token_storage.dart';
import 'package:gen/data/data_sources/remote/auth_remote_datasource.dart';
import 'package:gen/data/data_sources/remote/chat_remote_datasource.dart';
import 'package:gen/data/repositories/auth_repository_impl.dart';
import 'package:gen/data/repositories/chat_repository_impl.dart';
import 'package:gen/domain/repositories/auth_repository.dart';
import 'package:gen/domain/repositories/chat_repository.dart';
import 'package:gen/domain/usecases/auth/login_usecase.dart';
import 'package:gen/domain/usecases/auth/logout_usecase.dart';
import 'package:gen/domain/usecases/auth/refresh_token_usecase.dart';
import 'package:gen/domain/usecases/chat/connect_usecase.dart';
import 'package:gen/domain/usecases/chat/create_session_usecase.dart';
import 'package:gen/domain/usecases/chat/delete_session_usecase.dart';
import 'package:gen/domain/usecases/chat/get_session_messages_usecase.dart';
import 'package:gen/domain/usecases/chat/get_sessions_usecase.dart';
import 'package:gen/domain/usecases/chat/send_message_usecase.dart';
import 'package:gen/domain/usecases/chat/update_session_title_usecase.dart';
import 'package:gen/generated/grpc_pb/auth.pbgrpc.dart';
import 'package:gen/generated/grpc_pb/chat.pbgrpc.dart';
import 'package:gen/presentation/screens/auth/bloc/auth_bloc.dart';
import 'package:gen/presentation/screens/chat/bloc/chat_bloc.dart';
import 'package:get_it/get_it.dart';
import 'package:grpc/grpc.dart';

final sl = GetIt.instance;

Future<void> init() async {
  sl.registerLazySingleton<TokenStorage>(() => TokenStorage());

  sl.registerLazySingleton<AuthInterceptor>(
    () => AuthInterceptor(sl<TokenStorage>()),
  );

  sl.registerLazySingleton<ClientChannel>(() {
    return ClientChannel(
      '127.0.0.1',
      port: 50051,
      options: const ChannelOptions(
        credentials: ChannelCredentials.insecure(),
        idleTimeout: Duration(seconds: 30),
      ),
    );
  });

  sl.registerLazySingleton(() => ChatServiceClient(
        sl<ClientChannel>(),
        interceptors: [sl<AuthInterceptor>()],
      ));
  sl.registerLazySingleton(() => AuthServiceClient(
        sl<ClientChannel>(),
        interceptors: [sl<AuthInterceptor>()],
      ));

  sl.registerLazySingleton<IChatRemoteDataSource>(
    () => ChatRemoteDataSource(sl()),
  );

  sl.registerLazySingleton<IAuthRemoteDataSource>(
    () => AuthRemoteDataSource(sl()),
  );

  sl.registerLazySingleton<ChatRepository>(() => ChatRepositoryImpl(sl()));
  sl.registerLazySingleton<AuthRepository>(() => AuthRepositoryImpl(sl()));

  sl.registerFactory(() => ConnectUseCase(sl()));
  sl.registerFactory(() => SendMessageUseCase(sl()));
  sl.registerFactory(() => CreateSessionUseCase(sl()));
  sl.registerFactory(() => GetSessionsUseCase(sl()));
  sl.registerFactory(() => GetSessionMessagesUseCase(sl()));
  sl.registerFactory(() => DeleteSessionUseCase(sl()));
  sl.registerFactory(() => UpdateSessionTitleUseCase(sl()));

  sl.registerFactory(() => LoginUseCase(sl()));
  sl.registerFactory(() => RefreshTokenUseCase(sl()));
  sl.registerFactory(() => LogoutUseCase(sl()));

  sl.registerFactory(
    () => ChatBloc(
      connectUseCase: sl(),
      sendMessageUseCase: sl(),
      createSessionUseCase: sl(),
      getSessionsUseCase: sl(),
      getSessionMessagesUseCase: sl(),
      deleteSessionUseCase: sl(),
      updateSessionTitleUseCase: sl(),
    ),
  );

  sl.registerFactory(
    () => AuthBloc(
      loginUseCase: sl(),
      refreshTokenUseCase: sl(),
      logoutUseCase: sl(),
      tokenStorage: sl(),
    ),
  );
}
