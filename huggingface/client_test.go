package huggingface

import (
	"context"
	"errors"
	"testing"
)

func TestModelInfoCanceledContext(t *testing.T) {
	ctx, cancel := context.WithCancel(context.Background())
	cancel()

	c := NewClient("")
	_, err := c.ModelInfo(ctx, "some/repo")
	if err == nil {
		t.Fatal("ожидалась ошибка при отменённом контексте")
	}

	if !errors.Is(err, context.Canceled) {
		t.Fatalf("ожидался context.Canceled, получено: %v", err)
	}
}
