package agent

import (
	"context"

	contractsagent "github.com/jeanmolossi/ai-agent-cli/internal/contracts/agent"
	"github.com/tmc/langchaingo/llms"
)

type lcProvider struct {
	model llms.Model
}

var _ (contractsagent.LLMProvider) = (*lcProvider)(nil)

// Generate implements LLMProvider.
func (p *lcProvider) Generate(ctx context.Context, prompt string) (string, error) {
	return llms.GenerateFromSinglePrompt(ctx, p.model, prompt)
}

func (p *lcProvider) Model() llms.Model {
	return p.model
}
