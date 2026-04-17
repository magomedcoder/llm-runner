package service

import (
	"errors"
	"os"
	"path/filepath"
	"testing"

	"github.com/magomedcoder/gen-runner/domain"
)

func TestMergeGenParams_yamlDefaults(t *testing.T) {
	y := &ModelYAML{
		Parameter: &ModelYAMLParameter{
			Temperature: new(0.1),
			MaxTokens:   new(99),
			MinP:        new(0.07),
		},
	}
	out := MergeGenParams(nil, y)
	if out == nil || out.Temperature == nil || *out.Temperature != 0.1 {
		t.Fatalf("temperature из yaml: %+v", out)
	}

	if out.MaxTokens == nil || *out.MaxTokens != 99 {
		t.Fatalf("max_tokens из yaml: %+v", out)
	}
	if out.MinP == nil || *out.MinP != float32(0.07) {
		t.Fatalf("min_p из yaml: %+v", out)
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
	if out2.MinP == nil || *out2.MinP != float32(0.07) {
		t.Fatalf("min_p должен браться из yaml: %+v", out2)
	}

	reqMinP := float32(0.22)
	out3 := MergeGenParams(&domain.GenerationParams{MinP: &reqMinP}, y)
	if out3.MinP == nil || *out3.MinP != float32(0.22) {
		t.Fatalf("min_p из запроса должен перекрывать yaml: %+v", out3)
	}
}

func TestApplyModelYAMLMessages(t *testing.T) {
	cfg := &ModelYAML{
		Messages: []ModelYAMLMessage{
			{
				Role:    "user",
				Content: "u1",
			},
			{
				Role:    "assistant",
				Content: "a1",
			},
		},
	}
	norm := []*domain.AIChatMessage{
		domain.NewAIChatMessage(7, "sys", domain.AIChatMessageRoleSystem),
		domain.NewAIChatMessage(7, "hi", domain.AIChatMessageRoleUser),
	}
	out := ApplyModelYAMLMessages(norm, cfg)
	if len(out) != 4 {
		t.Fatalf("len=%d", len(out))
	}

	if out[0].Content != "sys" || out[1].Role != domain.AIChatMessageRoleUser || out[1].Content != "u1" {
		t.Fatalf("order: %+v %+v", out[0], out[1])
	}

	if out[3].Content != "hi" {
		t.Fatalf("last user: %+v", out[3])
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

func TestApplyModelYAMLSystem_existingSystemKeepsRuntimePolicyLast(t *testing.T) {
	norm := []*domain.AIChatMessage{
		domain.NewAIChatMessage(1, "runtime-policy", domain.AIChatMessageRoleSystem),
		domain.NewAIChatMessage(1, "answer in json", domain.AIChatMessageRoleUser),
	}

	cfg := &ModelYAML{
		System: "model-hint",
	}
	out := ApplyModelYAMLSystem(norm, cfg)
	if len(out) != 2 {
		t.Fatalf("len=%d", len(out))
	}

	if out[0].Role != domain.AIChatMessageRoleSystem {
		t.Fatalf("first must be system, got %s", out[0].Role)
	}

	if out[0].Content != "model-hint\n\nruntime-policy" {
		t.Fatalf("unexpected merged system content: %q", out[0].Content)
	}

	if out[1].Role != domain.AIChatMessageRoleUser || out[1].Content != "answer in json" {
		t.Fatalf("user instruction must remain unchanged and last: %+v", out[1])
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

func TestLoadManifestYAMLForShow_sidecar(t *testing.T) {
	dir := t.TempDir()
	_ = os.WriteFile(filepath.Join(dir, "M.gguf"), []byte{0}, 0o644)
	_ = os.WriteFile(filepath.Join(dir, "M.yaml"), []byte("system: x\n"), 0o644)
	p, raw, w, err := LoadManifestYAMLForShow(dir, "M")
	if err != nil || w != "M.gguf" || len(raw) == 0 {
		t.Fatalf("path=%q w=%q err=%v raw=%q", p, w, err, raw)
	}
}

func TestLoadManifestYAMLForShow_alias(t *testing.T) {
	dir := t.TempDir()
	_ = os.WriteFile(filepath.Join(dir, "B.gguf"), []byte{0}, 0o644)
	_ = os.WriteFile(filepath.Join(dir, "a.yaml"), []byte("from: B.gguf\n"), 0o644)
	p, raw, w, err := LoadManifestYAMLForShow(dir, "a")
	if err != nil || w != "B.gguf" || len(raw) == 0 {
		t.Fatalf("path=%q w=%q err=%v", p, w, err)
	}
}

func TestLoadManifestYAMLForShow_noYaml(t *testing.T) {
	dir := t.TempDir()
	_ = os.WriteFile(filepath.Join(dir, "M.gguf"), []byte{0}, 0o644)
	_, _, _, err := LoadManifestYAMLForShow(dir, "M")
	if !errors.Is(err, ErrNoSidecarManifest) {
		t.Fatalf("ожидался ErrNoSidecarManifest, получено %v", err)
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
