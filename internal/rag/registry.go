package rag

import (
	"fmt"

	contractsrag "github.com/jeanmolossi/ai-agent-cli/internal/contracts/rag"
	"github.com/spf13/viper"
)

var storeRegistry = map[string]func() (contractsrag.VectorStore, error){
	"local": newLocalStore,
}

func newVectorStore(cfg *Config) (contractsrag.VectorStore, error) {
	ctor, ok := storeRegistry[cfg.Provider]
	if !ok {
		return nil, fmt.Errorf("vector store %q not found", cfg.Provider)
	}

	return ctor()
}

func NewVectorStoreWithOptions(options ...RagOption) (contractsrag.VectorStore, error) {
	cfg := &Config{}

	for _, opt := range options {
		opt(cfg)
	}

	return newVectorStore(cfg)
}

func NewVectorStore() (contractsrag.VectorStore, error) {
	name := viper.GetString("rag.provider")
	if name == "" {
		name = "local"
	}

	return NewVectorStoreWithOptions(WithProvider(name))
}
