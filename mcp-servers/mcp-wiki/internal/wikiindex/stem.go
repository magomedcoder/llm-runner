package wikiindex

import (
	"strings"

	"github.com/kljensen/snowball"
)

func snowballLanguage(locale string) string {
	switch normalizeLocale(locale) {
	case "ru":
		return "russian"
	case "en":
		return "english"
	default:
		return ""
	}
}

func stemStatsLabel(locale string) string {
	switch normalizeLocale(locale) {
	case "ru":
		return "snowball_russian"
	case "en":
		return "snowball_english"
	default:
		return ""
	}
}

func applySnowballStem(locale string, terms []string) []string {
	lang := snowballLanguage(locale)
	if lang == "" || len(terms) == 0 {
		return terms
	}

	out := make([]string, 0, len(terms))
	for _, t := range terms {
		stemmed, err := snowball.Stem(t, lang, false)
		if err != nil || stemmed == "" {
			stemmed = t
		}

		if len([]rune(stemmed)) < 2 {
			continue
		}

		out = append(out, strings.ToLower(stemmed))
	}

	if len(out) == 0 {
		return nil
	}

	return out
}
