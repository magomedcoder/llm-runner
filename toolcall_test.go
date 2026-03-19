package runner

import (
	"testing"
)

func TestParseToolCalls_empty(t *testing.T) {
	if got := ParseToolCalls(""); len(got) != 0 {
		t.Errorf("ParseToolCalls(%q) = %v, ожидалось nil", "", got)
	}

	if got := ParseToolCalls("no json here"); len(got) != 0 {
		t.Errorf("ParseToolCalls(%q) = %v, ожидалось nil", "no json here", got)
	}
}

func TestParseToolCalls_single(t *testing.T) {
	content := `Some text {"name": "get_weather", "arguments": {"city": "Moscow"}} end`
	got := ParseToolCalls(content)
	if len(got) != 1 {
		t.Fatalf("ParseToolCalls: получено %d вызовов, ожидалось 1", len(got))
	}

	if got[0].Name != "get_weather" {
		t.Errorf("Name = %q, ожидалось get_weather", got[0].Name)
	}

	if got[0].Arguments != `{"city": "Moscow"}` {
		t.Errorf("Arguments = %q, ожидалось {\"city\": \"Moscow\"}", got[0].Arguments)
	}
}

func TestParseToolCalls_noName(t *testing.T) {
	content := `{"arguments": {}}`
	if got := ParseToolCalls(content); len(got) != 0 {
		t.Errorf("ParseToolCalls(%q) = %v, ожидалось nil (нет name)", content, got)
	}
}
