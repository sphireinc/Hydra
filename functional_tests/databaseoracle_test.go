//go:build integration
// +build integration

package functional_tests

import "testing"

func TestDatabaseOracle(t *testing.T) {
	db := mustOpenSQLDB(t, "godror", "HYDRA_ORACLE_DSN")
	defer db.Close()

	runHydrateScenariosSQL(t, db, "oracle")
	runFetchScenariosSQL(t, db, "oracle")
}
