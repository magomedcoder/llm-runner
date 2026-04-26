import 'package:flutter/material.dart';
import 'package:gen/domain/entities/rag_document_preview.dart';

class RagContextPreviewBanner extends StatefulWidget {
  const RagContextPreviewBanner({
    super.key,
    required this.preview,
    required this.onDismiss,
  });

  final RagDocumentPreview preview;
  final VoidCallback onDismiss;

  @override
  State<RagContextPreviewBanner> createState() => _RagContextPreviewBannerState();
}

String _modeTitle(String mode) {
  switch (mode) {
    case 'full_document':
      return 'Полный текст документа в промпте';
    case 'vector_rag':
      return 'Поиск по документу (фрагменты в промпте)';
    case 'vector_rag_deep':
      return 'Поиск по документу + deep map';
    default:
      if (mode.isEmpty) {
        return 'Контекст документа';
      }

      return 'Контекст: $mode';
  }
}

String _score(double s) {
  if (s == 0) {
    return '-';
  }

  final t = (s * 1000).round() / 1000;
  var out = t.toStringAsFixed(4);
  out = out.replaceFirst(RegExp(r'0+$'), '');
  out = out.replaceFirst(RegExp(r'\.$'), '');
  return out;
}

String? _collapsedSubtitle(RagDocumentPreview p) {
  final parts = <String>[];
  if (p.chunks.isNotEmpty) {
    parts.add('${p.chunks.length} фрагм. в промпте');
  }

  if (p.isFullDocument && p.fullDocumentExcerpt.isNotEmpty) {
    parts.add('есть превью текста');
  }

  if (p.summary.isNotEmpty) {
    final s = p.summary.trim();
    if (s.length <= 120) {
      parts.add(s);
    } else {
      parts.add('${s.substring(0, 117)}…');
    }
  }

  if (parts.isEmpty) {
    return null;
  }

  return parts.take(2).join(' · ');
}

class _RagContextPreviewBannerState extends State<RagContextPreviewBanner> {
  bool _expanded = false;

  @override
  void didUpdateWidget(covariant RagContextPreviewBanner oldWidget) {
    super.didUpdateWidget(oldWidget);
    if (oldWidget.preview != widget.preview) {
      _expanded = false;
    }
  }

  void _toggle() {
    setState(() => _expanded = !_expanded);
  }

  @override
  Widget build(BuildContext context) {
    final theme = Theme.of(context);
    final cs = theme.colorScheme;
    final p = widget.preview;
    final bg = cs.secondaryContainer;
    final fg = cs.onSecondaryContainer;
    final collapsedHint = _collapsedSubtitle(p);

    return Material(
      color: bg,
      elevation: 1,
      shadowColor: cs.shadow.withValues(alpha: 0.1),
      child: Padding(
        padding: const EdgeInsets.fromLTRB(4, 6, 4, 8),
        child: Column(
          crossAxisAlignment: CrossAxisAlignment.stretch,
          mainAxisSize: MainAxisSize.min,
          children: [
            Row(
              crossAxisAlignment: CrossAxisAlignment.start,
              children: [
                Padding(
                  padding: const EdgeInsets.only(top: 4, left: 4),
                  child: Icon(Icons.article_outlined, size: 22, color: fg),
                ),
                const SizedBox(width: 4),
                Expanded(
                  child: InkWell(
                    onTap: _toggle,
                    borderRadius: BorderRadius.circular(8),
                    child: Padding(
                      padding: const EdgeInsets.symmetric(
                        vertical: 4,
                        horizontal: 6,
                      ),
                      child: Column(
                        crossAxisAlignment: CrossAxisAlignment.start,
                        children: [
                          Row(
                            crossAxisAlignment: CrossAxisAlignment.start,
                            children: [
                              Expanded(
                                child: Text(
                                  'Как модель видит документ',
                                  style: theme.textTheme.titleSmall?.copyWith(
                                    color: fg,
                                    fontWeight: FontWeight.w600,
                                  ),
                                ),
                              ),
                              Icon(
                                _expanded
                                    ? Icons.expand_less_rounded
                                    : Icons.expand_more_rounded,
                                color: fg,
                                size: 24,
                              ),
                            ],
                          ),
                          const SizedBox(height: 2),
                          Text(
                            _modeTitle(p.mode),
                            style: theme.textTheme.labelMedium?.copyWith(
                              color: fg.withValues(alpha: 0.88),
                              fontWeight: FontWeight.w500,
                            ),
                          ),
                          if (!_expanded && collapsedHint != null) ...[
                            const SizedBox(height: 4),
                            Text(
                              collapsedHint,
                              maxLines: 3,
                              overflow: TextOverflow.ellipsis,
                              style: theme.textTheme.bodySmall?.copyWith(
                                color: fg.withValues(alpha: 0.78),
                                height: 1.32,
                              ),
                            ),
                          ],
                        ],
                      ),
                    ),
                  ),
                ),
                IconButton(
                  tooltip: 'Скрыть',
                  visualDensity: VisualDensity.compact,
                  icon: Icon(Icons.close_rounded, color: fg, size: 22),
                  onPressed: widget.onDismiss,
                ),
              ],
            ),
            if (_expanded) ...[
              const SizedBox(height: 4),
              Padding(
                padding: const EdgeInsets.only(left: 38, right: 10),
                child: Text(
                  'Ниже - то, как модель «видит» вложение при этой отправке (текст в промпте, не ответ).',
                  style: theme.textTheme.labelSmall?.copyWith(
                    color: fg.withValues(alpha: 0.82),
                    height: 1.35,
                  ),
                ),
              ),
              if (p.summary.isNotEmpty) ...[
                const SizedBox(height: 6),
                Padding(
                  padding: const EdgeInsets.only(left: 38, right: 10),
                  child: Text(
                    p.summary,
                    style: theme.textTheme.bodySmall?.copyWith(
                      color: fg.withValues(alpha: 0.92),
                      height: 1.38,
                    ),
                  ),
                ),
              ],
              if (!p.isFullDocument && (p.topK > 0 || p.neighborWindow > 0 || p.deepRagMapCalls > 0)) ...[
                const SizedBox(height: 6),
                Padding(
                  padding: const EdgeInsets.only(left: 38, right: 10),
                  child: Wrap(
                    spacing: 10,
                    runSpacing: 4,
                    children: [
                      if (p.topK > 0)
                        _MetaChip(label: 'top‑K: ${p.topK}', fg: fg),
                      if (p.neighborWindow > 0)
                        _MetaChip(
                          label: 'соседи: ±${p.neighborWindow}',
                          fg: fg,
                        ),
                      if (p.deepRagMapCalls > 0)
                        _MetaChip(
                          label: 'deep map: ${p.deepRagMapCalls}',
                          fg: fg,
                        ),
                    ],
                  ),
                ),
              ],
              if (p.isFullDocument && p.fullDocumentExcerpt.isNotEmpty) ...[
                const SizedBox(height: 10),
                Padding(
                  padding: const EdgeInsets.only(left: 38, right: 10),
                  child: Column(
                    crossAxisAlignment: CrossAxisAlignment.start,
                    children: [
                      Text(
                        'Начало извлечённого текста',
                        style: theme.textTheme.labelLarge?.copyWith(
                          color: fg,
                          fontWeight: FontWeight.w600,
                        ),
                      ),
                      const SizedBox(height: 4),
                      _PreviewTextBox(
                        fg: fg,
                        child: SelectableText(
                          p.fullDocumentExcerpt,
                          style: theme.textTheme.bodySmall?.copyWith(
                            color: fg.withValues(alpha: 0.88),
                            height: 1.42,
                            fontFamily: 'monospace',
                          ),
                        ),
                      ),
                    ],
                  ),
                ),
              ] else if (p.isFullDocument) ...[
                const SizedBox(height: 6),
                Padding(
                  padding: const EdgeInsets.only(left: 38, right: 10),
                  child: Text(
                    'В промпт уходит полный извлечённый текст файла (без разбиения на чанки для поиска).',
                    style: theme.textTheme.labelSmall?.copyWith(
                      color: fg.withValues(alpha: 0.85),
                      height: 1.35,
                    ),
                  ),
                ),
              ],
              if (p.chunks.isNotEmpty) ...[
                const SizedBox(height: 12),
                Padding(
                  padding: const EdgeInsets.only(left: 38, right: 10),
                  child: Text(
                    'Фрагменты, попавшие в промпт',
                    style: theme.textTheme.labelLarge?.copyWith(
                      color: fg,
                      fontWeight: FontWeight.w600,
                    ),
                  ),
                ),
                const SizedBox(height: 6),
                Padding(
                  padding: const EdgeInsets.only(left: 28, right: 8),
                  child: DecoratedBox(
                    decoration: BoxDecoration(
                      color: cs.surface.withValues(alpha: 0.35),
                      borderRadius: BorderRadius.circular(8),
                    ),
                    child: ConstrainedBox(
                      constraints: const BoxConstraints(maxHeight: 280),
                      child: Scrollbar(
                        thumbVisibility: true,
                        child: ListView.builder(
                          shrinkWrap: true,
                          padding: const EdgeInsets.symmetric(
                            vertical: 8,
                            horizontal: 10,
                          ),
                          itemCount: p.chunks.length,
                          itemBuilder: (context, i) {
                            final c = p.chunks[i];
                            return Padding(
                              padding: const EdgeInsets.only(bottom: 12),
                              child: Column(
                                crossAxisAlignment: CrossAxisAlignment.start,
                                children: [
                                  Wrap(
                                    spacing: 8,
                                    runSpacing: 4,
                                    crossAxisAlignment: WrapCrossAlignment.center,
                                    children: [
                                      Text(
                                        'Чанк #${c.chunkIndex}',
                                        style: theme.textTheme.labelMedium?.copyWith(
                                          color: fg,
                                          fontWeight: FontWeight.w700,
                                        ),
                                      ),
                                      Text(
                                        'релев. ${_score(c.score)}',
                                        style: theme.textTheme.labelSmall?.copyWith(color: fg),
                                      ),
                                      if (c.isNeighbor)
                                        Chip(
                                          label: Text(
                                            'рядом с хитом',
                                            style: theme.textTheme.labelSmall,
                                          ),
                                          visualDensity: VisualDensity.compact,
                                          padding: EdgeInsets.zero,
                                          materialTapTargetSize: MaterialTapTargetSize.shrinkWrap,
                                        ),
                                      if (c.pdfPageStart > 0)
                                        Text(
                                          c.pdfPageEnd > 0 && c.pdfPageEnd != c.pdfPageStart
                                              ? 'стр. ${c.pdfPageStart}–${c.pdfPageEnd}'
                                              : 'стр. ${c.pdfPageStart}',
                                          style: theme.textTheme.labelSmall?.copyWith(
                                            color: fg.withValues(alpha: 0.88),
                                          ),
                                        ),
                                    ],
                                  ),
                                  if (c.headingPath.isNotEmpty) ...[
                                    const SizedBox(height: 4),
                                    Text(
                                      c.headingPath,
                                      style: theme.textTheme.labelMedium?.copyWith(
                                        color: fg.withValues(alpha: 0.88),
                                        fontStyle: FontStyle.italic,
                                      ),
                                    ),
                                  ],
                                  if (c.excerpt.isNotEmpty) ...[
                                    const SizedBox(height: 6),
                                    SelectableText(
                                      c.excerpt,
                                      style: theme.textTheme.bodySmall?.copyWith(
                                        color: fg.withValues(alpha: 0.9),
                                        height: 1.4,
                                      ),
                                    ),
                                  ] else ...[
                                    const SizedBox(height: 4),
                                    Text(
                                      'Текст фрагмента недоступен в превью.',
                                      style: theme.textTheme.labelSmall?.copyWith(
                                        color: fg.withValues(alpha: 0.55),
                                        fontStyle: FontStyle.italic,
                                      ),
                                    ),
                                  ],
                                ],
                              ),
                            );
                          },
                        ),
                      ),
                    ),
                  ),
                ),
              ],
            ],
          ],
        ),
      ),
    );
  }
}

class _PreviewTextBox extends StatelessWidget {
  const _PreviewTextBox({required this.fg, required this.child});

  final Color fg;
  final Widget child;

  @override
  Widget build(BuildContext context) {
    return ConstrainedBox(
      constraints: const BoxConstraints(maxHeight: 160),
      child: DecoratedBox(
        decoration: BoxDecoration(
          color: Theme.of(context).colorScheme.surface.withValues(alpha: 0.4),
          borderRadius: BorderRadius.circular(8),
          border: Border.all(color: fg.withValues(alpha: 0.2)),
        ),
        child: Scrollbar(
          thumbVisibility: true,
          child: SingleChildScrollView(
            padding: const EdgeInsets.all(10),
            child: child,
          ),
        ),
      ),
    );
  }
}

class _MetaChip extends StatelessWidget {
  const _MetaChip({required this.label, required this.fg});

  final String label;
  final Color fg;

  @override
  Widget build(BuildContext context) {
    final theme = Theme.of(context);
    return Text(
      label,
      style: theme.textTheme.labelSmall?.copyWith(
        color: fg.withValues(alpha: 0.9),
      ),
    );
  }
}
