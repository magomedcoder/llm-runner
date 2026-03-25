package service

import (
	"os"
	"path/filepath"
	"testing"

	"github.com/magomedcoder/llm-runner/domain"
)

func TestMergeGenParams_yamlDefaults(t *testing.T) {
	y := &ModelYAML{
		Parameter: &ModelYAMLParameter{
			Temperature: ptrFloat64(0.1),
			MaxTokens:   ptrInt(99),
		},
	}
	out := MergeGenParams(nil, y)
	if out == nil || out.Temperature == nil || *out.Temperature != 0.1 {
		t.Fatalf("temperature из yaml: %+v", out)
	}
	if out.MaxTokens == nil || *out.MaxTokens != 99 {
		t.Fatalf("max_tokens из yaml: %+v", out)
	}

	reqT := float32(0.8)
	req := &domain.GenerationParams{Temperature: &reqT}
	out2 := MergeGenParams(req, y)
	if out2.Temperature == nil || *out2.Temperature != 0.8 {
		t.Fatalf("температура из запроса должна перекрывать yaml: %+v", out2)
	}

	if out2.MaxTokens == nil || *out2.MaxTokens != 99 {
		t.Fatalf("max_tokens должен браться из yaml: %+v", out2)
	}
}

func TestApplyModelYAMLSystem(t *testing.T) {
	norm := []*domain.AIChatMessage{
		domain.NewAIChatMessage(1, "hi", domain.AIChatMessageRoleUser),
	}

	cfg := &ModelYAML{System: "SYS"}
	out := ApplyModelYAMLSystem(norm, cfg)
	if len(out) != 2 || out[0].Role != domain.AIChatMessageRoleSystem || out[0].Content != "SYS" {
		t.Fatalf("ожидалось system + user, получено %+v", out)
	}
}

func TestResolveModelForInference_sidecar_template(t *testing.T) {
	dir := t.TempDir()
	_ = os.WriteFile(filepath.Join(dir, "T.gguf"), []byte{0}, 0o644)
	_ = os.WriteFile(filepath.Join(dir, "T.yaml"), []byte(`template: "TEST_CHAT_TEMPLATE_MARKER"`), 0o644)
	_, cfg, err := ResolveModelForInference(dir, "T")
	if err != nil || cfg == nil || cfg.Template != "TEST_CHAT_TEMPLATE_MARKER" {
		t.Fatalf("конфиг sidecar: %+v, ошибка %v", cfg, err)
	}
}

func TestResolveModelForInference_sidecar(t *testing.T) {
	dir := t.TempDir()
	_ = os.WriteFile(filepath.Join(dir, "M.gguf"), []byte{0}, 0o644)
	_ = os.WriteFile(filepath.Join(dir, "M.yaml"), []byte(`system: "X"
parameter:
  temperature: 0.3
`), 0o644)

	can, cfg, err := ResolveModelForInference(dir, "M")
	if err != nil || can != "M.gguf" || cfg == nil || cfg.System != "X" {
		t.Fatalf("канонический файл=%q конфиг=%+v ошибка=%v", can, cfg, err)
	}

	if cfg.From != "" {
		t.Fatal("у sidecar-файла поле from должно очищаться")
	}
}

func TestResolveModelForInference_manifest(t *testing.T) {
	dir := t.TempDir()
	_ = os.WriteFile(filepath.Join(dir, "Base-Q4.gguf"), []byte{0}, 0o644)
	_ = os.WriteFile(filepath.Join(dir, "alias.yaml"), []byte(`from: Base-Q4.gguf
system: "bot"
parameter:
  top_p: 0.95
`), 0o644)

	can, cfg, err := ResolveModelForInference(dir, "alias")
	if err != nil || can != "Base-Q4.gguf" || cfg == nil || cfg.System != "bot" {
		t.Fatalf("манифест: файл=%q конфиг=%+v ошибка=%v", can, cfg, err)
	}

	names, err := CatalogModelNames(dir)
	if err != nil {
		t.Fatal(err)
	}

	found := false
	for _, n := range names {
		if n == "alias" {
			found = true
		}
	}

	if !found {
		t.Fatalf("в каталоге должна быть запись alias, получено %v", names)
	}
}

func ptrFloat64(f float64) *float64 {
	return &f
}

func ptrInt(i int) *int {
	return &i
}
