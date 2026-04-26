package service

import (
	"sort"
	"strings"
	"unicode/utf8"
)

func trimTrailingStops(s string, stops []string) string {
	if s == "" || len(stops) == 0 {
		return s
	}
	st := append([]string(nil), stops...)
	sort.Slice(st, func(i, j int) bool {
		return len(st[i]) > len(st[j])
	})

	out := s
	changed := true
	for changed {
		changed = false
		for _, t := range st {
			if t != "" && strings.HasSuffix(out, t) {
				out = strings.TrimSuffix(out, t)
				changed = true
			}
		}
	}

	return strings.TrimRight(out, " \t\r\n")
}

func completionNoiseSuffixes() []string {
	return []string{
		chatMLImEnd,
		chatMLImEnd + "\n",
		chatMLImEnd + "\r\n",
		chatMLImEndFW,
		llamaEOT,
		llamaEOT + "\n",
		llamaEOT + "\r\n",
	}
}

func sanitizeCompletionSuffix(s string) string {
	return trimTrailingStops(s, completionNoiseSuffixes())
}

type stopStreamFilter struct {
	stops    []string
	maxRunes int
	buf      []rune
	emit     func(string)
}

func newStopStreamFilter(stops []string, emit func(string)) *stopStreamFilter {
	stops = MergeStopSequences(stops, completionNoiseSuffixes())
	var cleaned []string
	maxRunes := 0
	for _, s := range stops {
		s = strings.TrimSpace(s)
		if s == "" {
			continue
		}

		cleaned = append(cleaned, s)
		if n := utf8.RuneCountInString(s); n > maxRunes {
			maxRunes = n
		}
	}

	sort.Slice(cleaned, func(i, j int) bool {
		return len(cleaned[i]) > len(cleaned[j])
	})

	return &stopStreamFilter{
		stops:    cleaned,
		maxRunes: maxRunes,
		emit:     emit,
	}
}

func (f *stopStreamFilter) push(piece string) {
	if piece == "" {
		return
	}

	if f.maxRunes == 0 {
		f.emit(piece)
		return
	}

	f.buf = append(f.buf, []rune(piece)...)

	for len(f.buf) > f.maxRunes {
		emitCount := len(f.buf) - f.maxRunes
		chunk := string(f.buf[:emitCount])
		f.buf = f.buf[emitCount:]
		f.emit(chunk)
	}
}

func (f *stopStreamFilter) flush() {
	rest := trimTrailingStops(string(f.buf), f.stops)
	rest = sanitizeCompletionSuffix(rest)
	f.buf = nil
	if rest != "" {
		f.emit(rest)
	}
}
