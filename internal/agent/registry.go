package agent

import (
	"fmt"

	"github.com/spf13/viper"
	"github.com/tmc/langchaingo/llms/anthropic"
	"github.com/tmc/langchaingo/llms/ollama"
	"github.com/tmc/langchaingo/llms/openai"
)

var providerRegistry = map[string]func() (LLMProvider, error){
	"openai":    newOpenAIProvider,
	"anthropic": newAnthropicProvider,
	"ollama":    newOllamaProvider,
}

func NewLLMProvider() (LLMProvider, error) {
	name := viper.GetString("llm.provider")

	if name == "" {
		name = "ollama" // ollama as default
	}

	ctor, ok := providerRegistry[name]
	if !ok {
		return nil, fmt.Errorf("LLM provider %q not supported", name)
	}

	return ctor()
}

func newOpenAIProvider() (LLMProvider, error) {
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

func newAnthropicProvider() (LLMProvider, error) {
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

func newOllamaProvider() (LLMProvider, error) {
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
