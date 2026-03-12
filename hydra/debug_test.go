package hydra

import "testing"

func TestDebugDoesNotAffectBehavior(t *testing.T) {
	db := newTestSQLiteDB(t)
	defer db.Close()

	build := func() *hydratePerson {
		p := &hydratePerson{}
		p.Init(p)
		p.XDBTypeOverride = "sqlite"
		p.XTableNameOverride = "person"
		return p
	}

	debug = false
	first := build()
	if err := first.HydrateByPrimaryKey(db, 1); err != nil {
		t.Fatalf("debug off: %v", err)
	}

	debug = true
	second := build()
	if err := second.HydrateByPrimaryKey(db, 1); err != nil {
		t.Fatalf("debug on: %v", err)
	}
	debug = false

	if first.Name != second.Name || first.Email != second.Email || first.ID != second.ID {
		t.Fatalf("debug changed behavior: %#v %#v", first, second)
	}
}
