#!/usr/bin/env bash
set -euo pipefail

ROOT_DIR="$(cd "$(dirname "${BASH_SOURCE[0]}")" && pwd)"
COMPOSE_FILE="${ROOT_DIR}/functional_tests/docker-compose.yml"
PROJECT_NAME="hydra-functional-tests"

cleanup() {
  echo
  echo "==> Cleaning up functional test containers"
  docker compose -p "${PROJECT_NAME}" -f "${COMPOSE_FILE}" down -v --remove-orphans || true
}

trap cleanup EXIT

echo "==> Running Hydra functional test suite"
echo "==> Compose file: ${COMPOSE_FILE}"

docker compose \
  -p "${PROJECT_NAME}" \
  -f "${COMPOSE_FILE}" \
  up \
  --build \
  --abort-on-container-exit \
  --exit-code-from test-runner