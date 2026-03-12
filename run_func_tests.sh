#!/usr/bin/env bash
set -euo pipefail

ROOT_DIR="$(cd "$(dirname "${BASH_SOURCE[0]}")" && pwd)"
COMPOSE_FILE="${ROOT_DIR}/functional_tests/docker-compose.yml"
PROJECT_NAME="hydra-functional-tests"

cleanup() {
  local exit_code=$?
  echo
  echo "==> Cleaning up functional test containers"
  docker compose \
    -p "${PROJECT_NAME}" \
    -f "${COMPOSE_FILE}" \
    down -v --remove-orphans || true
  exit "${exit_code}"
}

trap cleanup EXIT

ENABLE_MSSQL="${ENABLE_MSSQL:-0}"
ENABLE_ORACLE="${ENABLE_ORACLE:-0}"

echo "==> Running Hydra functional test suite"
echo "==> Compose file: ${COMPOSE_FILE}"

compose_cmd=(
  docker compose
  -p "${PROJECT_NAME}"
  -f "${COMPOSE_FILE}"
)

if [[ "${ENABLE_MSSQL}" == "1" ]]; then
  echo "==> MSSQL profile enabled"
  : "${HYDRA_MSSQL_DSN:=sqlserver://sa:YourStrongPassword!123@mssql:1433?database=testdb&encrypt=disable}"
  export HYDRA_MSSQL_DSN
  compose_cmd+=(--profile mssql)
else
  echo "==> MSSQL profile disabled"
fi

if [[ "${ENABLE_ORACLE}" == "1" ]]; then
  echo "==> Oracle profile enabled"
  : "${HYDRA_ORACLE_DSN:=hydra/hydrapassword@oracle:1521/XEPDB1}"
  export HYDRA_ORACLE_DSN
  compose_cmd+=(--profile oracle)
else
  echo "==> Oracle profile disabled"
fi

compose_cmd+=(
  up
  --build
  --abort-on-container-exit
  --exit-code-from test-runner
)

"${compose_cmd[@]}"