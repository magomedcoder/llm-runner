package provider

import (
	"github.com/magomedcoder/llm-runner/config"
	"testing"
)

func TestNewTextProvider_llama_emptyPath(t *testing.T) {
	cfg := &config.Config{
		ModelPath: "",
	}

	_, err := NewTextProvider(cfg)
	if err == nil {
		t.Fatal("ожидалась ошибка при пустом model_path")
	}
}

func TestNewTextProvider_llama_withPath(t *testing.T) {
	cfg := &config.Config{
		ModelPath: "/models",
	}
	tp, err := NewTextProvider(cfg)
	if err != nil {
		t.Fatalf("NewTextProvider: %v", err)
	}

	if tp == nil {
		t.Fatal("ожидался непустой провайдер")
	}
}
