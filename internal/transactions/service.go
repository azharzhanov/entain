package transactions

import (
	"context"
	"fmt"
	"strconv"

	errors "entain/internal/error"
	"entain/schema/postgresql/dbs"

	"entain/internal/domain"
)

type service struct {
	transactionRepository domain.TransactionRepository
	accountRepository     domain.AccountRepository
	txRepository          domain.TxRepository
}

// NewService creates a new transactions service.
func NewService(
	transactionRepository domain.TransactionRepository,
	accountRepository domain.AccountRepository,
	txRepository domain.TxRepository,
) domain.TransactionService {
	var service domain.TransactionService
	{
		service = newBasicService(
			transactionRepository,
			accountRepository,
			txRepository,
		)
	}
	return service
}

// newBasicService returns a naive, stateless implementation of TransactionService.
func newBasicService(
	transactionRepository domain.TransactionRepository,
	accountRepository domain.AccountRepository,
	txRepository domain.TxRepository,
) domain.TransactionService {
	return &service{
		transactionRepository: transactionRepository,
		accountRepository:     accountRepository,
		txRepository:          txRepository,
	}
}

func (s *service) AggregateTransaction(
	ctx context.Context,
	userID uint64,
	sourceType string,
	transaction *domain.TransactionWrite,
) error {
	// Validate inputs
	err := s.validateAggregateTransaction(
		ctx,
		userID,
		sourceType,
		transaction,
	)
	if err != nil {
		return err
	}

	// Start the transaction
	txErr := s.txRepository.WithTransaction(ctx, func(queries *dbs.Queries) error {
		return s.processTransaction(
			ctx,
			queries,
			userID,
			sourceType,
			transaction,
		)
	})
	return txErr
}

func (s *service) validateAggregateTransaction(
	ctx context.Context,
	userID uint64,
	sourceType string,
	transaction *domain.TransactionWrite,
) error {
	if userID <= 0 {
		return errors.NewErrInvalidArgument("user id required")
	}
	if sourceType == "" {
		return errors.NewErrInvalidArgument("source type required")
	}
	if transaction == nil {
		return errors.NewErrInvalidArgument("transaction required")
	}
	if transaction.State == "" {
		return errors.NewErrInvalidArgument("transaction state required")
	}
	if transaction.Amount == "" {
		return errors.NewErrInvalidArgument("transaction amount required")
	}
	if transaction.TransactionID == "" {
		return errors.NewErrInvalidArgument("transaction id required")
	}

	// Check if the transaction state is valid
	if transaction.State != "win" && transaction.State != "lose" {
		return errors.NewErrInvalidArgument("transaction state must be 'win' or 'lose'")
	}

	// Check if the user exists
	isUserExists, err := s.accountRepository.IsAccountExists(ctx, userID)
	if err != nil {
		return err
	}

	// If the user does not exist, return an error
	if !isUserExists {
		return errors.NewErrNotFound(
			fmt.Sprintf("user with id %d does not exist", userID),
		)
	}

	// Check if the transaction already exists
	isTransactionExists, err := s.transactionRepository.IsTransactionExists(
		ctx,
		transaction.TransactionID,
	)
	if err != nil {
		return err
	}

	// If the transaction already exists, return an error
	if isTransactionExists {
		return errors.NewErrAlreadyExist(
			fmt.Sprintf("transaction with id %s already exists", transaction.TransactionID),
		)
	}

	// Check if the transaction amount is a valid float
	amountFloat, err := strconv.ParseFloat(transaction.Amount, 64)
	if err != nil {
		return errors.NewErrInvalidArgument(
			fmt.Sprintf("invalid transaction amount: %s", transaction.Amount),
		)
	}
	if amountFloat <= 0 {
		return errors.NewErrInvalidArgument("transaction amount must be greater than zero")
	}

	return nil
}

func (s *service) processTransaction(
	ctx context.Context,
	queries *dbs.Queries,
	userID uint64,
	sourceType string,
	transaction *domain.TransactionWrite,
) error {
	// 1. Select balance for the user using FOR UPDATE clause
	currentBalance, err := s.accountRepository.TxGetBalance(
		ctx,
		queries,
		userID,
	)
	if err != nil {
		return err
	}

	// 2. Calculate the new balance based on the transaction state
	newBalance, err := calculateNewBalance(
		currentBalance,
		transaction,
	)
	if err != nil {
		return err
	}

	// 3. Update the user's balance
	err = s.accountRepository.TxUpdateBalance(
		ctx,
		queries,
		userID,
		newBalance,
	)
	if err != nil {
		return err
	}

	// 4. Insert the transaction into the database
	err = s.transactionRepository.TxCreateTransaction(
		ctx,
		queries,
		userID,
		sourceType,
		transaction,
	)
	if err != nil {
		return err
	}

	return nil
}

func calculateNewBalance(
	currentBalance float64,
	transaction *domain.TransactionWrite,
) (float64, error) {
	// Parse the transaction amount from string to float64
	amountFloat, err := strconv.ParseFloat(transaction.Amount, 64)
	if err != nil {
		return 0, errors.NewErrInvalidArgument(
			fmt.Sprintf("invalid transaction amount: %s", transaction.Amount),
		)
	}

	// If the transaction state is "win", increase the balance
	if transaction.State == "win" {
		return currentBalance + amountFloat, nil
	}

	// If the transaction state is "lose", ensure the balance does not go negative
	if currentBalance < amountFloat {
		return 0, errors.NewErrInvalidArgument("insufficient balance for transaction")
	}

	// Decrease balance
	return currentBalance - amountFloat, nil
}
