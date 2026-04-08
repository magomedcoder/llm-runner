class McpProbeResultEntity {
  const McpProbeResultEntity({
    required this.ok,
    required this.errorMessage,
    required this.protocolVersion,
    required this.serverName,
    required this.serverVersion,
    required this.instructions,
    required this.toolsSupported,
    required this.resourcesSupported,
    required this.promptsSupported,
  });

  final bool ok;
  final String errorMessage;
  final String protocolVersion;
  final String serverName;
  final String serverVersion;
  final String instructions;
  final bool toolsSupported;
  final bool resourcesSupported;
  final bool promptsSupported;
}
