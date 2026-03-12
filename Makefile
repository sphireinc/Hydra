SHELL := /usr/bin/env bash

PROJECT_NAME := hydra-functional-tests
FUNC_COMPOSE_FILE := functional_tests/docker-compose.yml

.PHONY: help test test-unit test-hydra test-func test-func-keep test-func-down

help:
	@echo "Available targets:"
	@echo "  make test         - run unit tests"
	@echo "  make test-unit    - run unit tests"
	@echo "  make test-hydra   - run hydra package tests"
	@echo "  make test-func    - run full functional test stack with docker compose"
	@echo "  make test-func-keep - run functional tests but keep containers afterward"
	@echo "  make test-func-down - tear down functional test containers"

test: test-unit

test-unit:
	go test ./...

test-hydra:
	go test ./hydra/...

test-func:
	./run_func_tests.sh

test-func-keep:
	docker compose \
		-p $(PROJECT_NAME) \
		-f $(FUNC_COMPOSE_FILE) \
		up \
		--build \
		--abort-on-container-exit \
		--exit-code-from test-runner

test-func-down:
	docker compose \
		-p $(PROJECT_NAME) \
		-f $(FUNC_COMPOSE_FILE) \
		down -v --remove-orphans