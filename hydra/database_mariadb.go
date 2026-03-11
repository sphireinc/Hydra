package hydra

import (
	"database/sql"

	_ "github.com/go-sql-driver/mysql"
)

func (h *Hydratable) fetchMariaDB(db *sql.DB, tableName string, whereClauses map[string]interface{}) (map[string]interface{}, error) {
	query, params := buildSelectQuery(tableName, whereClauses, questionPlaceholder)
	return queryFirstRowSQL(db, query, params)
}
