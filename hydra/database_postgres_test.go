package hydra

import "testing"

func TestFetchPostgresUsesDollarPlaceholders(t *testing.T) {
	query, params := buildSelectQuery("person", map[string]interface{}{"name": "Alice", "id": 1}, dollarPlaceholder)

	expectedQuery := "SELECT * FROM person WHERE id = $1 AND name = $2"
	if query != expectedQuery {
		t.Fatalf("unexpected query\nwant: %s\n got: %s", expectedQuery, query)
	}

	if len(params) != 2 || params[0] != 1 || params[1] != "Alice" {
		t.Fatalf("unexpected params: %#v", params)
	}
}

func TestFetchPostgresBuildsDeterministicQuery(t *testing.T) {
	whereClauses := map[string]interface{}{
		"zeta":   3,
		"alpha":  1,
		"middle": 2,
	}

	query, params := buildSelectQuery("person", whereClauses, dollarPlaceholder)

	expectedQuery := "SELECT * FROM person WHERE alpha = $1 AND middle = $2 AND zeta = $3"
	if query != expectedQuery {
		t.Fatalf("unexpected query\nwant: %s\n got: %s", expectedQuery, query)
	}

	if len(params) != 3 || params[0] != 1 || params[1] != 2 || params[2] != 3 {
		t.Fatalf("unexpected params: %#v", params)
	}
}
