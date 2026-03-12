package hydra

import "testing"

func TestFetchMSSQL(t *testing.T) {
	db := newTestSQLiteDB(t)
	defer db.Close()

	h := &Hydratable{}
	result, err := h.fetchMSSQL(db, "person", []string{"id", "name"}, map[string]interface{}{"id": 2})
	if err != nil {
		t.Fatalf("fetchMSSQL: %v", err)
	}
	if result["name"].(string) != "Bob" {
		t.Fatalf("unexpected result: %#v", result)
	}
}

func TestFetchMSSQLUsesAtPPlaceholders(t *testing.T) {
	query, params, err := buildSelectQuery("person", []string{"id", "name"}, map[string]interface{}{"id": 2, "name": "Bob"}, atPPlaceholder)
	if err != nil {
		t.Fatalf("buildSelectQuery: %v", err)
	}
	if query != "SELECT id, name FROM person WHERE id = @p1 AND name = @p2" {
		t.Fatalf("unexpected query: %s", query)
	}
	if len(params) != 2 || params[0] != 2 || params[1] != "Bob" {
		t.Fatalf("unexpected params: %#v", params)
	}
}
