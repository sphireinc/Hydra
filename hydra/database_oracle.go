package hydra

import (
	"context"
	"database/sql"

	_ "github.com/godror/godror"
)

func (h *Hydratable) fetchOracle(ctx context.Context, db *sql.DB, tableName string, columns []string, whereClauses map[string]interface{}) (map[string]interface{}, error) {
	query, params, err := buildSelectQuery(tableName, columns, whereClauses, colonPlaceholder)
	if err != nil {
		return nil, err
	}
	return queryFirstRowSQL(ctx, db, query, params)
}
