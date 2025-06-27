package domain

import (
	"context"

	"entain/schema/postgresql/dbs"
)

// TxRepository provides access to a storage with transaction support.
type TxRepository interface {
	// WithTransaction executes a function within a transaction context.
	WithTransaction(
		ctx context.Context,
		fn func(queries *dbs.Queries) error,
	) error
}
