# Functional Testing Suite

This folder contains a set of functional tests for the **Sphire Hydra** library. The tests are 
written in Go and use the `testify` package for assertions.

## Running the Tests

To run the tests, we will use Docker to spin up containers with the required databases. This ensures 
consistency across environments and makes it easy to test against multiple databases.

### Step 1: Clone the Repository

First, clone the Hydra repository:

```bash
git clone https://github.com/sphireinc/Hydra.git
cd Hydra/tests
```

### Step 2: Set Up Docker Environment

To run the functional tests, you need to spin up Docker containers for all supported databases (MySQL, PostgreSQL, SQLite, MSSQL, Oracle, MariaDB, CockroachDB). The docker-compose.yml file in this directory has already been configured to set up these databases.

Start the Docker environment with:

```bash
docker-compose up -d
```

This command will:

- Spin up the necessary databases.
- Set up each container with test data via the respective SQL migration files (./migrations directory).

### Step 3: Run the Tests

Once the Docker environment is running, you can run the tests using the following command:

```bash
go test -v ./tests
```

This will execute all tests across the various databases configured in the test suite.

### Step 4: Stop the Docker Environment

After running the tests, you can stop the Docker containers with:

```bash
docker-compose down
```

This will tear down the test environment.

## Test Structure

To add a new test, create a new file in the tests directory with a name that describes the test. For example, if you want to test the Hydrate method of the Person struct, you can create a file named TestHydrate.go.

In the test file, import the necessary packages and define the test cases. For example:

```go
package tests

import (
    "Hydrator/hydra"
    "database/sql"
    "github.com/stretchr/testify/assert"
    _ "github.com/go-sql-driver/mysql" 
)

type Person struct {
    Name        string `json:"name" hydra:"name"`
    Age         int    `json:"age" hydra:"age"`
    Email       string `json:"email" hydra:"email"`
    hydra.Hydratable
}

func createDBConnection() *sql.DB {
    db, _ := sql.Open("mysql", "user:password@tcp(mysql-db:3306)/testdb")
    return db
}

func TestHydratePerson(t *testing.T) {
    db := createDBConnection()
    p := &Person{}
    p.Init(p)

    whereClause := map[string]interface{}{"id": "1"}
    p.Hydrate(db, whereClause)

    assert.Equal(t, "John Doe", p.Name)
    assert.Equal(t, 30, p.Age)
    assert.Equal(t, "john.doe@example.com", p.Email)
}
```

## Test Databases

The following databases are spun up for functional tests:

- MySQL: Port 3306
- MariaDB: Port 3307
- PostgreSQL: Port 5432
- SQLite: Embedded in the container
- Microsoft SQL Server (MSSQL): Port 1433
- Oracle: Port 1521
- CockroachDB: Port 26257 (Admin UI on port 8080)

Each database is seeded with two tables (Person, Addresses) and 10 rows of sample data.

This ensures that all functional tests are run against multiple database systems, ensuring 
full compatibility of the Hydra library across different environments.