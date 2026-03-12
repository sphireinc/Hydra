package hydra

import (
	"errors"
	"reflect"
	"testing"
)

func TestBuildSelectQueryTableDriven(t *testing.T) {
	columns := []string{"active", "id", "name"}
	where := map[string]interface{}{"name": "Alice", "id": 7}

	tests := []struct {
		name        string
		formatter   placeholderFormatter
		expectQuery string
	}{
		{"question", questionPlaceholder, "SELECT active, id, name FROM person WHERE id = ? AND name = ?"},
		{"dollar", dollarPlaceholder, "SELECT active, id, name FROM person WHERE id = $1 AND name = $2"},
		{"colon", colonPlaceholder, "SELECT active, id, name FROM person WHERE id = :1 AND name = :2"},
		{"atp", atPPlaceholder, "SELECT active, id, name FROM person WHERE id = @p1 AND name = @p2"},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			query, params, err := buildSelectQuery("person", columns, where, tt.formatter)
			if err != nil {
				t.Fatalf("buildSelectQuery: %v", err)
			}
			if query != tt.expectQuery {
				t.Fatalf("unexpected query\nwant: %s\n got: %s", tt.expectQuery, query)
			}
			expectedParams := []interface{}{7, "Alice"}
			if !reflect.DeepEqual(params, expectedParams) {
				t.Fatalf("unexpected params\nwant: %#v\n got: %#v", expectedParams, params)
			}
		})
	}
}

func TestBuildSelectQueryDeterministicOrdering(t *testing.T) {
	query, params, err := buildSelectQuery("person", []string{"name", "id"}, map[string]interface{}{"zeta": 3, "alpha": 1, "middle": 2}, dollarPlaceholder)
	if err != nil {
		t.Fatalf("buildSelectQuery: %v", err)
	}
	if query != "SELECT name, id FROM person WHERE alpha = $1 AND middle = $2 AND zeta = $3" {
		t.Fatalf("unexpected query: %s", query)
	}
	if !reflect.DeepEqual(params, []interface{}{1, 2, 3}) {
		t.Fatalf("unexpected params: %#v", params)
	}
}

func TestBuildSelectQueryRequiresWhereClauses(t *testing.T) {
	_, _, err := buildSelectQuery("person", []string{"id"}, map[string]interface{}{}, questionPlaceholder)
	if !errors.Is(err, ErrEmptyWhereClause) {
		t.Fatalf("expected ErrEmptyWhereClause, got: %v", err)
	}
}

func TestBuildSelectQueryRejectsInvalidIdentifiers(t *testing.T) {
	_, _, err := buildSelectQuery("person;DROP", []string{"id"}, map[string]interface{}{"id": 1}, questionPlaceholder)
	if err == nil {
		t.Fatal("expected invalid table identifier error")
	}

	_, _, err = buildSelectQuery("person", []string{"id"}, map[string]interface{}{"id;DROP": 1}, questionPlaceholder)
	if err == nil {
		t.Fatal("expected invalid column identifier error")
	}
}

func TestScanSQLRowsFirst(t *testing.T) {
	db := newTestSQLiteDB(t)
	defer db.Close()

	rows, err := db.Query(`SELECT id, name, active FROM person ORDER BY id`)
	if err != nil {
		t.Fatalf("query rows: %v", err)
	}
	defer rows.Close()

	result, err := scanSQLRowsFirst(rows)
	if err != nil {
		t.Fatalf("scan: %v", err)
	}

	if result["id"].(int64) != 1 || result["name"].(string) != "Alice" || result["active"].(bool) != true {
		t.Fatalf("unexpected result: %#v", result)
	}
}

func TestScanSQLRowsFirstNoRows(t *testing.T) {
	db := newTestSQLiteDB(t)
	defer db.Close()

	rows, err := db.Query(`SELECT id FROM person WHERE id = -1`)
	if err != nil {
		t.Fatalf("query rows: %v", err)
	}
	defer rows.Close()

	_, err = scanSQLRowsFirst(rows)
	if !errors.Is(err, ErrNotFound) {
		t.Fatalf("expected ErrNotFound, got: %v", err)
	}
}
