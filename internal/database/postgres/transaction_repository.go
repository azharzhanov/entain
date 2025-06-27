package postgres

import (
	"context"

	"entain/internal/domain"
	"entain/schema/postgresql/dbs"
)

type transactionRepository struct {
	db *SqlcRepository
}

// NewTransactionRepository creates a new repository for transactions.
func NewTransactionRepository(db *SqlcRepository) domain.TransactionRepository {
	return &transactionRepository{
		db: db,
	}
}

func (r *transactionRepository) TxCreateTransaction(
	ctx context.Context,
	queries *dbs.Queries,
	userID uint64,
	sourceType string,
	transaction *domain.TransactionWrite,
) error {
	return queries.CreateTransaction(
		ctx,
		dbs.CreateTransactionParams{
			TransactionID: transaction.TransactionID,
			UserID:        int64(userID),
			SourceType:    dbs.SourceType(sourceType),
			State:         dbs.TransactionState(transaction.State),
			Amount:        transaction.Amount,
		},
	)
}

func (r *transactionRepository) IsTransactionExists(
	ctx context.Context,
	transactionID string,
) (bool, error) {
	return r.db.Queries().IsTransactionExists(
		ctx,
		transactionID,
	)
}
