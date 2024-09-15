package hydra

import (
	"database/sql"
	"fmt"
	"github.com/jackc/pgx/v5"
)

// Fetch hydrates the object with data from the database
// @param db The database connection
// @param tableName The name of the table to fetch data from
// @param whereClauses The where clauses to filter the data
// @return map[string]interface{} The hydrated data
// @return error The error if any occurred
func (h *Hydratable) Fetch(db any, tableName string, whereClauses map[string]interface{}) (map[string]interface{}, error) {
	p("Fetching data from database")
	switch db := db.(type) {
	case *sql.DB:
		switch h.XDBTypeOverride {
		case "sqlite":
			p("Fetching data from SQLite")
			return h.fetchSQLite(db, tableName, whereClauses)
		case "mssql":
			p("Fetching data from MSSQL")
			return h.fetchMSSQL(db, tableName, whereClauses)
		case "mariadb":
			p("Fetching data from MariaDB")
			return h.fetchMariaDB(db, tableName, whereClauses)
		case "oracle":
			p("Fetching data from Oracle")
			return h.fetchOracle(db, tableName, whereClauses)
		case "mysql":
		default:
			p("Fetching data from MySQL")
			return h.fetchMySQL(db, tableName, whereClauses)
		}
	case *pgx.Conn:
		switch h.XDBTypeOverride {
		case "cockroachdb":
			p("Fetching data from CockroachDB")
			return h.fetchCockroachDB(db, tableName, whereClauses)
		case "postgres":
		default:
			p("Fetching data from PostgreSQL")
			return h.fetchPostgres(db, tableName, whereClauses)
		}
	default:
		err := fmt.Errorf("unsupported database type: %T", db)
		p("Unsupported database type:", err)
		return nil, err
	}
	return nil, fmt.Errorf("unsupported database type: %T", db)
}
