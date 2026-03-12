#!/usr/bin/env bash
set -euo pipefail

echo "==> Hydra functional test runner starting"

required_env_vars=(
  HYDRA_MYSQL_DSN
  HYDRA_MARIADB_DSN
  HYDRA_POSTGRES_DSN
  HYDRA_MSSQL_DSN
  HYDRA_ORACLE_DSN
  HYDRA_COCKROACHDB_DSN
)

for var_name in "${required_env_vars[@]}"; do
  if [[ -z "${!var_name:-}" ]]; then
    echo "ERROR: required environment variable ${var_name} is not set"
    exit 1
  fi
done

echo "==> Verifying migration files exist"
migration_files=(
  /app/functional_tests/migrations/mysql.sql
  /app/functional_tests/migrations/mariadb.sql
  /app/functional_tests/migrations/postgres.sql
  /app/functional_tests/migrations/mssql.sql
  /app/functional_tests/migrations/oracle.sql
  /app/functional_tests/migrations/cockroachdb.sql
  /app/functional_tests/migrations/sqlite.sql
)

for file_path in "${migration_files[@]}"; do
  if [[ ! -f "${file_path}" ]]; then
    echo "ERROR: missing migration file ${file_path}"
    exit 1
  fi
done

echo "==> Starting unified integration suite"
echo "    SQLite functional tests run inside this same test-runner container"
echo "    Server-backed functional tests use the DSNs injected by docker compose"

go test -tags=integration -count=1 -v ./functional_tests/...