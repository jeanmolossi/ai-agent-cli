package contractsagent

import (
	"context"

	"github.com/tmc/langchaingo/llms"
)

// LLMProvider is the contract who any LLM should implements
type LLMProvider interface {
	// Generate receives a prompt and returns the model response
	Generate(ctx context.Context, prompt string) (string, error)

	// Model
	Model() llms.Model
}
