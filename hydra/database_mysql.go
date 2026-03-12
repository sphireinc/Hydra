package hydra

import (
	"context"
	"database/sql"

	_ "github.com/go-sql-driver/mysql"
)

func (h *Hydratable) fetchMySQL(ctx context.Context, db *sql.DB, tableName string, columns []string, whereClauses map[string]interface{}) (map[string]interface{}, error) {
	query, params, err := buildSelectQuery(tableName, columns, whereClauses, questionPlaceholder)
	if err != nil {
		return nil, err
	}
	return queryFirstRowSQL(ctx, db, query, params)
}
