package template

import "testing"

func TestLoadAllPresets_embeddedPresets(t *testing.T) {
	presets, err := loadAllPresets()
	if err != nil {
		t.Fatal(err)
	}

	if len(presets) < 5 {
		t.Fatalf("ожидалось несколько встроенных пресетов, получено %d", len(presets))
	}

	for _, p := range presets {
		if len(p.Bytes) == 0 {
			t.Fatalf("у пресета %q пустые байты шаблона .gotmpl", p.Name)
		}
	}
}
