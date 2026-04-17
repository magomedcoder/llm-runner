package runner

import (
	"encoding/json"
	"fmt"
	"regexp"
	"strings"
)

type cohereActionRow struct {
	ToolName   string          `json:"tool_name"`
	Parameters json.RawMessage `json:"parameters"`
}

var reActionJSON = regexp.MustCompile("(?is)(?:Action|Действие):\\s*" + "```" + `json\s*([\s\S]*?)` + "```")

func parseCohereActionList(blob string) ([]cohereActionRow, error) {
	blob = strings.TrimSpace(blob)
	if blob == "" {
		return nil, nil
	}

	var asSlice []cohereActionRow
	if err := json.Unmarshal([]byte(blob), &asSlice); err == nil {
		if len(asSlice) > 0 {
			return asSlice, nil
		}

		if strings.HasPrefix(strings.TrimSpace(blob), "[") {
			return nil, fmt.Errorf("пустой список вызовов инструментов")
		}
	}

	type legacyNameArgs struct {
		Name      string          `json:"name"`
		Arguments json.RawMessage `json:"arguments"`
	}
	var legacy legacyNameArgs
	if err := json.Unmarshal([]byte(blob), &legacy); err == nil && strings.TrimSpace(legacy.Name) != "" {
		args := legacy.Arguments
		if len(args) == 0 || string(args) == "null" {
			args = json.RawMessage(`{}`)
		}

		return []cohereActionRow{{
			ToolName:   strings.TrimSpace(legacy.Name),
			Parameters: args,
		}}, nil
	}

	type legacyToolParams struct {
		ToolName   string          `json:"tool_name"`
		Parameters json.RawMessage `json:"parameters"`
	}
	var tp legacyToolParams
	if err := json.Unmarshal([]byte(blob), &tp); err == nil && strings.TrimSpace(tp.ToolName) != "" {
		args := tp.Parameters
		if len(args) == 0 || string(args) == "null" {
			args = json.RawMessage(`{}`)
		}

		return []cohereActionRow{{
			ToolName:   strings.TrimSpace(tp.ToolName),
			Parameters: args,
		}}, nil
	}

	return nil, fmt.Errorf("неверный формат вызова инструментов (ожидается JSON-массив с tool_name/parameters или объект name/arguments)")
}

func toolActionRowsHaveNames(rows []cohereActionRow) bool {
	for _, r := range rows {
		if strings.TrimSpace(r.ToolName) != "" {
			return true
		}
	}

	return false
}

func extractCohereActionJSON(text string) string {
	m := reActionJSON.FindStringSubmatch(text)
	if len(m) < 2 {
		return ""
	}

	return strings.TrimSpace(m[1])
}

func extractFirstFencedToolArray(text string) string {
	s := text
	for len(s) > 0 {
		start := strings.Index(s, "```")
		if start < 0 {
			return ""
		}

		afterOpen := s[start+3:]
		bodyStart := 0
		if nl := strings.IndexByte(afterOpen, '\n'); nl >= 0 {
			first := strings.TrimSpace(afterOpen[:nl])
			if len(first) > 0 && !strings.ContainsAny(first, " \t") {
				bodyStart = nl + 1
			}
		}

		rest := afterOpen[bodyStart:]
		before, _, ok := strings.Cut(rest, "```")
		if !ok {
			return ""
		}

		raw := strings.TrimSpace(before)
		tr := strings.TrimSpace(raw)
		if strings.HasPrefix(tr, "[") || strings.HasPrefix(tr, "{") {
			if rows, err := parseCohereActionList(raw); err == nil && len(rows) > 0 && toolActionRowsHaveNames(rows) {
				return raw
			}
		}

		s = afterOpen
	}

	return ""
}

func extractFirstJSONArray(text string) string {
	_, after, ok := strings.Cut(text, "```json")
	if !ok {
		return ""
	}

	rest := after
	before, _, ok := strings.Cut(rest, "```")
	if !ok {
		return ""
	}

	raw := strings.TrimSpace(before)
	tr := strings.TrimSpace(raw)
	if !strings.HasPrefix(tr, "[") && !strings.HasPrefix(tr, "{") {
		return ""
	}

	if rows, err := parseCohereActionList(raw); err != nil || len(rows) == 0 || !toolActionRowsHaveNames(rows) {
		return ""
	}

	return raw
}

func extractLeadingJSONArray(text string) string {
	s := strings.TrimSpace(text)
	if len(s) == 0 || s[0] != '[' {
		return ""
	}

	depth := 0
	inString := false
	escape := false
	for i := 0; i < len(s); i++ {
		b := s[i]
		if escape {
			escape = false
			continue
		}

		if inString {
			if b == '\\' {
				escape = true
			} else if b == '"' {
				inString = false
			}

			continue
		}

		switch b {
		case '"':
			inString = true
		case '[':
			depth++
		case ']':
			depth--
			if depth == 0 {
				return s[:i+1]
			}
		}
	}

	return ""
}

func extractLeadingJSONObject(text string) string {
	s := strings.TrimSpace(text)
	if len(s) == 0 || s[0] != '{' {
		return ""
	}

	depth := 0
	inString := false
	escape := false
	for i := 0; i < len(s); i++ {
		b := s[i]
		if escape {
			escape = false
			continue
		}

		if inString {
			if b == '\\' {
				escape = true
			} else if b == '"' {
				inString = false
			}

			continue
		}

		switch b {
		case '"':
			inString = true
		case '{':
			depth++
		case '}':
			depth--
			if depth == 0 {
				return s[:i+1]
			}
		}
	}

	return ""
}

func extractEmbeddedJSONArray(text string) string {
	s := text
	for {
		idx := strings.Index(s, "[")
		if idx < 0 {
			return ""
		}

		sub := s[idx:]
		candidate := extractLeadingJSONArray(sub)
		if candidate != "" {
			rows, err := parseCohereActionList(candidate)
			if err == nil && len(rows) > 0 && toolActionRowsHaveNames(rows) {
				return candidate
			}
		}
		s = s[idx+1:]
	}
}

func extractEmbeddedJSONObject(text string) string {
	s := text
	for {
		idx := strings.Index(s, "{")
		if idx < 0 {
			return ""
		}

		sub := s[idx:]
		candidate := extractLeadingJSONObject(sub)
		if candidate != "" {
			rows, err := parseCohereActionList(candidate)
			if err == nil && len(rows) > 0 && toolActionRowsHaveNames(rows) {
				return candidate
			}
		}

		s = s[idx+1:]
	}
}

func ExtractToolActionBlob(text string) string {
	if s := extractCohereActionJSON(text); s != "" {
		return s
	}

	if s := extractFirstJSONArray(text); s != "" {
		return s
	}

	if s := extractFirstFencedToolArray(text); s != "" {
		return s
	}

	if s := extractLeadingJSONArray(text); s != "" {
		return s
	}

	if s := extractLeadingJSONObject(text); s != "" {
		if rows, err := parseCohereActionList(s); err == nil && len(rows) > 0 && toolActionRowsHaveNames(rows) {
			return s
		}
	}

	if s := extractEmbeddedJSONArray(text); s != "" {
		return s
	}

	return extractEmbeddedJSONObject(text)
}
