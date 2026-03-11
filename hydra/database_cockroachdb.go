package hydra

import (
	"context"

	"github.com/jackc/pgx/v5"
)

func (h *Hydratable) fetchCockroachDB(db *pgx.Conn, tableName string, whereClauses map[string]interface{}) (map[string]interface{}, error) {
	query, params := buildSelectQuery(tableName, whereClauses, dollarPlaceholder)
	return queryFirstRowPGX(context.Background(), db, query, params)
}
