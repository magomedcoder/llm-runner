package service

import (
	"strings"
	"testing"

	"github.com/magomedcoder/gen/llm-runner/domain"
	"github.com/magomedcoder/gen/llm-runner/template"
)

const sampleChatMLJinja = `{% for message in messages %}{{'<|im_start|>' + message['role'] + '\n' + message['content'] + '<|im_end|>' + '\n'}}{% endfor %}{% if add_generation_prompt %}{{ '<|im_start|>assistant\n' }}{% endif %}`

func TestBuildChatPromptGotmpl_chatml(t *testing.T) {
	msgs := []*domain.AIChatMessage{
		domain.NewAIChatMessage(0, "Привет", domain.AIChatMessageRoleUser),
	}

	p, err := BuildChatPromptGotmpl(sampleChatMLJinja, msgs, nil)
	if err != nil {
		t.Fatal(err)
	}

	if !strings.Contains(p, "<|im_start|>user") || !strings.Contains(p, "Привет") {
		t.Fatalf("неожиданный промпт: %q", p)
	}

	if !strings.Contains(p, "<|im_start|>assistant") {
		t.Fatalf("нет заголовка assistant в промпте: %q", p)
	}
}

func TestRenderMatchedPreset_toolsOptional_chatmlUnchanged(t *testing.T) {
	msgs := []*domain.AIChatMessage{
		domain.NewAIChatMessage(0, "hi", domain.AIChatMessageRoleUser),
	}
	p, err := template.Named(strings.TrimSpace(sampleChatMLJinja))
	if err != nil {
		t.Fatal(err)
	}

	base, err := RenderMatchedPreset(p, msgs, nil)
	if err != nil {
		t.Fatal(err)
	}

	gp := &domain.GenerationParams{
		Tools: []domain.Tool{{
			Name:           "get_weather",
			Description:    "weather",
			ParametersJSON: `{"type":"object","properties":{"city":{"type":"string"}}}`,
		}},
	}

	withTools, err := RenderMatchedPreset(p, msgs, gp)
	if err != nil {
		t.Fatal(err)
	}

	if base != withTools {
		t.Fatalf("шаблон chatml не должен меняться от .Tools; без=%q с инструментами=%q", base, withTools)
	}
}

func TestRenderMatchedPreset_sameAsBuildChatPromptGotmpl(t *testing.T) {
	msgs := []*domain.AIChatMessage{
		domain.NewAIChatMessage(0, "x", domain.AIChatMessageRoleUser),
	}

	a, err := BuildChatPromptGotmpl(sampleChatMLJinja, msgs, nil)
	if err != nil {
		t.Fatal(err)
	}

	p, err := template.Named(strings.TrimSpace(sampleChatMLJinja))
	if err != nil {
		t.Fatal(err)
	}

	b, err := RenderMatchedPreset(p, msgs, nil)
	if err != nil {
		t.Fatal(err)
	}

	if a != b {
		t.Fatalf("расхождение промптов\na=%q\nb=%q", a, b)
	}
}
