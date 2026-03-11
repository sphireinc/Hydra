package hydra

import "testing"

func TestFetchMariaDB(t *testing.T) {
	db := newTestSQLiteDB(t)
	defer db.Close()

	h := &Hydratable{}
	result, err := h.fetchMariaDB(db, "person", map[string]interface{}{"id": 2})
	if err != nil {
		t.Fatalf("fetchMariaDB should succeed: %v", err)
	}

	if got, ok := result["id"].(int64); !ok || got != 2 {
		t.Fatalf("unexpected id value: %#v", result["id"])
	}

	if got, ok := result["name"].(string); !ok || got != "Bob" {
		t.Fatalf("unexpected name value: %#v", result["name"])
	}
}

func TestFetchMariaDBNoRows(t *testing.T) {
	db := newTestSQLiteDB(t)
	defer db.Close()

	h := &Hydratable{}
	result, err := h.fetchMariaDB(db, "person", map[string]interface{}{"id": 999})
	if err != nil {
		t.Fatalf("fetchMariaDB should not error on no rows: %v", err)
	}

	if len(result) != 0 {
		t.Fatalf("expected empty result on no rows, got: %#v", result)
	}
}
