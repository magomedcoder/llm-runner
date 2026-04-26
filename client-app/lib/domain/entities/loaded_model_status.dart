import 'package:equatable/equatable.dart';

class LoadedModelStatus extends Equatable {
  final bool loaded;
  final String displayName;
  final String ggufBasename;

  const LoadedModelStatus({
    required this.loaded,
    this.displayName = '',
    this.ggufBasename = '',
  });

  @override
  List<Object?> get props => [loaded, displayName, ggufBasename];
}
