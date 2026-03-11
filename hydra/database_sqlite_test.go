package hydra

import "testing"

func TestFetchSQLite(t *testing.T) {
	db := newTestSQLiteDB(t)
	defer db.Close()

	h := &Hydratable{}
	result, err := h.fetchSQLite(db, "person", map[string]interface{}{"id": 1})
	if err != nil {
		t.Fatalf("fetchSQLite should succeed: %v", err)
	}

	if got, ok := result["id"].(int64); !ok || got != 1 {
		t.Fatalf("unexpected id value: %#v", result["id"])
	}

	if got, ok := result["name"].(string); !ok || got != "Alice" {
		t.Fatalf("unexpected name value: %#v", result["name"])
	}
}

func TestFetchSQLiteNoRows(t *testing.T) {
	db := newTestSQLiteDB(t)
	defer db.Close()

	h := &Hydratable{}
	result, err := h.fetchSQLite(db, "person", map[string]interface{}{"id": 999})
	if err != nil {
		t.Fatalf("fetchSQLite should not error on no rows: %v", err)
	}

	if len(result) != 0 {
		t.Fatalf("expected empty result on no rows, got: %#v", result)
	}
}
