package hydra

import "testing"

func TestFetchMSSQL(t *testing.T) {
	db := newTestSQLiteDB(t)
	defer db.Close()

	h := &Hydratable{}
	result, err := h.fetchMSSQL(db, "person", map[string]interface{}{"id": 2})
	if err != nil {
		t.Fatalf("fetchMSSQL should succeed: %v", err)
	}

	if got, ok := result["id"].(int64); !ok || got != 2 {
		t.Fatalf("unexpected id value: %#v", result["id"])
	}

	if got, ok := result["name"].(string); !ok || got != "Bob" {
		t.Fatalf("unexpected name value: %#v", result["name"])
	}
}

func TestFetchMSSQLUsesAtPPlaceholders(t *testing.T) {
	query, params := buildSelectQuery("person", map[string]interface{}{"name": "Bob", "id": 2}, atPPlaceholder)

	expectedQuery := "SELECT * FROM person WHERE id = @p1 AND name = @p2"
	if query != expectedQuery {
		t.Fatalf("unexpected query\nwant: %s\n got: %s", expectedQuery, query)
	}

	if len(params) != 2 || params[0] != 2 || params[1] != "Bob" {
		t.Fatalf("unexpected params: %#v", params)
	}
}
