package chunkers

import (
	"fmt"
	"os"
	"path/filepath"
	"strings"

	sitter "github.com/tree-sitter/go-tree-sitter"
	tsgo "github.com/tree-sitter/tree-sitter-go/bindings/go"
	tsjs "github.com/tree-sitter/tree-sitter-javascript/bindings/go"
	tspy "github.com/tree-sitter/tree-sitter-python/bindings/go"
)

// Chunk represents a semantic code chunk extracted from source.
type Chunk struct {
	FilePath  string // Absolute or relative path of the source file
	StartLine int    // 1-based start line of this chunk
	EndLine   int    // 1-based end line of this chunk
	Content   string // Raw source lines joined by "\n"
}

// Chunker splits source files into semantic chunks, per-language.
type Chunker struct {
	MaxLines int                        // maximum number of lines per chunk
	langs    map[string]*LanguageConfig // keyed by the file extension (without dot)
}

// LanguageConfig holds Tree-sitter language and node types to extract
type LanguageConfig struct {
	Language      *sitter.Language
	DeclNodeTypes map[string]struct{} // node types indicating top-level declarations
}

// New returns a Chunker configured with maxLines heuristic.
func New(maxLines int) *Chunker {
	return &Chunker{
		MaxLines: maxLines,
		langs: map[string]*LanguageConfig{
			"go": {
				Language: sitter.NewLanguage(tsgo.Language()),
				DeclNodeTypes: map[string]struct{}{
					"function_declaration": {},
					"method_declaration":   {},
				},
			},
			"js": {
				Language: sitter.NewLanguage(tsjs.Language()),
				DeclNodeTypes: map[string]struct{}{
					"function_declaration": {},
					"method_declaration":   {},
					"class_declaration":    {},
				},
			},
			"py": {
				Language: sitter.NewLanguage(tspy.Language()),
				DeclNodeTypes: map[string]struct{}{
					"function_definition": {},
					"class_definition":    {},
				},
			},
		},
	}
}

// SplitFile routes the file at path to the appropriate language handler.
// Returns ErrUnsupportedExtension if no handler is registered.
func (c *Chunker) SplitFile(path string) ([]Chunk, error) {
	ext := strings.TrimPrefix(filepath.Ext(path), ".")
	cfg, ok := c.langs[ext]
	if !ok {
		return nil, fmt.Errorf("unsupported file extension: %s", ext)
	}

	src, err := os.ReadFile(path)
	if err != nil {
		return nil, err
	}

	p := sitter.NewParser()
	defer p.Close()

	err = p.SetLanguage(cfg.Language)
	if err != nil {
		return nil, err
	}

	tree := p.Parse(src, nil)
	defer tree.Close()

	var chunks []Chunk

	var recurse func(node *sitter.Node)
	recurse = func(node *sitter.Node) {
		if _, isDecl := cfg.DeclNodeTypes[node.Kind()]; isDecl {
			start := node.StartPosition().Row + 1
			end := node.EndPosition().Row + 1
			content := extractLines(string(src), int(start), int(end))
			chunks = append(chunks, Chunk{
				FilePath:  path,
				StartLine: int(start),
				EndLine:   int(end),
				Content:   content,
			})
		}

		for child := node.NamedChild(0); child != nil; child = child.NextNamedSibling() {
			recurse(child)
		}
	}
	recurse(tree.RootNode())

	return c.splitByMaxLines(chunks), nil
}

// splitByMaxLines split chunks who exceeds c.MaxLines into child chunks until c.MaxLines lines.
func (c *Chunker) splitByMaxLines(chunks []Chunk) []Chunk {
	var out []Chunk
	for _, ch := range chunks {
		lines := strings.Split(ch.Content, "\n")
		if len(lines) <= c.MaxLines {
			out = append(out, ch)
			continue
		}

		for i := 0; i < len(lines); i += c.MaxLines {
			end := i + c.MaxLines

			if end > len(lines) {
				end = len(lines)
			}

			subContent := strings.Join(lines[i:end], "\n")
			out = append(out, Chunk{
				FilePath:  ch.FilePath,
				StartLine: ch.StartLine + i,
				EndLine:   ch.StartLine + end - 1,
				Content:   subContent,
			})
		}
	}

	return out
}

// extractLines return lines [start..end] (1-based) from src content.
func extractLines(src string, start, end int) string {
	all := strings.Split(src, "\n")
	if start < 1 {
		start = 1
	}

	if end > len(all) {
		end = len(all)
	}

	return strings.Join(all[start-1:end], "\n")
}
