package hydra

import (
	"database/sql"
	"strings"
	"testing"

	_ "github.com/mattn/go-sqlite3"
)

func TestFetchSQLiteDispatch(t *testing.T) {
	db := newTestSQLiteDB(t)
	defer db.Close()

	h := &Hydratable{XDBTypeOverride: "sqlite"}
	result, err := h.Fetch(db, "person", map[string]interface{}{"id": 1})
	if err != nil {
		t.Fatalf("fetch should succeed: %v", err)
	}

	if got, ok := result["name"].(string); !ok || got != "Alice" {
		t.Fatalf("unexpected name value: %#v", result["name"])
	}
}

func TestFetchDefaultsToMySQLForSQLDBWithoutOverride(t *testing.T) {
	db := newTestSQLiteDB(t)
	defer db.Close()

	h := &Hydratable{}
	result, err := h.Fetch(db, "person", map[string]interface{}{"id": 2})
	if err != nil {
		t.Fatalf("default sql.DB dispatch should succeed: %v", err)
	}

	if got, ok := result["name"].(string); !ok || got != "Bob" {
		t.Fatalf("unexpected name value: %#v", result["name"])
	}
}

func TestFetchUnsupportedDatabaseType(t *testing.T) {
	h := &Hydratable{}
	_, err := h.Fetch(struct{}{}, "person", map[string]interface{}{"id": 1})
	if err == nil {
		t.Fatal("expected unsupported database type error")
	}

	if !strings.Contains(err.Error(), "unsupported database type") {
		t.Fatalf("unexpected error: %v", err)
	}
}

func TestFetchSQLiteWithRealSQLDBType(t *testing.T) {
	db, err := sql.Open("sqlite3", ":memory:")
	if err != nil {
		t.Fatalf("open sqlite db: %v", err)
	}
	defer db.Close()

	if _, err := db.Exec(`CREATE TABLE person (id INTEGER PRIMARY KEY, name TEXT)`); err != nil {
		t.Fatalf("create table: %v", err)
	}
	if _, err := db.Exec(`INSERT INTO person (id, name) VALUES (1, 'Alice')`); err != nil {
		t.Fatalf("insert row: %v", err)
	}

	h := &Hydratable{XDBTypeOverride: "sqlite"}
	result, err := h.Fetch(db, "person", map[string]interface{}{"id": 1})
	if err != nil {
		t.Fatalf("fetch should succeed: %v", err)
	}

	if got, ok := result["id"].(int64); !ok || got != 1 {
		t.Fatalf("unexpected id value: %#v", result["id"])
	}
}
