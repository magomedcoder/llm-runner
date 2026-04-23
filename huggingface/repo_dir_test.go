package huggingface

import (
	"path/filepath"
	"testing"
)

func TestRepoDownloadSubdir(t *testing.T) {
	got := RepoDownloadSubdir("Qwen/Qwen3-8B-GGUF")
	want := filepath.Join("Qwen", "Qwen3-8B-GGUF")
	if got != want {
		t.Fatalf("RepoDownloadSubdir: %q, want %q", got, want)
	}

	if g := RepoDownloadSubdir(""); g != "unknown-repo" {
		t.Fatalf("empty: %q", g)
	}

	if g := RepoDownloadSubdir("///"); g != "unknown-repo" {
		t.Fatalf("slashes only: %q", g)
	}
}
