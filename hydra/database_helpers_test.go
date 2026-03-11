package hydra

import (
	"database/sql"
	"reflect"
	"testing"

	_ "github.com/mattn/go-sqlite3"
)

func newTestSQLiteDB(t *testing.T) *sql.DB {
	t.Helper()

	db, err := sql.Open("sqlite3", ":memory:")
	if err != nil {
		t.Fatalf("open sqlite db: %v", err)
	}

	statements := []string{
		`CREATE TABLE person (id INTEGER PRIMARY KEY, name TEXT, active BOOLEAN);`,
		`INSERT INTO person (id, name, active) VALUES (1, 'Alice', 1);`,
		`INSERT INTO person (id, name, active) VALUES (2, 'Bob', 0);`,
	}

	for _, stmt := range statements {
		if _, err := db.Exec(stmt); err != nil {
			db.Close()
			t.Fatalf("exec %q: %v", stmt, err)
		}
	}

	return db
}

func TestBuildSelectQuery(t *testing.T) {
	whereClauses := map[string]interface{}{
		"name": "Alice",
		"id":   7,
	}

	tests := []struct {
		name        string
		formatter   placeholderFormatter
		expectQuery string
	}{
		{
			name:        "question placeholders",
			formatter:   questionPlaceholder,
			expectQuery: "SELECT * FROM person WHERE id = ? AND name = ?",
		},
		{
			name:        "dollar placeholders",
			formatter:   dollarPlaceholder,
			expectQuery: "SELECT * FROM person WHERE id = $1 AND name = $2",
		},
		{
			name:        "colon placeholders",
			formatter:   colonPlaceholder,
			expectQuery: "SELECT * FROM person WHERE id = :1 AND name = :2",
		},
		{
			name:        "at-p placeholders",
			formatter:   atPPlaceholder,
			expectQuery: "SELECT * FROM person WHERE id = @p1 AND name = @p2",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			query, params := buildSelectQuery("person", whereClauses, tt.formatter)

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
	whereClauses := map[string]interface{}{
		"zeta":   3,
		"alpha":  1,
		"middle": 2,
	}

	query1, params1 := buildSelectQuery("person", whereClauses, dollarPlaceholder)
	query2, params2 := buildSelectQuery("person", whereClauses, dollarPlaceholder)

	if query1 != query2 {
		t.Fatalf("query should be deterministic\nfirst:  %s\nsecond: %s", query1, query2)
	}

	if !reflect.DeepEqual(params1, params2) {
		t.Fatalf("params should be deterministic\nfirst:  %#v\nsecond: %#v", params1, params2)
	}

	expectedQuery := "SELECT * FROM person WHERE alpha = $1 AND middle = $2 AND zeta = $3"
	if query1 != expectedQuery {
		t.Fatalf("unexpected ordered query\nwant: %s\n got: %s", expectedQuery, query1)
	}

	expectedParams := []interface{}{1, 2, 3}
	if !reflect.DeepEqual(params1, expectedParams) {
		t.Fatalf("unexpected ordered params\nwant: %#v\n got: %#v", expectedParams, params1)
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
		t.Fatalf("scan first row: %v", err)
	}

	if got, ok := result["id"].(int64); !ok || got != 1 {
		t.Fatalf("unexpected id value: %#v", result["id"])
	}

	if got, ok := result["name"].(string); !ok || got != "Alice" {
		t.Fatalf("unexpected name value: %#v", result["name"])
	}

	if got, ok := result["active"].(bool); !ok || !got {
		t.Fatalf("unexpected active value: %#v", result["active"])
	}
}

func TestScanSQLRowsFirstNoRows(t *testing.T) {
	db := newTestSQLiteDB(t)
	defer db.Close()

	rows, err := db.Query(`SELECT id, name, active FROM person WHERE id = -1`)
	if err != nil {
		t.Fatalf("query rows: %v", err)
	}
	defer rows.Close()

	result, err := scanSQLRowsFirst(rows)
	if err != nil {
		t.Fatalf("scan first row: %v", err)
	}

	if len(result) != 0 {
		t.Fatalf("expected empty result for no rows, got: %#v", result)
	}
}

func TestQueryFirstRowSQL(t *testing.T) {
	db := newTestSQLiteDB(t)
	defer db.Close()

	result, err := queryFirstRowSQL(db, `SELECT id, name, active FROM person WHERE id = ?`, []interface{}{2})
	if err != nil {
		t.Fatalf("query first row: %v", err)
	}

	if got, ok := result["id"].(int64); !ok || got != 2 {
		t.Fatalf("unexpected id value: %#v", result["id"])
	}

	if got, ok := result["name"].(string); !ok || got != "Bob" {
		t.Fatalf("unexpected name value: %#v", result["name"])
	}

	if got, ok := result["active"].(bool); !ok || got {
		t.Fatalf("unexpected active value: %#v", result["active"])
	}
}
