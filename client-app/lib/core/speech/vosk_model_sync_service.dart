import 'dart:io';

import 'package:archive/archive.dart';
import 'package:flutter/foundation.dart';
import 'package:gen/core/grpc_channel_manager.dart';
import 'package:gen/core/speech/local_vosk_dictation_service.dart';
import 'package:gen/data/data_sources/local/user_local_data_source.dart';
import 'package:gen/generated/grpc_pb/chat.pb.dart';
import 'package:grpc/grpc.dart';
import 'package:path/path.dart' as p;
import 'package:path_provider/path_provider.dart';

class VoskModelSyncService {
  VoskModelSyncService(
    this._channelManager,
    this._localDataSource,
    this._dictation,
  );

  final GrpcChannelManager _channelManager;
  final UserLocalDataSourceImpl _localDataSource;
  final LocalVoskDictationService _dictation;

  Future<Directory> _modelsParentDir() async {
    final root = await getApplicationSupportDirectory();
    final dir = Directory(p.join(root.path, 'vosk_models'));
    if (!dir.existsSync()) {
      await dir.create(recursive: true);
    }

    return dir;
  }

  bool _isValidModelDir(String path) {
    final conf = File(p.join(path, 'conf', 'model.conf'));
    return conf.existsSync();
  }

  String? _findLocalModelPath(Directory parent) {
    if (_isValidModelDir(parent.path)) {
      return parent.path;
    }

    final candidates = <String>[];
    for (final e in parent.listSync(followLinks: false)) {
      if (e is! Directory) {
        continue;
      }

      if (_isValidModelDir(e.path)) {
        candidates.add(e.path);
      }
    }

    candidates.sort();
    if (candidates.isNotEmpty) {
      return candidates.first;
    }

    return null;
  }

  bool _safeZipEntryName(String name, String parentCanonical) {
    if (name.isEmpty || name.contains('..')) {
      return false;
    }
    final resolved = p.normalize(p.join(parentCanonical, name));
    return p.isWithin(parentCanonical, resolved) || resolved == parentCanonical;
  }

  String? _archiveSharedTopLevelSegment(Archive archive) {
    String? top;
    for (final f in archive.files) {
      var name = f.name.replaceAll(r'\', '/');
      if (name.isEmpty) {
        continue;
      }

      name = name.replaceFirst(RegExp(r'/+$'), '');
      if (name.isEmpty) {
        continue;
      }

      final slash = name.indexOf('/');
      final head = slash < 0 ? name : name.substring(0, slash);
      if (head.isEmpty) {
        continue;
      }

      if (top == null) {
        top = head;
      } else if (top != head) {
        return null;
      }
    }

    return top;
  }

  String? _zipEntryRelativePath(String rawName, String? sharedTop) {
    var name = rawName.replaceAll(r'\', '/');
    if (name.isEmpty) {
      return null;
    }

    name = name.replaceFirst(RegExp(r'/+$'), '');
    if (name.isEmpty) {
      return '';
    }

    if (sharedTop != null) {
      if (name == sharedTop) {
        return '';
      }

      final prefix = '$sharedTop/';
      if (!name.startsWith(prefix)) {
        return null;
      }

      name = name.substring(prefix.length);
      if (name.isEmpty) {
        return '';
      }
    }

    if (name.contains('..')) {
      return null;
    }

    return name;
  }

  Future<void> prefetchIfLoggedIn() async {
    if (kIsWeb || !_dictation.isPlatformSupported) {
      return;
    }

    if (!_localDataSource.hasToken) {
      return;
    }

    try {
      await ensureModelPath();
    } catch (_) {}
  }

  Future<String?> resolveLocalModelPath() async {
    if (kIsWeb) {
      return null;
    }

    final saved = await _dictation.getSavedModelPath();
    if (saved != null && saved.isNotEmpty && _isValidModelDir(saved)) {
      return saved;
    }

    final parent = await _modelsParentDir();
    final local = _findLocalModelPath(parent);
    if (local != null) {
      await _dictation.saveModelPath(local);
      return local;
    }
    return null;
  }

  Future<bool> shouldDownloadFromServer() async {
    if (kIsWeb || !_dictation.isPlatformSupported) {
      return false;
    }

    if (await resolveLocalModelPath() != null) {
      return false;
    }

    return _localDataSource.hasToken;
  }

  Future<String?> ensureModelPath() async {
    if (kIsWeb) {
      return null;
    }

    final existing = await resolveLocalModelPath();
    if (existing != null) {
      return existing;
    }

    final parent = await _modelsParentDir();

    if (!_localDataSource.hasToken) {
      return null;
    }

    await _downloadAndExtract(parent);
    final after = _findLocalModelPath(parent);
    if (after == null) {
      throw StateError('После загрузки не найдена модель с conf/model.conf');
    }
    await _dictation.saveModelPath(after);
    return after;
  }

  Future<void> _downloadAndExtract(Directory parent) async {
    final client = _channelManager.chatClient;
    final acc = BytesAccumulator();
    final req = VoskModelDownloadRequest();
    try {
      await for (final chunk in client.downloadVoskModel(req)) {
        acc.add(chunk.data);
      }
    } on GrpcError catch (e) {
      if (e.code == StatusCode.notFound) {
        throw GrpcError.notFound('На сервере нет');
      }
      rethrow;
    }

    if (acc.isEmpty) {
      throw StateError('Пустой ответ при загрузке модели Vosk');
    }

    final archive = ZipDecoder().decodeBytes(acc.takeBytes(), verify: false);
    final sharedTop = _archiveSharedTopLevelSegment(archive);
    final parentCanon = p.canonicalize(parent.path);

    for (final file in archive.files) {
      final rel = _zipEntryRelativePath(file.name, sharedTop);
      if (rel == null) {
        continue;
      }

      if (rel.isEmpty) {
        continue;
      }

      if (!_safeZipEntryName(rel, parentCanon)) {
        continue;
      }

      final outPath = p.join(parent.path, rel);
      if (file.isFile) {
        final bytes = file.content;
        final f = File(outPath);
        await f.parent.create(recursive: true);
        await f.writeAsBytes(bytes);
      } else {
        await Directory(outPath).create(recursive: true);
      }
    }
  }
}

class BytesAccumulator {
  final _chunks = <Uint8List>[];
  int _len = 0;

  void add(List<int> data) {
    if (data.isEmpty) {
      return;
    }

    final u = data is Uint8List ? data : Uint8List.fromList(data);
    _chunks.add(u);
    _len += u.length;
  }

  bool get isEmpty => _len == 0;

  Uint8List takeBytes() {
    if (_chunks.length == 1) {
      return _chunks.single;
    }

    final out = Uint8List(_len);
    var o = 0;
    for (final c in _chunks) {
      out.setRange(o, o + c.length, c);
      o += c.length;
    }

    _chunks.clear();
    _len = 0;
    return out;
  }
}
