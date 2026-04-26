package usecase

import (
	"strings"
	"testing"
)

func TestSanitizeRewrittenQuery(t *testing.T) {
	got := sanitizeRewrittenQuery("  одна строка  ")
	if got != "одна строка" {
		t.Fatalf("получено %q", got)
	}

	got = sanitizeRewrittenQuery("первая строка\nвторая игнорируется")
	if got != "первая строка" {
		t.Fatalf("ожидалась первая строка, получено %q", got)
	}

	long := strings.Repeat("а", ragQueryRewriteMaxRunes+50)
	got = sanitizeRewrittenQuery(long)
	if len([]rune(got)) != ragQueryRewriteMaxRunes {
		t.Fatalf("лимит: %d", len([]rune(got)))
	}
}

func TestSanitizeHyDEPseudoDocument(t *testing.T) {
	got := sanitizeHyDEPseudoDocument("  гипотетический ответ  ")
	if got != "гипотетический ответ" {
		t.Fatalf("получено %q", got)
	}

	long := strings.Repeat("b", ragHyDEMaxRunes+10)
	got = sanitizeHyDEPseudoDocument(long)
	if len([]rune(got)) != ragHyDEMaxRunes {
		t.Fatalf("лимит: %d", len([]rune(got)))
	}
}
