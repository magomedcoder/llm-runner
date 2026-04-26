//go:build !llama

package service

import (
	"context"
	"fmt"

	"github.com/magomedcoder/gen/llm-runner/domain"
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

func WithGPULayers(int) LlamaOption {
	return func(s *LlamaService) {}
}

func WithMLock(bool) LlamaOption {
	return func(s *LlamaService) {}
}

func WithMMap(bool) LlamaOption {
	return func(s *LlamaService) {}
}

func WithMainGPU(string) LlamaOption {
	return func(s *LlamaService) {}
}

func WithTensorSplit(string) LlamaOption {
	return func(s *LlamaService) {}
}

func WithSilentLoading(bool) LlamaOption {
	return func(s *LlamaService) {}
}

func WithProgressCallback(bool) LlamaOption {
	return func(s *LlamaService) {}
}

func WithThreads(int) LlamaOption {
	return func(s *LlamaService) {}
}

func WithThreadsBatch(int) LlamaOption {
	return func(s *LlamaService) {}
}

func WithBatchSize(int) LlamaOption {
	return func(s *LlamaService) {}
}

func WithF16Memory(bool) LlamaOption {
	return func(s *LlamaService) {}
}

func WithKVCacheType(string) LlamaOption {
	return func(s *LlamaService) {}
}

func WithFlashAttn(string) LlamaOption {
	return func(s *LlamaService) {}
}

func WithPrefixCaching(bool) LlamaOption {
	return func(s *LlamaService) {}
}

func WithParallel(int) LlamaOption {
	return func(s *LlamaService) {}
}

func WithTopNSigma(float32) LlamaOption {
	return func(s *LlamaService) {}
}

func WithFrequencyPenalty(float32) LlamaOption {
	return func(s *LlamaService) {}
}

func WithPresencePenalty(float32) LlamaOption {
	return func(s *LlamaService) {}
}

func WithIgnoreEOS(bool) LlamaOption {
	return func(s *LlamaService) {}
}

func WithDRYMultiplier(float32) LlamaOption {
	return func(s *LlamaService) {}
}

func WithDRYBase(float32) LlamaOption {
	return func(s *LlamaService) {}
}

func WithDRYAllowedLength(int) LlamaOption {
	return func(s *LlamaService) {}
}

func WithDRYPenaltyLastN(int) LlamaOption {
	return func(s *LlamaService) {}
}

func WithDRYSequenceBreakers([]string) LlamaOption {
	return func(s *LlamaService) {}
}

func WithXTC(float32, float32) LlamaOption {
	return func(s *LlamaService) {}
}

func WithMirostat(int) LlamaOption {
	return func(s *LlamaService) {}
}

func WithMirostatTau(float32) LlamaOption {
	return func(s *LlamaService) {}
}

func WithMirostatEta(float32) LlamaOption {
	return func(s *LlamaService) {}
}

func WithTypicalP(float32) LlamaOption {
	return func(s *LlamaService) {}
}

func WithMinKeep(int) LlamaOption {
	return func(s *LlamaService) {}
}

func WithDynamicTemperature(float32, float32) LlamaOption {
	return func(s *LlamaService) {}
}

func WithNPrev(int) LlamaOption {
	return func(s *LlamaService) {}
}

func WithNProbs(int) LlamaOption {
	return func(s *LlamaService) {}
}

func WithDebugGeneration(bool) LlamaOption {
	return func(s *LlamaService) {}
}

func WithSpeculativeEnabled(bool) LlamaOption {
	return func(s *LlamaService) {}
}

func WithSpeculativeDraftModel(string) LlamaOption {
	return func(s *LlamaService) {}
}

func WithSpeculativeDraftTokens(int) LlamaOption {
	return func(s *LlamaService) {}
}

func WithTokenPipelineEnabled(bool) LlamaOption {
	return func(s *LlamaService) {}
}

func WithChatAPIEnabled(bool) LlamaOption {
	return func(s *LlamaService) {}
}

func WithChatStreamBufferSize(int) LlamaOption {
	return func(s *LlamaService) {}
}

func WithChatReasoningFormat(string) LlamaOption {
	return func(s *LlamaService) {}
}

func WithChatEnableThinking(bool) LlamaOption {
	return func(s *LlamaService) {}
}

func WithChatReasoningBudget(int) LlamaOption {
	return func(s *LlamaService) {}
}

func WithReinitLlamaLogging(bool) LlamaOption {
	return func(s *LlamaService) {}
}

func WithLogModelStats(bool) LlamaOption {
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

func (s *LlamaService) SendMessage(ctx context.Context, model string, messages []*domain.AIChatMessage, stopSequences []string, genParams *domain.GenerationParams) (chan domain.TextStreamChunk, error) {
	ch := make(chan domain.TextStreamChunk)
	close(ch)
	return ch, fmt.Errorf("llama отключена")
}

func (s *LlamaService) Embed(ctx context.Context, model string, text string) ([]float32, error) {
	return nil, fmt.Errorf("llama отключена")
}

func (s *LlamaService) EmbedBatch(ctx context.Context, model string, texts []string) ([][]float32, error) {
	return nil, fmt.Errorf("llama отключена")
}
