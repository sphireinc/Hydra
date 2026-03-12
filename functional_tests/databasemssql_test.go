//go:build integration
// +build integration

package functional_tests

import "testing"

func TestDatabaseMSSQL(t *testing.T) {
	db := mustOpenSQLDB(t, "sqlserver", "HYDRA_MSSQL_DSN")
	defer db.Close()

	runHydrateScenariosSQL(t, db, "mssql")
	runFetchScenariosSQL(t, db, "mssql")
}
