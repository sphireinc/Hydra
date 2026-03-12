package hydra

import "testing"

func TestFetchMySQL(t *testing.T) {
	db := newTestSQLiteDB(t)
	defer db.Close()

	h := &Hydratable{}
	result, err := h.fetchMySQL(db, "person", []string{"id", "name"}, map[string]interface{}{"id": 1})
	if err != nil {
		t.Fatalf("fetchMySQL: %v", err)
	}
	if result["name"].(string) != "Alice" {
		t.Fatalf("unexpected result: %#v", result)
	}
}

func TestFetchMySQLMultipleWhereClauses(t *testing.T) {
	db := newTestSQLiteDB(t)
	defer db.Close()

	h := &Hydratable{}
	result, err := h.fetchMySQL(db, "person", []string{"id", "name"}, map[string]interface{}{"id": 2, "name": "Bob"})
	if err != nil {
		t.Fatalf("fetchMySQL: %v", err)
	}
	if result["name"].(string) != "Bob" {
		t.Fatalf("unexpected result: %#v", result)
	}
}
