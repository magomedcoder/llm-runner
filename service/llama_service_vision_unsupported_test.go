//go:build llama

package service

import (
	"context"
	"strings"
	"testing"

	"github.com/magomedcoder/llm-runner/domain"
)

func TestSendMessage_VisionAttachmentUnsupported_FastFail(t *testing.T) {
	svc := NewLlamaService("")
	msgs := []*domain.AIChatMessage{
		{
			Role:              domain.AIChatMessageRoleUser,
			Content:           "опиши изображение",
			AttachmentName:    "image.png",
			AttachmentContent: []byte{1, 2, 3},
		},
	}

	_, err := svc.SendMessage(context.Background(), "any-model", msgs, nil, nil)
	if err == nil {
		t.Fatal("ожидалась ошибка о неподдерживаемых vision-вложениях")
	}

	if !strings.Contains(err.Error(), "vision-вложения не поддерживаются") {
		t.Fatalf("неожиданная ошибка: %v", err)
	}
}
