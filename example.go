package main

import (
	"database/sql"
	"errors"
	"fmt"

	_ "github.com/mattn/go-sqlite3"
	"github.com/sphireinc/Hydra/hydra"
)

type ExamplePerson struct {
	ID    int    `hydra:"id,pk"`
	Email string `hydra:"email,lookup"`
	Name  string `hydra:"name"`

	hydra.Hydratable
}

func (ExamplePerson) HydraTableName() string {
	return "person"
}

func runExample() error {
	db, err := sql.Open("sqlite3", ":memory:")
	if err != nil {
		return fmt.Errorf("open sqlite db: %w", err)
	}
	defer db.Close()

	_, err = db.Exec(`
		CREATE TABLE person (
			id INTEGER PRIMARY KEY,
			email TEXT NOT NULL,
			name TEXT NOT NULL
		);

		INSERT INTO person (id, email, name)
		VALUES
			(1, 'alice@example.com', 'Alice'),
			(2, 'bob@example.com', 'Bob');
	`)
	if err != nil {
		return fmt.Errorf("seed sqlite db: %w", err)
	}

	person := &ExamplePerson{}
	person.Init(person)
	person.XDBTypeOverride = "sqlite"

	if err := person.HydrateByPrimaryKey(db, 1); err != nil {
		return fmt.Errorf("hydrate by primary key: %w", err)
	}

	fmt.Printf("loaded by pk: id=%d email=%s name=%s\n", person.ID, person.Email, person.Name)

	lookup := &ExamplePerson{
		Email: "bob@example.com",
	}
	lookup.Init(lookup)
	lookup.XDBTypeOverride = "sqlite"

	if err := lookup.HydrateByLookup(db); err != nil {
		if errors.Is(err, hydra.ErrNotFound) {
			return fmt.Errorf("hydrate by lookup: %w", err)
		}
		return fmt.Errorf("hydrate by lookup: %w", err)
	}

	fmt.Printf("loaded by lookup: id=%d email=%s name=%s\n", lookup.ID, lookup.Email, lookup.Name)

	return nil
}
