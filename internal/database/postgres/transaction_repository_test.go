package postgres

import (
	"context"
	"testing"

	"entain/internal/domain"

	"github.com/stretchr/testify/require"
)

func TestTransactionRepository_TxCreateTransactionTr(t *testing.T) {
	// Setup database.
	runTestSetup(t, setupDatabase)
	defer purgeDatabase(t, testSQLDB)

	// Initialize database
	//
	sqlcRepository, err := NewSqlcRepository(testSQLDB)
	require.NoError(t, err, "cannot initialize sqlc repository")
	defer sqlcRepository.Close()

	var (
		ctx                       = context.Background()
		testTransactionRepository = NewTransactionRepository(sqlcRepository)
	)

	// Define test cases.
	tests := []struct {
		name        string
		userID      uint64
		sourceType  string
		transaction *domain.TransactionWrite
	}{
		{
			name:       "Create transaction for user with id=1",
			userID:     1,
			sourceType: "payment",
			transaction: &domain.TransactionWrite{
				State:         "win",
				Amount:        "10.0",
				TransactionID: "transaction-id-1",
			},
		},
		{
			name:       "Create transaction for user with id=2",
			userID:     2,
			sourceType: "payment",
			transaction: &domain.TransactionWrite{
				State:         "lose",
				Amount:        "10.0",
				TransactionID: "transaction-id-2",
			},
		},
	}

	// Run test cases.
	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			err := testTransactionRepository.TxCreateTransaction(
				ctx,
				sqlcRepository.Queries(),
				test.userID,
				test.sourceType,
				test.transaction,
			)
			require.NoError(t, err)
		})
	}
}

func TestTransactionRepository_IsTransactionExists(t *testing.T) {
	// Setup database.
	runTestSetup(t, setupDatabase)
	defer purgeDatabase(t, testSQLDB)

	// Initialize database
	//
	sqlcRepository, err := NewSqlcRepository(testSQLDB)
	require.NoError(t, err, "cannot initialize sqlc repository")
	defer sqlcRepository.Close()

	var (
		ctx                       = context.Background()
		testTransactionRepository = NewTransactionRepository(sqlcRepository)
	)

	// Define prepare tests
	prepareTests := []struct {
		name        string
		userID      uint64
		sourceType  string
		transaction *domain.TransactionWrite
	}{
		{
			name:       "Create transaction for user with id=1",
			userID:     1,
			sourceType: "payment",
			transaction: &domain.TransactionWrite{
				State:         "win",
				Amount:        "10.0",
				TransactionID: "transaction-id-1",
			},
		},
		{
			name:       "Create transaction for user with id=2",
			userID:     2,
			sourceType: "payment",
			transaction: &domain.TransactionWrite{
				State:         "lose",
				Amount:        "10.0",
				TransactionID: "transaction-id-2",
			},
		},
	}

	// Prepare transactions
	for _, prepareTest := range prepareTests {
		err := testTransactionRepository.TxCreateTransaction(
			ctx,
			sqlcRepository.Queries(),
			prepareTest.userID,
			prepareTest.sourceType,
			prepareTest.transaction,
		)
		require.NoError(t, err)
	}

	// Define tests
	tests := []struct {
		name          string
		transactionID string
		expected      bool
	}{
		{
			name:          "Check if transaction with id=transaction-id-1 exists",
			transactionID: "transaction-id-1",
			expected:      true,
		},
		{
			name:          "Check if transaction with id=transaction-id-2 exists",
			transactionID: "transaction-id-2",
			expected:      true,
		},
		{
			name:          "Check if transaction with id=transaction-id-3 exists",
			transactionID: "transaction-id-3",
			expected:      false,
		},
	}

	// Run tests
	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			exists, err := testTransactionRepository.IsTransactionExists(
				ctx,
				test.transactionID,
			)
			require.NoError(t, err, "cannot check if transaction exists")
			require.Equal(t, test.expected, exists, "transaction existence mismatch")
		})
	}
}
