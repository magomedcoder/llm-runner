package provider

import (
	"github.com/magomedcoder/gen/llm-runner/config"
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

func TestNewTextProvider_smoke_chatAndDebugOptions(t *testing.T) {
	nProbs := 8
	enableThinking := true
	reasoningBudget := 256
	typicalP := float32(0.95)
	minKeep := 1
	dynRange := float32(0.5)
	dynExp := float32(1.2)
	cfg := &config.Config{
		ModelPath:                  "/models",
		MaxContextTokens:           4096,
		ChatAPIEnabled:             true,
		ChatStreamBufferSize:       128,
		ChatReasoningFormat:        "deepseek",
		ChatEnableThinking:         &enableThinking,
		ChatReasoningBudget:        &reasoningBudget,
		NProbs:                     &nProbs,
		DebugGeneration:            true,
		TypicalP:                   &typicalP,
		MinKeep:                    &minKeep,
		DynamicTemperatureRange:    &dynRange,
		DynamicTemperatureExponent: &dynExp,
		ReinitLlamaLogging:         true,
		LogModelStats:              true,
	}

	tp, err := NewTextProvider(cfg)
	if err != nil {
		t.Fatalf("NewTextProvider smoke: %v", err)
	}
	if tp == nil {
		t.Fatal("ожидался непустой провайдер")
	}
}
