package agent

import (
	"context"
	"fmt"
	"log/slog"

	"github.com/spf13/viper"
	"github.com/tmc/langchaingo/embeddings"
)

type EmbedProvider interface {
	Embed(ctx context.Context, content string) ([][]float32, error)
}

type embedProvider struct {
	embedder embeddings.Embedder
}

func NewEmbeddingProvider() (EmbedProvider, error) {
	ragProvider := viper.GetString("rag.embed.provider")
	if ragProvider == "" {
		ragProvider = viper.GetString("llm.provider")
	}

	provider, err := NewLLMProviderWithOptions(WithProvider(ragProvider))
	if err != nil {
		return nil, err
	}

	client, ok := provider.Model().(embeddings.EmbedderClient)
	if !ok {
		return nil, fmt.Errorf("the client is not supported as embedding provider")
	}

	slog.Debug("starting embedding provider", slog.String("provider", ragProvider))

	embedder, err := embeddings.NewEmbedder(client)
	if err != nil {
		return nil, err
	}

	return &embedProvider{embedder: embedder}, nil
}

// Embed implements EmbedProvider.
func (e *embedProvider) Embed(ctx context.Context, content string) ([][]float32, error) {
	return e.embedder.EmbedDocuments(ctx, []string{content})
}
