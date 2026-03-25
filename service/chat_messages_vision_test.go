package service

import (
	"testing"

	"github.com/magomedcoder/llm-runner/domain"
)

func TestNormalizeChatMessages_imageOnlyPayload(t *testing.T) {
	m := &domain.AIChatMessage{
		Role:              domain.AIChatMessageRoleUser,
		Content:           "",
		AttachmentContent: []byte{1, 2, 3},
		AttachmentName:    "a.png",
	}

	out := NormalizeChatMessages([]*domain.AIChatMessage{m})
	if len(out) != 1 {
		t.Fatalf("ожидалось 1 сообщение, получено %d", len(out))
	}
}
