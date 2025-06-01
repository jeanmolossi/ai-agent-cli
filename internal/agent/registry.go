package agent

import (
	"fmt"

	contractsagent "github.com/jeanmolossi/ai-agent-cli/internal/contracts/agent"
	"github.com/spf13/viper"
	"github.com/tmc/langchaingo/llms/anthropic"
	"github.com/tmc/langchaingo/llms/ollama"
	"github.com/tmc/langchaingo/llms/openai"
)

var providerRegistry = map[string]func() (contractsagent.LLMProvider, error){
	"openai":    newOpenAIProvider,
	"anthropic": newAnthropicProvider,
	"ollama":    newOllamaProvider,
}

func newProvider(cfg *Config) (contractsagent.LLMProvider, error) {
	ctor, ok := providerRegistry[cfg.Provider]
	if !ok {
		return nil, fmt.Errorf("LLM provider %q not supported", cfg.Provider)
	}

	return ctor()
}

func NewLLMProviderWithOptions(options ...LLMOption) (contractsagent.LLMProvider, error) {
	cfg := &Config{}

	for _, option := range options {
		option(cfg)
	}

	return newProvider(cfg)
}

func NewLLMProvider() (contractsagent.LLMProvider, error) {
	name := viper.GetString("llm.provider")

	if name == "" {
		name = "ollama" // ollama as default
	}

	return NewLLMProviderWithOptions(
		WithProvider(name),
	)
}

func newOpenAIProvider() (contractsagent.LLMProvider, error) {
	key := viper.GetString("llm.openai.api_key")
	if key == "" {
		return nil, fmt.Errorf("llm.openai.api_key not found")
	}

	client, err := openai.New(
		openai.WithToken(key),
	)
	if err != nil {
		return nil, err
	}

	return &lcProvider{model: client}, nil
}

func newAnthropicProvider() (contractsagent.LLMProvider, error) {
	key := viper.GetString("llm.anthropic.api_key")
	if key == "" {
		return nil, fmt.Errorf("llm.anthropic.api_key not found")
	}

	client, err := anthropic.New(
		anthropic.WithToken(key),
	)
	if err != nil {
		return nil, err
	}

	return &lcProvider{model: client}, nil
}

func newOllamaProvider() (contractsagent.LLMProvider, error) {
	host := viper.GetString("llm.ollama.host")
	if host == "" {
		host = "http://localhost"
	}

	port := viper.GetInt("llm.ollama.port")
	if port == 0 {
		port = 11434 // default Ollama port
	}

	model := viper.GetString("llm.ollama.model")
	if model == "" {
		model = "gemma"
	}

	baseURL := fmt.Sprintf("%s:%d", host, port)

	client, err := ollama.New(
		ollama.WithServerURL(baseURL),
		ollama.WithModel(model),
	)
	if err != nil {
		return nil, err
	}

	return &lcProvider{model: client}, nil
}
