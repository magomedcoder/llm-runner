package runner

import (
	"encoding/json"
	"strings"

	"github.com/magomedcoder/llm-runner/domain"
)

func ParseToolCalls(content string) []*domain.ToolCall {
	content = strings.TrimSpace(content)
	if content == "" {
		return nil
	}
	start := strings.Index(content, "{")
	if start < 0 {
		return nil
	}
	depth := 0
	end := -1
loop:
	for i := start; i < len(content); i++ {
		switch content[i] {
		case '{':
			depth++
		case '}':
			depth--
			if depth == 0 {
				end = i + 1
				break loop
			}
		}
	}
	if end <= start {
		return nil
	}

	jsonStr := content[start:end]
	var raw struct {
		Name      string          `json:"name"`
		Arguments json.RawMessage `json:"arguments"`
	}

	if err := json.Unmarshal([]byte(jsonStr), &raw); err != nil || raw.Name == "" {
		return nil
	}
	argsStr := string(raw.Arguments)
	if argsStr == "" || argsStr == "null" {
		argsStr = "{}"
	}

	if len(raw.Arguments) > 0 && raw.Arguments[0] != '"' {
		argsStr = string(raw.Arguments)
	}

	return []*domain.ToolCall{{
		Id:        "",
		Name:      raw.Name,
		Arguments: argsStr,
	}}
}
