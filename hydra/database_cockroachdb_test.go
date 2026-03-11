package hydra

import "testing"

func TestFetchCockroachDBUsesDollarPlaceholders(t *testing.T) {
	query, params := buildSelectQuery("person", map[string]interface{}{"name": "Bob", "id": 2}, dollarPlaceholder)

	expectedQuery := "SELECT * FROM person WHERE id = $1 AND name = $2"
	if query != expectedQuery {
		t.Fatalf("unexpected query\nwant: %s\n got: %s", expectedQuery, query)
	}

	if len(params) != 2 || params[0] != 2 || params[1] != "Bob" {
		t.Fatalf("unexpected params: %#v", params)
	}
}

func TestFetchCockroachDBBuildsDeterministicQuery(t *testing.T) {
	whereClauses := map[string]interface{}{
		"tenant_id": 42,
		"id":        7,
	}

	query, params := buildSelectQuery("person", whereClauses, dollarPlaceholder)

	expectedQuery := "SELECT * FROM person WHERE id = $1 AND tenant_id = $2"
	if query != expectedQuery {
		t.Fatalf("unexpected query\nwant: %s\n got: %s", expectedQuery, query)
	}

	if len(params) != 2 || params[0] != 7 || params[1] != 42 {
		t.Fatalf("unexpected params: %#v", params)
	}
}
