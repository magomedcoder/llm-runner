import 'dart:typed_data';

import 'package:equatable/equatable.dart';

class SessionFileDownload extends Equatable {
  final String filename;
  final Uint8List content;

  const SessionFileDownload({
    required this.filename,
    required this.content,
  });

  @override
  List<Object?> get props => [filename, content];
}
