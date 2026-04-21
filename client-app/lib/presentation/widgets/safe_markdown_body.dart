import 'dart:convert';

import 'package:flutter/material.dart';
import 'package:flutter_markdown_plus/flutter_markdown_plus.dart';
import 'package:markdown/markdown.dart' as md;

class SafeMarkdownBody extends StatefulWidget {
  const SafeMarkdownBody({
    super.key,
    required this.data,
    required this.styleSheet,
    this.selectable = true,
    this.builders = const <String, MarkdownElementBuilder>{},
    this.extensionSet,
    this.fallbackStyle,
  });

  final String data;
  final MarkdownStyleSheet styleSheet;
  final bool selectable;
  final Map<String, MarkdownElementBuilder> builders;
  final md.ExtensionSet? extensionSet;
  final TextStyle? fallbackStyle;

  @override
  State<SafeMarkdownBody> createState() => _SafeMarkdownBodyState();
}

class _SafeMarkdownBodyState extends State<SafeMarkdownBody> {
  bool _parseOk = true;

  void _probe(String data) {
    try {
      final doc = md.Document(
        extensionSet: widget.extensionSet ?? md.ExtensionSet.gitHubFlavored,
        encodeHtml: false,
      );
      final lines = const LineSplitter().convert(data);
      doc.parseLines(lines);
      _parseOk = true;
    } on Object catch (_) {
      _parseOk = false;
    }
  }

  @override
  void initState() {
    super.initState();
    _probe(widget.data);
  }

  @override
  void didUpdateWidget(covariant SafeMarkdownBody oldWidget) {
    super.didUpdateWidget(oldWidget);
    if (widget.data != oldWidget.data ||
        widget.extensionSet != oldWidget.extensionSet) {
      _probe(widget.data);
    }
  }

  @override
  Widget build(BuildContext context) {
    if (!_parseOk) {
      return SelectableText(
        widget.data,
        style: widget.fallbackStyle ?? widget.styleSheet.p,
      );
    }
    return MarkdownBody(
      data: widget.data,
      selectable: widget.selectable,
      styleSheet: widget.styleSheet,
      builders: widget.builders,
      extensionSet: widget.extensionSet,
    );
  }
}
