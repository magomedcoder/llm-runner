//go:build llama
// +build llama

package llama

var DefaultModelOptions = ModelOptions{
	ContextSize:   512,
	Seed:          0,
	F16Memory:     false,
	MLock:         false,
	Embeddings:    false,
	MMap:          true,
	LowVRAM:       false,
	NBatch:        512,
	FreqRopeBase:  10000,
	FreqRopeScale: 1.0,
}

func NewModelOptions(opts ...ModelOption) ModelOptions {
	p := DefaultModelOptions
	for _, opt := range opts {
		opt(&p)
	}

	return p
}
