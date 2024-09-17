package hydra

import (
	"database/sql"
	"fmt"
	"strings"

	_ "github.com/go-sql-driver/mysql"
)

func (h *Hydratable) fetchMariaDB(db *sql.DB, tableName string, whereClauses map[string]interface{}) (map[string]interface{}, error) {
	// Build the base query with the table name
	query := fmt.Sprintf("SELECT * FROM %s WHERE ", tableName)
	p("Query:", query)

	// Prepare the values for the SQL parameters
	var params []interface{}
	var conditions []string

	// Build the WHERE clause dynamically
	for column, value := range whereClauses {
		// Append the column name to the condition (safe concatenation)
		conditions = append(conditions, fmt.Sprintf("%s = ?", column))
		// Append the value to the params slice
		params = append(params, value)
	}

	// Join all conditions with "AND" and append to the query
	query += strings.Join(conditions, " AND ")
	p("Query with conditions:", query)

	// Execute the query with the dynamic parameters
	rows, err := db.Query(query, params...)
	p("Rows:", rows)
	if err != nil {
		p("Error executing query:", err)
		return nil, err
	}
	defer rows.Close()

	// Assuming a single row for hydration
	result := make(map[string]interface{})
	columns, err := rows.Columns()
	p("Columns:", columns)
	if err != nil {
		p("Error getting columns:", err)
		return nil, err
	}

	values := make([]interface{}, len(columns))
	valuePtrs := make([]interface{}, len(columns))
	for i := range values {
		valuePtrs[i] = &values[i]
	}

	if rows.Next() {
		if err := rows.Scan(valuePtrs...); err != nil {
			return nil, err
		}

		for i, col := range columns {
			result[col] = values[i]
		}
	}

	p("Result:", result)
	return result, nil
}
