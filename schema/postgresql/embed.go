package postgresql

import (
	"database/sql"
	"embed"

	"github.com/pressly/goose/v3"
)

//go:embed migrations/*.sql
var migrations embed.FS

// MigrateUp - up
// migrations only for run in active-gateway docker-compose local development
// in staging and live migrations up inside GitHub Actions
func MigrateUp(db *sql.DB) error {
	// goose has a lot of dependencies with ClickHouse and other DB drivers
	goose.SetBaseFS(migrations)
	goose.SetTableName("schema_migrations")
	// PostgreSQL by default
	// goose.SetDialect("postgres")

	return goose.Up(db, "migrations")
}
