class McpServerEntity {
  const McpServerEntity({
    required this.id,
    required this.name,
    required this.enabled,
    required this.transport,
    required this.command,
    required this.args,
    required this.env,
    required this.url,
    required this.headers,
    required this.timeoutSeconds,
    this.ownerUserId = 0,
  });

  final int id;
  final String name;
  final bool enabled;
  final String transport;
  final String command;
  final List<String> args;
  final Map<String, String> env;
  final String url;
  final Map<String, String> headers;
  final int timeoutSeconds;
  final int ownerUserId;

  bool get isGlobal => ownerUserId == 0;

  McpServerEntity copyWith({
    int? id,
    String? name,
    bool? enabled,
    String? transport,
    String? command,
    List<String>? args,
    Map<String, String>? env,
    String? url,
    Map<String, String>? headers,
    int? timeoutSeconds,
    int? ownerUserId,
  }) {
    return McpServerEntity(
      id: id ?? this.id,
      name: name ?? this.name,
      enabled: enabled ?? this.enabled,
      transport: transport ?? this.transport,
      command: command ?? this.command,
      args: args ?? this.args,
      env: env ?? this.env,
      url: url ?? this.url,
      headers: headers ?? this.headers,
      timeoutSeconds: timeoutSeconds ?? this.timeoutSeconds,
      ownerUserId: ownerUserId ?? this.ownerUserId,
    );
  }
}
