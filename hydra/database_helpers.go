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

func questionPlaceholder(_ int) string   { return "?" }
func dollarPlaceholder(index int) string { return fmt.Sprintf("$%d", index) }
func colonPlaceholder(index int) string  { return fmt.Sprintf(":%d", index) }
func atPPlaceholder(index int) string    { return fmt.Sprintf("@p%d", index) }

func buildSelectQuery(tableName string, columns []string, whereClauses map[string]interface{}, format placeholderFormatter) (string, []interface{}, error) {
	if err := validateIdentifier(tableName); err != nil {
		return "", nil, err
	}
	if len(columns) == 0 {
		return "", nil, fmt.Errorf("at least one selected column is required")
	}
	for _, column := range columns {
		if err := validateIdentifier(column); err != nil {
			return "", nil, err
		}
	}
	if len(whereClauses) == 0 {
		return "", nil, ErrEmptyWhereClause
	}

	keys := make([]string, 0, len(whereClauses))
	for column := range whereClauses {
		if err := validateIdentifier(column); err != nil {
			return "", nil, err
		}
		keys = append(keys, column)
	}
	sort.Strings(keys)

	params := make([]interface{}, 0, len(keys))
	conditions := make([]string, 0, len(keys))
	for i, column := range keys {
		conditions = append(conditions, fmt.Sprintf("%s = %s", column, format(i+1)))
		params = append(params, whereClauses[column])
	}

	query := fmt.Sprintf(
		"SELECT %s FROM %s WHERE %s",
		strings.Join(columns, ", "),
		tableName,
		strings.Join(conditions, " AND "),
	)

	return query, params, nil
}

func queryFirstRowSQL(ctx context.Context, db *sql.DB, query string, params []interface{}) (map[string]interface{}, error) {
	rows, err := db.QueryContext(ctx, query, params...)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	return scanSQLRowsFirst(rows)
}

func scanSQLRowsFirst(rows *sql.Rows) (map[string]interface{}, error) {
	columns, err := rows.Columns()
	if err != nil {
		return nil, err
	}

	values := make([]interface{}, len(columns))
	valuePtrs := make([]interface{}, len(columns))
	for i := range values {
		valuePtrs[i] = &values[i]
	}

	if !rows.Next() {
		if err := rows.Err(); err != nil {
			return nil, err
		}
		return nil, ErrNotFound
	}

	if err := rows.Scan(valuePtrs...); err != nil {
		return nil, err
	}

	result := make(map[string]interface{}, len(columns))
	for i, col := range columns {
		result[col] = values[i]
	}

	if err := rows.Err(); err != nil {
		return nil, err
	}

	return result, nil
}

func queryFirstRowPGX(ctx context.Context, db *pgx.Conn, query string, params []interface{}) (map[string]interface{}, error) {
	rows, err := db.Query(ctx, query, params...)
	if err != nil {
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

	values := make([]interface{}, len(columns))
	valuePtrs := make([]interface{}, len(columns))
	for i := range values {
		valuePtrs[i] = &values[i]
	}

	if !rows.Next() {
		if err := rows.Err(); err != nil {
			return nil, err
		}
		return nil, ErrNotFound
	}

	if err := rows.Scan(valuePtrs...); err != nil {
		return nil, err
	}

	result := make(map[string]interface{}, len(columns))
	for i, col := range columns {
		result[col] = values[i]
	}

	if err := rows.Err(); err != nil {
		return nil, err
	}

	return result, nil
}
