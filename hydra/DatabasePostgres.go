package hydra

import (
	"context"
	"fmt"
	"github.com/jackc/pgx/v5"
	"strings"
)

func (h *Hydratable) fetchPostgres(db *pgx.Conn, tableName string, whereClauses map[string]interface{}) (map[string]interface{}, error) {
	// Build the base query with the table name
	query := fmt.Sprintf("SELECT * FROM %s WHERE ", tableName)
	p("Query:", query)

	// Prepare the values for the SQL parameters
	var params []interface{}
	var conditions []string
	paramIndex := 1 // Postgres uses $1, $2... for placeholders

	// Build the WHERE clause dynamically
	for column, value := range whereClauses {
		// Append the column name to the condition (Postgres uses positional params like $1, $2...)
		conditions = append(conditions, fmt.Sprintf("%s = $%d", column, paramIndex))
		// Append the value to the params slice
		params = append(params, value)
		paramIndex++
	}

	// Join all conditions with "AND" and append to the query
	query += strings.Join(conditions, " AND ")
	p("Query with conditions:", query)

	// Execute the query with the dynamic parameters
	rows, err := db.Query(context.Background(), query, params...)
	p("Rows:", rows)
	if err != nil {
		p("Error executing query:", err)
		return nil, err
	}
	defer rows.Close()

	// Assuming a single row for hydration
	result := make(map[string]interface{})
	fieldDescriptions := rows.FieldDescriptions()
	columns := make([]string, len(fieldDescriptions))
	for i, field := range fieldDescriptions {
		columns[i] = string(field.Name)
	}
	p("Columns:", columns)

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
