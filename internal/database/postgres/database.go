package postgres

import (
	"context"
	"database/sql"

	"entain/schema/postgresql/dbs"

	// Import the Postgres driver
	_ "github.com/lib/pq"
)

// SqlcRepository represents a PostgresSQL database connection and its associated queries.
type SqlcRepository struct {
	connection *sql.DB
	queries    *dbs.Queries
}

// NewSqlcRepository creates a new instance with the provided SQL connection.
func NewSqlcRepository(connection *sql.DB) (*SqlcRepository, error) {
	var queries, err = dbs.Prepare(context.Background(), connection)
	if err != nil {
		return nil, err
	}

	return &SqlcRepository{
		connection: connection,
		queries:    queries,
	}, nil
}

// Queries - returns the queries associated with the database.
func (r *SqlcRepository) Queries() *dbs.Queries {
	return r.queries
}

// Close - closes the database connection and releases any resources associated with it.
func (r *SqlcRepository) Close() error {
	return r.queries.Close()
}

// WithTransaction starts a transaction and executes the provided function within that transaction context.
func (r *SqlcRepository) WithTransaction(ctx context.Context, fn func(queries *dbs.Queries) error) error {
	return withTransaction(ctx, r.connection, r.queries, fn)
}

func withTransaction(ctx context.Context, db *sql.DB, queries *dbs.Queries, fn func(queries *dbs.Queries) error) (err error) {
	tx, err := db.BeginTx(ctx, nil)
	if err != nil {
		return
	}

	defer func() {
		if p := recover(); p != nil {
			// a panic occurred, rollback and repanic
			tx.Rollback()

			panic(p)
		} else if err != nil {
			// something went wrong, rollback
			tx.Rollback()
		} else {
			// all good, commit
			err = tx.Commit()
		}
	}()

	err = fn(queries.WithTx(tx))

	return err
}
