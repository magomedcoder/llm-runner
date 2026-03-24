package repository

import (
	"context"

	"github.com/magomedcoder/gen/api/pb/llmrunnerpb"
	"github.com/magomedcoder/gen/internal/domain"
	"github.com/magomedcoder/gen/internal/service"
)

type LLMRunnerRepository struct {
	client *service.LLMRunnerService
}

func NewLLMRunnerRepository(address string) (*LLMRunnerRepository, error) {
	client, err := service.NewLLMRunnerService(address, "")
	if err != nil {
		return nil, err
	}
	return &LLMRunnerRepository{client: client}, nil
}

func (r *LLMRunnerRepository) CheckConnection(ctx context.Context) (bool, error) {
	return r.client.CheckConnection(ctx)
}

func (r *LLMRunnerRepository) GetModels(ctx context.Context) ([]string, error) {
	return r.client.GetModels(ctx)
}

func (r *LLMRunnerRepository) SendMessage(
	ctx context.Context,
	sessionID int64,
	model string,
	messages []*domain.Message,
	stopSequences []string,
	timeoutSeconds int32,
	genParams *domain.GenerationParams,
) (chan string, error) {
	return r.client.SendMessage(ctx, sessionID, model, messages, stopSequences, timeoutSeconds, genParams)
}

func (r *LLMRunnerRepository) Close() error {
	return r.client.Close()
}

func (r *LLMRunnerRepository) GetGpuInfo(ctx context.Context) (*llmrunnerpb.GetGpuInfoResponse, error) {
	return r.client.GetGpuInfo(ctx)
}

func (r *LLMRunnerRepository) GetServerInfo(ctx context.Context) (*llmrunnerpb.ServerInfo, error) {
	return r.client.GetServerInfo(ctx)
}
