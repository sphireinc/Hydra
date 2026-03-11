package hydra

import (
	"database/sql"

	_ "github.com/mattn/go-sqlite3"
)

func (h *Hydratable) fetchSQLite(db *sql.DB, tableName string, whereClauses map[string]interface{}) (map[string]interface{}, error) {
	query, params := buildSelectQuery(tableName, whereClauses, questionPlaceholder)
	return queryFirstRowSQL(db, query, params)
}
