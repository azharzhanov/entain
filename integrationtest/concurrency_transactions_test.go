package integrationtest

import (
	"fmt"
	"net/url"
	"testing"

	client "entain/integrationtest/client"
	"entain/internal/account"
	"entain/internal/domain"

	"github.com/google/uuid"
	"github.com/stretchr/testify/require"
	"golang.org/x/sync/errgroup"
)

func TestConcurrencyWinTransactions(t *testing.T) {
	var (
		userID   = int64(1)
		expected = &domain.AccountRead{
			UserID:  1,
			Balance: "30.00",
		}
	)

	// Initialize the client with the base URL and request header
	cl := client.NewClient(
		integrationTestBaseURL,
		&client.RequestHeader{
			SourceType:  "game",
			ContentType: "application/json",
		},
	)

	var (
		eg errgroup.Group
	)

	// Concurrent POST requests to create transactions
	for i := 0; i < 30; i++ {
		eg.Go(func() error {
			// Generate a unique transaction ID for each request
			transactionID := uuid.NewString()

			// Do a POST request to create a transaction
			response, err := cl.Do(client.NewPostRequest(
				fmt.Sprintf(createTransactionPath, userID),
				domain.TransactionWrite{
					State:         "win",
					Amount:        "1.00",
					TransactionID: transactionID,
				},
			))
			if err != nil {
				return err
			}

			// Check if the response status code is 200 OK
			if response.StatusCode() != 200 {
				return fmt.Errorf("expected a 200 OK for transaction, got %d", response.StatusCode())
			}

			return nil
		})
	}

	// Wait for all goroutines to finish
	err := eg.Wait()
	require.NoError(t, err, "Expected no errors during concurrent transaction processing")

	// After all transactions are processed, retrieve the account to check balance
	response, err := cl.Do(client.NewGetRequest(
		fmt.Sprintf(getAccountPath, userID),
		url.Values{},
	))
	require.NoError(t, err)
	require.Equal(t, 200, response.StatusCode(), "Expected a 200 OK for account retrieval after transactions")

	result := account.GetAccountResponseBody{}
	require.NoError(t, response.Read(&result))
	require.Equal(t, expected, result.Data, "Expected account balance to be updated correctly after concurrent transactions")
}

func TestConcurrencyLoseTransactions(t *testing.T) {
	var (
		userID   = int64(1)
		expected = &domain.AccountRead{
			UserID:  1,
			Balance: "0.00",
		}
	)

	// Initialize the client with the base URL and request header
	cl := client.NewClient(
		integrationTestBaseURL,
		&client.RequestHeader{
			SourceType:  "game",
			ContentType: "application/json",
		},
	)

	var (
		eg errgroup.Group
	)

	// Concurrent POST requests to create transactions
	for i := 0; i < 30; i++ {
		eg.Go(func() error {
			// Generate a unique transaction ID for each request
			transactionID := uuid.NewString()

			// Do a POST request to create a transaction
			response, err := cl.Do(client.NewPostRequest(
				fmt.Sprintf(createTransactionPath, userID),
				domain.TransactionWrite{
					State:         "lose",
					Amount:        "1.00",
					TransactionID: transactionID,
				},
			))
			if err != nil {
				return err
			}

			// Check if the response status code is 200 OK
			if response.StatusCode() != 200 {
				return fmt.Errorf("expected a 200 OK for transaction, got %d", response.StatusCode())
			}

			return nil
		})
	}

	// Wait for all goroutines to finish
	err := eg.Wait()
	require.NoError(t, err, "Expected no errors during concurrent transaction processing")

	// After all transactions are processed, retrieve the account to check balance
	response, err := cl.Do(client.NewGetRequest(
		fmt.Sprintf(getAccountPath, userID),
		url.Values{},
	))
	require.NoError(t, err)
	require.Equal(t, 200, response.StatusCode(), "Expected a 200 OK for account retrieval after transactions")

	result := account.GetAccountResponseBody{}
	require.NoError(t, response.Read(&result))
	require.Equal(t, expected, result.Data, "Expected account balance to be updated correctly after concurrent transactions")
}
