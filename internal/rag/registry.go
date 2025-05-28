package rag

import (
	"fmt"

	"github.com/spf13/viper"
)

var storeRegistry = map[string]func() (VectorStore, error){
	"local": newLocalStore,
}

func newVectorStore(cfg *Config) (VectorStore, error) {
	ctor, ok := storeRegistry[cfg.Provider]
	if !ok {
		return nil, fmt.Errorf("vector store %q not found", cfg.Provider)
	}

	return ctor()
}

func NewVectorStoreWithOptions(options ...RagOption) (VectorStore, error) {
	cfg := &Config{}

	for _, opt := range options {
		opt(cfg)
	}

	return newVectorStore(cfg)
}

func NewVectorStore() (VectorStore, error) {
	name := viper.GetString("rag.provider")
	if name == "" {
		name = "local"
	}

	return NewVectorStoreWithOptions(WithProvider(name))
}
