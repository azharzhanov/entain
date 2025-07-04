// Code generated by sqlc. DO NOT EDIT.
// versions:
//   sqlc v1.23.0

package dbs

import (
	"context"
	"database/sql"
	"fmt"
)

type DBTX interface {
	ExecContext(context.Context, string, ...interface{}) (sql.Result, error)
	PrepareContext(context.Context, string) (*sql.Stmt, error)
	QueryContext(context.Context, string, ...interface{}) (*sql.Rows, error)
	QueryRowContext(context.Context, string, ...interface{}) *sql.Row
}

func New(db DBTX) *Queries {
	return &Queries{db: db}
}

func Prepare(ctx context.Context, db DBTX) (*Queries, error) {
	q := Queries{db: db}
	var err error
	if q.createTransactionStmt, err = db.PrepareContext(ctx, createTransaction); err != nil {
		return nil, fmt.Errorf("error preparing query CreateTransaction: %w", err)
	}
	if q.getAccountStmt, err = db.PrepareContext(ctx, getAccount); err != nil {
		return nil, fmt.Errorf("error preparing query GetAccount: %w", err)
	}
	if q.getBalanceForUpdateStmt, err = db.PrepareContext(ctx, getBalanceForUpdate); err != nil {
		return nil, fmt.Errorf("error preparing query GetBalanceForUpdate: %w", err)
	}
	if q.isAccountExistsStmt, err = db.PrepareContext(ctx, isAccountExists); err != nil {
		return nil, fmt.Errorf("error preparing query IsAccountExists: %w", err)
	}
	if q.isTransactionExistsStmt, err = db.PrepareContext(ctx, isTransactionExists); err != nil {
		return nil, fmt.Errorf("error preparing query IsTransactionExists: %w", err)
	}
	if q.updateBalanceStmt, err = db.PrepareContext(ctx, updateBalance); err != nil {
		return nil, fmt.Errorf("error preparing query UpdateBalance: %w", err)
	}
	return &q, nil
}

func (q *Queries) Close() error {
	var err error
	if q.createTransactionStmt != nil {
		if cerr := q.createTransactionStmt.Close(); cerr != nil {
			err = fmt.Errorf("error closing createTransactionStmt: %w", cerr)
		}
	}
	if q.getAccountStmt != nil {
		if cerr := q.getAccountStmt.Close(); cerr != nil {
			err = fmt.Errorf("error closing getAccountStmt: %w", cerr)
		}
	}
	if q.getBalanceForUpdateStmt != nil {
		if cerr := q.getBalanceForUpdateStmt.Close(); cerr != nil {
			err = fmt.Errorf("error closing getBalanceForUpdateStmt: %w", cerr)
		}
	}
	if q.isAccountExistsStmt != nil {
		if cerr := q.isAccountExistsStmt.Close(); cerr != nil {
			err = fmt.Errorf("error closing isAccountExistsStmt: %w", cerr)
		}
	}
	if q.isTransactionExistsStmt != nil {
		if cerr := q.isTransactionExistsStmt.Close(); cerr != nil {
			err = fmt.Errorf("error closing isTransactionExistsStmt: %w", cerr)
		}
	}
	if q.updateBalanceStmt != nil {
		if cerr := q.updateBalanceStmt.Close(); cerr != nil {
			err = fmt.Errorf("error closing updateBalanceStmt: %w", cerr)
		}
	}
	return err
}

func (q *Queries) exec(ctx context.Context, stmt *sql.Stmt, query string, args ...interface{}) (sql.Result, error) {
	switch {
	case stmt != nil && q.tx != nil:
		return q.tx.StmtContext(ctx, stmt).ExecContext(ctx, args...)
	case stmt != nil:
		return stmt.ExecContext(ctx, args...)
	default:
		return q.db.ExecContext(ctx, query, args...)
	}
}

func (q *Queries) query(ctx context.Context, stmt *sql.Stmt, query string, args ...interface{}) (*sql.Rows, error) {
	switch {
	case stmt != nil && q.tx != nil:
		return q.tx.StmtContext(ctx, stmt).QueryContext(ctx, args...)
	case stmt != nil:
		return stmt.QueryContext(ctx, args...)
	default:
		return q.db.QueryContext(ctx, query, args...)
	}
}

func (q *Queries) queryRow(ctx context.Context, stmt *sql.Stmt, query string, args ...interface{}) *sql.Row {
	switch {
	case stmt != nil && q.tx != nil:
		return q.tx.StmtContext(ctx, stmt).QueryRowContext(ctx, args...)
	case stmt != nil:
		return stmt.QueryRowContext(ctx, args...)
	default:
		return q.db.QueryRowContext(ctx, query, args...)
	}
}

type Queries struct {
	db                      DBTX
	tx                      *sql.Tx
	createTransactionStmt   *sql.Stmt
	getAccountStmt          *sql.Stmt
	getBalanceForUpdateStmt *sql.Stmt
	isAccountExistsStmt     *sql.Stmt
	isTransactionExistsStmt *sql.Stmt
	updateBalanceStmt       *sql.Stmt
}

func (q *Queries) WithTx(tx *sql.Tx) *Queries {
	return &Queries{
		db:                      tx,
		tx:                      tx,
		createTransactionStmt:   q.createTransactionStmt,
		getAccountStmt:          q.getAccountStmt,
		getBalanceForUpdateStmt: q.getBalanceForUpdateStmt,
		isAccountExistsStmt:     q.isAccountExistsStmt,
		isTransactionExistsStmt: q.isTransactionExistsStmt,
		updateBalanceStmt:       q.updateBalanceStmt,
	}
}
