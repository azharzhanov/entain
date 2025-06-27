package postgres

import (
	"context"
	"testing"

	"entain/internal/domain"

	"github.com/stretchr/testify/require"
)

func TestAccountRepository_IsAccountExists(t *testing.T) {
	// Setup database.
	runTestSetup(t, setupDatabase)
	defer purgeDatabase(t, testSQLDB)

	// Initialize database
	//
	sqlcRepository, err := NewSqlcRepository(testSQLDB)
	require.NoError(t, err, "cannot initialize sqlc repository")
	defer sqlcRepository.Close()

	var (
		ctx                   = context.Background()
		testAccountRepository = NewAccountRepository(sqlcRepository)
	)

	// Define test cases.
	tests := []struct {
		name    string
		userID  uint64
		isExist bool
	}{
		{
			name:    "Get user balance with id=1",
			userID:  1,
			isExist: true,
		},
		{
			name:    "Get user balance with id=2",
			userID:  4,
			isExist: false,
		},
	}

	// Run tests.
	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			result, err := testAccountRepository.IsAccountExists(ctx, test.userID)
			require.NoError(t, err, "cannot get user balance")
			require.NotNil(t, result, "result should not be nil")
			require.Equal(t, test.isExist, result, "user existence should match")
		})
	}
}

func TestAccountRepository_GetAccount(t *testing.T) {
	// Setup database.
	runTestSetup(t, setupDatabase)
	defer purgeDatabase(t, testSQLDB)

	// Initialize database
	//
	sqlcRepository, err := NewSqlcRepository(testSQLDB)
	require.NoError(t, err, "cannot initialize sqlc repository")
	defer sqlcRepository.Close()

	var (
		ctx                   = context.Background()
		testAccountRepository = NewAccountRepository(sqlcRepository)
	)

	// Define test cases.
	tests := []struct {
		name     string
		userID   uint64
		expected *domain.AccountRead
	}{
		{
			name:   "Get user balance with id=1",
			userID: 1,
			expected: &domain.AccountRead{
				UserID:  1,
				Balance: "0.00",
			},
		},
		{
			name:   "Get user balance with id=2",
			userID: 2,
			expected: &domain.AccountRead{
				UserID:  2,
				Balance: "0.00",
			},
		},
		{
			name:   "Get user balance with id=3",
			userID: 3,
			expected: &domain.AccountRead{
				UserID:  3,
				Balance: "0.00",
			},
		},
	}

	// Run tests.
	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			result, err := testAccountRepository.GetAccount(ctx, test.userID)
			require.NoError(t, err, "cannot account by user id")
			require.NotNil(t, result, "result should not be nil")
			require.Equal(t, test.expected, result, "account should match expected result")
		})
	}
}

func TestAccountRepository_TxGetBalance(t *testing.T) {
	// Setup database.
	runTestSetup(t, setupDatabase)
	defer purgeDatabase(t, testSQLDB)

	// Initialize database
	//
	sqlcRepository, err := NewSqlcRepository(testSQLDB)
	require.NoError(t, err, "cannot initialize sqlc repository")
	defer sqlcRepository.Close()

	var (
		ctx                   = context.Background()
		testAccountRepository = NewAccountRepository(sqlcRepository)
	)

	// Define test cases.
	tests := []struct {
		name    string
		userID  uint64
		balance float64
	}{
		{
			name:    "Get user balance with id=1",
			userID:  1,
			balance: 0.0,
		},
	}

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			result, err := testAccountRepository.TxGetBalance(
				ctx,
				sqlcRepository.Queries(),
				test.userID,
			)
			require.NoError(t, err, "cannot get user balance")
			require.Equal(t, test.balance, result, "balance should match")
		})
	}
}

func TestAccountRepository_TxUpdateBalance(t *testing.T) {
	// Setup database.
	runTestSetup(t, setupDatabase)
	defer purgeDatabase(t, testSQLDB)

	// Initialize database
	//
	sqlcRepository, err := NewSqlcRepository(testSQLDB)
	require.NoError(t, err, "cannot initialize sqlc repository")
	defer sqlcRepository.Close()

	var (
		ctx                   = context.Background()
		testAccountRepository = NewAccountRepository(sqlcRepository)
	)

	// Define test cases.
	tests := []struct {
		name     string
		userID   uint64
		balance  float64
		expected *domain.AccountRead
	}{
		{
			name:    "Update user balance with id=1",
			userID:  1,
			balance: 100.0,
			expected: &domain.AccountRead{
				UserID:  1,
				Balance: "100.00",
			},
		},
	}

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			err := testAccountRepository.TxUpdateBalance(
				ctx,
				sqlcRepository.Queries(),
				test.userID,
				test.balance,
			)
			require.NoError(t, err, "cannot update user balance")

			account, err := testAccountRepository.GetAccount(
				ctx,
				test.userID,
			)
			require.NoError(t, err, "cannot get user account")
			require.Equal(t, test.expected, account, "account should match expected result")
		})
	}
}
