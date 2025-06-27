package transactions

import (
	"context"
	"testing"

	"entain/internal/domain"
	"entain/internal/mocks"
	"entain/schema/postgresql/dbs"

	"github.com/stretchr/testify/require"
	"go.uber.org/mock/gomock"
)

func TestService_AggregateTransaction(t *testing.T) {
	var (
		ctx        = context.Background()
		dummyQuery = &dbs.Queries{}
	)

	// Setup mocks
	stubCtrl := gomock.NewController(t)
	defer stubCtrl.Finish()

	// Mock repository
	var (
		transactionRepository = mocks.NewMockTransactionRepository(stubCtrl)
		accountRepository     = mocks.NewMockAccountRepository(stubCtrl)
		txRepository          = mocks.NewMockTxRepository(stubCtrl)
	)

	// Define service
	svc := service{
		transactionRepository: transactionRepository,
		accountRepository:     accountRepository,
		txRepository:          txRepository,
	}

	t.Run("check_inputs", func(t *testing.T) {
		var (
			userID      = uint64(1)
			transaction = &domain.TransactionWrite{
				State:         "win",
				Amount:        "10.0",
				TransactionID: "transaction-id-123",
			}
		)

		// Mock expectations
		accountRepository.EXPECT().
			IsAccountExists(
				ctx,
				userID,
			).Return(true, nil).
			Times(1)

		transactionRepository.EXPECT().
			IsTransactionExists(
				ctx,
				transaction.TransactionID,
			).Return(false, nil).
			Times(1)

		// Test: Fail: user id required
		err := svc.AggregateTransaction(
			ctx,
			0,
			"test",
			transaction,
		)
		require.Error(t, err)

		// Test: Fail: source type required
		err = svc.AggregateTransaction(
			ctx,
			userID,
			"",
			transaction,
		)
		require.Error(t, err)

		// Test: Fail: transaction required
		err = svc.AggregateTransaction(
			ctx,
			userID,
			"test",
			nil,
		)
		require.Error(t, err)

		// Test: Fail: transaction state required
		err = svc.AggregateTransaction(
			ctx,
			userID,
			"test",
			&domain.TransactionWrite{
				State:         "",
				Amount:        "10.0",
				TransactionID: "transaction-id-123",
			},
		)
		require.Error(t, err)

		// Test: Fail: transaction amount required
		err = svc.AggregateTransaction(
			ctx,
			userID,
			"test",
			&domain.TransactionWrite{
				State:         "win",
				Amount:        "",
				TransactionID: "transaction-id-123",
			},
		)
		require.Error(t, err)

		// Test: Fail: transaction amount must be a positive number
		err = svc.AggregateTransaction(
			ctx,
			userID,
			"test",
			&domain.TransactionWrite{
				State:         "lose",
				Amount:        "-10.0",
				TransactionID: "transaction-id-123",
			},
		)
		require.Error(t, err)

		// Test: Fail: transaction id required
		err = svc.AggregateTransaction(
			ctx,
			userID,
			"test",
			&domain.TransactionWrite{
				State:         "win",
				Amount:        "10.0",
				TransactionID: "",
			},
		)
		require.Error(t, err)
	})

	t.Run("success_aggregate_transaction", func(t *testing.T) {
		var (
			userID         = uint64(1)
			sourceType     = "payment"
			currentBalance = 100.0
			newBalance     = 110.0 // Assuming the transaction is a win and adds 10.0 to the balance
			transaction    = &domain.TransactionWrite{
				State:         "win",
				Amount:        "10.0",
				TransactionID: "transaction-id-123",
			}
		)

		// Mock get balance
		accountRepository.EXPECT().
			TxGetBalance(
				ctx,
				dummyQuery,
				userID,
			).Return(currentBalance, nil)

		// Mock update balance
		accountRepository.EXPECT().
			TxUpdateBalance(
				ctx,
				dummyQuery,
				userID,
				newBalance,
			).Return(nil)

		// Mock create transaction
		transactionRepository.EXPECT().
			TxCreateTransaction(
				ctx,
				dummyQuery,
				userID,
				sourceType,
				transaction,
			).Return(nil)

		// Call the service method
		err := svc.processTransaction(
			ctx,
			dummyQuery,
			userID,
			sourceType,
			transaction,
		)
		require.NoError(t, err)
	})
}

func TestCalculateNewBalance(t *testing.T) {
	var (
		currentBalance = 100.0
		transaction    = &domain.TransactionWrite{
			State:         "win",
			Amount:        "10.0",
			TransactionID: "transaction-id-123",
		}
	)

	t.Run("calculate_win_transaction", func(t *testing.T) {
		newBalance, err := calculateNewBalance(currentBalance, transaction)
		require.NoError(t, err)
		require.Equal(t, 110.0, newBalance)
	})

	t.Run("calculate_lose_transaction", func(t *testing.T) {
		transaction.State = "lose"
		newBalance, err := calculateNewBalance(currentBalance, transaction)
		require.NoError(t, err)
		require.Equal(t, 90.0, newBalance)
	})
}
