package hydra

import "testing"

func TestFetchOracle(t *testing.T) {
	db := newTestSQLiteDB(t)
	defer db.Close()

	h := &Hydratable{}
	result, err := h.fetchOracle(db, "person", []string{"id", "name"}, map[string]interface{}{"id": 1})
	if err != nil {
		t.Fatalf("fetchOracle: %v", err)
	}
	if result["name"].(string) != "Alice" {
		t.Fatalf("unexpected result: %#v", result)
	}
}

func TestFetchOracleUsesColonPlaceholders(t *testing.T) {
	query, params, err := buildSelectQuery("person", []string{"id", "name"}, map[string]interface{}{"id": 1, "name": "Alice"}, colonPlaceholder)
	if err != nil {
		t.Fatalf("buildSelectQuery: %v", err)
	}
	if query != "SELECT id, name FROM person WHERE id = :1 AND name = :2" {
		t.Fatalf("unexpected query: %s", query)
	}
	if len(params) != 2 || params[0] != 1 || params[1] != "Alice" {
		t.Fatalf("unexpected params: %#v", params)
	}
}
