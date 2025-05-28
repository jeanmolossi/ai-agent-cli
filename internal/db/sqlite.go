package db

import (
	"database/sql"
	"path/filepath"

	_ "github.com/mattn/go-sqlite3"
)

func OpenLocal(dataDir string) (*sql.DB, error) {
	dbPath := filepath.Join(dataDir, "agent_memory.db")
	return sql.Open("sqlite3", dbPath)
}
