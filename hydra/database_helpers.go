package hydra

import (
	"context"
	"database/sql"
	"fmt"
	"sort"
	"strings"

	"github.com/jackc/pgx/v5"
)

type placeholderFormatter func(index int) string

func questionPlaceholder(_ int) string {
	return "?"
}

func dollarPlaceholder(index int) string {
	return fmt.Sprintf("$%d", index)
}

func colonPlaceholder(index int) string {
	return fmt.Sprintf(":%d", index)
}

func atPPlaceholder(index int) string {
	return fmt.Sprintf("@p%d", index)
}

func buildSelectQuery(tableName string, whereClauses map[string]interface{}, format placeholderFormatter) (string, []interface{}) {
	query := fmt.Sprintf("SELECT * FROM %s WHERE ", tableName)
	p("Query:", query)

	keys := make([]string, 0, len(whereClauses))
	for column := range whereClauses {
		keys = append(keys, column)
	}
	sort.Strings(keys)

	params := make([]interface{}, 0, len(keys))
	conditions := make([]string, 0, len(keys))

	for i, column := range keys {
		conditions = append(conditions, fmt.Sprintf("%s = %s", column, format(i+1)))
		params = append(params, whereClauses[column])
	}

	query += strings.Join(conditions, " AND ")
	p("Query with conditions:", query)
	p("Query params:", params)

	return query, params
}

func queryFirstRowSQL(db *sql.DB, query string, params []interface{}) (map[string]interface{}, error) {
	rows, err := db.Query(query, params...)
	p("Rows:", rows)
	if err != nil {
		p("Error executing query:", err)
		return nil, err
	}
	defer rows.Close()

	return scanSQLRowsFirst(rows)
}

func scanSQLRowsFirst(rows *sql.Rows) (map[string]interface{}, error) {
	columns, err := rows.Columns()
	p("Columns:", columns)
	if err != nil {
		p("Error getting columns:", err)
		return nil, err
	}

	result := make(map[string]interface{})
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

	if err := rows.Err(); err != nil {
		return nil, err
	}

	p("Result:", result)
	return result, nil
}

func queryFirstRowPGX(ctx context.Context, db *pgx.Conn, query string, params []interface{}) (map[string]interface{}, error) {
	rows, err := db.Query(ctx, query, params...)
	p("Rows:", rows)
	if err != nil {
		p("Error executing query:", err)
		return nil, err
	}
	defer rows.Close()

	return scanPGXRowsFirst(rows)
}

func scanPGXRowsFirst(rows pgx.Rows) (map[string]interface{}, error) {
	fieldDescriptions := rows.FieldDescriptions()
	columns := make([]string, len(fieldDescriptions))
	for i, field := range fieldDescriptions {
		columns[i] = field.Name
	}
	p("Columns:", columns)

	result := make(map[string]interface{})
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

	if err := rows.Err(); err != nil {
		return nil, err
	}

	p("Result:", result)
	return result, nil
}
