import 'package:equatable/equatable.dart';
import 'package:gen/domain/entities/message.dart';

class SessionMessagesPage extends Equatable {
  final List<Message> messages;
  final bool hasMoreOlder;

  const SessionMessagesPage({
    required this.messages,
    required this.hasMoreOlder,
  });

  @override
  List<Object?> get props => [messages, hasMoreOlder];
}
