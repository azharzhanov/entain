package postgres

import (
	"context"
	"database/sql"
	"fmt"

	"entain/internal/domain"
	errors "entain/internal/error"
	"entain/schema/postgresql/dbs"
)

type accountRepository struct {
	db *SqlcRepository
}

// NewAccountRepository creates a new instance of account repository.
func NewAccountRepository(db *SqlcRepository) domain.AccountRepository {
	return &accountRepository{
		db: db,
	}
}

func (r *accountRepository) IsAccountExists(
	ctx context.Context,
	userID uint64,
) (bool, error) {
	return r.db.Queries().IsAccountExists(ctx, int64(userID))
}

func (r *accountRepository) GetAccount(
	ctx context.Context,
	userID uint64,
) (*domain.AccountRead, error) {
	row, err := r.db.Queries().GetAccount(ctx, int64(userID))
	if err != nil {
		if err == sql.ErrNoRows {
			return nil, errors.NewErrNotFound(
				fmt.Sprintf("user with id=%d not found", userID),
			)
		}
		return nil, err
	}

	return &domain.AccountRead{
		UserID:  uint64(row.ID),
		Balance: row.Balance,
	}, nil
}

func (r *accountRepository) TxGetBalance(
	ctx context.Context,
	queries *dbs.Queries,
	userID uint64,
) (float64, error) {
	return queries.GetBalanceForUpdate(
		ctx,
		int64(userID),
	)
}

func (r *accountRepository) TxUpdateBalance(
	ctx context.Context,
	queries *dbs.Queries,
	userID uint64,
	balance float64,
) error {
	return queries.UpdateBalance(
		ctx,
		dbs.UpdateBalanceParams{
			NewBalance: balance,
			ID:         int64(userID),
		},
	)
}
