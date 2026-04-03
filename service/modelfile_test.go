package service

import (
	"os"
	"path/filepath"
	"strings"
	"testing"
)

func TestParseModelfile_basic(t *testing.T) {
	const src = `
# comment
FROM base.gguf
SYSTEM You are concise.
PARAMETER temperature 0.2
PARAMETER num_ctx 8192
`
	mf, err := ParseModelfile(strings.NewReader(src))
	if err != nil {
		t.Fatal(err)
	}

	if mf.From != "base.gguf" {
		t.Fatalf("from=%q", mf.From)
	}

	if mf.System != "You are concise." {
		t.Fatalf("system=%q", mf.System)
	}

	if mf.Parameter == nil || mf.Parameter.Temperature == nil || *mf.Parameter.Temperature != 0.2 {
		t.Fatalf("parameter=%+v", mf.Parameter)
	}

	if mf.NumCtx == nil || *mf.NumCtx != 8192 {
		t.Fatalf("num_ctx=%+v", mf.NumCtx)
	}
}

func TestParseModelfile_multilineSystem(t *testing.T) {
	src := "FROM m.gguf\nSYSTEM \"\"\"\nline1\nline2\n\"\"\"\n"
	mf, err := ParseModelfile(strings.NewReader(src))
	if err != nil {
		t.Fatal(err)
	}

	if mf.System != "line1\nline2" {
		t.Fatalf("system=%q", mf.System)
	}
}

func TestParseModelfile_templateTripleOneLine(t *testing.T) {
	src := "FROM m.gguf\nTEMPLATE \"\"\"{{ .NotUsed }}\"\"\"\n"
	mf, err := ParseModelfile(strings.NewReader(src))
	if err != nil {
		t.Fatal(err)
	}

	if mf.Template != "{{ .NotUsed }}" {
		t.Fatalf("template=%q", mf.Template)
	}
}

func TestWriteModelManifest_alias(t *testing.T) {
	dir := t.TempDir()
	if err := os.WriteFile(filepath.Join(dir, "Base.gguf"), []byte{0}, 0o644); err != nil {
		t.Fatal(err)
	}

	cfg := &ModelYAML{
		From:   "Base.gguf",
		System: "x",
		Parameter: &ModelYAMLParameter{
			TopP: new(0.9),
		},
	}

	if err := WriteModelManifest(dir, "alias", cfg, false); err != nil {
		t.Fatal(err)
	}

	can, out, err := ResolveModelForInference(dir, "alias")
	if err != nil || can != "Base.gguf" || out == nil || out.System != "x" {
		t.Fatalf("resolve: can=%q cfg=%+v err=%v", can, out, err)
	}
}

func TestWriteModelManifest_sidecarOmitsFrom(t *testing.T) {
	dir := t.TempDir()
	if err := os.WriteFile(filepath.Join(dir, "Same.gguf"), []byte{0}, 0o644); err != nil {
		t.Fatal(err)
	}

	cfg := &ModelYAML{
		From:   "Same.gguf",
		System: "s",
	}
	if err := WriteModelManifest(dir, "Same", cfg, false); err != nil {
		t.Fatal(err)
	}

	raw, err := os.ReadFile(filepath.Join(dir, "Same.yaml"))
	if err != nil {
		t.Fatal(err)
	}

	if strings.Contains(string(raw), "from:") {
		t.Fatalf("sidecar не должен содержать from, получено:\n%s", raw)
	}
}

func TestParseModelfile_sampling(t *testing.T) {
	src := `
FROM b.gguf
PARAMETER repeat_last_n 128
PARAMETER repeat_penalty 1.15
PARAMETER seed 1
PARAMETER min_p 0.05
`
	mf, err := ParseModelfile(strings.NewReader(src))
	if err != nil {
		t.Fatal(err)
	}

	if mf.Parameter == nil || mf.Parameter.RepeatLastN == nil || *mf.Parameter.RepeatLastN != 128 {
		t.Fatalf("param=%+v", mf.Parameter)
	}
}

func TestParseModelfile_adapterUnsupported(t *testing.T) {
	src := "FROM b.gguf\nADAPTER ./lora.gguf\n"
	_, err := ParseModelfile(strings.NewReader(src))
	if err == nil || !strings.Contains(err.Error(), "ADAPTER не поддерживается") {
		t.Fatalf("expected unsupported ADAPTER error, got %v", err)
	}
}

func TestParseModelfile_messageAndStop(t *testing.T) {
	src := `
FROM m.gguf
PARAMETER stop "</think>"
PARAMETER stop <|eot|>
MESSAGE user Is Paris in France?
MESSAGE assistant yes
MESSAGE user Is London in France?
MESSAGE assistant no
REQUIRES 0.99.0
`
	mf, err := ParseModelfile(strings.NewReader(src))
	if err != nil {
		t.Fatal(err)
	}

	if len(mf.Stop) != 2 || mf.Stop[0] != "</think>" || mf.Stop[1] != "<|eot|>" {
		t.Fatalf("stop=%#v", mf.Stop)
	}

	if len(mf.Messages) != 4 {
		t.Fatalf("messages=%d", len(mf.Messages))
	}

	if mf.Messages[0].Role != "user" || mf.Messages[1].Content != "yes" {
		t.Fatalf("msg0=%+v msg1=%+v", mf.Messages[0], mf.Messages[1])
	}
}

func TestParseModelfile_messageMultiline(t *testing.T) {
	src := "FROM x.gguf\nMESSAGE user \"\"\"\nhi\nthere\n\"\"\"\n"
	mf, err := ParseModelfile(strings.NewReader(src))
	if err != nil {
		t.Fatal(err)
	}

	if len(mf.Messages) != 1 || mf.Messages[0].Content != "hi\nthere" {
		t.Fatalf("%+v", mf.Messages)
	}
}

func TestWriteModelManifest_conflictStemGguf(t *testing.T) {
	dir := t.TempDir()
	_ = os.WriteFile(filepath.Join(dir, "A.gguf"), []byte{0}, 0o644)
	_ = os.WriteFile(filepath.Join(dir, "B.gguf"), []byte{0}, 0o644)
	cfg := &ModelYAML{
		From: "B.gguf",
	}
	err := WriteModelManifest(dir, "A", cfg, false)
	if err == nil || !strings.Contains(err.Error(), "уже есть") {
		t.Fatalf("ожидалась ошибка конфликта, получено %v", err)
	}
}
