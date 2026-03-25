import 'package:equatable/equatable.dart';
import 'package:gen/domain/entities/gpu_info.dart';
import 'package:gen/domain/entities/loaded_model_status.dart';
import 'package:gen/domain/entities/server_info.dart';

class RunnerInfo extends Equatable {
  final String address;
  final String name;
  final bool enabled;
  final bool connected;
  final List<GpuInfo> gpus;
  final ServerInfo? serverInfo;
  final LoadedModelStatus? loadedModel;

  const RunnerInfo({
    required this.address,
    this.name = '',
    required this.enabled,
    this.connected = false,
    this.gpus = const [],
    this.serverInfo,
    this.loadedModel,
  });

  @override
  List<Object?> get props => [
    address,
    name,
    enabled,
    connected,
    gpus,
    serverInfo,
    loadedModel,
  ];
}
