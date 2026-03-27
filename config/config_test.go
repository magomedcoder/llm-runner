package config

import (
	"os"
	"path/filepath"
	"testing"
)

func TestLoad(t *testing.T) {
	dir := t.TempDir()
	configPath := filepath.Join(dir, "runner.yaml")
	content := `core:
  host: "127.0.0.1"
  port: 50051
listen:
  host: "127.0.0.1"
  port: 50052
log_level: "info"
model_path: "./models"
default_model: "test-model"
`
	if err := os.WriteFile(configPath, []byte(content), 0o644); err != nil {
		t.Fatalf("запись тестового конфига: %v", err)
	}

	t.Setenv("LLM_RUNNER_CONFIG", configPath)

	cfg, err := Load()
	if err != nil {
		t.Fatalf("загрузка конфигурации: %v", err)
	}

	if cfg == nil {
		t.Fatal("конфиг не должен быть nil")
	}

	if cfg.CoreAddr() == "" || cfg.ListenAddr() == "" {
		t.Error("адреса должны быть заданы")
	}
	if cfg.CoreAddr() != "127.0.0.1:50051" {
		t.Errorf("CoreAddr: ожидалось 127.0.0.1:50051, получено %q", cfg.CoreAddr())
	}
	if cfg.ListenAddr() != "127.0.0.1:50052" {
		t.Errorf("ListenAddr: ожидалось 127.0.0.1:50052, получено %q", cfg.ListenAddr())
	}

	if cfg.DefaultModel != "test-model" {
		t.Errorf("default_model: ожидалось test-model, получено %q", cfg.DefaultModel)
	}
}

func TestLoad_ChatAPIReasoningAndDebugFields(t *testing.T) {
	dir := t.TempDir()
	configPath := filepath.Join(dir, "runner.yaml")
	content := `core:
  host: "127.0.0.1"
  port: 50051
listen:
  host: "127.0.0.1"
  port: 50052
model_path: "./models"
default_model: "test-model"
chat_api_enabled: true
chat_stream_buffer_size: 256
chat_reasoning_format: "deepseek"
chat_enable_thinking: true
chat_reasoning_budget: 512
n_probs: 8
debug_generation: true
reinit_llama_logging: true
log_model_stats: true
`
	if err := os.WriteFile(configPath, []byte(content), 0o644); err != nil {
		t.Fatalf("запись тестового конфига: %v", err)
	}
	t.Setenv("LLM_RUNNER_CONFIG", configPath)

	cfg, err := Load()
	if err != nil {
		t.Fatalf("загрузка конфигурации: %v", err)
	}

	if !cfg.ChatAPIEnabled {
		t.Fatal("chat_api_enabled должен быть true")
	}

	if cfg.ChatReasoningFormat != "deepseek" {
		t.Fatalf("chat_reasoning_format: %q", cfg.ChatReasoningFormat)
	}

	if cfg.ChatEnableThinking == nil || !*cfg.ChatEnableThinking {
		t.Fatal("chat_enable_thinking должен быть true")
	}

	if cfg.ChatReasoningBudget == nil || *cfg.ChatReasoningBudget != 512 {
		t.Fatalf("chat_reasoning_budget: %+v", cfg.ChatReasoningBudget)
	}

	if cfg.NProbs == nil || *cfg.NProbs != 8 {
		t.Fatalf("n_probs: %+v", cfg.NProbs)
	}

	if !cfg.DebugGeneration {
		t.Fatal("debug_generation должен быть true")
	}

	if !cfg.ReinitLlamaLogging || !cfg.LogModelStats {
		t.Fatal("observability flags должны быть true")
	}
}
