package llama

import (
	gocontext "context"
	"fmt"
	"runtime"
	"runtime/cgo"
	"sync"
	"unsafe"
)

/*
#include "wrapper.h"
#include <stdlib.h>
*/
import "C"

type Context struct {
	contextPtr unsafe.Pointer // указатель llama_wrapper_context_t*
	model      *Model
	config     contextConfig
	mu         sync.RWMutex
	closed     bool
}

func (c *Context) Close() error {
	c.mu.Lock()
	defer c.mu.Unlock()

	if c.closed {
		return nil
	}

	runtime.SetFinalizer(c, nil)

	if c.contextPtr != nil {
		C.llama_wrapper_context_free(c.contextPtr)
		c.contextPtr = nil
	}

	c.closed = true
	return nil
}

func (c *Context) Tokenize(text string) ([]int32, error) {
	c.mu.RLock()
	defer c.mu.RUnlock()

	if c.closed {
		return nil, fmt.Errorf("контекст закрыт")
	}

	cText := C.CString(text)
	defer C.free(unsafe.Pointer(cText))

	var tokensPtr *C.int
	var count C.int

	C.llama_wrapper_tokenize_alloc(c.contextPtr, cText, &tokensPtr, &count)

	if tokensPtr != nil {
		defer C.llama_wrapper_free_tokens(tokensPtr)
	}

	if count < 0 || tokensPtr == nil {
		return nil, fmt.Errorf("ошибка токенизации: %s", C.GoString(C.llama_wrapper_last_error()))
	}

	tokens := (*[1 << 30]C.int)(unsafe.Pointer(tokensPtr))[:count:count]
	result := make([]int32, count)
	for i := 0; i < int(count); i++ {
		result[i] = int32(tokens[i])
	}

	return result, nil
}

func (c *Context) GetCachedTokenCount() (int, error) {
	c.mu.RLock()
	defer c.mu.RUnlock()

	if c.closed {
		return 0, fmt.Errorf("контекст закрыт")
	}

	count := int(C.llama_wrapper_get_cached_token_count(c.contextPtr))
	if count < 0 {
		return 0, fmt.Errorf("не удалось получить количество кэшированных токенов: %s", C.GoString(C.llama_wrapper_last_error()))
	}

	return count, nil
}

func (c *Context) GetEmbeddings(text string) ([]float32, error) {
	c.mu.RLock()
	defer c.mu.RUnlock()

	if c.closed {
		return nil, fmt.Errorf("контекст закрыт")
	}

	cText := C.CString(text)
	defer C.free(unsafe.Pointer(cText))

	maxEmbeddings := 4096
	embeddings := make([]C.float, maxEmbeddings)

	count := C.llama_wrapper_embeddings(c.contextPtr, cText, &embeddings[0], C.int(maxEmbeddings))
	if count < 0 {
		return nil, fmt.Errorf("ошибка генерации эмбеддинга: %s", C.GoString(C.llama_wrapper_last_error()))
	}

	result := make([]float32, count)
	for i := 0; i < int(count); i++ {
		result[i] = float32(embeddings[i])
	}

	return result, nil
}

func (c *Context) GetEmbeddingsBatch(texts []string) ([][]float32, error) {
	c.mu.RLock()
	defer c.mu.RUnlock()

	if c.closed {
		return nil, fmt.Errorf("контекст закрыт")
	}

	if len(texts) == 0 {
		return nil, fmt.Errorf("не переданы тексты")
	}

	nEmbd := int(C.llama_wrapper_model_n_embd(c.model.modelPtr))
	if nEmbd <= 0 {
		return nil, fmt.Errorf("некорректная размерность эмбеддинга: %d", nEmbd)
	}

	cTexts := make([]*C.char, len(texts))
	for i, text := range texts {
		cTexts[i] = C.CString(text)
	}
	defer func() {
		for i := range cTexts {
			C.free(unsafe.Pointer(cTexts[i]))
		}
	}()

	outputSize := len(texts) * nEmbd
	cEmbeddings := make([]C.float, outputSize)

	count := C.llama_wrapper_embeddings_batch(
		c.contextPtr,
		(**C.char)(unsafe.Pointer(&cTexts[0])),
		C.int(len(texts)),
		&cEmbeddings[0],
		C.int(nEmbd),
	)

	if count < 0 {
		return nil, fmt.Errorf("ошибка пакетной генерации эмбеддингов: %s", C.GoString(C.llama_wrapper_last_error()))
	}

	if int(count) != len(texts) {
		return nil, fmt.Errorf("несовпадение количества эмбеддингов: ожидалось %d, получено %d", len(texts), count)
	}

	result := make([][]float32, len(texts))
	for i := 0; i < len(texts); i++ {
		result[i] = make([]float32, nEmbd)
		for j := 0; j < nEmbd; j++ {
			result[i][j] = float32(cEmbeddings[i*nEmbd+j])
		}
	}

	return result, nil
}

func (c *Context) Generate(prompt string, opts ...GenerateOption) (string, error) {
	config := defaultGenerateConfig
	for _, opt := range opts {
		opt(&config)
	}

	return c.generateWithConfig(prompt, config, nil)
}

func (c *Context) GenerateStream(prompt string, callback func(token string) bool, opts ...GenerateOption) error {
	config := defaultGenerateConfig
	for _, opt := range opts {
		opt(&config)
	}

	_, err := c.generateWithConfig(prompt, config, callback)
	return err
}

func (c *Context) GenerateChannel(ctx gocontext.Context, prompt string, opts ...GenerateOption) (<-chan string, <-chan error) {
	tokenChan := make(chan string, 10)
	errChan := make(chan error, 1)

	go func() {
		defer close(tokenChan)
		defer close(errChan)

		err := c.GenerateStream(prompt, func(token string) bool {
			select {
			case <-ctx.Done():
				return false
			case tokenChan <- token:
				return true
			}
		}, opts...)

		if err != nil {
			errChan <- err
		}
	}()

	return tokenChan, errChan
}

func (c *Context) GenerateWithTokens(tokens []int32, opts ...GenerateOption) (string, error) {
	config := defaultGenerateConfig
	for _, opt := range opts {
		opt(&config)
	}

	return c.generateWithTokensAndConfig(tokens, config, nil)
}

func (c *Context) GenerateWithTokensStream(tokens []int32, callback func(token string) bool, opts ...GenerateOption) error {
	config := defaultGenerateConfig
	for _, opt := range opts {
		opt(&config)
	}

	_, err := c.generateWithTokensAndConfig(tokens, config, callback)
	return err
}

func (c *Context) GenerateWithDraft(prompt string, draft *Context, opts ...GenerateOption) (string, error) {
	config := defaultGenerateConfig
	for _, opt := range opts {
		opt(&config)
	}

	return c.generateWithDraftAndConfig(prompt, draft, config, nil)
}

func (c *Context) GenerateWithDraftStream(prompt string, draft *Context, callback func(token string) bool, opts ...GenerateOption) error {
	config := defaultGenerateConfig
	for _, opt := range opts {
		opt(&config)
	}

	_, err := c.generateWithDraftAndConfig(prompt, draft, config, callback)
	return err
}

func (c *Context) GenerateWithDraftChannel(ctx gocontext.Context, prompt string, draft *Context, opts ...GenerateOption) (<-chan string, <-chan error) {
	tokenChan := make(chan string, 10)
	errChan := make(chan error, 1)

	go func() {
		defer close(tokenChan)
		defer close(errChan)

		err := c.GenerateWithDraftStream(prompt, draft, func(token string) bool {
			select {
			case <-ctx.Done():
				return false
			case tokenChan <- token:
				return true
			}
		}, opts...)

		if err != nil {
			errChan <- err
		}
	}()

	return tokenChan, errChan
}

func (c *Context) Chat(ctx gocontext.Context, messages []ChatMessage, opts ChatOptions) (*ChatResponse, error) {
	return c.model.chatWithContext(ctx, c, messages, opts)
}

func (c *Context) ChatStream(ctx gocontext.Context, messages []ChatMessage, opts ChatOptions) (<-chan ChatDelta, <-chan error) {
	return c.model.chatStreamWithContext(ctx, c, messages, opts)
}

//export goTokenCallback
func goTokenCallback(handle C.uintptr_t, token *C.char) C.bool {
	h := cgo.Handle(handle)
	callback := h.Value().(func(string) bool)
	return C.bool(callback(C.GoString(token)))
}

func findCommonPrefix(a, b []int32) int {
	commonLen := 0
	for i := 0; i < len(a) && i < len(b); i++ {
		if a[i] != b[i] {
			break
		}

		commonLen++
	}

	return commonLen
}

func (c *Context) generateWithConfig(prompt string, config generateConfig, callback func(string) bool) (string, error) {
	c.mu.Lock()
	defer c.mu.Unlock()

	if c.closed {
		return "", fmt.Errorf("контекст закрыт")
	}

	if c.model == nil {
		return "", fmt.Errorf("модель закрыта")
	}
	c.model.mu.RLock()
	modelClosed := c.model.closed
	c.model.mu.RUnlock()
	if modelClosed {
		return "", fmt.Errorf("модель закрыта")
	}

	cPrompt := C.CString(prompt)
	defer C.free(unsafe.Pointer(cPrompt))

	var cStopWords **C.char
	var stopWordsCount C.int

	if len(config.stopWords) > 0 {
		stopWordsCount = C.int(len(config.stopWords))
		cStopWordsArray := make([]*C.char, len(config.stopWords))
		for i, word := range config.stopWords {
			cStopWordsArray[i] = C.CString(word)
		}
		defer func() {
			for _, ptr := range cStopWordsArray {
				C.free(unsafe.Pointer(ptr))
			}
		}()
		cStopWords = (**C.char)(unsafe.Pointer(&cStopWordsArray[0]))
	}

	var handle cgo.Handle
	var callbackHandle C.uintptr_t
	if callback != nil {
		handle = cgo.NewHandle(callback)
		callbackHandle = C.uintptr_t(handle)
		defer handle.Delete()
	}

	var cDryBreakers **C.char
	var dryBreakersCount C.int
	if len(config.drySequenceBreakers) > 0 {
		dryBreakersCount = C.int(len(config.drySequenceBreakers))
		cDryBreakersArray := make([]*C.char, len(config.drySequenceBreakers))
		for i, breaker := range config.drySequenceBreakers {
			cDryBreakersArray[i] = C.CString(breaker)
		}
		defer func() {
			for _, ptr := range cDryBreakersArray {
				C.free(unsafe.Pointer(ptr))
			}
		}()
		cDryBreakers = (**C.char)(unsafe.Pointer(&cDryBreakersArray[0]))
	}

	params := C.llama_wrapper_generate_params{
		prompt:                      cPrompt,
		max_tokens:                  C.int(config.maxTokens),
		temperature:                 C.float(config.temperature),
		top_k:                       C.int(config.topK),
		top_p:                       C.float(config.topP),
		min_p:                       C.float(config.minP),
		typ_p:                       C.float(config.typP),
		top_n_sigma:                 C.float(config.topNSigma),
		penalty_last_n:              C.int(config.penaltyLastN),
		penalty_repeat:              C.float(config.penaltyRepeat),
		penalty_freq:                C.float(config.penaltyFreq),
		penalty_present:             C.float(config.penaltyPresent),
		dry_multiplier:              C.float(config.dryMultiplier),
		dry_base:                    C.float(config.dryBase),
		dry_allowed_length:          C.int(config.dryAllowedLength),
		dry_penalty_last_n:          C.int(config.dryPenaltyLastN),
		dry_sequence_breakers:       cDryBreakers,
		dry_sequence_breakers_count: dryBreakersCount,
		dynatemp_range:              C.float(config.dynatempRange),
		dynatemp_exponent:           C.float(config.dynatempExponent),
		xtc_probability:             C.float(config.xtcProbability),
		xtc_threshold:               C.float(config.xtcThreshold),
		mirostat:                    C.int(config.mirostat),
		mirostat_tau:                C.float(config.mirostatTau),
		mirostat_eta:                C.float(config.mirostatEta),
		n_prev:                      C.int(config.nPrev),
		n_probs:                     C.int(config.nProbs),
		min_keep:                    C.int(config.minKeep),
		seed:                        C.int(config.seed),
		stop_words:                  cStopWords,
		stop_words_count:            stopWordsCount,
		callback_handle:             callbackHandle,
		ignore_eos:                  C.bool(config.ignoreEOS),
		debug:                       C.bool(config.debug),
	}

	cResult := C.llama_wrapper_generate(c.contextPtr, params)
	if cResult == nil {
		return "", fmt.Errorf("ошибка генерации: %s", C.GoString(C.llama_wrapper_last_error()))
	}

	result := C.GoString(cResult)
	C.llama_wrapper_free_result(cResult)

	return result, nil
}

func (c *Context) generateWithTokensAndConfig(tokens []int32, config generateConfig, callback func(string) bool) (string, error) {
	c.mu.Lock()
	defer c.mu.Unlock()

	if c.closed {
		return "", fmt.Errorf("контекст закрыт")
	}

	if len(tokens) == 0 {
		return "", fmt.Errorf("токены не переданы")
	}

	cTokens := make([]C.int, len(tokens))
	for i, token := range tokens {
		cTokens[i] = C.int(token)
	}

	var cStopWords **C.char
	var stopWordsCount C.int

	if len(config.stopWords) > 0 {
		stopWordsCount = C.int(len(config.stopWords))
		cStopWordsArray := make([]*C.char, len(config.stopWords))
		for i, word := range config.stopWords {
			cStopWordsArray[i] = C.CString(word)
		}

		defer func() {
			for _, ptr := range cStopWordsArray {
				C.free(unsafe.Pointer(ptr))
			}
		}()

		cStopWords = (**C.char)(unsafe.Pointer(&cStopWordsArray[0]))
	}

	var handle cgo.Handle
	var callbackHandle C.uintptr_t
	if callback != nil {
		handle = cgo.NewHandle(callback)
		callbackHandle = C.uintptr_t(handle)
		defer handle.Delete()
	}

	var cDryBreakers **C.char
	var dryBreakersCount C.int
	if len(config.drySequenceBreakers) > 0 {
		dryBreakersCount = C.int(len(config.drySequenceBreakers))
		cDryBreakersArray := make([]*C.char, len(config.drySequenceBreakers))
		for i, breaker := range config.drySequenceBreakers {
			cDryBreakersArray[i] = C.CString(breaker)
		}

		defer func() {
			for _, ptr := range cDryBreakersArray {
				C.free(unsafe.Pointer(ptr))
			}
		}()

		cDryBreakers = (**C.char)(unsafe.Pointer(&cDryBreakersArray[0]))
	}

	params := C.llama_wrapper_generate_params{
		prompt:                      nil, // Не используется при генерации по токенам
		max_tokens:                  C.int(config.maxTokens),
		temperature:                 C.float(config.temperature),
		top_k:                       C.int(config.topK),
		top_p:                       C.float(config.topP),
		min_p:                       C.float(config.minP),
		typ_p:                       C.float(config.typP),
		top_n_sigma:                 C.float(config.topNSigma),
		penalty_last_n:              C.int(config.penaltyLastN),
		penalty_repeat:              C.float(config.penaltyRepeat),
		penalty_freq:                C.float(config.penaltyFreq),
		penalty_present:             C.float(config.penaltyPresent),
		dry_multiplier:              C.float(config.dryMultiplier),
		dry_base:                    C.float(config.dryBase),
		dry_allowed_length:          C.int(config.dryAllowedLength),
		dry_penalty_last_n:          C.int(config.dryPenaltyLastN),
		dry_sequence_breakers:       cDryBreakers,
		dry_sequence_breakers_count: dryBreakersCount,
		dynatemp_range:              C.float(config.dynatempRange),
		dynatemp_exponent:           C.float(config.dynatempExponent),
		xtc_probability:             C.float(config.xtcProbability),
		xtc_threshold:               C.float(config.xtcThreshold),
		mirostat:                    C.int(config.mirostat),
		mirostat_tau:                C.float(config.mirostatTau),
		mirostat_eta:                C.float(config.mirostatEta),
		n_prev:                      C.int(config.nPrev),
		n_probs:                     C.int(config.nProbs),
		min_keep:                    C.int(config.minKeep),
		seed:                        C.int(config.seed),
		stop_words:                  cStopWords,
		stop_words_count:            stopWordsCount,
		callback_handle:             callbackHandle,
		ignore_eos:                  C.bool(config.ignoreEOS),
		debug:                       C.bool(config.debug),
	}

	cResult := C.llama_wrapper_generate_with_tokens(
		c.contextPtr,
		&cTokens[0],
		C.int(len(tokens)),
		C.int(0), // длина префикса: префиксное кэширование для этой функции не используется
		params,
	)

	if cResult == nil {
		return "", fmt.Errorf("ошибка генерации по токенам: %s", C.GoString(C.llama_wrapper_last_error()))
	}

	result := C.GoString(cResult)
	C.llama_wrapper_free_result(cResult)

	return result, nil
}

func (c *Context) generateWithDraftAndConfig(prompt string, draft *Context, config generateConfig, callback func(string) bool) (string, error) {
	c.mu.Lock()
	defer c.mu.Unlock()

	if c.closed {
		return "", fmt.Errorf("контекст закрыт")
	}

	if c.model == nil {
		return "", fmt.Errorf("модель закрыта")
	}
	c.model.mu.RLock()
	modelClosed := c.model.closed
	c.model.mu.RUnlock()
	if modelClosed {
		return "", fmt.Errorf("модель закрыта")
	}

	draft.mu.RLock()
	if draft.closed {
		draft.mu.RUnlock()
		return "", fmt.Errorf("draft-контекст закрыт")
	}
	draftPtr := draft.contextPtr
	draft.mu.RUnlock()

	if draft.model == nil {
		return "", fmt.Errorf("draft-модель закрыта")
	}
	draft.model.mu.RLock()
	draftModelClosed := draft.model.closed
	draft.model.mu.RUnlock()
	if draftModelClosed {
		return "", fmt.Errorf("draft-модель закрыта")
	}

	cPrompt := C.CString(prompt)
	defer C.free(unsafe.Pointer(cPrompt))

	var cStopWords **C.char
	var stopWordsCount C.int

	if len(config.stopWords) > 0 {
		stopWordsCount = C.int(len(config.stopWords))
		cStopWordsArray := make([]*C.char, len(config.stopWords))
		for i, word := range config.stopWords {
			cStopWordsArray[i] = C.CString(word)
		}

		defer func() {
			for _, ptr := range cStopWordsArray {
				C.free(unsafe.Pointer(ptr))
			}
		}()

		cStopWords = (**C.char)(unsafe.Pointer(&cStopWordsArray[0]))
	}

	var handle cgo.Handle
	var callbackHandle C.uintptr_t
	if callback != nil {
		handle = cgo.NewHandle(callback)
		callbackHandle = C.uintptr_t(handle)
		defer handle.Delete()
	}

	var cDryBreakers **C.char
	var dryBreakersCount C.int
	if len(config.drySequenceBreakers) > 0 {
		dryBreakersCount = C.int(len(config.drySequenceBreakers))
		cDryBreakersArray := make([]*C.char, len(config.drySequenceBreakers))
		for i, breaker := range config.drySequenceBreakers {
			cDryBreakersArray[i] = C.CString(breaker)
		}

		defer func() {
			for _, ptr := range cDryBreakersArray {
				C.free(unsafe.Pointer(ptr))
			}
		}()

		cDryBreakers = (**C.char)(unsafe.Pointer(&cDryBreakersArray[0]))
	}

	params := C.llama_wrapper_generate_params{
		prompt:                      cPrompt,
		max_tokens:                  C.int(config.maxTokens),
		temperature:                 C.float(config.temperature),
		top_k:                       C.int(config.topK),
		top_p:                       C.float(config.topP),
		min_p:                       C.float(config.minP),
		typ_p:                       C.float(config.typP),
		top_n_sigma:                 C.float(config.topNSigma),
		penalty_last_n:              C.int(config.penaltyLastN),
		penalty_repeat:              C.float(config.penaltyRepeat),
		penalty_freq:                C.float(config.penaltyFreq),
		penalty_present:             C.float(config.penaltyPresent),
		dry_multiplier:              C.float(config.dryMultiplier),
		dry_base:                    C.float(config.dryBase),
		dry_allowed_length:          C.int(config.dryAllowedLength),
		dry_penalty_last_n:          C.int(config.dryPenaltyLastN),
		dry_sequence_breakers:       cDryBreakers,
		dry_sequence_breakers_count: dryBreakersCount,
		dynatemp_range:              C.float(config.dynatempRange),
		dynatemp_exponent:           C.float(config.dynatempExponent),
		xtc_probability:             C.float(config.xtcProbability),
		xtc_threshold:               C.float(config.xtcThreshold),
		mirostat:                    C.int(config.mirostat),
		mirostat_tau:                C.float(config.mirostatTau),
		mirostat_eta:                C.float(config.mirostatEta),
		n_prev:                      C.int(config.nPrev),
		n_probs:                     C.int(config.nProbs),
		min_keep:                    C.int(config.minKeep),
		seed:                        C.int(config.seed),
		stop_words:                  cStopWords,
		stop_words_count:            stopWordsCount,
		callback_handle:             callbackHandle,
		ignore_eos:                  C.bool(config.ignoreEOS),
		debug:                       C.bool(config.debug),
	}

	cResult := C.llama_wrapper_generate_draft(
		c.contextPtr,
		draftPtr,
		params,
	)

	if cResult == nil {
		return "", fmt.Errorf("ошибка draft-генерации: %s", C.GoString(C.llama_wrapper_last_error()))
	}

	result := C.GoString(cResult)
	C.llama_wrapper_free_result(cResult)

	return result, nil
}
