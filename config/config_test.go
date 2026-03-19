package config

import (
	"os"
	"path/filepath"
	"testing"
)

func TestLoad(t *testing.T) {
	wd, err := os.Getwd()
	if err != nil {
		t.Fatalf("Getwd: %v", err)
	}

	projectRoot := filepath.Join(wd, "..", "..")
	configPath := filepath.Join(projectRoot, "configs", "config.yaml")
	os.Setenv("LLM_RUNNER_CONFIG", configPath)
	defer os.Unsetenv("LLM_RUNNER_CONFIG")

	cfg, err := Load()
	if err != nil {
		t.Fatalf("Load: %v", err)
	}

	if cfg == nil {
		t.Fatal("конфиг не должен быть nil")
	}

	if cfg.CoreAddr == "" || cfg.ListenAddr == "" {
		t.Error("адреса должны быть заданы")
	}
}
