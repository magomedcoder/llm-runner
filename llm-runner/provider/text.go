package provider

import (
	"context"

	"github.com/magomedcoder/gen/llm-runner/domain"
)

type Text struct {
	backend TextBackend
}

func NewText(backend TextBackend) *Text {
	return &Text{
		backend: backend,
	}
}

func (t *Text) CheckConnection(ctx context.Context) (bool, error) {
	return t.backend.CheckConnection(ctx)
}

func (t *Text) WarmDefaultModel(ctx context.Context, model string) error {
	return t.backend.WarmDefaultModel(ctx, model)
}

func (t *Text) GetModels(ctx context.Context) ([]string, error) {
	return t.backend.GetModels(ctx)
}

func (t *Text) GetLoadedModel(ctx context.Context) (loaded bool, ggufBasename, displayName string, err error) {
	return t.backend.GetLoadedModel(ctx)
}

func (t *Text) UnloadModel(ctx context.Context) error {
	return t.backend.UnloadModel(ctx)
}

func (t *Text) SendMessage(ctx context.Context, sessionId int64, model string, messages []*domain.AIChatMessage, stopSequences []string, genParams *domain.GenerationParams) (chan domain.TextStreamChunk, error) {
	return t.backend.SendMessage(ctx, model, messages, stopSequences, genParams)
}

func (t *Text) Embed(ctx context.Context, model string, text string) ([]float32, error) {
	return t.backend.Embed(ctx, model, text)
}

func (t *Text) EmbedBatch(ctx context.Context, model string, texts []string) ([][]float32, error) {
	return t.backend.EmbedBatch(ctx, model, texts)
}
