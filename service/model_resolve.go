package service

import (
	"fmt"
	"os"
	"path/filepath"
	"sort"
	"strings"
	"unicode"
)

func DisplayModelName(filename string) string {
	base := filepath.Base(filename)
	if len(base) > 5 && strings.EqualFold(base[len(base)-5:], ".gguf") {
		return base[:len(base)-5]
	}

	return base
}

func ListGGUFBasenames(dir string) ([]string, error) {
	entries, err := os.ReadDir(dir)
	if err != nil {
		return nil, err
	}

	var out []string
	for _, e := range entries {
		if e.IsDir() {
			continue
		}

		name := e.Name()
		if strings.EqualFold(filepath.Ext(name), ".gguf") {
			out = append(out, name)
		}
	}

	sort.Strings(out)

	return out, nil
}

func SplitModelRef(userInput string) (name, tag string) {
	s := strings.TrimSpace(userInput)
	if s == "" {
		return "", ""
	}

	if strings.Contains(s, "/") || strings.Contains(s, string(filepath.Separator)) {
		return s, ""
	}

	i := strings.LastIndex(s, ":")
	if i <= 0 || i+1 >= len(s) {
		return s, ""
	}

	return strings.TrimSpace(s[:i]), strings.TrimSpace(s[i+1:])
}

func reasonableTag(tag string) bool {
	if tag == "" || len(tag) > 48 {
		return false
	}

	for _, r := range tag {
		if unicode.IsLetter(r) || unicode.IsDigit(r) || r == '.' || r == '_' {
			continue
		}

		return false
	}

	return true
}

func splitStemToTagged(stem string) (base, tag string, ok bool) {
	if stem == "" {
		return "", "", false
	}

	i := strings.LastIndex(stem, "-")
	if i <= 0 || i+1 >= len(stem) {
		return "", "", false
	}

	base, tag = stem[:i], stem[i+1:]
	if base == "" || !reasonableTag(tag) {
		return "", "", false
	}

	return base, tag, true
}

func CatalogModelNames(dir string) ([]string, error) {
	files, err := ListGGUFBasenames(dir)
	if err != nil {
		return nil, err
	}

	seen := make(map[string]struct{}, len(files)*2)
	for _, f := range files {
		stem := DisplayModelName(f)
		if stem == "" {
			continue
		}

		seen[stem] = struct{}{}
		if b, tg, ok := splitStemToTagged(stem); ok {
			seen[b+":"+tg] = struct{}{}
		}
	}

	if err := addManifestCatalogEntries(dir, seen); err != nil {
		return nil, err
	}

	out := make([]string, 0, len(seen))
	for k := range seen {
		out = append(out, k)
	}
	sort.Strings(out)

	return out, nil
}

func SortedDisplayModelNames(dir string) ([]string, error) {
	return CatalogModelNames(dir)
}

func ResolveGGUFFile(modelsDir, userInput string) (canonical string, err error) {
	raw := strings.TrimSpace(userInput)
	if raw == "" {
		return "", fmt.Errorf("пустое имя модели")
	}

	refName, refTag := SplitModelRef(raw)
	if refTag != "" && !strings.EqualFold(refTag, "latest") {
		candName := fmt.Sprintf("%s-%s.gguf", refName, refTag)
		if c, err := matchGGUFBasename(modelsDir, candName); err == nil && c != "" {
			return c, nil
		}

		return "", fmt.Errorf("модель %q не найдена", userInput)
	}

	base := filepath.Base(raw)
	if base == "." || base == string(filepath.Separator) {
		return "", fmt.Errorf("некорректное имя модели")
	}

	lookup := filepath.Base(refName)
	if lookup == "." || lookup == string(filepath.Separator) {
		lookup = base
	}

	try := filepath.Join(modelsDir, lookup)
	if st, e := os.Stat(try); e == nil && !st.IsDir() {
		if strings.EqualFold(filepath.Ext(lookup), ".gguf") {
			return filepath.Base(lookup), nil
		}
	}

	if filepath.Ext(lookup) == "" {
		cand := lookup + ".gguf"
		try2 := filepath.Join(modelsDir, cand)
		if _, e := os.Stat(try2); e == nil {
			return cand, nil
		}
	}

	files, err := ListGGUFBasenames(modelsDir)
	if err != nil {
		return "", err
	}

	for _, name := range files {
		if name == lookup {
			return name, nil
		}
	}

	for _, name := range files {
		if DisplayModelName(name) == lookup {
			return name, nil
		}
	}

	for _, name := range files {
		if strings.EqualFold(name, lookup) {
			return name, nil
		}
	}

	for _, name := range files {
		if strings.EqualFold(DisplayModelName(name), lookup) {
			return name, nil
		}
	}

	return "", fmt.Errorf("модель %q не найдена", userInput)
}

func matchGGUFBasename(modelsDir, wantBase string) (string, error) {
	files, err := ListGGUFBasenames(modelsDir)
	if err != nil {
		return "", err
	}

	for _, name := range files {
		if name == wantBase {
			return name, nil
		}
	}

	for _, name := range files {
		if strings.EqualFold(name, wantBase) {
			return name, nil
		}
	}

	return "", fmt.Errorf("в каталоге моделей %q нет файла .gguf с базовым именем %q (проверьте имя файла и расширение .gguf)", modelsDir, wantBase)
}
