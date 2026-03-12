//go:build integration
// +build integration

package functional_tests

import (
	"context"
	"database/sql"
	"errors"
	"os"
	"strings"
	"testing"

	"github.com/jackc/pgx/v5"

	"github.com/sphireinc/Hydra/hydra"

	_ "github.com/denisenkom/go-mssqldb"
	_ "github.com/go-sql-driver/mysql"
	_ "github.com/godror/godror"
	_ "github.com/mattn/go-sqlite3"
)

func mustOpenSQLDB(t *testing.T, driverName string, envVars ...string) *sql.DB {
	t.Helper()

	dsn := firstEnv(envVars...)
	if dsn == "" {
		t.Skipf("skipping: missing DSN env var, tried %v", envVars)
	}

	db, err := sql.Open(driverName, dsn)
	if err != nil {
		t.Fatalf("sql.Open(%s): %v", driverName, err)
	}

	if err := db.Ping(); err != nil {
		db.Close()
		t.Fatalf("db.Ping(%s): %v", driverName, err)
	}

	return db
}

func mustOpenPGXConn(t *testing.T, envVars ...string) *pgx.Conn {
	t.Helper()

	dsn := firstEnv(envVars...)
	if dsn == "" {
		t.Skipf("skipping: missing DSN env var, tried %v", envVars)
	}

	conn, err := pgx.Connect(context.Background(), dsn)
	if err != nil {
		t.Fatalf("pgx.Connect: %v", err)
	}

	if err := conn.Ping(context.Background()); err != nil {
		conn.Close(context.Background())
		t.Fatalf("pgx.Ping: %v", err)
	}

	return conn
}

func firstEnv(keys ...string) string {
	for _, key := range keys {
		if value := strings.TrimSpace(os.Getenv(key)); value != "" {
			return value
		}
	}
	return ""
}

func mustCreateSQLiteFunctionalDB(t *testing.T) *sql.DB {
	t.Helper()

	db, err := sql.Open("sqlite3", ":memory:")
	if err != nil {
		t.Fatalf("sql.Open(sqlite3): %v", err)
	}

	if err := applySQLiteMigration(t, db, "migrations/sqlite.sql"); err != nil {
		db.Close()
		t.Fatalf("applySQLiteMigration: %v", err)
	}

	return db
}

func applySQLiteMigration(t *testing.T, db *sql.DB, path string) error {
	t.Helper()

	content, err := os.ReadFile(path)
	if err != nil {
		return err
	}

	var statements []string
	var builder strings.Builder

	for _, line := range strings.Split(string(content), "\n") {
		trimmed := strings.TrimSpace(line)
		if strings.HasPrefix(trimmed, "--") || trimmed == "" {
			continue
		}
		builder.WriteString(line)
		builder.WriteString("\n")
	}

	for _, stmt := range strings.Split(builder.String(), ";") {
		trimmed := strings.TrimSpace(stmt)
		if trimmed == "" {
			continue
		}
		statements = append(statements, trimmed)
	}

	for _, stmt := range statements {
		if _, err := db.Exec(stmt); err != nil {
			return err
		}
	}

	return nil
}

func validateSeedAssumptionsSQL(t *testing.T, db *sql.DB) {
	t.Helper()

	var count int
	if err := db.QueryRow(`SELECT COUNT(*) FROM Person`).Scan(&count); err != nil {
		t.Fatalf("count Person rows: %v", err)
	}
	if count != 10 {
		t.Fatalf("expected 10 Person rows, got %d", count)
	}

	if err := db.QueryRow(`SELECT COUNT(*) FROM Addresses`).Scan(&count); err != nil {
		t.Fatalf("count Addresses rows: %v", err)
	}
	if count != 10 {
		t.Fatalf("expected 10 Addresses rows, got %d", count)
	}

	var firstName, lastName, sex string
	if err := db.QueryRow(`SELECT first_name, last_name, sex FROM Person WHERE id = 1`).Scan(&firstName, &lastName, &sex); err != nil {
		t.Fatalf("seed person row: %v", err)
	}
	if firstName != "John" || lastName != "Doe" || sex != "M" {
		t.Fatalf("unexpected seeded person row: %q %q %q", firstName, lastName, sex)
	}

	var address1, address2, city, state, postalCode, country string
	var userID int
	if err := db.QueryRow(`SELECT user_id, address_1, address_2, city, state, postal_code, country FROM Addresses WHERE id = 1`).Scan(
		&userID, &address1, &address2, &city, &state, &postalCode, &country,
	); err != nil {
		t.Fatalf("seed address row: %v", err)
	}
	if userID != 1 || address1 != "123 Main St" || address2 != "Apt 4" || city != "New York" || state != "NY" || postalCode != "10001" || country != "USA" {
		t.Fatalf("unexpected seeded address row: userID=%d address1=%q address2=%q city=%q state=%q postalCode=%q country=%q",
			userID, address1, address2, city, state, postalCode, country)
	}
}

func validateSeedAssumptionsPGX(t *testing.T, conn *pgx.Conn) {
	t.Helper()

	var count int
	if err := conn.QueryRow(context.Background(), `SELECT COUNT(*) FROM Person`).Scan(&count); err != nil {
		t.Fatalf("count Person rows: %v", err)
	}
	if count != 10 {
		t.Fatalf("expected 10 Person rows, got %d", count)
	}

	if err := conn.QueryRow(context.Background(), `SELECT COUNT(*) FROM Addresses`).Scan(&count); err != nil {
		t.Fatalf("count Addresses rows: %v", err)
	}
	if count != 10 {
		t.Fatalf("expected 10 Addresses rows, got %d", count)
	}

	var firstName, lastName, sex string
	if err := conn.QueryRow(context.Background(), `SELECT first_name, last_name, sex FROM Person WHERE id = 1`).Scan(&firstName, &lastName, &sex); err != nil {
		t.Fatalf("seed person row: %v", err)
	}
	if firstName != "John" || lastName != "Doe" || sex != "M" {
		t.Fatalf("unexpected seeded person row: %q %q %q", firstName, lastName, sex)
	}

	var address1, address2, city, state, postalCode, country string
	var userID int
	if err := conn.QueryRow(context.Background(), `SELECT user_id, address_1, address_2, city, state, postal_code, country FROM Addresses WHERE id = 1`).Scan(
		&userID, &address1, &address2, &city, &state, &postalCode, &country,
	); err != nil {
		t.Fatalf("seed address row: %v", err)
	}
	if userID != 1 || address1 != "123 Main St" || address2 != "Apt 4" || city != "New York" || state != "NY" || postalCode != "10001" || country != "USA" {
		t.Fatalf("unexpected seeded address row: userID=%d address1=%q address2=%q city=%q state=%q postalCode=%q country=%q",
			userID, address1, address2, city, state, postalCode, country)
	}
}

func assertHydratedPerson(t *testing.T, person *Person) {
	t.Helper()

	if person.ID != 1 {
		t.Fatalf("expected person.ID = 1, got %d", person.ID)
	}
	if person.FirstName != "John" {
		t.Fatalf("expected person.FirstName = John, got %q", person.FirstName)
	}
	if person.LastName != "Doe" {
		t.Fatalf("expected person.LastName = Doe, got %q", person.LastName)
	}
	if person.Sex != "M" {
		t.Fatalf("expected person.Sex = M, got %q", person.Sex)
	}
}

func assertHydratedAddress(t *testing.T, address *Address) {
	t.Helper()

	if address.ID != 1 {
		t.Fatalf("expected address.ID = 1, got %d", address.ID)
	}
	if address.UserID != 1 {
		t.Fatalf("expected address.UserID = 1, got %d", address.UserID)
	}
	if address.Address1 != "123 Main St" {
		t.Fatalf("expected address.Address1 = 123 Main St, got %q", address.Address1)
	}
	if address.Address2 != "Apt 4" {
		t.Fatalf("expected address.Address2 = Apt 4, got %q", address.Address2)
	}
	if address.City != "New York" {
		t.Fatalf("expected address.City = New York, got %q", address.City)
	}
	if address.State != "NY" {
		t.Fatalf("expected address.State = NY, got %q", address.State)
	}
	if address.PostalCode != "10001" {
		t.Fatalf("expected address.PostalCode = 10001, got %q", address.PostalCode)
	}
	if address.Country != "USA" {
		t.Fatalf("expected address.Country = USA, got %q", address.Country)
	}
}

func runHydrateScenariosSQL(t *testing.T, db *sql.DB, dbOverride string) {
	t.Helper()

	validateSeedAssumptionsSQL(t, db)

	person := &Person{}
	person.Init(person)
	person.XDBTypeOverride = dbOverride

	if err := person.Hydrate(db, map[string]interface{}{"id": 1}); err != nil {
		t.Fatalf("hydrate person: %v", err)
	}
	assertHydratedPerson(t, person)

	address := &Address{}
	address.Init(address)
	address.XDBTypeOverride = dbOverride

	if err := address.Hydrate(db, map[string]interface{}{"id": 1}); err != nil {
		t.Fatalf("hydrate address: %v", err)
	}
	assertHydratedAddress(t, address)

	missing := &Person{}
	missing.Init(missing)
	missing.XDBTypeOverride = dbOverride

	err := missing.Hydrate(db, map[string]interface{}{"id": 9999})
	if !errors.Is(err, hydra.ErrNotFound) {
		t.Fatalf("expected hydra.ErrNotFound, got %v", err)
	}
}

func runHydrateScenariosPGX(t *testing.T, conn *pgx.Conn, dbOverride string) {
	t.Helper()

	validateSeedAssumptionsPGX(t, conn)

	person := &Person{}
	person.Init(person)
	person.XDBTypeOverride = dbOverride

	if err := person.Hydrate(conn, map[string]interface{}{"id": 1}); err != nil {
		t.Fatalf("hydrate person: %v", err)
	}
	assertHydratedPerson(t, person)

	address := &Address{}
	address.Init(address)
	address.XDBTypeOverride = dbOverride

	if err := address.Hydrate(conn, map[string]interface{}{"id": 1}); err != nil {
		t.Fatalf("hydrate address: %v", err)
	}
	assertHydratedAddress(t, address)

	missing := &Person{}
	missing.Init(missing)
	missing.XDBTypeOverride = dbOverride

	err := missing.Hydrate(conn, map[string]interface{}{"id": 9999})
	if !errors.Is(err, hydra.ErrNotFound) {
		t.Fatalf("expected hydra.ErrNotFound, got %v", err)
	}
}

func runFetchScenariosSQL(t *testing.T, db *sql.DB, dbOverride string) {
	t.Helper()

	validateSeedAssumptionsSQL(t, db)

	h := &hydra.Hydratable{XDBTypeOverride: dbOverride}

	personColumns := []string{"id", "first_name", "last_name", "sex"}
	personRow, err := h.Fetch(db, "Person", personColumns, map[string]interface{}{"id": 1})
	if err != nil {
		t.Fatalf("fetch person: %v", err)
	}

	if personRow["id"] == nil || personRow["first_name"] == nil || personRow["last_name"] == nil || personRow["sex"] == nil {
		t.Fatalf("unexpected fetched person row: %#v", personRow)
	}

	addressColumns := []string{"id", "user_id", "address_1", "address_2", "city", "state", "postal_code", "country"}
	addressRow, err := h.Fetch(db, "Addresses", addressColumns, map[string]interface{}{"id": 1})
	if err != nil {
		t.Fatalf("fetch address: %v", err)
	}

	if addressRow["id"] == nil || addressRow["user_id"] == nil || addressRow["address_1"] == nil || addressRow["city"] == nil {
		t.Fatalf("unexpected fetched address row: %#v", addressRow)
	}

	_, err = h.Fetch(db, "Person", personColumns, map[string]interface{}{"id": 9999})
	if !errors.Is(err, hydra.ErrNotFound) {
		t.Fatalf("expected hydra.ErrNotFound, got %v", err)
	}
}

func runFetchScenariosPGX(t *testing.T, conn *pgx.Conn, dbOverride string) {
	t.Helper()

	validateSeedAssumptionsPGX(t, conn)

	h := &hydra.Hydratable{XDBTypeOverride: dbOverride}

	personColumns := []string{"id", "first_name", "last_name", "sex"}
	personRow, err := h.Fetch(conn, "Person", personColumns, map[string]interface{}{"id": 1})
	if err != nil {
		t.Fatalf("fetch person: %v", err)
	}

	if personRow["id"] == nil || personRow["first_name"] == nil || personRow["last_name"] == nil || personRow["sex"] == nil {
		t.Fatalf("unexpected fetched person row: %#v", personRow)
	}

	addressColumns := []string{"id", "user_id", "address_1", "address_2", "city", "state", "postal_code", "country"}
	addressRow, err := h.Fetch(conn, "Addresses", addressColumns, map[string]interface{}{"id": 1})
	if err != nil {
		t.Fatalf("fetch address: %v", err)
	}

	if addressRow["id"] == nil || addressRow["user_id"] == nil || addressRow["address_1"] == nil || addressRow["city"] == nil {
		t.Fatalf("unexpected fetched address row: %#v", addressRow)
	}

	_, err = h.Fetch(conn, "Person", personColumns, map[string]interface{}{"id": 9999})
	if !errors.Is(err, hydra.ErrNotFound) {
		t.Fatalf("expected hydra.ErrNotFound, got %v", err)
	}
}
