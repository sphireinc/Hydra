package hydra

import "testing"

func TestFetchCockroachDBUsesDollarPlaceholders(t *testing.T) {
	query, params, err := buildSelectQuery("person", []string{"id", "name"}, map[string]interface{}{"id": 2, "name": "Bob"}, dollarPlaceholder)
	if err != nil {
		t.Fatalf("buildSelectQuery: %v", err)
	}
	if query != "SELECT id, name FROM person WHERE id = $1 AND name = $2" {
		t.Fatalf("unexpected query: %s", query)
	}
	if len(params) != 2 || params[0] != 2 || params[1] != "Bob" {
		t.Fatalf("unexpected params: %#v", params)
	}
}

func TestFetchCockroachDBBuildsDeterministicQuery(t *testing.T) {
	query, params, err := buildSelectQuery("person", []string{"id", "tenant_id"}, map[string]interface{}{"tenant_id": 42, "id": 7}, dollarPlaceholder)
	if err != nil {
		t.Fatalf("buildSelectQuery: %v", err)
	}
	if query != "SELECT id, tenant_id FROM person WHERE id = $1 AND tenant_id = $2" {
		t.Fatalf("unexpected query: %s", query)
	}
	if len(params) != 2 || params[0] != 7 || params[1] != 42 {
		t.Fatalf("unexpected params: %#v", params)
	}
}
