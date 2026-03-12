package hydra

import (
	"errors"
	"testing"
)

func TestFetchMariaDB(t *testing.T) {
	db := newTestSQLiteDB(t)
	defer db.Close()

	h := &Hydratable{}
	result, err := h.fetchMariaDB(db, "person", []string{"id", "name"}, map[string]interface{}{"id": 2})
	if err != nil {
		t.Fatalf("fetchMariaDB: %v", err)
	}
	if result["name"].(string) != "Bob" {
		t.Fatalf("unexpected result: %#v", result)
	}
}

func TestFetchMariaDBNoRows(t *testing.T) {
	db := newTestSQLiteDB(t)
	defer db.Close()

	h := &Hydratable{}
	_, err := h.fetchMariaDB(db, "person", []string{"id"}, map[string]interface{}{"id": 999})
	if !errors.Is(err, ErrNotFound) {
		t.Fatalf("expected ErrNotFound, got: %v", err)
	}
}
