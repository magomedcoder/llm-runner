package llm_runner

import (
	"encoding/json"
	"strings"
	"testing"
)

func TestExtractToolActionBlob_rawPrefix(t *testing.T) {
	blob := ExtractToolActionBlob(`[{"tool_name":"apply_spreadsheet","parameters":{"operations_json":"[]"}}]`)

	var rows []cohereActionRow
	if err := json.Unmarshal([]byte(blob), &rows); err != nil || len(rows) != 1 || rows[0].ToolName != "apply_spreadsheet" {
		t.Fatalf("blob=%q err=%v rows=%v", blob, err, rows)
	}
}

func TestExtractToolActionBlob_legacyNameArgumentsObject(t *testing.T) {
	text := `{"name":"apply_spreadsheet","arguments":{"operations_json":"[]"}}`
	blob := ExtractToolActionBlob(text)
	rows, err := parseCohereActionList(blob)
	if err != nil || len(rows) != 1 || rows[0].ToolName != "apply_spreadsheet" {
		t.Fatalf("blob=%q err=%v rows=%+v", blob, err, rows)
	}
}

func TestExtractCohereActionJSON(t *testing.T) {
	text := `Краткое рассуждение здесь.

Действие: ` + "```json\n[\n  {\"tool_name\": \"apply_spreadsheet\", \"parameters\": {\"operations_json\": \"[]\"}}\n]\n```"
	got := extractCohereActionJSON(text)
	if !strings.Contains(got, "apply_spreadsheet") {
		t.Fatalf("ожидался JSON с apply_spreadsheet, получено: %q", got)
	}
}
