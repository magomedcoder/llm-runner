package internal

import "github.com/magomedcoder/gen/internal/domain"

func generationParamsFromWire(g *GenerationWire) *domain.GenerationParams {
	if g == nil {
		return nil
	}

	if g.Temperature == nil && g.MaxTokens == nil {
		return nil
	}

	out := &domain.GenerationParams{}
	if g.Temperature != nil {
		t := *g.Temperature
		out.Temperature = &t
	}

	if g.MaxTokens != nil && *g.MaxTokens > 0 {
		out.MaxTokens = g.MaxTokens
	}

	return out
}
