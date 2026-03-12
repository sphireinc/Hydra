//go:build integration
// +build integration

package functional_tests

import (
	"context"
	"testing"
)

func TestDatabaseCockroachDB(t *testing.T) {
	conn := mustOpenPGXConn(t, "HYDRA_COCKROACHDB_DSN")
	defer conn.Close(context.Background())

	runHydrateScenariosPGX(t, conn, "cockroachdb")
	runFetchScenariosPGX(t, conn, "cockroachdb")
}
