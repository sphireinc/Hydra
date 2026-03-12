package hydra

import (
	"context"
	"database/sql"
	"fmt"

	"github.com/jackc/pgx/v5"
)

type fetchRoute string

const (
	fetchRouteMySQL       fetchRoute = "mysql"
	fetchRouteSQLite      fetchRoute = "sqlite"
	fetchRouteMSSQL       fetchRoute = "mssql"
	fetchRouteMariaDB     fetchRoute = "mariadb"
	fetchRouteOracle      fetchRoute = "oracle"
	fetchRoutePostgres    fetchRoute = "postgres"
	fetchRouteCockroachDB fetchRoute = "cockroachdb"
)

func (h *Hydratable) resolveFetchRoute(db any) (fetchRoute, error) {
	switch db.(type) {
	case *sql.DB:
		switch h.XDBTypeOverride {
		case "sqlite":
			return fetchRouteSQLite, nil
		case "mssql":
			return fetchRouteMSSQL, nil
		case "mariadb":
			return fetchRouteMariaDB, nil
		case "oracle":
			return fetchRouteOracle, nil
		case "", "mysql":
			return fetchRouteMySQL, nil
		default:
			return "", fmt.Errorf("unsupported sql.DB override: %s", h.XDBTypeOverride)
		}
	case *pgx.Conn:
		switch h.XDBTypeOverride {
		case "cockroachdb":
			return fetchRouteCockroachDB, nil
		case "", "postgres":
			return fetchRoutePostgres, nil
		default:
			return "", fmt.Errorf("unsupported pgx override: %s", h.XDBTypeOverride)
		}
	default:
		return "", fmt.Errorf("unsupported database type: %T", db)
	}
}

func (h *Hydratable) Fetch(db any, tableName string, columns []string, whereClauses map[string]interface{}) (map[string]interface{}, error) {
	return h.FetchContext(context.Background(), db, tableName, columns, whereClauses)
}

func (h *Hydratable) FetchContext(ctx context.Context, db any, tableName string, columns []string, whereClauses map[string]interface{}) (map[string]interface{}, error) {
	route, err := h.resolveFetchRoute(db)
	if err != nil {
		return nil, err
	}

	switch typed := db.(type) {
	case *sql.DB:
		switch route {
		case fetchRouteSQLite:
			return h.fetchSQLite(ctx, typed, tableName, columns, whereClauses)
		case fetchRouteMSSQL:
			return h.fetchMSSQL(ctx, typed, tableName, columns, whereClauses)
		case fetchRouteMariaDB:
			return h.fetchMariaDB(ctx, typed, tableName, columns, whereClauses)
		case fetchRouteOracle:
			return h.fetchOracle(ctx, typed, tableName, columns, whereClauses)
		case fetchRouteMySQL:
			return h.fetchMySQL(ctx, typed, tableName, columns, whereClauses)
		}
	case *pgx.Conn:
		switch route {
		case fetchRouteCockroachDB:
			return h.fetchCockroachDB(ctx, typed, tableName, columns, whereClauses)
		case fetchRoutePostgres:
			return h.fetchPostgres(ctx, typed, tableName, columns, whereClauses)
		}
	}

	return nil, fmt.Errorf("unsupported database type: %T", db)
}
