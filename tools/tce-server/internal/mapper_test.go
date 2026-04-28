package internal

import "testing"

func TestBuildEditorPrompt_Empty(t *testing.T) {
	t.Parallel()

	if got := buildEditorPrompt(nil); got != "" {
		t.Fatalf("ожидался пустой prompt для nil editor, получено %q", got)
	}

	editor := &EditorContext{}
	if got := buildEditorPrompt(editor); got != "" {
		t.Fatalf("ожидался пустой prompt для пустого editor, получено %q", got)
	}
}

func TestBuildEditorPrompt_Full(t *testing.T) {
	t.Parallel()

	line := 12
	col := 4
	editor := &EditorContext{
		Path:         "src/main.rs",
		Language:     "rust",
		Snippet:      "fn main() {}",
		CursorLine:   &line,
		CursorColumn: &col,
	}

	got := buildEditorPrompt(editor)
	want := "Editor context:\npath: src/main.rs\nlanguage: rust\ncursor: 12:4\nsnippet:\nfn main() {}"
	if got != want {
		t.Fatalf("неожиданный prompt:\nожидалось: %q\nполучено:  %q", want, got)
	}
}

func TestMapGenerateParams(t *testing.T) {
	t.Parallel()

	if got := mapGenerateParams(nil); got != nil {
		t.Fatalf("для nil-входа ожидался nil, получено %#v", got)
	}

	maxTokens := 1024
	temp := 0.2
	got := mapGenerateParams(&GenerateParams{
		MaxTokens:   &maxTokens,
		Temperature: &temp,
	})

	if got == nil {
		t.Fatal("ожидался не-nil результат")
	}

	if got.MaxTokens == nil || *got.MaxTokens != 1024 {
		t.Fatalf("неожиданный max_tokens: %#v", got.MaxTokens)
	}

	if got.Temperature == nil || *got.Temperature != float32(0.2) {
		t.Fatalf("неожиданная temperature: %#v", got.Temperature)
	}
}

func TestMakeRunnerMessages_IncludesSystemAndEditor(t *testing.T) {
	t.Parallel()

	system := "Ты помощник по коду."
	line := 0
	col := 1
	editor := &EditorContext{
		Path:         "src/main.go",
		Language:     "go",
		Snippet:      "package main",
		CursorLine:   &line,
		CursorColumn: &col,
	}
	input := []ChatMessage{
		{Role: "user", Content: "Привет"},
	}

	got := makeRunnerMessages(system, input, editor)
	if len(got) != 3 {
		t.Fatalf("ожидалось 3 сообщения, получено %d", len(got))
	}

	if got[0].Role != "system" || got[0].Content != system {
		t.Fatalf("неожиданное системное сообщение: %#v", got[0])
	}

	if got[1].Role != "system" {
		t.Fatalf("ожидался контекст editor как системное сообщение, получена роль=%q", got[1].Role)
	}
	
	if got[2].Role != "user" || got[2].Content != "Привет" {
		t.Fatalf("неожиданное пользовательское сообщение: %#v", got[2])
	}
}
