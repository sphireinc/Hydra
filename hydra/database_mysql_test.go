package hydra

import "testing"

func TestFetchMySQL(t *testing.T) {
	db := newTestSQLiteDB(t)
	defer db.Close()

	h := &Hydratable{}
	result, err := h.fetchMySQL(db, "person", map[string]interface{}{"id": 1})
	if err != nil {
		t.Fatalf("fetchMySQL should succeed: %v", err)
	}

	if got, ok := result["id"].(int64); !ok || got != 1 {
		t.Fatalf("unexpected id value: %#v", result["id"])
	}

	if got, ok := result["name"].(string); !ok || got != "Alice" {
		t.Fatalf("unexpected name value: %#v", result["name"])
	}
}

func TestFetchMySQLMultipleWhereClauses(t *testing.T) {
	db := newTestSQLiteDB(t)
	defer db.Close()

	h := &Hydratable{}
	result, err := h.fetchMySQL(db, "person", map[string]interface{}{"name": "Bob", "id": 2})
	if err != nil {
		t.Fatalf("fetchMySQL should succeed with multiple where clauses: %v", err)
	}

	if got, ok := result["name"].(string); !ok || got != "Bob" {
		t.Fatalf("unexpected name value: %#v", result["name"])
	}
}
