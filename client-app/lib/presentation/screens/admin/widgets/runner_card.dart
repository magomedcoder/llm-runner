import 'package:flutter/material.dart';
import 'package:gen/domain/entities/runner_info.dart';
import 'package:gen/presentation/screens/admin/widgets/runner_card_header.dart';
import 'package:gen/presentation/screens/admin/widgets/runner_gpu_section.dart';
import 'package:gen/presentation/screens/admin/widgets/runner_loaded_model_section.dart';
import 'package:gen/presentation/screens/admin/widgets/runner_server_info_section.dart';
import 'package:gen/presentation/screens/admin/widgets/runner_status.dart';

class RunnerCard extends StatelessWidget {
  final RunnerInfo runner;
  final VoidCallback onToggleEnabled;
  final String? defaultModel;
  final ValueChanged<String>? onDefaultModelChanged;

  const RunnerCard({
    super.key,
    required this.runner,
    required this.onToggleEnabled,
    this.defaultModel,
    this.onDefaultModelChanged,
  });

  @override
  Widget build(BuildContext context) {
    final status = runnerStatusFromRunner(runner);
    final hasServerInfo = runner.serverInfo != null;
    final hasGpus = runner.gpus.isNotEmpty;
    final hasLoadedModel = runner.loadedModel != null;

    return Card(
      clipBehavior: Clip.antiAlias,
      child: Padding(
        padding: const EdgeInsets.all(14),
        child: Column(
          crossAxisAlignment: CrossAxisAlignment.stretch,
          mainAxisSize: MainAxisSize.min,
          children: [
            RunnerCardHeader(
              runner: runner,
              status: status,
              onToggleEnabled: onToggleEnabled,
            ),
            if (hasLoadedModel) ...[
              const SizedBox(height: 12),
              const Divider(height: 1),
              const SizedBox(height: 8),
              RunnerLoadedModelSection(status: runner.loadedModel!),
            ],
            if (hasServerInfo) ...[
              const SizedBox(height: 12),
              const Divider(height: 1),
              const SizedBox(height: 8),
              RunnerServerInfoSection(
                serverInfo: runner.serverInfo!,
                defaultModel: defaultModel,
                onDefaultModelChanged: onDefaultModelChanged,
              ),
            ],
            if (hasGpus) ...[
              const SizedBox(height: 12),
              const Divider(height: 1),
              const SizedBox(height: 8),
              RunnerGpuSection(gpus: runner.gpus),
            ],
          ],
        ),
      ),
    );
  }
}
