package hydra

import (
	"context"

	"github.com/jackc/pgx/v5"
)

func (h *Hydratable) fetchPostgres(ctx context.Context, db *pgx.Conn, tableName string, columns []string, whereClauses map[string]interface{}) (map[string]interface{}, error) {
	query, params, err := buildSelectQuery(tableName, columns, whereClauses, dollarPlaceholder)
	if err != nil {
		return nil, err
	}
	return queryFirstRowPGX(ctx, db, query, params)
}
