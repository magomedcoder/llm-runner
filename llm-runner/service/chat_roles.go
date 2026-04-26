package service

import "github.com/magomedcoder/gen/llm-runner/domain"

func ChatRoleString(role domain.AIChatMessageRole) string {
	switch role {
	case domain.AIChatMessageRoleSystem:
		return "system"
	case domain.AIChatMessageRoleAssistant:
		return "assistant"
	case domain.AIChatMessageRoleTool:
		return "tool"
	default:
		return "user"
	}
}
