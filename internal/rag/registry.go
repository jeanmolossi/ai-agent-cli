package rag

import (
	"fmt"

	"github.com/spf13/viper"
)

var storeRegistry = map[string]func() (VectorStore, error){
	"local": newLocalStore,
}

func NewVectorStore() (VectorStore, error) {
	name := viper.GetString("rag.provider")
	if name == "" {
		name = "local"
	}

	ctor, ok := storeRegistry[name]
	if !ok {
		return nil, fmt.Errorf("vector store %q not found", name)
	}

	return ctor()
}
