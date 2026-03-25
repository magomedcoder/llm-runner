package provider

import (
	"context"
	"fmt"
	"strings"

	"github.com/magomedcoder/llm-runner/config"
	"github.com/magomedcoder/llm-runner/domain"
	"github.com/magomedcoder/llm-runner/service"
)

type TextBackend interface {
	CheckConnection(ctx context.Context) (bool, error)

	WarmDefaultModel(ctx context.Context, model string) error

	GetModels(ctx context.Context) ([]string, error)

	GetLoadedModel(ctx context.Context) (loaded bool, ggufBasename, displayName string, err error)

	UnloadModel(ctx context.Context) error

	SendMessage(ctx context.Context, model string, messages []*domain.AIChatMessage, stopSequences []string, genParams *domain.GenerationParams) (chan string, error)

	Embed(ctx context.Context, model string, text string) ([]float32, error)
}

type TextProvider interface {
	CheckConnection(ctx context.Context) (bool, error)

	WarmDefaultModel(ctx context.Context, model string) error

	GetModels(ctx context.Context) ([]string, error)

	GetLoadedModel(ctx context.Context) (loaded bool, ggufBasename, displayName string, err error)

	UnloadModel(ctx context.Context) error

	SendMessage(ctx context.Context, sessionId int64, model string, messages []*domain.AIChatMessage, stopSequences []string, genParams *domain.GenerationParams) (chan string, error)

	Embed(ctx context.Context, model string, text string) ([]float32, error)
}

func NewTextProvider(cfg *config.Config) (TextProvider, error) {
	if cfg.ModelPath == "" {
		return nil, fmt.Errorf("укажите model_path")
	}

	var opts []service.LlamaOption
	nCtx := cfg.MaxContextTokens
	if nCtx <= 0 {
		nCtx = 4096
	}

	opts = append(opts, service.WithLlamaNCtx(nCtx))
	if cfg.MaxContextTokens > 0 {
		opts = append(opts, service.WithMaxContextTokens(cfg.MaxContextTokens))
	}

	opts = append(opts, service.WithEmbeddings(true))
	if strings.TrimSpace(cfg.MmprojPath) != "" {
		opts = append(opts, service.WithMmprojPath(cfg.MmprojPath))
	}
	svc := service.NewLlamaService(cfg.ModelPath, opts...)

	return NewText(svc), nil
}
