package rag

// SearchResult assoc a chunk with their similarity
type SearchResult struct {
	ID      string
	Content string
	Score   float64
}
