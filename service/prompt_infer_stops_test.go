package service

import (
	"slices"
	"testing"
)

func TestInferStopSequencesFromPrompt_chatML(t *testing.T) {
	p := "x" + chatMLImStart + "assistant\nhi" + chatMLImEnd
	got := inferStopSequencesFromPrompt(p)
	if !slices.Contains(got, chatMLImEnd) || !slices.Contains(got, chatMLImStart) {
		t.Fatalf("ChatML стопы: получено %v", got)
	}
}

func TestInferStopSequencesFromPrompt_llama(t *testing.T) {
	p := "x" + llamaHdrStart + "assistant" + llamaHdrEnd + llamaEOT
	got := inferStopSequencesFromPrompt(p)
	if !slices.Contains(got, llamaEOT) {
		t.Fatalf("Llama стопы: получено %v", got)
	}
}

func TestInferStopSequencesFromPrompt_plain(t *testing.T) {
	if len(inferStopSequencesFromPrompt("User: hi\nAssistant:")) != 0 {
		t.Fatal("для простого промпта не должно выводиться стоп-последовательностей")
	}
}
