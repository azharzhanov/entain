package domain

import (
	"context"

	"entain/schema/postgresql/dbs"
)

// TransactionWrite represents read transaction model.
type TransactionWrite struct {
	// State indicates the state of the transaction.
	// Possible values for `state` field are: (`win`, `lose`)
	//
	// `win` - request must increase the user balance
	// `lose` - request must decrease user balance
	State string `json:"state"`

	// Amount indicates the amount of the transaction.
	//
	Amount string `json:"amount"`

	// TransactionID is the unique identifier for the transaction.
	//
	TransactionID string `json:"transactionId"`
}

// TransactionRepository provides access to a storage.
type TransactionRepository interface {
	// TxCreateTransaction creates a transaction for a user.
	TxCreateTransaction(
		ctx context.Context,
		queries *dbs.Queries,
		userID uint64,
		sourceType string,
		transaction *TransactionWrite,
	) error

	// IsTransactionExists checks if a transaction exists for a user.
	IsTransactionExists(
		ctx context.Context,
		transactionID string,
	) (bool, error)
}

// TransactionService provides access to a business logic.
type TransactionService interface {
	// AggregateTransaction aggregates a transaction for a user.
	AggregateTransaction(
		ctx context.Context,
		userID uint64,
		sourceType string,
		transaction *TransactionWrite,
	) error
}
