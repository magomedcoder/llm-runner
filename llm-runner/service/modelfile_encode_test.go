package service

import (
	"strings"
	"testing"
)

func TestEncodeModelfile_roundtrip(t *testing.T) {
	temperature := 0.3
	maxTok := 100
	nctx := 2048
	cfg := &ModelYAML{
		From:   "Base.gguf",
		System: "be short",
		Stop:   []string{"<|eot|>", "end"},
		NumCtx: &nctx,
		Parameter: &ModelYAMLParameter{
			Temperature: &temperature,
			MaxTokens:   &maxTok,
		},
		Messages: []ModelYAMLMessage{
			{
				Role:    "user",
				Content: "hi",
			},
			{
				Role:    "assistant",
				Content: "yo",
			},
		},
	}

	s, err := EncodeModelfile(cfg, "")
	if err != nil {
		t.Fatal(err)
	}

	out, err := ParseModelfile(strings.NewReader(s))
	if err != nil {
		t.Fatalf("parse encoded: %v\n%s", err, s)
	}

	if out.From != "Base.gguf" || out.System != "be short" {
		t.Fatalf("from/system: %+v", out)
	}

	if out.NumCtx == nil || *out.NumCtx != 2048 {
		t.Fatalf("num_ctx %+v", out.NumCtx)
	}

	if len(out.Stop) != 2 {
		t.Fatalf("stop %#v", out.Stop)
	}

	if len(out.Messages) != 2 {
		t.Fatalf("messages %d", len(out.Messages))
	}

	if out.Parameter == nil || out.Parameter.MaxTokens == nil || *out.Parameter.MaxTokens != 100 {
		t.Fatalf("num_predict %+v", out.Parameter)
	}
}

func TestEncodeModelfile_sidecarFromWeights(t *testing.T) {
	cfg := &ModelYAML{
		System: "x",
	}
	s, err := EncodeModelfile(cfg, "M.gguf")
	if err != nil {
		t.Fatal(err)
	}

	if !strings.Contains(s, "FROM M.gguf") {
		t.Fatal(s)
	}
}

func TestParseModelYAMLData_EncodeBasic(t *testing.T) {
	raw := []byte("from: B.gguf\n")
	cfg, err := ParseModelYAMLData(raw)
	if err != nil {
		t.Fatal(err)
	}

	s, err := EncodeModelfile(cfg, "")
	if err != nil {
		t.Fatal(err)
	}

	if !strings.Contains(s, "FROM B.gguf") {
		t.Fatal(s)
	}
}
