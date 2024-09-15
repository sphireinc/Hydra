# Sphire Hydra

![Build Passing](https://github.com/sphireinc/Hydra/actions/workflows/go.yml/badge.svg)

<div align="center">
    <img src="logo.jpg" width="400px" />
</div>

Sphire Hydra is a Go library designed to dynamically hydrate Go structs with data from various 
databases. The library supports multiple databases, including MySQL, PostgreSQL, SQLite, 
Microsoft SQL Server, Oracle, MariaDB, and CockroachDB. Using reflection and `hydra` tags, 
it automatically fills struct fields with data fetched from database queries.

## Features

- **Automatic Hydration**: Automatically populates Go structs with data fetched from databases using reflection.
- **Multiple Database Support**: Supports MySQL, PostgreSQL, SQLite, Microsoft SQL Server, Oracle, MariaDB, and CockroachDB.
- **Flexible Queries**: Allows dynamic construction of SQL `WHERE` clauses.
- **Type Safety**: Ensures proper type conversions between database values and Go struct fields.
- **Easily Extendable**: Easily extendable to support more databases in the future.

## Installation

To install Sphire Hydra, use `go get`:

```bash
go get github.com/sphireinc/Hydra
```

## Supported Databases

- MySQL 
- PostgreSQL 
- SQLite
- Microsoft SQL Server
- Oracle
- MariaDB
- CockroachDB

## Usage

### Struct Definition

Define your structs with hydra tags to map the struct fields to the corresponding database columns:

```go
type Person struct {
    Name        string `json:"name" hydra:"name"`
    Age         int    `json:"age" hydra:"age"`
    Email       string `json:"email" hydra:"email"`
    hydra.Hydratable
}
```

### Hydration

To hydrate a struct, initialize the struct and then call the Hydrate method, which automatically fetches the
data from the database and populates the fields:

```go
package main

import (
    "database/sql"
    "github.com/sphireinc/Hydra"
    _ "github.com/go-sql-driver/mysql"
)

type Person struct {
    Name        string `json:"name" hydra:"name"`
    Age         int    `json:"age" hydra:"age"`
    Email       string `json:"email" hydra:"email"`
    hydra.Hydratable
}

func main() {
    // Initialize a MySQL database connection
    db, _ := sql.Open("mysql", "user:password@/dbname")

    // Create an addressable Person instance
    p := &Person{}
    p.Init(p)

    // Call Hydrate to populate the struct with data from the database
    p.Hydrate(db, map[string]interface{}{"id": "U6"})

    // Print the hydrated struct
    fmt.Printf("Hydrated person: %+v\n", p)
}
```

# Extensibility

Sphire Hydra is designed to be easily extensible. You can add support for additional databases by implementing a 
fetch function specific to the database’s query syntax and integrating it with the existing hydration process.

# Contributing

We welcome contributions! Feel free to open an issue or submit a pull request to improve the library.
