package postgres

import (
	"context"
	"database/sql"
	"fmt"
	"log"
	"os"
	"testing"
	"time"

	"entain/pkg/database/postgres"
	"entain/schema/postgresql"

	"github.com/docker/go-connections/nat"
	"github.com/stretchr/testify/require"
	"github.com/testcontainers/testcontainers-go"
	"github.com/testcontainers/testcontainers-go/wait"

	_ "github.com/lib/pq"
)

var (
	testSQLDB *sql.DB
)

func TestMain(m *testing.M) {
	var exitCode int

	runTestContainer(func() {
		exitCode = m.Run()
	})

	os.Exit(exitCode)
}

func runTestContainer(testRunner func()) {
	const (
		databaseName = "test"
		username     = "admin"
		password     = "test"
		exposedPort  = "5432/tcp"
	)

	// Define if we will use reaper to clean up resources.
	//
	skipReaper := os.Getenv("TESTCONTAINERS_RYUK_DISABLED") != ""
	log.Printf("flag SkipReaper = %v", skipReaper)

	// Setup Database.
	//
	ctx := context.Background()
	req := testcontainers.ContainerRequest{
		Image:        "postgres:14.17", // PostgreSQL 14.17 on x86_64-pc-linux-gnu, compiled by Debian clang version 12.0.1, 64-bit
		ExposedPorts: []string{exposedPort},
		WaitingFor: wait.ForSQL(exposedPort, "postgres", func(host string, port nat.Port) string {
			return fmt.Sprintf(
				"user=%s password=%s host=%s port=%s dbname=%s sslmode=disable",
				username, password, host, port.Port(), databaseName,
			)
		}).WithStartupTimeout(2 * time.Minute),
		Env: map[string]string{
			"POSTGRES_DB":       databaseName,
			"POSTGRES_USER":     username,
			"POSTGRES_PASSWORD": password,
		},
		SkipReaper: skipReaper, // if you have problems with reaper rights just skip this
	}

	// Start test container.
	//
	container, err := testcontainers.GenericContainer(ctx, testcontainers.GenericContainerRequest{
		ContainerRequest: req,
		Started:          true,
	})
	if err != nil {
		log.Fatalf("Failed to start container %v", err)
	}

	defer func() { _ = container.Terminate(ctx) }()

	// Setup test dependencies.
	//
	ip, err := container.Host(ctx)
	if err != nil {
		log.Fatalf("Failed to start container host %v", err)
	}
	port, err := container.MappedPort(ctx, nat.Port(exposedPort))
	if err != nil {
		log.Fatalf("Failed to start container port %v", err)
	}
	postgresAddress := fmt.Sprintf("%s:%s", ip, port.Port())
	fmt.Printf("postgresql address: %s\n", postgresAddress)

	postgresDSN := fmt.Sprintf("user=%s password=%s host=%s port=%s dbname=%s sslmode=disable", username, password, ip, port.Port(), databaseName)
	fmt.Printf("postgresql DSN: %s\n", postgresDSN)

	// Initialize database connection.
	//
	testSQLDB, err = postgres.NewConnection(postgresDSN)
	if err != nil {
		log.Fatalf("Cannot establish connection with database %v", err)
	}
	if err := testSQLDB.Ping(); err != nil {
		log.Fatalf("Failed to ping database %v", err)
	}
	defer func() { _ = testSQLDB.Close() }()

	testRunner()
}

func setupDatabase(t testing.TB) error {
	testMigrate(t)

	return nil
}

func testMigrate(t testing.TB) {
	t.Helper()

	require.NoError(t, postgresql.MigrateUp(testSQLDB), "migrations failed")
}

func runTestSetup(t *testing.T, setup func(t testing.TB) error) {
	err := setup(t)
	require.NoError(t, err, "Database setup failed")
}

func purgeDatabase(t testing.TB, db *sql.DB) {
	requireDatabaseStatements(t, db, "DROP SCHEMA public CASCADE;")
	requireDatabaseStatements(t, db, "CREATE SCHEMA public;")
}

func requireDatabaseStatements(t testing.TB, db *sql.DB, queries ...string) {
	t.Helper()

	for _, query := range queries {
		_, err := db.Exec(query)
		require.NoError(t, err, fmt.Sprintf("cannot exec %q", query))
	}
}
