package hydra

import "testing"

func TestFetchPostgresUsesDollarPlaceholders(t *testing.T) {
	query, params, err := buildSelectQuery("person", []string{"id", "name"}, map[string]interface{}{"id": 1, "name": "Alice"}, dollarPlaceholder)
	if err != nil {
		t.Fatalf("buildSelectQuery: %v", err)
	}
	if query != "SELECT id, name FROM person WHERE id = $1 AND name = $2" {
		t.Fatalf("unexpected query: %s", query)
	}
	if len(params) != 2 || params[0] != 1 || params[1] != "Alice" {
		t.Fatalf("unexpected params: %#v", params)
	}
}

func TestFetchPostgresBuildsDeterministicQuery(t *testing.T) {
	query, params, err := buildSelectQuery("person", []string{"id", "name"}, map[string]interface{}{"zeta": 3, "alpha": 1, "middle": 2}, dollarPlaceholder)
	if err != nil {
		t.Fatalf("buildSelectQuery: %v", err)
	}
	if query != "SELECT id, name FROM person WHERE alpha = $1 AND middle = $2 AND zeta = $3" {
		t.Fatalf("unexpected query: %s", query)
	}
	if len(params) != 3 || params[0] != 1 || params[1] != 2 || params[2] != 3 {
		t.Fatalf("unexpected params: %#v", params)
	}
}
