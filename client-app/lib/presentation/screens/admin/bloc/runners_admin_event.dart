import 'package:equatable/equatable.dart';

sealed class RunnersAdminEvent extends Equatable {
  const RunnersAdminEvent();

  @override
  List<Object?> get props => [];
}

class RunnersAdminLoadRequested extends RunnersAdminEvent {
  const RunnersAdminLoadRequested();
}

class RunnersAdminSetEnabledRequested extends RunnersAdminEvent {
  final String address;
  final bool enabled;

  const RunnersAdminSetEnabledRequested({
    required this.address,
    required this.enabled,
  });

  @override
  List<Object?> get props => [address, enabled];
}

class RunnersAdminClearError extends RunnersAdminEvent {
  const RunnersAdminClearError();
}

class RunnersAdminDefaultRunnerChanged extends RunnersAdminEvent {
  final String? address;

  const RunnersAdminDefaultRunnerChanged(this.address);

  @override
  List<Object?> get props => [address];
}

class RunnersAdminDefaultModelChanged extends RunnersAdminEvent {
  final String runnerAddress;
  final String? model;

  const RunnersAdminDefaultModelChanged({
    required this.runnerAddress,
    required this.model,
  });

  @override
  List<Object?> get props => [runnerAddress, model];
}
