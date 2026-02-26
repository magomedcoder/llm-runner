import 'package:flutter/material.dart';
import 'package:flutter_bloc/flutter_bloc.dart';
import 'package:gen/core/injector.dart' as di;
import 'package:gen/presentation/screens/admin/bloc/runners_admin_bloc.dart';
import 'package:gen/presentation/screens/admin/bloc/runners_admin_event.dart';
import 'package:gen/presentation/screens/admin/bloc/runners_admin_state.dart';
import 'package:gen/domain/entities/runner_info.dart';

String _runnerStatusText(RunnerInfo runner) {
  if (!runner.enabled) return 'Отключён';
  if (runner.connected) return 'Подключён';
  return 'Ожидание подключения';
}

Color _runnerStatusColor(BuildContext context, RunnerInfo runner) {
  if (runner.connected) return Theme.of(context).colorScheme.primary;
  if (runner.enabled) return Theme.of(context).colorScheme.tertiary;
  return Theme.of(context).colorScheme.outline;
}

class RunnersAdminScreen extends StatefulWidget {
  const RunnersAdminScreen({super.key});

  @override
  State<RunnersAdminScreen> createState() => _RunnersAdminScreenState();
}

class _RunnersAdminScreenState extends State<RunnersAdminScreen> {
  late final RunnersAdminBloc _bloc;

  @override
  void initState() {
    super.initState();
    _bloc = di.sl<RunnersAdminBloc>()..add(const RunnersAdminLoadRequested());
  }

  @override
  void dispose() {
    _bloc.close();
    super.dispose();
  }

  @override
  Widget build(BuildContext context) {
    return BlocProvider.value(
      value: _bloc,
      child: Scaffold(
        appBar: AppBar(
          title: const Text('Раннеры'),
          actions: [
            IconButton(
              icon: const Icon(Icons.refresh),
              onPressed: () => _bloc.add(const RunnersAdminLoadRequested()),
            ),
          ],
        ),
        body: BlocConsumer<RunnersAdminBloc, RunnersAdminState>(
          listener: (context, state) {
            if (state.error != null) {
              ScaffoldMessenger.of(context).showSnackBar(
                SnackBar(
                  content: Text(state.error!),
                  backgroundColor: Theme.of(context).colorScheme.error,
                  behavior: SnackBarBehavior.floating,
                ),
              );
              _bloc.add(const RunnersAdminClearError());
            }
          },
          builder: (context, state) {
            if (state.isLoading && state.runners.isEmpty) {
              return const Center(child: CircularProgressIndicator());
            }
            if (state.runners.isEmpty) {
              return const Center(
                child: Text('Нет зарегистрированных раннеров'),
              );
            }
            return ListView.builder(
              padding: const EdgeInsets.all(16),
              itemCount: state.runners.length,
              itemBuilder: (context, index) {
                final r = state.runners[index];
                final isConnected = r.enabled && r.connected;
                return Card(
                  margin: const EdgeInsets.only(bottom: 8),
                  child: ListTile(
                    leading: Icon(
                      r.enabled ? Icons.link : Icons.link_off,
                      color: isConnected
                          ? Theme.of(context).colorScheme.primary
                          : Theme.of(context).colorScheme.outline,
                    ),
                    title: Text(
                      r.address,
                      style: const TextStyle(fontFamily: 'monospace'),
                    ),
                    subtitle: Text(
                      _runnerStatusText(r),
                      style: Theme.of(context).textTheme.bodySmall?.copyWith(
                            color: _runnerStatusColor(context, r),
                          ),
                    ),
                    trailing: Switch(
                      value: r.enabled,
                      onChanged: (value) {
                        _bloc.add(RunnersAdminSetEnabledRequested(
                          address: r.address,
                          enabled: value,
                        ));
                      },
                    ),
                  ),
                );
              },
            );
          },
        ),
      ),
    );
  }
}
