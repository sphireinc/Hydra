package hydra

import (
	"errors"
	"testing"
)

func TestFetchSQLite(t *testing.T) {
	db := newTestSQLiteDB(t)
	defer db.Close()

	h := &Hydratable{}
	result, err := h.fetchSQLite(db, "person", []string{"id", "name"}, map[string]interface{}{"id": 1})
	if err != nil {
		t.Fatalf("fetchSQLite: %v", err)
	}
	if result["id"].(int64) != 1 || result["name"].(string) != "Alice" {
		t.Fatalf("unexpected result: %#v", result)
	}
}

func TestFetchSQLiteNoRows(t *testing.T) {
	db := newTestSQLiteDB(t)
	defer db.Close()

	h := &Hydratable{}
	_, err := h.fetchSQLite(db, "person", []string{"id", "name"}, map[string]interface{}{"id": 999})
	if !errors.Is(err, ErrNotFound) {
		t.Fatalf("expected ErrNotFound, got: %v", err)
	}
}
