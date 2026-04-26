import 'package:gen/domain/entities/session_file_download.dart';
import 'package:gen/domain/repositories/chat_repository.dart';

class GetSessionFileUseCase {
  final ChatRepository _repository;

  GetSessionFileUseCase(this._repository);

  Future<SessionFileDownload> call({
    required int sessionId,
    required int fileId,
  }) {
    return _repository.getSessionFile(sessionId: sessionId, fileId: fileId);
  }
}
