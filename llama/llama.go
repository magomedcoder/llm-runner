//go:build llama
// +build llama

package llama

// #cgo CXXFLAGS: -I${SRCDIR}/llama.cpp_lib/common -I${SRCDIR}/llama.cpp_lib/include -I${SRCDIR}/llama.cpp_lib/ggml/include -I${SRCDIR}/llama.cpp_lib -std=c++17
// #cgo LDFLAGS: -L${SRCDIR}/ -lllama -lm -lstdc++
// #cgo linux LDFLAGS: -fopenmp
// #include "llama.h"
// #include <stdlib.h>
import "C"
import (
	"fmt"
	"os"
	"strings"
	"sync"
	"unsafe"
)

type LLama struct {
	state       unsafe.Pointer
	embeddings  bool
	contextSize int
}

func New(model string, opts ...ModelOption) (*LLama, error) {
	mo := NewModelOptions(opts...)
	modelPath := C.CString(model)
	defer C.free(unsafe.Pointer(modelPath))

	loraBase := C.CString(mo.LoraBase)
	defer C.free(unsafe.Pointer(loraBase))

	loraAdapter := C.CString(mo.LoraAdapter)
	defer C.free(unsafe.Pointer(loraAdapter))

	mainGPU := C.CString(mo.MainGPU)
	defer C.free(unsafe.Pointer(mainGPU))
	tensorSplit := C.CString(mo.TensorSplit)
	defer C.free(unsafe.Pointer(tensorSplit))

	var mmprojArg *C.char
	if mp := strings.TrimSpace(mo.MmprojPath); mp != "" {
		mmprojArg = C.CString(mp)
		defer C.free(unsafe.Pointer(mmprojArg))
	}

	result := C.load_model(
		modelPath,
		C.int(mo.ContextSize),
		C.int(mo.Seed),
		C.bool(mo.F16Memory),
		C.bool(mo.MLock),
		C.bool(mo.Embeddings),
		C.bool(mo.MMap),
		C.bool(mo.LowVRAM),
		C.int(mo.NGPULayers),
		C.int(mo.NBatch),
		mainGPU,
		tensorSplit,
		C.bool(mo.NUMA),
		C.float(mo.FreqRopeBase),
		C.float(mo.FreqRopeScale),
		loraAdapter, loraBase,
		mmprojArg,
	)

	if result == nil {
		return nil, fmt.Errorf("не удалось загрузить модель")
	}

	ll := &LLama{
		state:       result,
		contextSize: mo.ContextSize,
		embeddings:  mo.Embeddings,
	}
	return ll, nil
}

func (l *LLama) Free() {
	C.llama_binding_free_model(l.state)
}

type ModelInfo struct {
	VocabSize     int
	ContextLength int
	EmbeddingSize int
	LayerCount    int
	ModelSize     int64
	ParamCount    int64
	Description   string
}

func (l *LLama) GetModelInfo() ModelInfo {
	descBuf := make([]byte, 256)
	C.get_model_description(l.state, (*C.char)(unsafe.Pointer(&descBuf[0])), C.int(len(descBuf)))

	return ModelInfo{
		VocabSize:     int(C.get_model_n_vocab(l.state)),
		ContextLength: int(C.get_model_n_ctx_train(l.state)),
		EmbeddingSize: int(C.get_model_n_embd(l.state)),
		LayerCount:    int(C.get_model_n_layer(l.state)),
		ModelSize:     int64(C.get_model_size(l.state)),
		ParamCount:    int64(C.get_model_n_params(l.state)),
		Description:   string(descBuf[:cStrLen(descBuf)]),
	}
}

func (l *LLama) GetChatTemplate(name string) string {
	buf := make([]byte, 4096)
	var cName *C.char
	if name != "" {
		cName = C.CString(name)
		defer C.free(unsafe.Pointer(cName))
	}

	ret := C.get_model_chat_template(l.state, cName, (*C.char)(unsafe.Pointer(&buf[0])), C.int(len(buf)))
	if ret <= 0 {
		return ""
	}

	return string(buf[:ret])
}

func (l *LLama) ApplyChatTemplate(tmplOverride string, roles []string, contents []string, addGenerationPrompt bool) (string, error) {
	n := len(roles)
	if n == 0 || n != len(contents) {
		return "", fmt.Errorf("шаблон чата: нужны пары роль/содержимое одинаковой длины")
	}

	rolePtrs := make([]*C.char, n)
	contentPtrs := make([]*C.char, n)
	defer func() {
		for i := 0; i < n; i++ {
			C.free(unsafe.Pointer(rolePtrs[i]))
			C.free(unsafe.Pointer(contentPtrs[i]))
		}
	}()

	for i := 0; i < n; i++ {
		rolePtrs[i] = C.CString(roles[i])
		contentPtrs[i] = C.CString(contents[i])
	}

	var tmplC *C.char
	if tmplOverride != "" {
		tmplC = C.CString(tmplOverride)
		defer C.free(unsafe.Pointer(tmplC))
	}

	need := int(C.apply_chat_template(
		l.state,
		tmplC,
		(**C.char)(unsafe.Pointer(&rolePtrs[0])),
		(**C.char)(unsafe.Pointer(&contentPtrs[0])),
		C.size_t(n),
		C.bool(addGenerationPrompt),
		(*C.char)(nil),
		C.int(0),
	))
	if need < 0 {
		return "", fmt.Errorf("шаблон чата: модель не распознала шаблон или ошибка (%d)", need)
	}

	buf := make([]byte, need+1)
	got := int(C.apply_chat_template(
		l.state,
		tmplC,
		(**C.char)(unsafe.Pointer(&rolePtrs[0])),
		(**C.char)(unsafe.Pointer(&contentPtrs[0])),
		C.size_t(n),
		C.bool(addGenerationPrompt),
		(*C.char)(unsafe.Pointer(&buf[0])),
		C.int(len(buf)),
	))
	if got < 0 {
		return "", fmt.Errorf("шаблон чата: ошибка при повторном применении (%d)", got)
	}

	return string(buf[:cStrLen(buf)]), nil
}

func cStrLen(b []byte) int {
	for i, v := range b {
		if v == 0 {
			return i
		}
	}
	return len(b)
}

func (l *LLama) LoadState(state string) error {
	d := C.CString(state)
	w := C.CString("rb")
	result := C.load_state(l.state, d, w)

	defer C.free(unsafe.Pointer(d))
	defer C.free(unsafe.Pointer(w))

	if result != 0 {
		return fmt.Errorf("ошибка при загрузке состояния")
	}

	return nil
}

func (l *LLama) SaveState(dst string) error {
	d := C.CString(dst)
	w := C.CString("wb")

	C.save_state(l.state, d, w)

	defer C.free(unsafe.Pointer(d))
	defer C.free(unsafe.Pointer(w))

	_, err := os.Stat(dst)
	return err
}

func (l *LLama) TokenEmbeddings(tokens []int, opts ...PredictOption) ([]float32, error) {
	if !l.embeddings {
		return []float32{}, fmt.Errorf("модель загружена без поддержки эмбеддингов")
	}

	po := NewPredictOptions(opts...)

	outSize := po.Tokens
	if po.Tokens == 0 {
		outSize = 9999999
	}

	floats := make([]float32, outSize)

	myArray := (*C.int)(C.malloc(C.size_t(len(tokens)) * C.sizeof_int))

	for i, v := range tokens {
		(*[1<<31 - 1]int32)(unsafe.Pointer(myArray))[i] = int32(v)
	}

	params := C.llama_allocate_params(
		C.CString(""),
		C.int(po.Seed),
		C.int(po.Threads),
		C.int(po.Tokens),
		C.int(po.TopK),
		C.float(po.TopP),
		C.float(po.MinP),
		C.float(po.Temperature),
		C.float(po.Penalty),
		C.int(po.Repeat),
		C.bool(po.IgnoreEOS),
		C.bool(po.F16KV),
		C.int(po.Batch),
		C.int(po.NKeep),
		nil,
		C.int(0),
		C.float(po.TailFreeSamplingZ),
		C.float(po.TypicalP),
		C.float(po.FrequencyPenalty),
		C.float(po.PresencePenalty),
		C.int(po.Mirostat),
		C.float(po.MirostatETA),
		C.float(po.MirostatTAU),
		C.bool(po.PenalizeNL),
		C.CString(po.LogitBias),
		C.CString(po.PathPromptCache),
		C.bool(po.PromptCacheAll),
		C.bool(po.MLock),
		C.bool(po.MMap),
		C.CString(po.MainGPU),
		C.CString(po.TensorSplit),
		C.bool(po.PromptCacheRO),
		C.CString(po.Grammar),
		C.float(po.RopeFreqBase),
		C.float(po.RopeFreqScale),
		C.int(po.NDraft),
		C.float(po.XTCProbability),
		C.float(po.XTCThreshold),
		C.float(po.DRYMultiplier),
		C.float(po.DRYBase),
		C.int(po.DRYAllowedLength),
		C.int(po.DRYPenaltyLastN),
		C.float(po.TopNSigma),
	)
	ret := C.get_token_embeddings(params, l.state, myArray, C.int(len(tokens)), (*C.float)(&floats[0]))
	if ret != 0 {
		return floats, fmt.Errorf("ошибка вывода эмбеддингов")
	}

	return floats, nil
}

func (l *LLama) Embeddings(text string, opts ...PredictOption) ([]float32, error) {
	if !l.embeddings {
		return []float32{}, fmt.Errorf("модель загружена без поддержки эмбеддингов")
	}

	po := NewPredictOptions(opts...)

	input := C.CString(text)
	if po.Tokens == 0 {
		po.Tokens = 99999999
	}

	floats := make([]float32, po.Tokens)
	reverseCount := len(po.StopPrompts)
	reversePrompt := make([]*C.char, reverseCount)
	var pass **C.char
	for i, s := range po.StopPrompts {
		cs := C.CString(s)
		reversePrompt[i] = cs
		pass = &reversePrompt[0]
	}

	params := C.llama_allocate_params(
		input,
		C.int(po.Seed),
		C.int(po.Threads),
		C.int(po.Tokens),
		C.int(po.TopK),
		C.float(po.TopP),
		C.float(po.MinP),
		C.float(po.Temperature),
		C.float(po.Penalty),
		C.int(po.Repeat),
		C.bool(po.IgnoreEOS),
		C.bool(po.F16KV),
		C.int(po.Batch),
		C.int(po.NKeep),
		pass,
		C.int(reverseCount),
		C.float(po.TailFreeSamplingZ),
		C.float(po.TypicalP),
		C.float(po.FrequencyPenalty),
		C.float(po.PresencePenalty),
		C.int(po.Mirostat),
		C.float(po.MirostatETA),
		C.float(po.MirostatTAU),
		C.bool(po.PenalizeNL),
		C.CString(po.LogitBias),
		C.CString(po.PathPromptCache),
		C.bool(po.PromptCacheAll),
		C.bool(po.MLock),
		C.bool(po.MMap),
		C.CString(po.MainGPU),
		C.CString(po.TensorSplit),
		C.bool(po.PromptCacheRO),
		C.CString(po.Grammar),
		C.float(po.RopeFreqBase),
		C.float(po.RopeFreqScale),
		C.int(po.NDraft),
		C.float(po.XTCProbability),
		C.float(po.XTCThreshold),
		C.float(po.DRYMultiplier),
		C.float(po.DRYBase),
		C.int(po.DRYAllowedLength),
		C.int(po.DRYPenaltyLastN),
		C.float(po.TopNSigma),
	)

	ret := C.get_embeddings(params, l.state, (*C.float)(&floats[0]))
	if ret != 0 {
		return floats, fmt.Errorf("ошибка вывода эмбеддингов")
	}

	return floats, nil
}

func (l *LLama) Predict(text string, opts ...PredictOption) (string, error) {
	po := NewPredictOptions(opts...)

	if po.TokenCallback != nil {
		setCallback(l.state, po.TokenCallback)
	}

	input := C.CString(text)
	if po.Tokens == 0 {
		po.Tokens = 99999999
	}
	out := make([]byte, po.Tokens)

	reverseCount := len(po.StopPrompts)
	reversePrompt := make([]*C.char, reverseCount)
	var pass **C.char
	for i, s := range po.StopPrompts {
		cs := C.CString(s)
		reversePrompt[i] = cs
		pass = &reversePrompt[0]
	}

	params := C.llama_allocate_params(
		input,
		C.int(po.Seed),
		C.int(po.Threads),
		C.int(po.Tokens),
		C.int(po.TopK),
		C.float(po.TopP),
		C.float(po.MinP),
		C.float(po.Temperature),
		C.float(po.Penalty),
		C.int(po.Repeat),
		C.bool(po.IgnoreEOS),
		C.bool(po.F16KV),
		C.int(po.Batch),
		C.int(po.NKeep),
		pass,
		C.int(reverseCount),
		C.float(po.TailFreeSamplingZ),
		C.float(po.TypicalP),
		C.float(po.FrequencyPenalty),
		C.float(po.PresencePenalty),
		C.int(po.Mirostat),
		C.float(po.MirostatETA),
		C.float(po.MirostatTAU),
		C.bool(po.PenalizeNL),
		C.CString(po.LogitBias),
		C.CString(po.PathPromptCache),
		C.bool(po.PromptCacheAll),
		C.bool(po.MLock),
		C.bool(po.MMap),
		C.CString(po.MainGPU),
		C.CString(po.TensorSplit),
		C.bool(po.PromptCacheRO),
		C.CString(po.Grammar),
		C.float(po.RopeFreqBase),
		C.float(po.RopeFreqScale),
		C.int(po.NDraft),
		C.float(po.XTCProbability),
		C.float(po.XTCThreshold),
		C.float(po.DRYMultiplier),
		C.float(po.DRYBase),
		C.int(po.DRYAllowedLength),
		C.int(po.DRYPenaltyLastN),
		C.float(po.TopNSigma),
	)
	ret := C.llama_predict(params, l.state, (*C.char)(unsafe.Pointer(&out[0])), C.bool(po.DebugMode))
	if ret != 0 {
		return "", fmt.Errorf("ошибка вывода")
	}

	res := C.GoString((*C.char)(unsafe.Pointer(&out[0])))

	res = strings.TrimPrefix(res, " ")
	res = strings.TrimPrefix(res, text)
	res = strings.TrimPrefix(res, "\n")

	for _, s := range po.StopPrompts {
		if s == "" {
			continue
		}
		for strings.HasSuffix(res, s) {
			res = strings.TrimSuffix(res, s)
		}
	}

	C.llama_free_params(params)

	if po.TokenCallback != nil {
		setCallback(l.state, nil)
	}

	return res, nil
}

func MTMDDefaultMarker() string {
	return C.GoString(C.llama_mtmd_default_marker())
}

func (l *LLama) HasMTMD() bool {
	if l == nil || l.state == nil {
		return false
	}

	return bool(C.llama_has_mtmd(l.state))
}

func (l *LLama) PredictMTMD(text string, images [][]byte, opts ...PredictOption) (string, error) {
	if len(images) == 0 {
		return l.Predict(text, opts...)
	}
	po := NewPredictOptions(opts...)

	if po.TokenCallback != nil {
		setCallback(l.state, po.TokenCallback)
	}

	input := C.CString(text)
	if po.Tokens == 0 {
		po.Tokens = 99999999
	}
	out := make([]byte, po.Tokens)

	reverseCount := len(po.StopPrompts)
	reversePrompt := make([]*C.char, reverseCount)
	var pass **C.char
	for i, s := range po.StopPrompts {
		cs := C.CString(s)
		reversePrompt[i] = cs
		pass = &reversePrompt[0]
	}

	ptrSlice := make([]*C.uchar, len(images))
	lenSlice := make([]C.size_t, len(images))
	for i, img := range images {
		if len(img) == 0 {
			for j := 0; j < i; j++ {
				C.free(unsafe.Pointer(ptrSlice[j]))
			}
			return "", fmt.Errorf("llama: пустые байты изображения %d", i)
		}

		cpy := C.CBytes(img)
		ptrSlice[i] = (*C.uchar)(cpy)
		lenSlice[i] = C.size_t(len(img))
	}

	defer func() {
		for _, p := range ptrSlice {
			if p != nil {
				C.free(unsafe.Pointer(p))
			}
		}
	}()

	params := C.llama_allocate_params(
		input,
		C.int(po.Seed),
		C.int(po.Threads),
		C.int(po.Tokens),
		C.int(po.TopK),
		C.float(po.TopP),
		C.float(po.MinP),
		C.float(po.Temperature),
		C.float(po.Penalty),
		C.int(po.Repeat),
		C.bool(po.IgnoreEOS),
		C.bool(po.F16KV),
		C.int(po.Batch),
		C.int(po.NKeep),
		pass,
		C.int(reverseCount),
		C.float(po.TailFreeSamplingZ),
		C.float(po.TypicalP),
		C.float(po.FrequencyPenalty),
		C.float(po.PresencePenalty),
		C.int(po.Mirostat),
		C.float(po.MirostatETA),
		C.float(po.MirostatTAU),
		C.bool(po.PenalizeNL),
		C.CString(po.LogitBias),
		C.CString(po.PathPromptCache),
		C.bool(po.PromptCacheAll),
		C.bool(po.MLock),
		C.bool(po.MMap),
		C.CString(po.MainGPU),
		C.CString(po.TensorSplit),
		C.bool(po.PromptCacheRO),
		C.CString(po.Grammar),
		C.float(po.RopeFreqBase),
		C.float(po.RopeFreqScale),
		C.int(po.NDraft),
		C.float(po.XTCProbability),
		C.float(po.XTCThreshold),
		C.float(po.DRYMultiplier),
		C.float(po.DRYBase),
		C.int(po.DRYAllowedLength),
		C.int(po.DRYPenaltyLastN),
		C.float(po.TopNSigma),
	)
	ret := C.llama_predict_mtmd(params, l.state,
		(**C.uchar)(unsafe.Pointer(&ptrSlice[0])),
		(*C.size_t)(unsafe.Pointer(&lenSlice[0])),
		C.int(len(images)),
		(*C.char)(unsafe.Pointer(&out[0])),
		C.bool(po.DebugMode),
	)
	if ret != 0 {
		C.llama_free_params(params)
		if po.TokenCallback != nil {
			setCallback(l.state, nil)
		}
		return "", fmt.Errorf("ошибка vision-вывода (mtmd)")
	}

	res := C.GoString((*C.char)(unsafe.Pointer(&out[0])))

	res = strings.TrimPrefix(res, " ")
	res = strings.TrimPrefix(res, "\n")

	for _, s := range po.StopPrompts {
		if s == "" {
			continue
		}
		for strings.HasSuffix(res, s) {
			res = strings.TrimSuffix(res, s)
		}
	}

	C.llama_free_params(params)

	if po.TokenCallback != nil {
		setCallback(l.state, nil)
	}

	return res, nil
}

func (l *LLama) TokenizeString(text string, opts ...PredictOption) (int32, []int32, error) {
	po := NewPredictOptions(opts...)

	input := C.CString(text)
	if po.Tokens == 0 {
		po.Tokens = 4096
	}

	out := make([]C.int, po.Tokens)

	var fakeDblPtr **C.char

	params := C.llama_allocate_params(
		input,
		C.int(po.Seed),
		C.int(po.Threads),
		C.int(po.Tokens),
		C.int(po.TopK),
		C.float(po.TopP),
		C.float(po.MinP),
		C.float(po.Temperature),
		C.float(po.Penalty),
		C.int(po.Repeat),
		C.bool(po.IgnoreEOS),
		C.bool(po.F16KV),
		C.int(po.Batch),
		C.int(po.NKeep),
		fakeDblPtr,
		C.int(0),
		C.float(po.TailFreeSamplingZ),
		C.float(po.TypicalP),
		C.float(po.FrequencyPenalty),
		C.float(po.PresencePenalty),
		C.int(po.Mirostat),
		C.float(po.MirostatETA),
		C.float(po.MirostatTAU),
		C.bool(po.PenalizeNL),
		C.CString(po.LogitBias),
		C.CString(po.PathPromptCache),
		C.bool(po.PromptCacheAll),
		C.bool(po.MLock),
		C.bool(po.MMap),
		C.CString(po.MainGPU),
		C.CString(po.TensorSplit),
		C.bool(po.PromptCacheRO),
		C.CString(po.Grammar),
		C.float(po.RopeFreqBase),
		C.float(po.RopeFreqScale),
		C.int(po.NDraft),
		C.float(po.XTCProbability),
		C.float(po.XTCThreshold),
		C.float(po.DRYMultiplier),
		C.float(po.DRYBase),
		C.int(po.DRYAllowedLength),
		C.int(po.DRYPenaltyLastN),
		C.float(po.TopNSigma),
	)

	tokRet := C.llama_tokenize_string(params, l.state, (*C.int)(unsafe.Pointer(&out[0])))

	if tokRet < 0 {
		return int32(tokRet), []int32{}, fmt.Errorf("токенизация (llama_tokenize_string) вернула отрицательное число %d", tokRet)
	}

	gTokRet := int32(tokRet)

	gLenOut := min(len(out), int(gTokRet))

	goSlice := make([]int32, gLenOut)
	for i := 0; i < gLenOut; i++ {
		goSlice[i] = int32(out[i])
	}

	return gTokRet, goSlice, nil
}

func (l *LLama) SetTokenCallback(callback func(token string) bool) {
	setCallback(l.state, callback)
}

var (
	m         sync.RWMutex
	callbacks = map[uintptr]func(string) bool{}
)

//export tokenCallback
func tokenCallback(statePtr unsafe.Pointer, token *C.char) C.uchar {
	m.RLock()
	defer m.RUnlock()

	if callback, ok := callbacks[uintptr(statePtr)]; ok {
		if callback(C.GoString(token)) {
			return 1
		}
		return 0
	}

	return 1
}

func setCallback(statePtr unsafe.Pointer, callback func(string) bool) {
	m.Lock()
	defer m.Unlock()

	if callback == nil {
		delete(callbacks, uintptr(statePtr))
	} else {
		callbacks[uintptr(statePtr)] = callback
	}
}
