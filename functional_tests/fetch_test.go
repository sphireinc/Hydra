//go:build integration
// +build integration

package functional_tests

import "testing"

func TestFetchSQLite(t *testing.T) {
	db := mustCreateSQLiteFunctionalDB(t)
	defer db.Close()

	runFetchScenariosSQL(t, db, "sqlite")
}
