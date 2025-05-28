package rag

import (
	"fmt"
	"io/fs"
	"log/slog"
	"os"
	"path/filepath"
	"strings"

	"github.com/spf13/viper"
)

var defaultIgnores = []string{
	".git",
	".vscode",
	".docker",
	".idea",
	"node_modules",
	"vendor",
}

func ScanAndIndex(roots ...string) error {
	ignores := viper.GetStringSlice("rag.ignore")
	if len(ignores) == 0 {
		slog.Debug(
			"Nenhum diretório para ignorar, utilizando o padrão...",
			slog.String("paths", strings.Join(ignores, "|")),
		)

		ignores = defaultIgnores
	}

	store, err := NewVectorStore()
	if err != nil {
		return err
	}

	store.Drop()
	store.Load()

	for _, root := range roots {
		if err := filepath.Walk(root, func(path string, info fs.FileInfo, err error) error {
			if err != nil {
				return err
			}

			if info.IsDir() {
				for _, ig := range ignores {
					if info.Name() == ig {
						return filepath.SkipDir
					}
				}

				return nil
			}

			if !isTextFile(path) {
				return nil
			}

			data, _ := os.ReadFile(path)
			chunks := chunkText(string(data), viper.GetInt("rag.local.chunk_size"))
			for i, c := range chunks {
				id := fmt.Sprintf("%s#%d", path, i)
				store.Add(id, c)
			}

			return nil
		}); err != nil {
			return fmt.Errorf("can not index %s: %w", root, err)
		}
	}

	return store.Persist()
}

var exts = []string{".go", ".md", ".yaml", ".yml", ".txt", ".json"}

func isTextFile(path string) bool {
	for _, e := range exts {
		if strings.HasSuffix(path, e) {
			return true
		}
	}

	return false
}

func chunkText(text string, size int) []string {
	var chunks []string

	for len(text) > size {
		chunks = append(chunks, text[:size])
		text = text[size:]
	}

	if len(text) > 0 {
		chunks = append(chunks, text)
	}

	return chunks
}
