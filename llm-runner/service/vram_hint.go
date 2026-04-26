//go:build llama

package service

import "fmt"

func roughKVCacheMiBUpperBound(nLayer, nCtx, nEmbd int) float64 {
	if nLayer <= 0 || nCtx <= 0 || nEmbd <= 0 {
		return 0
	}

	const bytesFP16 = 2
	bytes := int64(nLayer) * int64(nCtx) * int64(nEmbd) * 2 * bytesFP16

	return float64(bytes) / (1024 * 1024)
}

func formatVRAMHint(modelSizeBytes int64, nCtx, nLayer, nEmbd int) string {
	kv := roughKVCacheMiBUpperBound(nLayer, nCtx, nEmbd)
	wGiB := float64(modelSizeBytes) / (1024 * 1024 * 1024)

	return fmt.Sprintf("оценка VRAM: тензоры модели ~%.2f GiB; KV при n_ctx=%d - грубый верх ~%.0f MiB (GQA обычно меньше); OOM -> gpu_layers↓ или max_context_tokens↓ или другая квантовка", wGiB, nCtx, kv)
}
