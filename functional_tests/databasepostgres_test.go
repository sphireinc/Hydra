//go:build integration
// +build integration

package functional_tests

import (
	"context"
	"testing"
)

func TestDatabasePostgres(t *testing.T) {
	conn := mustOpenPGXConn(t, "HYDRA_POSTGRES_DSN")
	defer conn.Close(context.Background())

	runHydrateScenariosPGX(t, conn, "postgres")
	runFetchScenariosPGX(t, conn, "postgres")
}
