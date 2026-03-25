package template

import (
	"strings"
	"testing"
)

func mergeUniqueStrings(a, b []string) []string {
	seen := make(map[string]struct{}, len(a)+len(b))
	var out []string
	for _, s := range a {
		if _, ok := seen[s]; ok {
			continue
		}

		seen[s] = struct{}{}
		out = append(out, s)
	}

	for _, s := range b {
		if _, ok := seen[s]; ok {
			continue
		}

		seen[s] = struct{}{}
		out = append(out, s)
	}

	return out
}

func TestMergeUniqueStrings(t *testing.T) {
	got := mergeUniqueStrings([]string{"a", "b"}, []string{"b", "c"})
	if strings.Join(got, ",") != "a,b,c" {
		t.Fatalf("ожидалось a,b,c без дубликатов, получено %v", got)
	}
}

func TestTemplatesOnce_chatmlStopwords(t *testing.T) {
	presets, err := templatesOnce()
	if err != nil {
		t.Fatal(err)
	}

	var found bool
	for _, p := range presets {
		if p.Name != "chatml" || p.Parameters == nil {
			continue
		}

		if containsAll(p.Parameters.Stop, "<|im_start|>", "<|im_end|>") {
			found = true
			break
		}
	}

	if !found {
		t.Fatal("пресет chatml должен содержать стоп-слова из presets/chatml.json")
	}
}

func containsAll(hay []string, needles ...string) bool {
	set := make(map[string]struct{}, len(hay))
	for _, s := range hay {
		set[s] = struct{}{}
	}

	for _, n := range needles {
		if _, ok := set[n]; !ok {
			return false
		}
	}

	return true
}
