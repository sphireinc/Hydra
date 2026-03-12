# Hydra package guide

## Import

```go
import "github.com/sphireinc/Hydra/hydra"
```

## Overview

Hydra populates tagged struct fields from a single database row.

Typical flow:

1. define a struct with `hydra` tags 
2. embed `hydra.Hydratable`
3. initialize with `Init`
4. hydrate using one of the supported APIs

## Basic example

```go
package main

import (
	"database/sql"
	"log"

	"github.com/sphireinc/Hydra/hydra"
	_ "github.com/mattn/go-sqlite3"
)

type Person struct {
	ID    int    `hydra:"id,pk"`
	Email string `hydra:"email,lookup"`
	Name  string `hydra:"name"`

	hydra.Hydratable
}

func (Person) HydraTableName() string {
	return "person"
}

func main() {
	db, err := sql.Open("sqlite3", ":memory:")
	if err != nil {
		log.Fatal(err)
	}
	defer db.Close()

	_, err = db.Exec(`
		CREATE TABLE person (
			id INTEGER PRIMARY KEY,
			email TEXT NOT NULL,
			name TEXT NOT NULL
		);

		INSERT INTO person (id, email, name)
		VALUES (1, 'alice@example.com', 'Alice');
	`)
	if err != nil {
		log.Fatal(err)
	}

	person := &Person{}
	person.Init(person)
	person.XDBTypeOverride = "sqlite"

	if err := person.HydrateByPrimaryKey(db, 1); err != nil {
		log.Fatal(err)
	}
}
```

## Core contract

### Initialization

Hydra expects an addressable struct pointer.

```go
person := &Person{}
person.Init(person)
```

Calling `Init` on a non-pointer struct value is not the intended usage.

### Table name resolution

Hydra resolves the table name in this order:

1. `XTableNameOverride`
2. `HydraTableName() string`
3. lowercase struct type name

### Supported handle types

Hydra supports these database handle types:

- `*sql.DB`
    - mysql
    - mariadb
    - sqlite
    - mssql
    - oracle
- `*pgx.Conn`
    - postgres
    - cockroachdb

Use `XDBTypeOverride` to route the fetcher.

Examples:

```go
obj.XDBTypeOverride = "sqlite"
obj.XDBTypeOverride = "mysql"
obj.XDBTypeOverride = "postgres"
obj.XDBTypeOverride = "cockroachdb"
```

### No matching row

If no matching row exists, Hydra returns:

```go
hydra.ErrNotFound
```

### Empty filters

If hydration is attempted with an empty `where` map, Hydra returns:

```go
hydra.ErrEmptyWhereClause
```

### Identifier safety

Hydra validates table names and column names before building SQL.

Only simple identifiers are accepted.

## Tags

### Standard mapping

```go
type Person struct {
	ID   int    `hydra:"id"`
	Name string `hydra:"name"`
}
```

### Primary key tags

```go
type Person struct {
	ID int `hydra:"id,pk"`
}
```

This enables:

```go
err := person.HydrateByPrimaryKey(db, 1)
```

### Lookup tags

```go
type Person struct {
	Email string `hydra:"email,lookup"`
}
```

This enables:

```go
person := &Person{Email: "alice@example.com"}
person.Init(person)
err := person.HydrateByLookup(db)
```

## Supported field conversion behavior

Hydra supports built-in conversion for:

- `string`
- `bool`
- signed integers
- unsigned integers
- floats
- pointers to supported types
- interfaces
- directly assignable/convertible values for matching struct, slice, array, and map types

`NULL` is supported for:

- pointers
- slices
- maps
- interfaces

Attempting to place `NULL` into a non-nullable scalar field returns an error.

## Custom converters

Hydra supports two converter models.

### Field-level converter

If the field implements:

```go
HydraConvert(src any) error
```

Hydra uses it before built-in primitive conversion.

Example:

```go
type RFC3339Time struct {
	time.Time
}

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
		return fmt.Errorf("unsupported value %T", src)
	}
}
```

### Parent-level converters

A struct can provide converters keyed by field name or column name:

```go
HydraConverters() map[string]hydra.HydraFieldConverter
```

This is useful for:

- enums
- JSON columns
- nullable wrappers
- custom map/slice decoding

## Context-aware APIs

Hydra exposes context-aware calls:

```go
HydrateContext(ctx, db, whereClauses)
HydrateByPrimaryKeyContext(ctx, db, value)
HydrateByLookupContext(ctx, db)
FetchContext(ctx, db, tableName, columns, whereClauses)
```

Use them when the caller needs cancellation or deadlines.

## Error handling

Common errors include:

- `hydra.ErrNotInitialized`
- `hydra.ErrEmptyWhereClause`
- `hydra.ErrNotFound`

Example:

```go
if err := person.HydrateByPrimaryKey(db, 42); err != nil {
	switch {
	case errors.Is(err, hydra.ErrNotFound):
		// row missing
	case errors.Is(err, hydra.ErrNotInitialized):
		// Init was not called
	default:
		// other conversion, validation, or database error
	}
}
```

## Recommended usage patterns

### Explicit where clauses

```go
person := &Person{}
person.Init(person)
person.XDBTypeOverride = "sqlite"

err := person.Hydrate(db, map[string]interface{}{
	"id": 1,
})
```

### Primary key lookup

```go
person := &Person{}
person.Init(person)
person.XDBTypeOverride = "sqlite"

err := person.HydrateByPrimaryKey(db, 1)
```

### Lookup-field hydration

```go
person := &Person{
	Email: "alice@example.com",
}
person.Init(person)
person.XDBTypeOverride = "sqlite"

err := person.HydrateByLookup(db)
```

## Testing flow

Unit tests:

```bash
make test
make test-hydra
```

Full Docker-backed functional suite:

```bash
make test-func
```
