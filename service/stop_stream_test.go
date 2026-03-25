package service

import (
	"strings"
	"testing"
)

const testImEnd = "<|im_end|>"

func TestTrimTrailingStops(t *testing.T) {
	got := trimTrailingStops("hello"+testImEnd, []string{testImEnd})
	if got != "hello" {
		t.Fatalf("обрезка стопа: получено %q, ожидалось hello", got)
	}
}

func TestStopStreamFilter_fullStopBufferedThenDropped(t *testing.T) {
	var b strings.Builder
	f := newStopStreamFilter([]string{testImEnd}, func(s string) { b.WriteString(s) })
	f.push(testImEnd)
	f.flush()

	if strings.Contains(b.String(), testImEnd) {
		t.Fatalf("стоп-последовательность не должна попадать в вывод: %q", b.String())
	}
}

func TestStopStreamFilter_textThenStop(t *testing.T) {
	var b strings.Builder
	f := newStopStreamFilter([]string{testImEnd}, func(s string) { b.WriteString(s) })
	f.push("Привет!")
	f.push(testImEnd)
	f.flush()
	got := b.String()
	if strings.Contains(got, testImEnd) {
		t.Fatalf("в выводе не должно быть стоп-токена: %q", got)
	}

	if !strings.Contains(got, "Привет") {
		t.Fatalf("текст до стопа потерян: %q", got)
	}
}

func TestStopStreamFilter_russianNotSplit(t *testing.T) {
	var b strings.Builder
	f := newStopStreamFilter([]string{testImEnd}, func(s string) { b.WriteString(s) })
	for _, r := range "Привет" {
		f.push(string(r))
	}

	f.push(testImEnd)
	f.flush()
	got := b.String()
	if strings.Contains(got, testImEnd) {
		t.Fatalf("побайтово: в выводе не должно быть стоп-токена: %q", got)
	}

	if got != "Привет" {
		t.Fatalf("ожидалось «Привет», получено %q", got)
	}
}
