//go:build llama

package service

import (
	"strings"
	"testing"

	"github.com/magomedcoder/gen/llm-runner/domain"
)

func TestApplyResponseFormatPrompt_NoParams(t *testing.T) {
	in := "промпт"
	out := applyResponseFormatPrompt(in, nil)
	if out != in {
		t.Fatalf("промпт должен остаться без изменений, получено %q", out)
	}
}

func TestApplyResponseFormatPrompt_JsonObject_DefaultGrammar(t *testing.T) {
	in := "Ответ:"
	p := &domain.GenerationParams{
		ResponseFormat: &domain.ResponseFormat{
			Type: "json_object",
		},
	}

	out := applyResponseFormatPrompt(in, p)
	if !strings.Contains(out, "Return ONLY a single valid JSON object") {
		t.Fatalf("отсутствует заголовок json-ограничения: %q", out)
	}

	if !strings.Contains(out, DefaultJSONObjectGrammar) {
		t.Fatalf("должна быть включена грамматика по умолчанию")
	}
}

func TestApplyResponseFormatPrompt_JsonObject_CustomSchema(t *testing.T) {
	in := "Ответ:"
	schema := "root ::= \"ok\""
	p := &domain.GenerationParams{
		ResponseFormat: &domain.ResponseFormat{
			Type:   "json_object",
			Schema: &schema,
		},
	}

	out := applyResponseFormatPrompt(in, p)
	if !strings.Contains(out, schema) {
		t.Fatalf("должна быть включена пользовательская схема: %q", out)
	}
}

func TestApplyResponseFormatPrompt_NonJsonObject(t *testing.T) {
	in := "Ответ:"
	p := &domain.GenerationParams{
		ResponseFormat: &domain.ResponseFormat{
			Type: "text",
		},
	}

	out := applyResponseFormatPrompt(in, p)
	if out != in {
		t.Fatalf("промпт должен остаться без изменений, получено %q", out)
	}
}
