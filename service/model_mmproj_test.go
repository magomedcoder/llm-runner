package service

import (
	"os"
	"path/filepath"
	"testing"
)

func TestResolveMmprojPath_auto(t *testing.T) {
	dir := t.TempDir()
	_ = os.WriteFile(filepath.Join(dir, "M.gguf"), nil, 0o644)
	_ = os.WriteFile(filepath.Join(dir, "M-mmproj.gguf"), nil, 0o644)

	p, err := ResolveMmprojPath(dir, "M.gguf", "")
	if err != nil || filepath.Base(p) != "M-mmproj.gguf" {
		t.Fatalf("авто mmproj: путь %q, ошибка %v", p, err)
	}
}

func TestResolveMmprojPath_override(t *testing.T) {
	dir := t.TempDir()
	_ = os.WriteFile(filepath.Join(dir, "proj.gguf"), nil, 0o644)

	p, err := ResolveMmprojPath(dir, "x.gguf", "proj.gguf")
	if err != nil || filepath.Base(p) != "proj.gguf" {
		t.Fatalf("переопределение mmproj: путь %q, ошибка %v", p, err)
	}
}

func TestResolveMmprojPath_missingOverride(t *testing.T) {
	dir := t.TempDir()
	_, err := ResolveMmprojPath(dir, "x.gguf", "nope.gguf")
	if err == nil {
		t.Fatal("ожидалась ошибка для отсутствующего mmproj")
	}
}
