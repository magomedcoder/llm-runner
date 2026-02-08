package repository

import (
	"context"

	"github.com/magomedcoder/gen/internal/domain"
	"github.com/magomedcoder/gen/internal/service"
)

type LLMRunnerRepository struct {
	client *service.LLMRunnerService
}

// NewLLMRunnerRepository создаёт репозиторий для llm-runner (gRPC по llmrunner.proto).
// address — адрес llm-runner, например "localhost:50052".
func NewLLMRunnerRepository(address, model string) (*LLMRunnerRepository, error) {
	client, err := service.NewLLMRunnerService(address, model)
	if err != nil {
		return nil, err
	}
	return &LLMRunnerRepository{client: client}, nil
}

func (r *LLMRunnerRepository) CheckConnection(ctx context.Context) (bool, error) {
	return r.client.CheckConnection(ctx)
}

func (r *LLMRunnerRepository) SendMessage(ctx context.Context, sessionID string, messages []*domain.Message) (chan string, error) {
	return r.client.SendMessage(ctx, sessionID, messages)
}

// Close закрывает соединение с llm-runner.
func (r *LLMRunnerRepository) Close() error {
	return r.client.Close()
}
