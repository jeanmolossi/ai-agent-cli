package agent

import (
	"context"

	"github.com/tmc/langchaingo/llms"
)

type lcProvider struct {
	model llms.Model
}

var _ (LLMProvider) = (*lcProvider)(nil)

// Generate implements LLMProvider.
func (p *lcProvider) Generate(ctx context.Context, prompt string) (string, error) {
	return llms.GenerateFromSinglePrompt(ctx, p.model, prompt)
}

func (p *lcProvider) Model() llms.Model {
	return p.model
}
