import 'package:gen/domain/entities/chat_stream_chunk.dart';
import 'package:gen/domain/entities/message.dart';
import 'package:gen/domain/repositories/chat_repository.dart';

class SendMessageUseCase {
  final ChatRepository repository;

  SendMessageUseCase(this.repository);

  Stream<ChatStreamChunk> call(
    int sessionId,
    Message message,
  ) {
    return repository.sendMessage(sessionId, message);
  }
}
