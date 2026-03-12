package hydra

import (
	"context"
	"database/sql"
	"errors"
	"reflect"
	"testing"

	"github.com/jackc/pgx/v5"
)

func TestResolveFetchRoute(t *testing.T) {
	tests := []struct {
		name      string
		h         Hydratable
		db        any
		expect    fetchRoute
		expectErr bool
	}{
		{"sql default mysql", Hydratable{}, &sql.DB{}, fetchRouteMySQL, false},
		{"sql sqlite override", Hydratable{XDBTypeOverride: "sqlite"}, &sql.DB{}, fetchRouteSQLite, false},
		{"sql oracle override", Hydratable{XDBTypeOverride: "oracle"}, &sql.DB{}, fetchRouteOracle, false},
		{"sql unsupported override", Hydratable{XDBTypeOverride: "postgres"}, &sql.DB{}, "", true},
		{"pgx default postgres", Hydratable{}, (*pgx.Conn)(nil), fetchRoutePostgres, false},
		{"pgx cockroach override", Hydratable{XDBTypeOverride: "cockroachdb"}, (*pgx.Conn)(nil), fetchRouteCockroachDB, false},
		{"unsupported db type", Hydratable{}, struct{}{}, "", true},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := tt.h.resolveFetchRoute(tt.db)
			if tt.expectErr {
				if err == nil {
					t.Fatal("expected error")
				}
				return
			}
			if err != nil {
				t.Fatalf("resolveFetchRoute: %v", err)
			}
			if got != tt.expect {
				t.Fatalf("unexpected route: %s", got)
			}
		})
	}
}

func TestFetchSQLiteDispatchAndContext(t *testing.T) {
	db := newTestSQLiteDB(t)
	defer db.Close()

	h := &Hydratable{XDBTypeOverride: "sqlite"}
	result, err := h.FetchContext(context.Background(), db, "person", []string{"id", "name"}, map[string]interface{}{"id": 1})
	if err != nil {
		t.Fatalf("FetchContext: %v", err)
	}
	if result["name"].(string) != "Alice" {
		t.Fatalf("unexpected result: %#v", result)
	}
}

func TestBuildSelectQuery(t *testing.T) {
	columns := []string{"active", "id", "name"}
	where := map[string]interface{}{"name": "Alice", "id": 7}

	tests := []struct {
		name      string
		formatter placeholderFormatter
		expect    string
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
			if query != tt.expect {
				t.Fatalf("unexpected query: %s", query)
			}
			if !reflect.DeepEqual(params, []interface{}{7, "Alice"}) {
				t.Fatalf("unexpected params: %#v", params)
			}
		})
	}
}

func TestBuildSelectQueryErrors(t *testing.T) {
	if _, _, err := buildSelectQuery("person", []string{"id"}, map[string]interface{}{}, questionPlaceholder); !errors.Is(err, ErrEmptyWhereClause) {
		t.Fatalf("expected ErrEmptyWhereClause, got: %v", err)
	}

	if _, _, err := buildSelectQuery("person;drop", []string{"id"}, map[string]interface{}{"id": 1}, questionPlaceholder); err == nil {
		t.Fatal("expected invalid identifier error")
	}
}

func TestMetadataCacheAndTagParsing(t *testing.T) {
	type sample struct {
		Hydratable
		ID    int    `hydra:"id,pk"`
		Email string `hydra:"email,lookup"`
		Name  string `hydra:"name"`
	}

	meta1, err := getTypeMetadata(&sample{})
	if err != nil {
		t.Fatalf("getTypeMetadata: %v", err)
	}
	meta2, err := getTypeMetadata(&sample{})
	if err != nil {
		t.Fatalf("getTypeMetadata: %v", err)
	}

	if meta1 != meta2 {
		t.Fatal("expected cached metadata pointer reuse")
	}
	if len(meta1.PrimaryKeys) != 1 || meta1.PrimaryKeys[0].Column != "id" {
		t.Fatalf("unexpected pk metadata: %#v", meta1.PrimaryKeys)
	}
	if len(meta1.LookupFields) != 1 || meta1.LookupFields[0].Column != "email" {
		t.Fatalf("unexpected lookup metadata: %#v", meta1.LookupFields)
	}
}
