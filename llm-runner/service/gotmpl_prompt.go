package service

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"strings"

	"github.com/magomedcoder/gen/llm-runner/domain"
	"github.com/magomedcoder/gen/llm-runner/template"
)

func domainToolsToTemplate(tools []domain.Tool) template.Tools {
	if len(tools) == 0 {
		return nil
	}

	out := make(template.Tools, 0, len(tools))
	for _, dt := range tools {
		t := template.Tool{
			Type: "function",
			Function: template.ToolFunction{
				Name:        dt.Name,
				Description: dt.Description,
				Parameters: template.ToolFunctionParameters{
					Type: "object",
				},
			},
		}

		s := strings.TrimSpace(dt.ParametersJSON)
		if s != "" {
			var p template.ToolFunctionParameters
			if err := json.Unmarshal([]byte(s), &p); err == nil {
				t.Function.Parameters = p
			}
		}

		out = append(out, t)
	}

	return out
}

func RenderMatchedPreset(preset *template.MatchedPreset, norm []*domain.AIChatMessage, genParams *domain.GenerationParams) (string, error) {
	if preset == nil {
		return "", fmt.Errorf("пресет не задан (nil)")
	}

	raw, err := io.ReadAll(preset.Reader())
	if err != nil {
		return "", err
	}

	tmpl, err := template.Parse(strings.TrimSpace(string(raw)))
	if err != nil {
		return "", fmt.Errorf("пресет %q: разбор шаблона: %w", preset.Name, err)
	}

	msgs := make([]template.Message, 0, len(norm))
	for _, m := range norm {
		tm := template.Message{
			Role:       ChatRoleString(m.Role),
			Content:    m.Content,
			ToolCallID: m.ToolCallID,
			ToolName:   m.ToolName,
		}

		if strings.TrimSpace(m.ToolCallsJSON) != "" {
			calls, err := parseToolCallsJSON(m.ToolCallsJSON)
			if err != nil {
				return "", fmt.Errorf("разбор tool_calls_json: %w", err)
			}
			tm.ToolCalls = calls
		}

		msgs = append(msgs, tm)
	}

	var tools template.Tools
	if genParams != nil {
		tools = domainToolsToTemplate(genParams.Tools)
	}

	var buf bytes.Buffer
	if err := tmpl.Execute(&buf, template.Values{
		Messages: msgs,
		Tools:    tools,
	}); err != nil {
		return "", fmt.Errorf("пресет %q: выполнение шаблона: %w", preset.Name, err)
	}

	return buf.String(), nil
}

func BuildChatPromptGotmpl(chatTemplateJinja string, norm []*domain.AIChatMessage, genParams *domain.GenerationParams) (string, error) {
	j := strings.TrimSpace(chatTemplateJinja)
	if j == "" {
		return "", fmt.Errorf("пресет: пустой chat_template у модели")
	}

	p, err := template.Named(j)
	if err != nil {
		return "", err
	}

	return RenderMatchedPreset(p, norm, genParams)
}
