package rag

import (
	"context"
	"database/sql"
	"encoding/json"
	"fmt"
	"log/slog"
	"os"
	"path/filepath"
	"sort"

	"github.com/jeanmolossi/ai-agent-cli/internal/agent"
	contractsagent "github.com/jeanmolossi/ai-agent-cli/internal/contracts/agent"
	contractsrag "github.com/jeanmolossi/ai-agent-cli/internal/contracts/rag"
	"github.com/jeanmolossi/ai-agent-cli/pkg/similarity"
	_ "github.com/mattn/go-sqlite3"
)

type localStore struct {
	db       *sql.DB
	embedder contractsagent.EmbedProvider
}

var (
	_ contractsrag.VectorStore  = (*localStore)(nil)
	_ contractsrag.SearchResult = (*searchResult)(nil)

	dataDir = filepath.Join(".", ".ai-agent-cli")
)

func newLocalStore() (contractsrag.VectorStore, error) {
	err := os.MkdirAll(dataDir, 0o755)
	if err != nil {
		return nil, err
	}

	db, err := sql.Open("sqlite3", filepath.Join(dataDir, "rag.db"))
	if err != nil {
		return nil, err
	}

	_, _ = db.Exec(`
        CREATE TABLE IF NOT EXISTS docs (
            id TEXT PRIMARY KEY,
            content TEXT,
            embedding BLOB
        )
    `)

	ep, err := agent.NewEmbeddingProvider()
	if err != nil {
		return nil, err
	}

	return &localStore{db: db, embedder: ep}, nil
}

// Add implements VectorStore.
func (l *localStore) Add(id string, content string) error {
	vec, err := l.embedder.Embed(context.Background(), content)
	if err != nil {
		return err
	}

	blob, _ := json.Marshal(vec)
	_, err = l.db.Exec(
		`INSERT OR REPLACE INTO docs (id, content, embedding) VALUES (?,?,?)`,
		id, content, blob,
	)

	return err
}

// Search implements VectorStore.
func (l *localStore) Search(query string, topK int) ([]contractsrag.SearchResult, error) {
	qvec, err := l.embedder.Embed(context.Background(), query)
	if err != nil {
		return nil, err
	}

	rows, err := l.db.Query(`SELECT id, content, embedding FROM docs`)
	if err != nil {
		return nil, err
	}

	//nolint:errcheck
	defer rows.Close()

	var results []contractsrag.SearchResult

	for rows.Next() {
		var id, content string
		var blob []byte

		_ = rows.Scan(&id, &content, &blob)

		var vec [][]float32
		_ = json.Unmarshal(blob, &vec)

		score, err := similarity.CosSimilarity(qvec, vec)
		if err != nil {
			return nil, err
		}

		slog.Debug("similarity retrieved", slog.String("id", id), slog.Float64("similarity", score*100))

		results = append(results, &searchResult{id, content, score})
	}

	sort.Slice(results, func(i, j int) bool {
		return results[i].Score() > results[j].Score()
	})

	if len(results) > topK {
		results = results[:topK]
	}

	slog.Debug(fmt.Sprintf("still with %d results", len(results)), slog.Int("topK", topK))
	return results, nil
}

// Load implements VectorStore.
func (l *localStore) Load() error {
	db, err := sql.Open("sqlite3", filepath.Join(dataDir, "rag.db"))
	if err != nil {
		return err
	}

	l.db = db

	return nil
}

func (l *localStore) Drop() error {
	_, err := l.db.Exec(`DELETE FROM docs`)
	return err
}

// Persist implements VectorStore.
func (l *localStore) Persist() error {
	return nil
}

func (l *localStore) Close() error {
	return l.db.Close()
}

// search result

type searchResult struct {
	id      string
	content string
	score   float64
}

// Content implements contractsrag.SearchResult.
func (s *searchResult) Content() string { return s.content }

// ID implements contractsrag.SearchResult.
func (s *searchResult) ID() string { return s.id }

// Score implements contractsrag.SearchResult.
func (s *searchResult) Score() float64 { return s.score }
