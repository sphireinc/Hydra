package hydra

import "errors"

var (
	ErrNotInitialized   = errors.New("not initialized, cannot hydrate")
	ErrEmptyWhereClause = errors.New("at least one where clause is required")
	ErrNotFound         = errors.New("no matching row found")
)
