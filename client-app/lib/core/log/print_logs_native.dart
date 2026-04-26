import 'package:gen/core/log/logs.dart';

extension LogEventPrint on LogEvent {
  static const _ansiReset = '\u001b[0m';
  static const _ansiRed = '\u001b[31m';
  static const _ansiGreen = '\u001b[32m';
  static const _ansiYellow = '\u001b[33m';
  static const _ansiBlue = '\u001b[34m';

  String get _levelTag {
    switch (level) {
      case Level.debug:
        return 'DEBUG';
      case Level.verbose:
        return 'VERBOSE';
      case Level.info:
        return 'INFO';
      case Level.warn:
        return 'WARN';
      case Level.error:
        return 'ERROR';
    }
  }

  String _colorize(String text, String ansi) => '$ansi$text$_ansiReset';

  void printOut({bool useColor = true}) {
    final timeStr = time.toIso8601String();
    var logsStr = '$_levelTag $timeStr $message';
    if (exception != null) {
      logsStr += ' | $exception';
      if (stackTrace != null) {
        logsStr += '\n$stackTrace';
      }
    }
    if (useColor) {
      switch (level) {
        case Level.debug:
          logsStr = _colorize(logsStr, _ansiBlue);
          break;
        case Level.info:
          logsStr = _colorize(logsStr, _ansiGreen);
          break;
        case Level.warn:
          logsStr = _colorize(logsStr, _ansiYellow);
          break;
        case Level.error:
          logsStr = _colorize(logsStr, _ansiRed);
          break;
        case Level.verbose:
          break;
      }
    }
    print('[Gen] $logsStr');
  }
}
