package hydra

import (
	"database/sql"

	_ "github.com/go-sql-driver/mysql"
)

// fetchMySQL fetches data from a MySQL database using a query
// @param db The database connection
// @param params The parameters to pass to the query
// @return map[string]interface{} The hydrated data
// @return error The error if any occurred
func (h *Hydratable) fetchMySQL(db *sql.DB, tableName string, whereClauses map[string]interface{}) (map[string]interface{}, error) {
	query, params := buildSelectQuery(tableName, whereClauses, questionPlaceholder)
	return queryFirstRowSQL(db, query, params)
}
