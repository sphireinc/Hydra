package hydra

import (
	"context"
	"errors"
	"reflect"
	"strings"
	"testing"
	"time"
)

type RFC3339Time struct{ time.Time }

func (t *RFC3339Time) HydraConvert(src any) error {
	switch v := src.(type) {
	case string:
		parsed, err := time.Parse(time.RFC3339, v)
		if err != nil {
			return err
		}
		t.Time = parsed
		return nil
	case []byte:
		parsed, err := time.Parse(time.RFC3339, string(v))
		if err != nil {
			return err
		}
		t.Time = parsed
		return nil
	default:
		return errors.New("unsupported time input")
	}
}

type hydratePerson struct {
	Hydratable
	ID        int         `hydra:"id,pk"`
	Name      string      `hydra:"name"`
	Email     string      `hydra:"email,lookup"`
	Active    bool        `hydra:"active"`
	Score     float64     `hydra:"score"`
	Nickname  *string     `hydra:"nickname"`
	CreatedAt RFC3339Time `hydra:"created_at"`
}

type hydratePersonWithProvider struct {
	Hydratable
	ID      int               `hydra:"id,pk"`
	Meta    map[string]string `hydra:"meta"`
	Created string            `hydra:"created_at"`
}

func (p *hydratePersonWithProvider) HydraConverters() map[string]HydraFieldConverter {
	return map[string]HydraFieldConverter{
		"meta": func(src any) (any, error) {
			var s string
			switch v := src.(type) {
			case string:
				s = v
			case []byte:
				s = string(v)
			default:
				return nil, errors.New("unsupported meta input")
			}
			return map[string]string{"raw": s}, nil
		},
	}
}

type unsupportedHydrateField struct {
	Hydratable
	Meta chan int `hydra:"meta"`
}

func TestHydrateRequiresInit(t *testing.T) {
	var p hydratePerson
	err := p.Hydrate(nil, map[string]interface{}{"id": 1})
	if !errors.Is(err, ErrNotInitialized) {
		t.Fatalf("expected ErrNotInitialized, got: %v", err)
	}
}

func TestHydrateSuccessAndPrimaryKeyLookup(t *testing.T) {
	db := newTestSQLiteDB(t)
	defer db.Close()

	p := &hydratePerson{}
	p.Init(p)
	p.XDBTypeOverride = "sqlite"
	p.XTableNameOverride = "person"

	if err := p.HydrateByPrimaryKey(db, 1); err != nil {
		t.Fatalf("hydrate by primary key: %v", err)
	}

	if p.ID != 1 || p.Name != "Alice" || p.Email != "alice@example.com" || !p.Active || p.Score != 9.5 {
		t.Fatalf("unexpected hydrated struct: %#v", p)
	}
	if p.Nickname != nil {
		t.Fatalf("expected nil nickname, got: %#v", p.Nickname)
	}
	if p.CreatedAt.IsZero() {
		t.Fatal("expected CreatedAt to be populated via custom converter")
	}
}

func TestHydrateByLookup(t *testing.T) {
	db := newTestSQLiteDB(t)
	defer db.Close()

	p := &hydratePerson{Email: "bob@example.com"}
	p.Init(p)
	p.XDBTypeOverride = "sqlite"
	p.XTableNameOverride = "person"

	if err := p.HydrateByLookupContext(context.Background(), db); err != nil {
		t.Fatalf("hydrate by lookup: %v", err)
	}

	if p.ID != 2 || p.Name != "Bob" {
		t.Fatalf("unexpected hydrated struct: %#v", p)
	}
}

func TestHydrateWithProviderConverter(t *testing.T) {
	db := newTestSQLiteDB(t)
	defer db.Close()

	p := &hydratePersonWithProvider{}
	p.Init(p)
	p.XDBTypeOverride = "sqlite"
	p.XTableNameOverride = "person"

	if err := p.HydrateByPrimaryKey(db, 1); err != nil {
		t.Fatalf("hydrate by primary key: %v", err)
	}

	if p.Meta["raw"] != "alpha" {
		t.Fatalf("unexpected provider-converted meta: %#v", p.Meta)
	}
}

func TestHydrateEmptyWhereClauses(t *testing.T) {
	db := newTestSQLiteDB(t)
	defer db.Close()

	p := &hydratePerson{}
	p.Init(p)
	p.XDBTypeOverride = "sqlite"
	p.XTableNameOverride = "person"

	err := p.Hydrate(db, map[string]interface{}{})
	if !errors.Is(err, ErrEmptyWhereClause) {
		t.Fatalf("expected ErrEmptyWhereClause, got: %v", err)
	}
}

func TestHydrateNotFound(t *testing.T) {
	db := newTestSQLiteDB(t)
	defer db.Close()

	p := &hydratePerson{}
	p.Init(p)
	p.XDBTypeOverride = "sqlite"
	p.XTableNameOverride = "person"

	err := p.Hydrate(db, map[string]interface{}{"id": 999})
	if !errors.Is(err, ErrNotFound) {
		t.Fatalf("expected ErrNotFound, got: %v", err)
	}
}

func TestAssignFieldValueConversions(t *testing.T) {
	target := struct {
		Name   string
		Active bool
		ID     int
		Count  uint
		Score  float64
	}{}

	v := reflect.ValueOf(&target).Elem()
	meta := fieldMetadata{Name: "Name", Column: "name"}

	if err := assignFieldValue(nil, meta, v.FieldByName("Name"), []byte("Alice")); err != nil {
		t.Fatalf("string: %v", err)
	}
	if err := assignFieldValue(nil, meta, v.FieldByName("Active"), "true"); err != nil {
		t.Fatalf("bool: %v", err)
	}
	if err := assignFieldValue(nil, meta, v.FieldByName("ID"), int64(7)); err != nil {
		t.Fatalf("int: %v", err)
	}
	if err := assignFieldValue(nil, meta, v.FieldByName("Count"), "12"); err != nil {
		t.Fatalf("uint: %v", err)
	}
	if err := assignFieldValue(nil, meta, v.FieldByName("Score"), []byte("9.5")); err != nil {
		t.Fatalf("float: %v", err)
	}

	if target.Name != "Alice" || !target.Active || target.ID != 7 || target.Count != 12 || target.Score != 9.5 {
		t.Fatalf("unexpected target: %#v", target)
	}
}

func TestAssignFieldValueNilPointer(t *testing.T) {
	target := struct{ Nickname *string }{}
	v := reflect.ValueOf(&target).Elem().FieldByName("Nickname")

	if err := assignFieldValue(nil, fieldMetadata{Name: "Nickname", Column: "nickname"}, v, nil); err != nil {
		t.Fatalf("nil pointer assign: %v", err)
	}
	if target.Nickname != nil {
		t.Fatalf("expected nil nickname, got: %#v", target.Nickname)
	}
}

func TestHydrateUnsupportedTypeReturnsError(t *testing.T) {
	db := newTestSQLiteDB(t)
	defer db.Close()

	p := &unsupportedHydrateField{}
	p.Init(p)
	p.XDBTypeOverride = "sqlite"
	p.XTableNameOverride = "person"

	err := p.Hydrate(db, map[string]interface{}{"id": 1})
	if err == nil || !strings.Contains(err.Error(), "field Meta") {
		t.Fatalf("expected field context error, got: %v", err)
	}
}
