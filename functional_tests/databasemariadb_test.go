//go:build integration
// +build integration

package functional_tests

import "testing"

func TestDatabaseMariaDB(t *testing.T) {
	db := mustOpenSQLDB(t, "mysql", "HYDRA_MARIADB_DSN")
	defer db.Close()

	runHydrateScenariosSQL(t, db, "mariadb")
	runFetchScenariosSQL(t, db, "mariadb")
}
