// This is a generated file - do not edit.
//
// Generated from chat.proto.

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

import 'chat.pb.dart' as $0;

export 'chat.pb.dart';

@$pb.GrpcServiceName('chat.ChatService')
class ChatServiceClient extends $grpc.Client {
  /// The hostname for this service.
  static const $core.String defaultHost = '';

  /// OAuth scopes needed for the client.
  static const $core.List<$core.String> oauthScopes = [
    '',
  ];

  ChatServiceClient(super.channel, {super.options, super.interceptors});

  $grpc.ResponseFuture<$0.ConnectionResponse> checkConnection(
    $0.Empty request, {
    $grpc.CallOptions? options,
  }) {
    return $createUnaryCall(_$checkConnection, request, options: options);
  }

  $grpc.ResponseStream<$0.ChatResponse> sendMessage(
    $0.SendMessageRequest request, {
    $grpc.CallOptions? options,
  }) {
    return $createStreamingCall(
        _$sendMessage, $async.Stream.fromIterable([request]),
        options: options);
  }

  $grpc.ResponseFuture<$0.ChatSession> createSession(
    $0.CreateSessionRequest request, {
    $grpc.CallOptions? options,
  }) {
    return $createUnaryCall(_$createSession, request, options: options);
  }

  $grpc.ResponseFuture<$0.ChatSession> getSession(
    $0.GetSessionRequest request, {
    $grpc.CallOptions? options,
  }) {
    return $createUnaryCall(_$getSession, request, options: options);
  }

  $grpc.ResponseFuture<$0.ListSessionsResponse> listSessions(
    $0.ListSessionsRequest request, {
    $grpc.CallOptions? options,
  }) {
    return $createUnaryCall(_$listSessions, request, options: options);
  }

  $grpc.ResponseFuture<$0.GetSessionMessagesResponse> getSessionMessages(
    $0.GetSessionMessagesRequest request, {
    $grpc.CallOptions? options,
  }) {
    return $createUnaryCall(_$getSessionMessages, request, options: options);
  }

  $grpc.ResponseFuture<$0.Empty> deleteSession(
    $0.DeleteSessionRequest request, {
    $grpc.CallOptions? options,
  }) {
    return $createUnaryCall(_$deleteSession, request, options: options);
  }

  $grpc.ResponseFuture<$0.ChatSession> updateSessionTitle(
    $0.UpdateSessionTitleRequest request, {
    $grpc.CallOptions? options,
  }) {
    return $createUnaryCall(_$updateSessionTitle, request, options: options);
  }

  // method descriptors

  static final _$checkConnection =
      $grpc.ClientMethod<$0.Empty, $0.ConnectionResponse>(
          '/chat.ChatService/CheckConnection',
          ($0.Empty value) => value.writeToBuffer(),
          $0.ConnectionResponse.fromBuffer);
  static final _$sendMessage =
      $grpc.ClientMethod<$0.SendMessageRequest, $0.ChatResponse>(
          '/chat.ChatService/SendMessage',
          ($0.SendMessageRequest value) => value.writeToBuffer(),
          $0.ChatResponse.fromBuffer);
  static final _$createSession =
      $grpc.ClientMethod<$0.CreateSessionRequest, $0.ChatSession>(
          '/chat.ChatService/CreateSession',
          ($0.CreateSessionRequest value) => value.writeToBuffer(),
          $0.ChatSession.fromBuffer);
  static final _$getSession =
      $grpc.ClientMethod<$0.GetSessionRequest, $0.ChatSession>(
          '/chat.ChatService/GetSession',
          ($0.GetSessionRequest value) => value.writeToBuffer(),
          $0.ChatSession.fromBuffer);
  static final _$listSessions =
      $grpc.ClientMethod<$0.ListSessionsRequest, $0.ListSessionsResponse>(
          '/chat.ChatService/ListSessions',
          ($0.ListSessionsRequest value) => value.writeToBuffer(),
          $0.ListSessionsResponse.fromBuffer);
  static final _$getSessionMessages = $grpc.ClientMethod<
          $0.GetSessionMessagesRequest, $0.GetSessionMessagesResponse>(
      '/chat.ChatService/GetSessionMessages',
      ($0.GetSessionMessagesRequest value) => value.writeToBuffer(),
      $0.GetSessionMessagesResponse.fromBuffer);
  static final _$deleteSession =
      $grpc.ClientMethod<$0.DeleteSessionRequest, $0.Empty>(
          '/chat.ChatService/DeleteSession',
          ($0.DeleteSessionRequest value) => value.writeToBuffer(),
          $0.Empty.fromBuffer);
  static final _$updateSessionTitle =
      $grpc.ClientMethod<$0.UpdateSessionTitleRequest, $0.ChatSession>(
          '/chat.ChatService/UpdateSessionTitle',
          ($0.UpdateSessionTitleRequest value) => value.writeToBuffer(),
          $0.ChatSession.fromBuffer);
}

@$pb.GrpcServiceName('chat.ChatService')
abstract class ChatServiceBase extends $grpc.Service {
  $core.String get $name => 'chat.ChatService';

  ChatServiceBase() {
    $addMethod($grpc.ServiceMethod<$0.Empty, $0.ConnectionResponse>(
        'CheckConnection',
        checkConnection_Pre,
        false,
        false,
        ($core.List<$core.int> value) => $0.Empty.fromBuffer(value),
        ($0.ConnectionResponse value) => value.writeToBuffer()));
    $addMethod($grpc.ServiceMethod<$0.SendMessageRequest, $0.ChatResponse>(
        'SendMessage',
        sendMessage_Pre,
        false,
        true,
        ($core.List<$core.int> value) =>
            $0.SendMessageRequest.fromBuffer(value),
        ($0.ChatResponse value) => value.writeToBuffer()));
    $addMethod($grpc.ServiceMethod<$0.CreateSessionRequest, $0.ChatSession>(
        'CreateSession',
        createSession_Pre,
        false,
        false,
        ($core.List<$core.int> value) =>
            $0.CreateSessionRequest.fromBuffer(value),
        ($0.ChatSession value) => value.writeToBuffer()));
    $addMethod($grpc.ServiceMethod<$0.GetSessionRequest, $0.ChatSession>(
        'GetSession',
        getSession_Pre,
        false,
        false,
        ($core.List<$core.int> value) => $0.GetSessionRequest.fromBuffer(value),
        ($0.ChatSession value) => value.writeToBuffer()));
    $addMethod(
        $grpc.ServiceMethod<$0.ListSessionsRequest, $0.ListSessionsResponse>(
            'ListSessions',
            listSessions_Pre,
            false,
            false,
            ($core.List<$core.int> value) =>
                $0.ListSessionsRequest.fromBuffer(value),
            ($0.ListSessionsResponse value) => value.writeToBuffer()));
    $addMethod($grpc.ServiceMethod<$0.GetSessionMessagesRequest,
            $0.GetSessionMessagesResponse>(
        'GetSessionMessages',
        getSessionMessages_Pre,
        false,
        false,
        ($core.List<$core.int> value) =>
            $0.GetSessionMessagesRequest.fromBuffer(value),
        ($0.GetSessionMessagesResponse value) => value.writeToBuffer()));
    $addMethod($grpc.ServiceMethod<$0.DeleteSessionRequest, $0.Empty>(
        'DeleteSession',
        deleteSession_Pre,
        false,
        false,
        ($core.List<$core.int> value) =>
            $0.DeleteSessionRequest.fromBuffer(value),
        ($0.Empty value) => value.writeToBuffer()));
    $addMethod(
        $grpc.ServiceMethod<$0.UpdateSessionTitleRequest, $0.ChatSession>(
            'UpdateSessionTitle',
            updateSessionTitle_Pre,
            false,
            false,
            ($core.List<$core.int> value) =>
                $0.UpdateSessionTitleRequest.fromBuffer(value),
            ($0.ChatSession value) => value.writeToBuffer()));
  }

  $async.Future<$0.ConnectionResponse> checkConnection_Pre(
      $grpc.ServiceCall $call, $async.Future<$0.Empty> $request) async {
    return checkConnection($call, await $request);
  }

  $async.Future<$0.ConnectionResponse> checkConnection(
      $grpc.ServiceCall call, $0.Empty request);

  $async.Stream<$0.ChatResponse> sendMessage_Pre($grpc.ServiceCall $call,
      $async.Future<$0.SendMessageRequest> $request) async* {
    yield* sendMessage($call, await $request);
  }

  $async.Stream<$0.ChatResponse> sendMessage(
      $grpc.ServiceCall call, $0.SendMessageRequest request);

  $async.Future<$0.ChatSession> createSession_Pre($grpc.ServiceCall $call,
      $async.Future<$0.CreateSessionRequest> $request) async {
    return createSession($call, await $request);
  }

  $async.Future<$0.ChatSession> createSession(
      $grpc.ServiceCall call, $0.CreateSessionRequest request);

  $async.Future<$0.ChatSession> getSession_Pre($grpc.ServiceCall $call,
      $async.Future<$0.GetSessionRequest> $request) async {
    return getSession($call, await $request);
  }

  $async.Future<$0.ChatSession> getSession(
      $grpc.ServiceCall call, $0.GetSessionRequest request);

  $async.Future<$0.ListSessionsResponse> listSessions_Pre(
      $grpc.ServiceCall $call,
      $async.Future<$0.ListSessionsRequest> $request) async {
    return listSessions($call, await $request);
  }

  $async.Future<$0.ListSessionsResponse> listSessions(
      $grpc.ServiceCall call, $0.ListSessionsRequest request);

  $async.Future<$0.GetSessionMessagesResponse> getSessionMessages_Pre(
      $grpc.ServiceCall $call,
      $async.Future<$0.GetSessionMessagesRequest> $request) async {
    return getSessionMessages($call, await $request);
  }

  $async.Future<$0.GetSessionMessagesResponse> getSessionMessages(
      $grpc.ServiceCall call, $0.GetSessionMessagesRequest request);

  $async.Future<$0.Empty> deleteSession_Pre($grpc.ServiceCall $call,
      $async.Future<$0.DeleteSessionRequest> $request) async {
    return deleteSession($call, await $request);
  }

  $async.Future<$0.Empty> deleteSession(
      $grpc.ServiceCall call, $0.DeleteSessionRequest request);

  $async.Future<$0.ChatSession> updateSessionTitle_Pre($grpc.ServiceCall $call,
      $async.Future<$0.UpdateSessionTitleRequest> $request) async {
    return updateSessionTitle($call, await $request);
  }

  $async.Future<$0.ChatSession> updateSessionTitle(
      $grpc.ServiceCall call, $0.UpdateSessionTitleRequest request);
}
