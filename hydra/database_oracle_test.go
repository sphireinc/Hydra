package hydra

import "testing"

func TestFetchOracle(t *testing.T) {
	db := newTestSQLiteDB(t)
	defer db.Close()

	h := &Hydratable{}
	result, err := h.fetchOracle(db, "person", map[string]interface{}{"id": 1})
	if err != nil {
		t.Fatalf("fetchOracle should succeed: %v", err)
	}

	if got, ok := result["id"].(int64); !ok || got != 1 {
		t.Fatalf("unexpected id value: %#v", result["id"])
	}

	if got, ok := result["name"].(string); !ok || got != "Alice" {
		t.Fatalf("unexpected name value: %#v", result["name"])
	}
}

func TestFetchOracleUsesColonPlaceholders(t *testing.T) {
	query, params := buildSelectQuery("person", map[string]interface{}{"name": "Alice", "id": 1}, colonPlaceholder)

	expectedQuery := "SELECT * FROM person WHERE id = :1 AND name = :2"
	if query != expectedQuery {
		t.Fatalf("unexpected query\nwant: %s\n got: %s", expectedQuery, query)
	}

	if len(params) != 2 || params[0] != 1 || params[1] != "Alice" {
		t.Fatalf("unexpected params: %#v", params)
	}
}
