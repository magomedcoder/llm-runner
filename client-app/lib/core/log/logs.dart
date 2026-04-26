import 'package:gen/core/log/print_logs_native.dart';

enum Level { debug, verbose, info, warn, error }

class LogEvent {
  final Level level;
  final String message;
  final DateTime time;
  final Object? exception;
  final StackTrace? stackTrace;

  LogEvent({
    required this.level,
    required this.message,
    DateTime? time,
    this.exception,
    this.stackTrace,
  }) : time = time ?? DateTime.now();
}

class Logs {
  static final Logs _instance = Logs._();
  static Logs get instance => _instance;

  factory Logs() => _instance;

  Logs._();

  final List<LogEvent> _events = [];
  final List<void Function(LogEvent)> _listeners = [];

  List<LogEvent> get events => List.unmodifiable(_events);

  void addLogEvent(LogEvent event) {
    _events.add(event);
    event.printOut();
    for (final fn in _listeners) {
      fn(event);
    }
  }

  void addListener(void Function(LogEvent) fn) {
    _listeners.add(fn);
  }

  void removeListener(void Function(LogEvent) fn) {
    _listeners.remove(fn);
  }

  void d(String message, {Object? exception, StackTrace? stackTrace}) {
    addLogEvent(LogEvent(
      level: Level.debug,
      message: message,
      exception: exception,
      stackTrace: stackTrace,
    ));
  }

  void v(String message, {Object? exception, StackTrace? stackTrace}) {
    addLogEvent(LogEvent(
      level: Level.verbose,
      message: message,
      exception: exception,
      stackTrace: stackTrace,
    ));
  }

  void i(String message, {Object? exception, StackTrace? stackTrace}) {
    addLogEvent(LogEvent(
      level: Level.info,
      message: message,
      exception: exception,
      stackTrace: stackTrace,
    ));
  }

  void w(String message, {Object? exception, StackTrace? stackTrace}) {
    addLogEvent(LogEvent(
      level: Level.warn,
      message: message,
      exception: exception,
      stackTrace: stackTrace,
    ));
  }

  void e(String message, {Object? exception, StackTrace? stackTrace}) {
    addLogEvent(LogEvent(
      level: Level.error,
      message: message,
      exception: exception,
      stackTrace: stackTrace,
    ));
  }
}
