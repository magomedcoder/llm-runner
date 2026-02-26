import 'package:equatable/equatable.dart';

class RunnerInfo extends Equatable {
  final String address;
  final bool enabled;
  final bool connected;

  const RunnerInfo({
    required this.address,
    required this.enabled,
    this.connected = false,
  });

  @override
  List<Object?> get props => [address, enabled, connected];
}
