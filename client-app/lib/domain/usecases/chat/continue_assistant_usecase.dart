import 'package:gen/domain/entities/chat_stream_chunk.dart';
import 'package:gen/domain/repositories/chat_repository.dart';

class ContinueAssistantUseCase {
  final ChatRepository repository;

  ContinueAssistantUseCase(this.repository);

  Stream<ChatStreamChunk> call(int sessionId, int assistantMessageId) {
    return repository.continueAssistantResponse(sessionId, assistantMessageId);
  }
}
