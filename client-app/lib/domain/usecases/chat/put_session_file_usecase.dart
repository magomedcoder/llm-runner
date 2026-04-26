import 'package:gen/domain/repositories/chat_repository.dart';

class PutSessionFileUseCase {
  final ChatRepository _repository;

  PutSessionFileUseCase(this._repository);

  Future<int> call({
    required int sessionId,
    required String filename,
    required List<int> content,
    int ttlSeconds = 0,
  }) {
    return _repository.putSessionFile(
      sessionId: sessionId,
      filename: filename,
      content: content,
      ttlSeconds: ttlSeconds,
    );
  }
}
