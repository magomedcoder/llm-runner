package template

import (
	"bytes"
	"embed"
	"encoding/json"
	"errors"
	"fmt"
	"path"
	"strings"
)

const maxPresetMatchDistance = 1000

//go:embed presets/*.gotmpl
//go:embed presets/*.json
var presetsFS embed.FS

type presetFile struct {
	Stop         []string `json:"stop"`
	Fingerprints []string `json:"fingerprints"`
}

func loadPresetEntries() ([]*MatchedPreset, error) {
	entries, err := presetsFS.ReadDir("presets")
	if err != nil {
		return nil, fmt.Errorf("каталог пресетов: %w", err)
	}

	var out []*MatchedPreset
	for _, e := range entries {
		if e.IsDir() || !strings.HasSuffix(e.Name(), ".gotmpl") {
			continue
		}

		stem := strings.TrimSuffix(e.Name(), ".gotmpl")
		raw, err := presetsFS.ReadFile(path.Join("presets", stem+".json"))
		if err != nil {
			continue
		}

		var cfg presetFile
		if err := json.Unmarshal(raw, &cfg); err != nil {
			return nil, fmt.Errorf("presets/%s.json: %w", stem, err)
		}

		if len(cfg.Fingerprints) == 0 {
			continue
		}

		for _, fp := range cfg.Fingerprints {
			fp = strings.TrimSpace(fp)
			if len(fp) < 80 {
				continue
			}

			out = append(out, &MatchedPreset{
				Name:       stem,
				IndexJinja: fp,
			})
		}
	}

	if len(out) == 0 {
		return nil, errors.New("пресеты: в presets/*.json нет записей с непустым полем fingerprint (отпечаток шаблона)")
	}

	return out, nil
}

func hydratePresetAssets(templates []*MatchedPreset) error {
	type stopOnly struct {
		Stop []string `json:"stop"`
	}

	gotmplBytes := map[string][]byte{}
	paramsByName := map[string]*struct {
		Stop []string `json:"stop"`
	}{}

	for _, t := range templates {
		if len(t.Bytes) > 0 {
			continue
		}

		if b, ok := gotmplBytes[t.Name]; ok {
			t.Bytes = b
			t.Parameters = paramsByName[t.Name]
			continue
		}

		bts, err := presetsFS.ReadFile("presets/" + t.Name + ".gotmpl")
		if err != nil {
			return fmt.Errorf("presets/%s.gotmpl: %w", t.Name, err)
		}
		bts = bytes.ReplaceAll(bts, []byte("\r\n"), []byte("\n"))
		gotmplBytes[t.Name] = bts

		var params *struct {
			Stop []string `json:"stop"`
		}
		raw, err := presetsFS.ReadFile("presets/" + t.Name + ".json")
		if err == nil {
			var so stopOnly
			if err := json.Unmarshal(raw, &so); err != nil {
				return fmt.Errorf("presets/%s.json: %w", t.Name, err)
			}

			if len(so.Stop) > 0 {
				params = &struct {
					Stop []string `json:"stop"`
				}{Stop: so.Stop}
			}
		}
		paramsByName[t.Name] = params

		t.Bytes = bts
		t.Parameters = params
	}

	return nil
}
