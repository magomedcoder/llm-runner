package service

import (
	"fmt"
	"os"
	"path/filepath"
	"strings"

	"github.com/magomedcoder/llm-runner/domain"
	"gopkg.in/yaml.v3"
)

type ModelYAML struct {
	From      string              `yaml:"from"`
	System    string              `yaml:"system"`
	Template  string              `yaml:"template"`
	Parameter *ModelYAMLParameter `yaml:"parameter"`
}

type ModelYAMLParameter struct {
	Temperature *float64 `yaml:"temperature"`
	TopP        *float64 `yaml:"top_p"`
	TopK        *int     `yaml:"top_k"`
	MaxTokens   *int     `yaml:"max_tokens"`
}

func parseModelYAML(data []byte) (*ModelYAML, error) {
	var cfg ModelYAML
	if err := yaml.Unmarshal(data, &cfg); err != nil {
		return nil, err
	}
	return &cfg, nil
}

func yamlPathCandidates(modelsDir, stem string) []string {
	return []string{
		filepath.Join(modelsDir, stem+".yaml"),
		filepath.Join(modelsDir, stem+".yml"),
	}
}

func readYAMLByStem(modelsDir, stem string) ([]byte, error) {
	for _, p := range yamlPathCandidates(modelsDir, stem) {
		data, err := os.ReadFile(p)
		if err == nil {
			return data, nil
		}
	}
	return nil, os.ErrNotExist
}

func ManifestYAMLStemForRef(userInput string) string {
	name, tag := SplitModelRef(userInput)
	if tag != "" && !strings.EqualFold(tag, "latest") {
		return fmt.Sprintf("%s-%s", name, tag)
	}

	return name
}

func ResolveModelForInference(modelsDir, userInput string) (canonicalGGUF string, cfg *ModelYAML, err error) {
	if strings.TrimSpace(modelsDir) == "" {
		return "", nil, fmt.Errorf("каталог моделей не задан")
	}

	canonical, errGGUF := ResolveGGUFFile(modelsDir, userInput)
	if errGGUF == nil {
		stem := strings.TrimSuffix(canonical, ".gguf")
		data, e := readYAMLByStem(modelsDir, stem)
		if e != nil {
			return canonical, nil, nil
		}

		cfg, e := parseModelYAML(data)
		if e != nil {
			return "", nil, fmt.Errorf("YAML %s: %w", stem, e)
		}

		if strings.TrimSpace(cfg.From) != "" {
			cfg.From = ""
		}

		return canonical, cfg, nil
	}

	stem := ManifestYAMLStemForRef(userInput)
	data, err := readYAMLByStem(modelsDir, stem)
	if err != nil {
		return "", nil, errGGUF
	}

	cfg, err = parseModelYAML(data)
	if err != nil {
		return "", nil, fmt.Errorf("манифест %s.yaml: %w", stem, err)
	}

	from := strings.TrimSpace(cfg.From)
	if from == "" {
		return "", nil, fmt.Errorf("%w (в манифесте %s.yaml укажите from: базовый .gguf)", errGGUF, stem)
	}

	canonical, err = ResolveGGUFFile(modelsDir, from)
	if err != nil {
		return "", nil, fmt.Errorf("from в манифесте %q: %w", from, err)
	}

	return canonical, cfg, nil
}

func MergeGenParams(req *domain.GenerationParams, yamlCfg *ModelYAML) *domain.GenerationParams {
	if yamlCfg == nil || yamlCfg.Parameter == nil {
		return req
	}

	p := yamlCfg.Parameter
	out := domain.GenerationParams{}
	if req != nil {
		out.Temperature = req.Temperature
		out.MaxTokens = req.MaxTokens
		out.TopK = req.TopK
		out.TopP = req.TopP
		out.ResponseFormat = req.ResponseFormat
		if len(req.Tools) > 0 {
			out.Tools = append([]domain.Tool(nil), req.Tools...)
		}
	}

	if out.Temperature == nil && p.Temperature != nil {
		v := float32(*p.Temperature)
		out.Temperature = &v
	}

	if out.TopP == nil && p.TopP != nil {
		v := float32(*p.TopP)
		out.TopP = &v
	}

	if out.TopK == nil && p.TopK != nil {
		v := int32(*p.TopK)
		out.TopK = &v
	}

	if out.MaxTokens == nil && p.MaxTokens != nil {
		v := int32(*p.MaxTokens)
		out.MaxTokens = &v
	}

	return &out
}

func ApplyModelYAMLSystem(norm []*domain.AIChatMessage, yamlCfg *ModelYAML) []*domain.AIChatMessage {
	if yamlCfg == nil {
		return norm
	}

	extra := strings.TrimSpace(yamlCfg.System)
	if extra == "" {
		return norm
	}

	if len(norm) == 0 {
		return []*domain.AIChatMessage{domain.NewAIChatMessage(0, extra, domain.AIChatMessageRoleSystem)}
	}

	sid := norm[0].SessionId
	if norm[0].Role == domain.AIChatMessageRoleSystem {
		merged := extra + "\n\n" + norm[0].Content
		cp := *norm[0]
		cp.Content = merged
		out := make([]*domain.AIChatMessage, len(norm))
		out[0] = &cp
		copy(out[1:], norm[1:])
		return out
	}

	out := make([]*domain.AIChatMessage, 0, len(norm)+1)
	out = append(out, domain.NewAIChatMessage(sid, extra, domain.AIChatMessageRoleSystem))
	out = append(out, norm...)

	return out
}

func addManifestCatalogEntries(dir string, seen map[string]struct{}) error {
	entries, err := os.ReadDir(dir)
	if err != nil {
		return err
	}

	for _, e := range entries {
		if e.IsDir() {
			continue
		}

		name := e.Name()
		low := strings.ToLower(name)
		var stem string
		switch {
		case strings.HasSuffix(low, ".yaml"):
			stem = name[:len(name)-5]
		case strings.HasSuffix(low, ".yml"):
			stem = name[:len(name)-4]
		default:
			continue
		}

		if stem == "" {
			continue
		}

		gguf := filepath.Join(dir, stem+".gguf")
		if _, err := os.Stat(gguf); err == nil {
			continue
		}

		data, err := os.ReadFile(filepath.Join(dir, name))
		if err != nil {
			continue
		}

		cfg, err := parseModelYAML(data)
		if err != nil || strings.TrimSpace(cfg.From) == "" {
			continue
		}

		seen[stem] = struct{}{}
		if b, tg, ok := splitStemToTagged(stem); ok {
			seen[b+":"+tg] = struct{}{}
		}
	}

	return nil
}
