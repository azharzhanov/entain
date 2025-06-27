package domain

import (
	"context"

	"entain/schema/postgresql/dbs"
)

// AccountRead represents the account information of a user.
type AccountRead struct {
	// UserID is the unique identifier of the user.
	//
	UserID uint64 `json:"userId"`

	// Balance is the current balance of the user.
	//
	Balance string `json:"balance"`
}

// AccountRepository provides access to a storage.
type AccountRepository interface {
	IsAccountExists(
		ctx context.Context,
		userID uint64,
	) (bool, error)

	GetAccount(
		ctx context.Context,
		userID uint64,
	) (*AccountRead, error)

	TxGetBalance(
		ctx context.Context,
		queries *dbs.Queries,
		userID uint64,
	) (float64, error)

	TxUpdateBalance(
		ctx context.Context,
		queries *dbs.Queries,
		userID uint64,
		balance float64,
	) error
}

// AccountService provides access to a business logic.
type AccountService interface {
	// GetAccount retrieves the account information of a user by their ID.
	GetAccount(
		ctx context.Context,
		userID uint64,
	) (*AccountRead, error)
}
