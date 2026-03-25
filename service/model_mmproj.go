package service

import (
	"fmt"
	"os"
	"path/filepath"
	"strings"
)

func mmprojBasenamesToTry(canonicalGGUF string) []string {
	base := filepath.Base(strings.TrimSpace(canonicalGGUF))
	stem := strings.TrimSuffix(base, ".gguf")
	if stem == "" || stem == base {
		return nil
	}

	return []string{
		stem + "-mmproj.gguf",
		stem + ".mmproj.gguf",
	}
}

func ResolveMmprojPath(modelsDir, canonicalGGUF, override string) (string, error) {
	modelsDir = strings.TrimSpace(modelsDir)
	o := strings.TrimSpace(override)
	if o != "" {
		var p string
		if filepath.IsAbs(o) {
			p = o
		} else {
			if modelsDir == "" {
				return "", fmt.Errorf("mmproj_path: пустой model_path для относительного пути %q", o)
			}

			p = filepath.Join(modelsDir, o)
		}

		st, err := os.Stat(p)
		if err != nil || st.IsDir() {
			return "", fmt.Errorf("mmproj_path: файл не найден: %s", p)
		}

		return p, nil
	}

	if modelsDir == "" || canonicalGGUF == "" {
		return "", nil
	}

	for _, name := range mmprojBasenamesToTry(canonicalGGUF) {
		p := filepath.Join(modelsDir, name)
		st, err := os.Stat(p)
		if err == nil && !st.IsDir() {
			return p, nil
		}
	}

	return "", nil
}
