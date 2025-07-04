// Code generated by sqlc. DO NOT EDIT.
// versions:
//   sqlc v1.23.0
// source: transactions.sql

package dbs

import (
	"context"
)

const createTransaction = `-- name: CreateTransaction :exec
INSERT INTO transactions (transaction_id, user_id, source_type, state, amount, created_at)
VALUES ($1, $2, $3, $4, $5, NOW())
`

type CreateTransactionParams struct {
	TransactionID string
	UserID        int64
	SourceType    SourceType
	State         TransactionState
	Amount        string
}

func (q *Queries) CreateTransaction(ctx context.Context, arg CreateTransactionParams) error {
	_, err := q.exec(ctx, q.createTransactionStmt, createTransaction,
		arg.TransactionID,
		arg.UserID,
		arg.SourceType,
		arg.State,
		arg.Amount,
	)
	return err
}

const isTransactionExists = `-- name: IsTransactionExists :one
SELECT EXISTS (
    SELECT
    FROM transactions
    WHERE transaction_id = $1
) AS exists
`

func (q *Queries) IsTransactionExists(ctx context.Context, transactionID string) (bool, error) {
	row := q.queryRow(ctx, q.isTransactionExistsStmt, isTransactionExists, transactionID)
	var exists bool
	err := row.Scan(&exists)
	return exists, err
}
