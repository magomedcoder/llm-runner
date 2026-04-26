import 'package:gen/core/server_config.dart';

({String host, int port})? parseServerAddress(String raw) {
  final s = raw.trim();
  if (s.isEmpty) {
    return null;
  }

  if (s.startsWith('[')) {
    final close = s.indexOf(']');
    if (close <= 0) {
      return null;
    }

    final host = s.substring(1, close);
    if (host.isEmpty) {
      return null;
    }

    if (close == s.length - 1) {
      return (host: host, port: ServerConfig.defaultPort);
    }

    if (close + 1 >= s.length || s[close + 1] != ':') {
      return null;
    }

    final port = int.tryParse(s.substring(close + 2));
    if (port == null || port < 1 || port > 65535) {
      return null;
    }

    return (host: host, port: port);
  }

  final uri = Uri.tryParse(s);
  if (uri != null && uri.hasScheme && uri.host.isNotEmpty) {
    return (
      host: uri.host,
      port: uri.hasPort ? uri.port : ServerConfig.defaultPort,
    );
  }

  final lastColon = s.lastIndexOf(':');
  if (lastColon > 0) {
    final after = s.substring(lastColon + 1);
    final port = int.tryParse(after);
    if (port != null && port >= 1 && port <= 65535) {
      final before = s.substring(0, lastColon);
      if (!before.contains(':') && before.isNotEmpty) {
        return (host: before, port: port);
      }
    }
  }

  return (host: s, port: ServerConfig.defaultPort);
}

String formatServerAddressForField(ServerConfig config) {
  if (config.host.isEmpty) {
    return '';
  }

  if (config.port == ServerConfig.defaultPort) {
    return config.host;
  }

  return '${config.host}:${config.port}';
}

String get serverAddressInputHint => 'хост или хост:порт';

String? validateServerAddressInput(String? value) {
  if (value == null || value.trim().isEmpty) {
    return 'Введите хост';
  }

  if (parseServerAddress(value) == null) {
    return 'Неверный формат (например: example.com:${ServerConfig.defaultPort})';
  }

  return null;
}
