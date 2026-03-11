package hydra

import (
	"database/sql"

	_ "github.com/godror/godror"
)

func (h *Hydratable) fetchOracle(db *sql.DB, tableName string, whereClauses map[string]interface{}) (map[string]interface{}, error) {
	query, params := buildSelectQuery(tableName, whereClauses, colonPlaceholder)
	return queryFirstRowSQL(db, query, params)
}
