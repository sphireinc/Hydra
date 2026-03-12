# Sphire Hydra

Sphire Hydra hydrates Go structs from database rows using `hydra` tags.

## Documentation

- [Hydra package guide](code/hydra.md)

## Quick commands

```bash
go test ./...
make test
make test-hydra
make test-func
```

## What to read first

The package guide covers:
- table-name resolution
- override behavior
- supported database handle types
- not-found behavior
- supported field types and conversions
- custom converters
- primary-key and lookup tags
- context-aware APIs

