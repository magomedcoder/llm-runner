// This is a generated file - do not edit.
//
// Generated from chat.proto.

// @dart = 3.3

// ignore_for_file: annotate_overrides, camel_case_types, comment_references
// ignore_for_file: constant_identifier_names
// ignore_for_file: curly_braces_in_flow_control_structures
// ignore_for_file: deprecated_member_use_from_same_package, library_prefixes
// ignore_for_file: non_constant_identifier_names, prefer_relative_imports
// ignore_for_file: unused_import

import 'dart:convert' as $convert;
import 'dart:core' as $core;
import 'dart:typed_data' as $typed_data;

@$core.Deprecated('Use emptyDescriptor instead')
const Empty$json = {
  '1': 'Empty',
};

/// Descriptor for `Empty`. Decode as a `google.protobuf.DescriptorProto`.
final $typed_data.Uint8List emptyDescriptor =
    $convert.base64Decode('CgVFbXB0eQ==');

@$core.Deprecated('Use connectionResponseDescriptor instead')
const ConnectionResponse$json = {
  '1': 'ConnectionResponse',
  '2': [
    {'1': 'is_connected', '3': 1, '4': 1, '5': 8, '10': 'isConnected'},
  ],
};

/// Descriptor for `ConnectionResponse`. Decode as a `google.protobuf.DescriptorProto`.
final $typed_data.Uint8List connectionResponseDescriptor = $convert.base64Decode(
    'ChJDb25uZWN0aW9uUmVzcG9uc2USIQoMaXNfY29ubmVjdGVkGAEgASgIUgtpc0Nvbm5lY3RlZA'
    '==');

@$core.Deprecated('Use chatMessageDescriptor instead')
const ChatMessage$json = {
  '1': 'ChatMessage',
  '2': [
    {'1': 'id', '3': 1, '4': 1, '5': 9, '10': 'id'},
    {'1': 'content', '3': 2, '4': 1, '5': 9, '10': 'content'},
    {'1': 'role', '3': 3, '4': 1, '5': 9, '10': 'role'},
    {'1': 'created_at', '3': 4, '4': 1, '5': 3, '10': 'createdAt'},
  ],
};

/// Descriptor for `ChatMessage`. Decode as a `google.protobuf.DescriptorProto`.
final $typed_data.Uint8List chatMessageDescriptor = $convert.base64Decode(
    'CgtDaGF0TWVzc2FnZRIOCgJpZBgBIAEoCVICaWQSGAoHY29udGVudBgCIAEoCVIHY29udGVudB'
    'ISCgRyb2xlGAMgASgJUgRyb2xlEh0KCmNyZWF0ZWRfYXQYBCABKANSCWNyZWF0ZWRBdA==');

@$core.Deprecated('Use chatResponseDescriptor instead')
const ChatResponse$json = {
  '1': 'ChatResponse',
  '2': [
    {'1': 'id', '3': 1, '4': 1, '5': 9, '10': 'id'},
    {'1': 'content', '3': 2, '4': 1, '5': 9, '10': 'content'},
    {'1': 'role', '3': 3, '4': 1, '5': 9, '10': 'role'},
    {'1': 'created_at', '3': 4, '4': 1, '5': 3, '10': 'createdAt'},
    {'1': 'done', '3': 5, '4': 1, '5': 8, '10': 'done'},
  ],
};

/// Descriptor for `ChatResponse`. Decode as a `google.protobuf.DescriptorProto`.
final $typed_data.Uint8List chatResponseDescriptor = $convert.base64Decode(
    'CgxDaGF0UmVzcG9uc2USDgoCaWQYASABKAlSAmlkEhgKB2NvbnRlbnQYAiABKAlSB2NvbnRlbn'
    'QSEgoEcm9sZRgDIAEoCVIEcm9sZRIdCgpjcmVhdGVkX2F0GAQgASgDUgljcmVhdGVkQXQSEgoE'
    'ZG9uZRgFIAEoCFIEZG9uZQ==');

@$core.Deprecated('Use sendMessageRequestDescriptor instead')
const SendMessageRequest$json = {
  '1': 'SendMessageRequest',
  '2': [
    {'1': 'session_id', '3': 1, '4': 1, '5': 9, '10': 'sessionId'},
    {
      '1': 'messages',
      '3': 2,
      '4': 3,
      '5': 11,
      '6': '.chat.ChatMessage',
      '10': 'messages'
    },
  ],
};

/// Descriptor for `SendMessageRequest`. Decode as a `google.protobuf.DescriptorProto`.
final $typed_data.Uint8List sendMessageRequestDescriptor = $convert.base64Decode(
    'ChJTZW5kTWVzc2FnZVJlcXVlc3QSHQoKc2Vzc2lvbl9pZBgBIAEoCVIJc2Vzc2lvbklkEi0KCG'
    '1lc3NhZ2VzGAIgAygLMhEuY2hhdC5DaGF0TWVzc2FnZVIIbWVzc2FnZXM=');

@$core.Deprecated('Use chatSessionDescriptor instead')
const ChatSession$json = {
  '1': 'ChatSession',
  '2': [
    {'1': 'id', '3': 1, '4': 1, '5': 9, '10': 'id'},
    {'1': 'title', '3': 2, '4': 1, '5': 9, '10': 'title'},
    {'1': 'created_at', '3': 3, '4': 1, '5': 3, '10': 'createdAt'},
    {'1': 'updated_at', '3': 4, '4': 1, '5': 3, '10': 'updatedAt'},
  ],
};

/// Descriptor for `ChatSession`. Decode as a `google.protobuf.DescriptorProto`.
final $typed_data.Uint8List chatSessionDescriptor = $convert.base64Decode(
    'CgtDaGF0U2Vzc2lvbhIOCgJpZBgBIAEoCVICaWQSFAoFdGl0bGUYAiABKAlSBXRpdGxlEh0KCm'
    'NyZWF0ZWRfYXQYAyABKANSCWNyZWF0ZWRBdBIdCgp1cGRhdGVkX2F0GAQgASgDUgl1cGRhdGVk'
    'QXQ=');

@$core.Deprecated('Use createSessionRequestDescriptor instead')
const CreateSessionRequest$json = {
  '1': 'CreateSessionRequest',
  '2': [
    {'1': 'title', '3': 1, '4': 1, '5': 9, '10': 'title'},
  ],
};

/// Descriptor for `CreateSessionRequest`. Decode as a `google.protobuf.DescriptorProto`.
final $typed_data.Uint8List createSessionRequestDescriptor =
    $convert.base64Decode(
        'ChRDcmVhdGVTZXNzaW9uUmVxdWVzdBIUCgV0aXRsZRgBIAEoCVIFdGl0bGU=');

@$core.Deprecated('Use getSessionRequestDescriptor instead')
const GetSessionRequest$json = {
  '1': 'GetSessionRequest',
  '2': [
    {'1': 'session_id', '3': 1, '4': 1, '5': 9, '10': 'sessionId'},
  ],
};

/// Descriptor for `GetSessionRequest`. Decode as a `google.protobuf.DescriptorProto`.
final $typed_data.Uint8List getSessionRequestDescriptor = $convert.base64Decode(
    'ChFHZXRTZXNzaW9uUmVxdWVzdBIdCgpzZXNzaW9uX2lkGAEgASgJUglzZXNzaW9uSWQ=');

@$core.Deprecated('Use listSessionsRequestDescriptor instead')
const ListSessionsRequest$json = {
  '1': 'ListSessionsRequest',
  '2': [
    {'1': 'page', '3': 1, '4': 1, '5': 5, '10': 'page'},
    {'1': 'page_size', '3': 2, '4': 1, '5': 5, '10': 'pageSize'},
  ],
};

/// Descriptor for `ListSessionsRequest`. Decode as a `google.protobuf.DescriptorProto`.
final $typed_data.Uint8List listSessionsRequestDescriptor = $convert.base64Decode(
    'ChNMaXN0U2Vzc2lvbnNSZXF1ZXN0EhIKBHBhZ2UYASABKAVSBHBhZ2USGwoJcGFnZV9zaXplGA'
    'IgASgFUghwYWdlU2l6ZQ==');

@$core.Deprecated('Use listSessionsResponseDescriptor instead')
const ListSessionsResponse$json = {
  '1': 'ListSessionsResponse',
  '2': [
    {
      '1': 'sessions',
      '3': 1,
      '4': 3,
      '5': 11,
      '6': '.chat.ChatSession',
      '10': 'sessions'
    },
    {'1': 'total', '3': 2, '4': 1, '5': 5, '10': 'total'},
    {'1': 'page', '3': 3, '4': 1, '5': 5, '10': 'page'},
    {'1': 'page_size', '3': 4, '4': 1, '5': 5, '10': 'pageSize'},
  ],
};

/// Descriptor for `ListSessionsResponse`. Decode as a `google.protobuf.DescriptorProto`.
final $typed_data.Uint8List listSessionsResponseDescriptor = $convert.base64Decode(
    'ChRMaXN0U2Vzc2lvbnNSZXNwb25zZRItCghzZXNzaW9ucxgBIAMoCzIRLmNoYXQuQ2hhdFNlc3'
    'Npb25SCHNlc3Npb25zEhQKBXRvdGFsGAIgASgFUgV0b3RhbBISCgRwYWdlGAMgASgFUgRwYWdl'
    'EhsKCXBhZ2Vfc2l6ZRgEIAEoBVIIcGFnZVNpemU=');

@$core.Deprecated('Use getSessionMessagesRequestDescriptor instead')
const GetSessionMessagesRequest$json = {
  '1': 'GetSessionMessagesRequest',
  '2': [
    {'1': 'session_id', '3': 1, '4': 1, '5': 9, '10': 'sessionId'},
    {'1': 'page', '3': 2, '4': 1, '5': 5, '10': 'page'},
    {'1': 'page_size', '3': 3, '4': 1, '5': 5, '10': 'pageSize'},
  ],
};

/// Descriptor for `GetSessionMessagesRequest`. Decode as a `google.protobuf.DescriptorProto`.
final $typed_data.Uint8List getSessionMessagesRequestDescriptor =
    $convert.base64Decode(
        'ChlHZXRTZXNzaW9uTWVzc2FnZXNSZXF1ZXN0Eh0KCnNlc3Npb25faWQYASABKAlSCXNlc3Npb2'
        '5JZBISCgRwYWdlGAIgASgFUgRwYWdlEhsKCXBhZ2Vfc2l6ZRgDIAEoBVIIcGFnZVNpemU=');

@$core.Deprecated('Use getSessionMessagesResponseDescriptor instead')
const GetSessionMessagesResponse$json = {
  '1': 'GetSessionMessagesResponse',
  '2': [
    {
      '1': 'messages',
      '3': 1,
      '4': 3,
      '5': 11,
      '6': '.chat.ChatMessage',
      '10': 'messages'
    },
    {'1': 'total', '3': 2, '4': 1, '5': 5, '10': 'total'},
    {'1': 'page', '3': 3, '4': 1, '5': 5, '10': 'page'},
    {'1': 'page_size', '3': 4, '4': 1, '5': 5, '10': 'pageSize'},
  ],
};

/// Descriptor for `GetSessionMessagesResponse`. Decode as a `google.protobuf.DescriptorProto`.
final $typed_data.Uint8List getSessionMessagesResponseDescriptor =
    $convert.base64Decode(
        'ChpHZXRTZXNzaW9uTWVzc2FnZXNSZXNwb25zZRItCghtZXNzYWdlcxgBIAMoCzIRLmNoYXQuQ2'
        'hhdE1lc3NhZ2VSCG1lc3NhZ2VzEhQKBXRvdGFsGAIgASgFUgV0b3RhbBISCgRwYWdlGAMgASgF'
        'UgRwYWdlEhsKCXBhZ2Vfc2l6ZRgEIAEoBVIIcGFnZVNpemU=');

@$core.Deprecated('Use deleteSessionRequestDescriptor instead')
const DeleteSessionRequest$json = {
  '1': 'DeleteSessionRequest',
  '2': [
    {'1': 'session_id', '3': 1, '4': 1, '5': 9, '10': 'sessionId'},
  ],
};

/// Descriptor for `DeleteSessionRequest`. Decode as a `google.protobuf.DescriptorProto`.
final $typed_data.Uint8List deleteSessionRequestDescriptor = $convert.base64Decode(
    'ChREZWxldGVTZXNzaW9uUmVxdWVzdBIdCgpzZXNzaW9uX2lkGAEgASgJUglzZXNzaW9uSWQ=');

@$core.Deprecated('Use updateSessionTitleRequestDescriptor instead')
const UpdateSessionTitleRequest$json = {
  '1': 'UpdateSessionTitleRequest',
  '2': [
    {'1': 'session_id', '3': 1, '4': 1, '5': 9, '10': 'sessionId'},
    {'1': 'title', '3': 2, '4': 1, '5': 9, '10': 'title'},
  ],
};

/// Descriptor for `UpdateSessionTitleRequest`. Decode as a `google.protobuf.DescriptorProto`.
final $typed_data.Uint8List updateSessionTitleRequestDescriptor =
    $convert.base64Decode(
        'ChlVcGRhdGVTZXNzaW9uVGl0bGVSZXF1ZXN0Eh0KCnNlc3Npb25faWQYASABKAlSCXNlc3Npb2'
        '5JZBIUCgV0aXRsZRgCIAEoCVIFdGl0bGU=');
