//go:build integration
// +build integration

package functional_tests

import "testing"

func TestHydrateSQLite(t *testing.T) {
	db := mustCreateSQLiteFunctionalDB(t)
	defer db.Close()

	runHydrateScenariosSQL(t, db, "sqlite")
}
