//go:build !llama

package service

import (
	"context"
	"fmt"

	"github.com/magomedcoder/llm-runner/domain"
)

type LlamaService struct{}

type LlamaOption func(*LlamaService)

func NewLlamaService(modelPath string, opts ...LlamaOption) *LlamaService {
	return &LlamaService{}
}

func WithMaxContextTokens(n int) LlamaOption {
	return func(s *LlamaService) {}
}

func WithLlamaNCtx(n int) LlamaOption {
	return func(s *LlamaService) {}
}

func WithEmbeddings(enable bool) LlamaOption {
	return func(s *LlamaService) {}
}

func WithMmprojPath(string) LlamaOption {
	return func(s *LlamaService) {}
}

func (s *LlamaService) WarmDefaultModel(ctx context.Context, model string) error {
	return fmt.Errorf("llama отключена")
}

func (s *LlamaService) CheckConnection(ctx context.Context) (bool, error) {
	return false, fmt.Errorf("llama отключена")
}

func (s *LlamaService) GetModels(ctx context.Context) ([]string, error) {
	return nil, fmt.Errorf("llama отключена")
}

func (s *LlamaService) GetLoadedModel(ctx context.Context) (loaded bool, ggufBasename, displayName string, err error) {
	return false, "", "", nil
}

func (s *LlamaService) UnloadModel(ctx context.Context) error {
	return nil
}

func (s *LlamaService) SendMessage(ctx context.Context, model string, messages []*domain.AIChatMessage, stopSequences []string, genParams *domain.GenerationParams) (chan string, error) {
	ch := make(chan string)
	close(ch)
	return ch, fmt.Errorf("llama отключена")
}

func (s *LlamaService) Embed(ctx context.Context, model string, text string) ([]float32, error) {
	return nil, fmt.Errorf("llama отключена")
}
