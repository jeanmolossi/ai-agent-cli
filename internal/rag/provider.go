package rag

// SearchResult assoc a chunk with their similarity
type SearchResult struct {
	ID      string
	Content string
	Score   float64
}

// VectorStore is the contract to index and seek vectors.
type VectorStore interface {
	Add(id, content string) error
	Search(query string, topK int) ([]SearchResult, error)
	Persist() error
	Load() error
	Drop() error
	Close() error
}
