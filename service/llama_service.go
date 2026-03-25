//go:build llama

package service

import (
	"context"
	"fmt"
	"os"
	"path/filepath"
	"strings"
	"sync"

	"github.com/magomedcoder/llm-runner/domain"
	"github.com/magomedcoder/llm-runner/llama"
	"github.com/magomedcoder/llm-runner/template"
)

const defaultChunkSize = 128

type LlamaService struct {
	modelsDir        string
	currentModelName string
	mmprojPathConfig string
	chunkSize        int
	predictOpts      []llama.PredictOption
	mu               sync.RWMutex
	model            *llama.LLama
	maxContextTokens int
	llamaNCtx        int
	enableEmbeddings bool
}

type LlamaOption func(*LlamaService)

func WithChunkSize(n int) LlamaOption {
	return func(s *LlamaService) {
		if n > 0 {
			s.chunkSize = n
		}
	}
}

func WithPredictOptions(opts ...llama.PredictOption) LlamaOption {
	return func(s *LlamaService) {
		s.predictOpts = opts
	}
}

func WithMaxContextTokens(n int) LlamaOption {
	return func(s *LlamaService) {
		if n > 0 {
			s.maxContextTokens = n
		}
	}
}

func WithLlamaNCtx(n int) LlamaOption {
	return func(s *LlamaService) {
		if n > 0 {
			s.llamaNCtx = n
		}
	}
}

func WithEmbeddings(enable bool) LlamaOption {
	return func(s *LlamaService) {
		s.enableEmbeddings = enable
	}
}

func WithMmprojPath(path string) LlamaOption {
	return func(s *LlamaService) {
		s.mmprojPathConfig = strings.TrimSpace(path)
	}
}

func NewLlamaService(modelPath string, opts ...LlamaOption) *LlamaService {
	modelsDir := modelPath
	if modelPath != "" {
		if info, err := os.Stat(modelPath); err == nil && !info.IsDir() {
			modelsDir = filepath.Dir(modelPath)
		}
	}

	s := &LlamaService{
		modelsDir: modelsDir,
		chunkSize: defaultChunkSize,
		llamaNCtx: 4096,
	}

	for _, opt := range opts {
		opt(s)
	}

	if s.chunkSize <= 0 {
		s.chunkSize = defaultChunkSize
	}

	return s
}

func (s *LlamaService) applyModelChatTemplate(norm []*domain.AIChatMessage, addGen bool, visionInject bool, mediaMarker string, chatTemplateOverride string) (string, error) {
	if s.model == nil {
		return "", fmt.Errorf("llama: модель не загружена")
	}

	roles := make([]string, len(norm))
	contents := make([]string, len(norm))
	for i, m := range norm {
		roles[i] = ChatRoleString(m.Role)
		contents[i] = FormatContentForChatTemplateWithVision(m, visionInject, mediaMarker)
	}

	return s.model.ApplyChatTemplate(strings.TrimSpace(chatTemplateOverride), roles, contents, addGen)
}

func (s *LlamaService) resolveChatPrompt(norm []*domain.AIChatMessage, genParams *domain.GenerationParams, visionInject bool, mediaMarker string, chatTemplateOverride string) (prompt string, presetStops []string, err error) {
	jinja := strings.TrimSpace(s.model.GetChatTemplate(""))
	var matched *template.MatchedPreset
	if jinja != "" {
		if p, e := template.Named(jinja); e == nil {
			matched = p
			presetStops = template.PresetStopSequences(p)
		}
	}

	if p, e := s.applyModelChatTemplate(norm, true, visionInject, mediaMarker, chatTemplateOverride); e == nil && p != "" {
		return p, presetStops, nil
	}

	if matched != nil && !visionInject && strings.TrimSpace(chatTemplateOverride) == "" {
		text, e := RenderMatchedPreset(matched, norm, genParams)
		if e != nil {
			return "", presetStops, fmt.Errorf("llama: пресет %q: %w", matched.Name, e)
		}

		if strings.TrimSpace(text) != "" {
			return text, presetStops, nil
		}
	}

	return fallbackPlainChatPrompt(norm, genParams, visionInject, mediaMarker), presetStops, nil
}

func (s *LlamaService) ensureModel(modelName string) (*ModelYAML, string, error) {
	s.mu.Lock()
	defer s.mu.Unlock()

	if s.modelsDir == "" {
		return nil, "", fmt.Errorf("llama: путь к папке с моделями не задан")
	}

	if strings.TrimSpace(modelName) == "" {
		return nil, "", fmt.Errorf("llama: укажите модель (доступные: %s)", strings.Join(s.modelDisplayNamesLocked(), ", "))
	}

	resolved, yamlCfg, err := ResolveModelForInference(s.modelsDir, modelName)
	if err != nil {
		return nil, "", fmt.Errorf("llama: %w (доступные: %s)", err, strings.Join(s.modelDisplayNamesLocked(), ", "))
	}

	mmprojPath, err := ResolveMmprojPath(s.modelsDir, resolved, s.mmprojPathConfig)
	if err != nil {
		return nil, "", fmt.Errorf("llama: %w", err)
	}

	fullPath := filepath.Join(s.modelsDir, resolved)
	if s.model != nil && s.currentModelName == resolved {
		return yamlCfg, mmprojPath, nil
	}

	if s.model != nil {
		s.model.Free()
		s.model = nil
		s.currentModelName = ""
	}

	var modelOpts []llama.ModelOption
	if s.enableEmbeddings {
		modelOpts = append(modelOpts, llama.EnableEmbeddings)
	}
	nCtx := s.llamaNCtx
	if nCtx <= 0 {
		nCtx = 4096
	}
	modelOpts = append(modelOpts, llama.SetContext(nCtx))
	if strings.TrimSpace(mmprojPath) != "" {
		modelOpts = append(modelOpts, llama.WithMmproj(mmprojPath))
	}

	m, err := llama.New(fullPath, modelOpts...)
	if err != nil {
		return nil, "", fmt.Errorf("llama: не удалось загрузить модель %q: %w", DisplayModelName(resolved), err)
	}

	s.model = m
	s.currentModelName = resolved
	return yamlCfg, mmprojPath, nil
}

func (s *LlamaService) modelDisplayNamesLocked() []string {
	if s.modelsDir == "" {
		return nil
	}

	names, err := SortedDisplayModelNames(s.modelsDir)
	if err != nil {
		return nil
	}
	return names
}

func (s *LlamaService) WarmDefaultModel(ctx context.Context, model string) error {
	if err := ctx.Err(); err != nil {
		return err
	}

	_, _, err := s.ensureModel(strings.TrimSpace(model))
	return err
}

func (s *LlamaService) CheckConnection(ctx context.Context) (bool, error) {
	s.mu.RLock()
	defer s.mu.RUnlock()
	files, err := ListGGUFBasenames(s.modelsDir)
	if err != nil || len(files) == 0 {
		return false, fmt.Errorf("llama: нет моделей в папке %q", s.modelsDir)
	}
	return true, nil
}

func (s *LlamaService) GetModels(ctx context.Context) ([]string, error) {
	s.mu.RLock()
	defer s.mu.RUnlock()
	out := s.modelDisplayNamesLocked()
	if out == nil {
		return []string{}, nil
	}

	return out, nil
}

func (s *LlamaService) GetLoadedModel(ctx context.Context) (loaded bool, ggufBasename, displayName string, err error) {
	s.mu.RLock()
	defer s.mu.RUnlock()
	if s.model == nil || strings.TrimSpace(s.currentModelName) == "" {
		return false, "", "", nil
	}

	return true, s.currentModelName, DisplayModelName(s.currentModelName), nil
}

func (s *LlamaService) UnloadModel(ctx context.Context) error {
	s.mu.Lock()
	defer s.mu.Unlock()
	if s.model != nil {
		s.model.Free()
		s.model = nil
		s.currentModelName = ""
	}

	return nil
}

func (s *LlamaService) SendMessage(ctx context.Context, model string, messages []*domain.AIChatMessage, stopSequences []string, genParams *domain.GenerationParams) (chan string, error) {
	yamlCfg, mmprojPath, err := s.ensureModel(model)
	if err != nil {
		return nil, err
	}

	norm := NormalizeChatMessages(messages)
	norm = ApplyModelYAMLSystem(norm, yamlCfg)
	if len(norm) == 0 {
		return nil, fmt.Errorf("llama: пустой список сообщений (нет текста content)")
	}

	if err := errIfVisionAttachments(norm, mmprojPath, s.currentModelName, s.modelsDir); err != nil {
		return nil, err
	}

	hasVision := MessagesHaveVisionAttachments(norm)
	if hasVision && !s.model.HasMTMD() {
		return nil, fmt.Errorf("llama: vision-запрос, но mmproj не загружен (проверьте путь и сборку с libmtmd)")
	}

	genMerged := MergeGenParams(genParams, yamlCfg)

	visionInject := hasVision
	mediaMarker := llama.MTMDDefaultMarker()
	chatTmpl := ""
	if yamlCfg != nil {
		chatTmpl = yamlCfg.Template
	}
	prompt, presetStops, err := s.resolveChatPrompt(norm, genMerged, visionInject, mediaMarker, chatTmpl)
	if err != nil {
		return nil, err
	}
	stopForPredict := MergeStopSequences(stopSequences, presetStops)
	stopForPredict = MergeStopSequences(stopForPredict, inferStopSequencesFromPrompt(prompt))

	if s.maxContextTokens > 0 {
		approxTokens := len(prompt)/4 + 1
		if approxTokens > s.maxContextTokens {
			return nil, fmt.Errorf("llama: контекст слишком велик (≈%d токенов, лимит %d)", approxTokens, s.maxContextTokens)
		}
	}

	out := make(chan string, 32)
	go func() {
		defer close(out)
		opts := make([]llama.PredictOption, 0, len(s.predictOpts)+6)
		opts = append(opts, s.predictOpts...)
		if genMerged != nil {
			if genMerged.Temperature != nil {
				opts = append(opts, llama.SetTemperature(*genMerged.Temperature))
			}

			if genMerged.MaxTokens != nil && *genMerged.MaxTokens > 0 {
				opts = append(opts, llama.SetTokens(int(*genMerged.MaxTokens)))
			}

			if genMerged.TopK != nil && *genMerged.TopK > 0 {
				opts = append(opts, llama.SetTopK(int(*genMerged.TopK)))
			}

			if genMerged.TopP != nil {
				opts = append(opts, llama.SetTopP(*genMerged.TopP))
			}

			if genMerged.ResponseFormat != nil && genMerged.ResponseFormat.Type == "json_object" {
				grammar := DefaultJSONObjectGrammar
				if genMerged.ResponseFormat.Schema != nil && *genMerged.ResponseFormat.Schema != "" {
					grammar = *genMerged.ResponseFormat.Schema
				}

				if grammar != "" {
					opts = append(opts, llama.WithGrammar(grammar))
				}
			}
		}

		if len(stopForPredict) > 0 {
			opts = append(opts, llama.SetStopWords(stopForPredict...))
		}

		emitToken := func(token string) {
			if token == "" {
				return
			}
			select {
			case <-ctx.Done():
			case out <- token:
			}
		}

		streamFilter := newStopStreamFilter(stopForPredict, emitToken)

		opts = append(opts, llama.SetTokenCallback(func(token string) bool {
			select {
			case <-ctx.Done():
				return false
			default:
				if token == "" {
					return true
				}
				streamFilter.push(token)
				return true
			}
		}))
		imgs := CollectVisionImageBytes(norm)
		s.mu.Lock()
		var err error
		if len(imgs) > 0 {
			_, err = s.model.PredictMTMD(prompt, imgs, opts...)
		} else {
			_, err = s.model.Predict(prompt, opts...)
		}
		s.mu.Unlock()
		streamFilter.flush()
		if err != nil {
			return
		}
	}()
	return out, nil
}

func (s *LlamaService) Embed(ctx context.Context, model string, text string) ([]float32, error) {
	if _, _, err := s.ensureModel(model); err != nil {
		return nil, err
	}

	s.mu.Lock()
	defer s.mu.Unlock()
	if s.model == nil {
		return nil, fmt.Errorf("llama: модель не загружена")
	}

	return s.model.Embeddings(text, s.predictOpts...)
}
