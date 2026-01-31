package runner

import (
	"sync"
	"time"
)

type InferenceMetrics struct {
	mu               sync.RWMutex
	lastTokens       int64
	lastLatencyMs    float64
	lastTokensPerSec float64
	lastAt           time.Time
}

func NewInferenceMetrics() *InferenceMetrics {
	return &InferenceMetrics{}
}

func (m *InferenceMetrics) Record(tokens int64, duration time.Duration) {
	if tokens < 0 {
		tokens = 0
	}

	sec := duration.Seconds()
	var tps float64
	if sec > 0 {
		tps = float64(tokens) / sec
	}

	m.mu.Lock()
	defer m.mu.Unlock()
	m.lastTokens = tokens
	m.lastLatencyMs = duration.Seconds() * 1000
	m.lastTokensPerSec = tps
	m.lastAt = time.Now()
}

func (m *InferenceMetrics) Get() (tokens int64, latencyMs float64, tokensPerSec float64) {
	m.mu.RLock()
	defer m.mu.RUnlock()

	return m.lastTokens, m.lastLatencyMs, m.lastTokensPerSec
}
