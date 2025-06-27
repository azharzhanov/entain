package postgres

import (
	"database/sql"
)

// NewConnection - creates a new database connection
func NewConnection(dsn string) (*sql.DB, error) {
	db, err := sql.Open("postgres", dsn)
	if err != nil {
		return nil, err
	}

	// Ping the database to ensure the connection is established
	if err = db.Ping(); err != nil {
		return nil, err
	}

	return db, nil
}
