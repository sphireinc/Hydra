//go:build integration
// +build integration

package functional_tests

import "testing"

func TestDatabaseMySQL(t *testing.T) {
	db := mustOpenSQLDB(t, "mysql", "HYDRA_MYSQL_DSN")
	defer db.Close()

	runHydrateScenariosSQL(t, db, "mysql")
	runFetchScenariosSQL(t, db, "mysql")
}
