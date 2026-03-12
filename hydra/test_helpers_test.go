package hydra

import (
	"database/sql"
	"testing"

	_ "github.com/mattn/go-sqlite3"
)

func newTestSQLiteDB(t *testing.T) *sql.DB {
	t.Helper()

	db, err := sql.Open("sqlite3", ":memory:")
	if err != nil {
		t.Fatalf("open sqlite db: %v", err)
	}

	statements := []string{
		`CREATE TABLE person (id INTEGER PRIMARY KEY, name TEXT, email TEXT, active BOOLEAN, score REAL, nickname TEXT, meta TEXT, created_at TEXT);`,
		`INSERT INTO person (id, name, email, active, score, nickname, meta, created_at) VALUES (1, 'Alice', 'alice@example.com', 1, 9.5, NULL, 'alpha', '2026-01-02T03:04:05Z');`,
		`INSERT INTO person (id, name, email, active, score, nickname, meta, created_at) VALUES (2, 'Bob', 'bob@example.com', 0, 7.25, 'Bobby', 'beta', '2026-01-03T03:04:05Z');`,
	}

	for _, stmt := range statements {
		if _, err := db.Exec(stmt); err != nil {
			db.Close()
			t.Fatalf("exec %q: %v", stmt, err)
		}
	}

	return db
}
