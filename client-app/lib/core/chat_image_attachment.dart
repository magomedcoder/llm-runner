import 'package:gen/domain/entities/message.dart';

String? guessImageMimeFromFilename(String? filename) {
  final n = filename?.trim().toLowerCase() ?? '';
  if (n.isEmpty) {
    return null;
  }

  if (n.endsWith('.png')) {
    return 'image/png';
  }

  if (n.endsWith('.jpg') || n.endsWith('.jpeg')) {
    return 'image/jpeg';
  }

  if (n.endsWith('.webp')) {
    return 'image/webp';
  }

  if (n.endsWith('.gif')) {
    return 'image/gif';
  }

  return null;
}

bool messageEligibleForChatImageThumb(Message m) {
  if (m.role != MessageRole.user) {
    return false;
  }

  final mime = m.attachmentMime?.trim().toLowerCase() ?? '';
  if (mime.startsWith('image/')) {
    return true;
  }

  if (guessImageMimeFromFilename(m.attachmentFileName) != null) {
    return true;
  }

  return false;
}

bool messageHasImageBytesOrFileRef(Message m) {
  final hasBytes = m.attachmentContent != null && m.attachmentContent!.isNotEmpty;
  final fid = m.attachmentFileId;
  final hasId = fid != null && fid > 0;
  return hasBytes || hasId;
}
