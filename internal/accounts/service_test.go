package accounts

import (
	"context"
	"testing"

	"entain/internal/domain"
	"entain/internal/mocks"

	"github.com/stretchr/testify/require"
	"go.uber.org/mock/gomock"
)

func TestService_GetAccount(t *testing.T) {
	var (
		ctx     = context.Background()
		balance = &domain.AccountRead{
			UserID:  1,
			Balance: "10.0",
		}
	)

	// Setup mocks
	stubCtrl := gomock.NewController(t)
	defer stubCtrl.Finish()

	// Mock repository
	var (
		accountRepository = mocks.NewMockAccountRepository(stubCtrl)
	)

	// Mock expectations
	accountRepository.EXPECT().
		IsAccountExists(ctx, uint64(1)).
		Return(true, nil).
		Times(1)

	accountRepository.EXPECT().
		GetAccount(ctx, uint64(1)).
		Return(balance, nil).
		Times(1)

	// Define tests
	type arguments struct {
		userID uint64
	}

	type result struct {
		balance *domain.AccountRead
	}

	tests := []struct {
		name        string
		arguments   arguments
		expected    result
		expectError bool
	}{
		{
			name: "Success: valid user id",
			arguments: arguments{
				userID: 1,
			},
			expected: result{
				balance: balance,
			},
			expectError: false,
		},
		{
			name: "Fail: invalid user id",
			arguments: arguments{
				userID: 0,
			},
			expected: result{
				balance: nil,
			},
			expectError: true,
		},
	}

	// Define service and run tests
	service := newBasicService(
		accountRepository,
	)
	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			args := test.arguments
			expected := test.expected

			result, err := service.GetAccount(
				ctx,
				args.userID,
			)
			if !test.expectError {
				require.NoError(t, err)
				require.Equal(t, expected.balance, result)
			} else {
				require.Error(t, err, "expected an error but got nil")
			}
		})
	}
}
