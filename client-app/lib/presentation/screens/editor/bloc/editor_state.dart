import 'package:equatable/equatable.dart';
import 'package:gen/generated/grpc_pb/editor.pb.dart' as grpc;

class EditorState extends Equatable {
  final bool isLoading;
  final String documentText;
  final List<String> undoStack;
  final List<String> redoStack;
  final String? selectedRunner;
  final grpc.TransformType type;
  final bool preserveMarkdown;
  final String? error;
  final int documentVersion;

  const EditorState({
    this.isLoading = false,
    this.documentText = '',
    this.undoStack = const [],
    this.redoStack = const [],
    this.selectedRunner,
    this.type = grpc.TransformType.TRANSFORM_TYPE_FIX,
    this.preserveMarkdown = false,
    this.error,
    this.documentVersion = 0,
  });

  bool get canUndo => undoStack.isNotEmpty;

  bool get canRedo => redoStack.isNotEmpty;

  EditorState copyWith({
    bool? isLoading,
    String? documentText,
    List<String>? undoStack,
    List<String>? redoStack,
    String? selectedRunner,
    bool clearSelectedRunner = false,
    grpc.TransformType? type,
    bool? preserveMarkdown,
    String? error,
    bool clearError = false,
    int? documentVersion,
  }) {
    return EditorState(
      isLoading: isLoading ?? this.isLoading,
      documentText: documentText ?? this.documentText,
      undoStack: undoStack ?? this.undoStack,
      redoStack: redoStack ?? this.redoStack,
      selectedRunner: clearSelectedRunner
        ? null
        : (selectedRunner ?? this.selectedRunner),
      type: type ?? this.type,
      preserveMarkdown: preserveMarkdown ?? this.preserveMarkdown,
      error: clearError ? null : (error ?? this.error),
      documentVersion: documentVersion ?? this.documentVersion,
    );
  }

  @override
  List<Object?> get props => [
    isLoading,
    documentText,
    undoStack,
    redoStack,
    selectedRunner,
    type,
    preserveMarkdown,
    error,
    documentVersion,
  ];
}
