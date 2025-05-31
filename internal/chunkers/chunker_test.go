package chunkers_test

import (
	"os"
	"path/filepath"
	"strings"
	"testing"

	"github.com/jeanmolossi/ai-agent-cli/internal/chunkers"
	"github.com/stretchr/testify/assert"
)

func writeFile(t *testing.T, dir, name, content string) string {
	t.Helper()
	path := filepath.Join(dir, name)
	if err := os.WriteFile(path, []byte(content), 0o644); err != nil {
		t.Fatalf("falha ao escrever arquivo %s: %v", path, err)
	}
	return path
}

func TestChunker_SplitFile(t *testing.T) {
	// Cria um diretório temporário para colocar todos os arquivos de teste.
	baseDir := t.TempDir()

	// Definições de conteúdo de exemplo para cada linguagem.
	goSmall := `package main

func Small() {
	println("small")
}

func Another() {
	println("small too")
}
`
	// Função com 8 linhas de corpo, excedendo MaxLines = 5.
	goLarge := `package main

func Large() {
	// line1
	// line2
	// line3
	// line4
	// line5
	// line6
	println("end")
}
`

	jsTwoFuncs := `function foo() {
    console.log("foo small");
}

function bar() {
    console.log("bar small");
}
`

	pyClassMethods := `class Exemple:
    def method1(self):
        print("method1")

    def method2(self):
        print("method2")
`

	// Cria arquivos temporários para cada linguagem.
	goSmallPath := writeFile(t, baseDir, "small.go", goSmall)
	goLargePath := writeFile(t, baseDir, "large.go", goLarge)
	jsPath := writeFile(t, baseDir, "twofuncs.js", jsTwoFuncs)
	pyPath := writeFile(t, baseDir, "class.py", pyClassMethods)

	t.Run("Go: small and large funcs", func(t *testing.T) {
		// maxLines low (5) to force chunking of Large func
		c := chunkers.New(5)

		// 1) Check the file with two small funcs -> should return exactly 2 chunks
		chunksSmall, err := c.SplitFile(goSmallPath)
		assert.NoError(t, err)
		assert.Len(t, chunksSmall, 2, "expected 2 chunks (two small funcs)")

		// Each chunk chunkSmall[i].Content should have "func Small" or "func Another"
		foundSmall := 0
		for _, ch := range chunksSmall {
			if strings.Contains(ch.Content, "func Small") || strings.Contains(ch.Content, "func Another") {
				foundSmall++
			}
		}
		assert.Equal(t, 2, foundSmall, "each small func should appears lonely")

		// 2) verify the file with Large func -> should return multiple chunks (>1)
		chunksLarge, err := c.SplitFile(goLargePath)
		assert.NoError(t, err)
		// The Large func have 8 lines body, MaxLines=5, then expected 2 chunks:
		//    - 1st chunk: lines 3 to 7 (5 lines of content)
		//    - 2nd chunk: lines 8 to 10 (ramaining body + closing)
		assert.Len(t, chunksLarge, 2, "expected 2 chunks for the Large func")

		// Looks for the first chunk if contains the comment "// line" and the second "println"
		assert.True(t, strings.Contains(chunksLarge[0].Content, "// line"), "chunk1 should have the first line commented")
		assert.True(t, strings.Contains(chunksLarge[1].Content, "println"), "chunk2 should have a println call at the end")
	})

	t.Run("JavaScript: two small funcs", func(t *testing.T) {
		// MaxLines bigger enough (10) to do not chunk internal body func
		c := chunkers.New(10)

		chunksJS, err := c.SplitFile(jsPath)
		assert.NoError(t, err)

		// Should have 2 chunks: one for "function foo" another one for "function bar"
		assert.Len(t, chunksJS, 2, "expected 2 chunks in JS (foo & bar)")

		foundJS := 0
		for _, ch := range chunksJS {
			if strings.Contains(ch.Content, "function foo") || strings.Contains(ch.Content, "function bar") {
				foundJS++
			}
		}
		assert.Equal(t, 2, foundJS, "each JS function should be a single chunk")
	})

	t.Run("Python: class and methods", func(t *testing.T) {
		// MaxLines enough (10) to do not chunk each block
		c := chunkers.New(10)

		chunksPY, err := c.SplitFile(pyPath)
		assert.NoError(t, err)

		// The class contains 2 methods, but our parser extracts:
		//   - 1 chunk to "class Example"
		//   - 1 chunk to "def method1"
		//   - 1 chunk to "def method2"
		assert.Len(t, chunksPY, 3, "expected 3 Python chunks (class + 2 methods)")

		foundPY := 0
		for _, ch := range chunksPY {
			if strings.Contains(ch.Content, "class Exemple") ||
				strings.Contains(ch.Content, "def method1") ||
				strings.Contains(ch.Content, "def method2") {
				foundPY++
			}
		}
		assert.Equal(t, 3, foundPY, "each definition (class & methods) should be a chunk")
	})
}
