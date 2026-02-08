// This is a generated file - do not edit.
//
// Generated from auth.proto.

// @dart = 3.3

// ignore_for_file: annotate_overrides, camel_case_types, comment_references
// ignore_for_file: constant_identifier_names
// ignore_for_file: curly_braces_in_flow_control_structures
// ignore_for_file: deprecated_member_use_from_same_package, library_prefixes
// ignore_for_file: non_constant_identifier_names, prefer_relative_imports

import 'dart:async' as $async;
import 'dart:core' as $core;

import 'package:grpc/service_api.dart' as $grpc;
import 'package:protobuf/protobuf.dart' as $pb;

import 'auth.pb.dart' as $0;

export 'auth.pb.dart';

@$pb.GrpcServiceName('auth.AuthService')
class AuthServiceClient extends $grpc.Client {
  /// The hostname for this service.
  static const $core.String defaultHost = '';

  /// OAuth scopes needed for the client.
  static const $core.List<$core.String> oauthScopes = [
    '',
  ];

  AuthServiceClient(super.channel, {super.options, super.interceptors});

  $grpc.ResponseFuture<$0.LoginResponse> login(
    $0.LoginRequest request, {
    $grpc.CallOptions? options,
  }) {
    return $createUnaryCall(_$login, request, options: options);
  }

  $grpc.ResponseFuture<$0.RefreshTokenResponse> refreshToken(
    $0.RefreshTokenRequest request, {
    $grpc.CallOptions? options,
  }) {
    return $createUnaryCall(_$refreshToken, request, options: options);
  }

  $grpc.ResponseFuture<$0.LogoutResponse> logout(
    $0.LogoutRequest request, {
    $grpc.CallOptions? options,
  }) {
    return $createUnaryCall(_$logout, request, options: options);
  }

  // method descriptors

  static final _$login = $grpc.ClientMethod<$0.LoginRequest, $0.LoginResponse>(
      '/auth.AuthService/Login',
      ($0.LoginRequest value) => value.writeToBuffer(),
      $0.LoginResponse.fromBuffer);
  static final _$refreshToken =
      $grpc.ClientMethod<$0.RefreshTokenRequest, $0.RefreshTokenResponse>(
          '/auth.AuthService/RefreshToken',
          ($0.RefreshTokenRequest value) => value.writeToBuffer(),
          $0.RefreshTokenResponse.fromBuffer);
  static final _$logout =
      $grpc.ClientMethod<$0.LogoutRequest, $0.LogoutResponse>(
          '/auth.AuthService/Logout',
          ($0.LogoutRequest value) => value.writeToBuffer(),
          $0.LogoutResponse.fromBuffer);
}

@$pb.GrpcServiceName('auth.AuthService')
abstract class AuthServiceBase extends $grpc.Service {
  $core.String get $name => 'auth.AuthService';

  AuthServiceBase() {
    $addMethod($grpc.ServiceMethod<$0.LoginRequest, $0.LoginResponse>(
        'Login',
        login_Pre,
        false,
        false,
        ($core.List<$core.int> value) => $0.LoginRequest.fromBuffer(value),
        ($0.LoginResponse value) => value.writeToBuffer()));
    $addMethod(
        $grpc.ServiceMethod<$0.RefreshTokenRequest, $0.RefreshTokenResponse>(
            'RefreshToken',
            refreshToken_Pre,
            false,
            false,
            ($core.List<$core.int> value) =>
                $0.RefreshTokenRequest.fromBuffer(value),
            ($0.RefreshTokenResponse value) => value.writeToBuffer()));
    $addMethod($grpc.ServiceMethod<$0.LogoutRequest, $0.LogoutResponse>(
        'Logout',
        logout_Pre,
        false,
        false,
        ($core.List<$core.int> value) => $0.LogoutRequest.fromBuffer(value),
        ($0.LogoutResponse value) => value.writeToBuffer()));
  }

  $async.Future<$0.LoginResponse> login_Pre(
      $grpc.ServiceCall $call, $async.Future<$0.LoginRequest> $request) async {
    return login($call, await $request);
  }

  $async.Future<$0.LoginResponse> login(
      $grpc.ServiceCall call, $0.LoginRequest request);

  $async.Future<$0.RefreshTokenResponse> refreshToken_Pre(
      $grpc.ServiceCall $call,
      $async.Future<$0.RefreshTokenRequest> $request) async {
    return refreshToken($call, await $request);
  }

  $async.Future<$0.RefreshTokenResponse> refreshToken(
      $grpc.ServiceCall call, $0.RefreshTokenRequest request);

  $async.Future<$0.LogoutResponse> logout_Pre(
      $grpc.ServiceCall $call, $async.Future<$0.LogoutRequest> $request) async {
    return logout($call, await $request);
  }

  $async.Future<$0.LogoutResponse> logout(
      $grpc.ServiceCall call, $0.LogoutRequest request);
}
