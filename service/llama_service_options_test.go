//go:build llama

package service

import (
	"testing"

	"github.com/magomedcoder/llm-runner/domain"
	"github.com/magomedcoder/llm-runner/llama"
)

func TestMapChatReasoningFormat(t *testing.T) {
	cases := []struct {
		name string
		in   string
		want llama.ReasoningFormat
	}{
		{
			name: "none по умолчанию",
			in:   "",
			want: llama.ReasoningFormatNone,
		},
		{
			name: "режим auto",
			in:   "auto",
			want: llama.ReasoningFormatAuto,
		},
		{
			name: "устаревший deepseek",
			in:   "deepseek-legacy",
			want: llama.ReasoningFormatDeepSeekLegacy,
		},
		{
			name: "режим deepseek",
			in:   "deepseek",
			want: llama.ReasoningFormatDeepSeek,
		},
		{
			name: "неизвестный формат",
			in:   "abc",
			want: llama.ReasoningFormatNone,
		},
	}

	for _, tc := range cases {
		t.Run(tc.name, func(t *testing.T) {
			got := mapChatReasoningFormat(tc.in)
			if got != tc.want {
				t.Fatalf("mapChatReasoningFormat(%q)=%v, ожидалось=%v", tc.in, got, tc.want)
			}
		})
	}
}

func TestRequiresGeneratePipeline_MinP(t *testing.T) {
	minP := float32(0.1)
	gen := &domain.GenerationParams{
		MinP: &minP,
	}

	if !requiresGeneratePipeline(gen, nil, &LlamaService{}) {
		t.Fatal("ожидался переход на generate pipeline для min_p")
	}
}

func TestRequiresGeneratePipeline_NProbsAndDebug(t *testing.T) {
	nProbs := 5
	s := &LlamaService{
		nProbs:          &nProbs,
		debugGeneration: true,
	}

	if !requiresGeneratePipeline(nil, nil, s) {
		t.Fatal("ожидался переход на generate pipeline для n_probs/debug_generation")
	}
}

func TestRequiresGeneratePipeline_ReasoningOnlyDoesNotFallback(t *testing.T) {
	s := &LlamaService{
		chatAPIEnabled:      true,
		chatReasoningFormat: "deepseek",
	}

	if requiresGeneratePipeline(nil, nil, s) {
		t.Fatal("не ожидался переход на generate pipeline, когда заданы только chat reasoning опции")
	}
}

func TestRequiresGeneratePipeline_ChatAPIWithReasoningAndAdvancedSamplingFallsBack(t *testing.T) {
	nProbs := 4
	s := &LlamaService{
		chatAPIEnabled:      true,
		chatReasoningFormat: "deepseek",
		nProbs:              &nProbs,
	}

	if !requiresGeneratePipeline(nil, nil, s) {
		t.Fatal("ожидался переход на generate pipeline, когда включен chat api с расширенными generate-опциями")
	}
}
