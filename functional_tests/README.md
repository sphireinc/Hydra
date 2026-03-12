# Functional Testing Suite

This directory contains executable integration tests for Hydra.

## How the suite is separated from unit tests

These tests are behind the `integration` build tag, so they do not run during normal unit-test execution.

Run them explicitly with:

```bash
go test -tags=integration -v ./functional_tests/...
```

## What the functional suite covers

The suite performs real end-to-end checks for:
- hydrating a Person by ID
- hydrating an Address by ID
- verifying expected field values from seeded data
- verifying ErrNotFound when a row is missing
- validating that the migration/seed data matches the test assumptions

## SQLite

SQLite is fully local and uses the actual migrations/sqlite.sql file to create and seed an in-memory database.

This makes the SQLite functional tests executable without Docker or external services.

## External databases

The other functional tests are real integration tests, but they only run when the appropriate DSN environment 
variable is present. Otherwise they are skipped.

Supported environment variables:
- HYDRA_MYSQL_DSN
- HYDRA_MARIADB_DSN
- HYDRA_MSSQL_DSN
- HYDRA_ORACLE_DSN
- HYDRA_POSTGRES_DSN
- HYDRA_COCKROACHDB_DSN

Examples:

```bash
HYDRA_MYSQL_DSN='testuser:testpassword@tcp(127.0.0.1:3306)/testdb' \
go test -tags=integration -v ./functional_tests -run TestDatabaseMySQL


HYDRA_POSTGRES_DSN='postgres://testuser:testpassword@127.0.0.1:5433/testdb?sslmode=disable' \
go test -tags=integration -v ./functional_tests -run TestDatabasePostgres
```


## Notes on table naming

The migrations create Person and Addresses tables. The functional test structs implement `HydraTableName()` so the
hydrator targets those exact table names.

Seed expectations validated by the suite

The suite verifies that:
- Person contains 10 rows
- Addresses contains 10 rows
- Person(id=1) is John Doe, sex M
- Addresses(id=1) belongs to user 1 and is 123 Main St, Apt 4, New York, NY, 10001, USA