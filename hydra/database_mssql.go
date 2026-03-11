package hydra

import (
	"database/sql"

	_ "github.com/denisenkom/go-mssqldb"
)

func (h *Hydratable) fetchMSSQL(db *sql.DB, tableName string, whereClauses map[string]interface{}) (map[string]interface{}, error) {
	query, params := buildSelectQuery(tableName, whereClauses, atPPlaceholder)
	return queryFirstRowSQL(db, query, params)
}
