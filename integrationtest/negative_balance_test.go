package integrationtest

import (
	"fmt"
	"testing"

	client "entain/integrationtest/client"

	"github.com/google/uuid"
	"github.com/stretchr/testify/require"
)

func TestNegativeBalance(t *testing.T) {
	var (
		userID         = int64(3)
		transactionID  = uuid.NewString()
		requestPayload = map[string]string{
			"state":         "win",
			"amount":        "-10.00",
			"transactionId": transactionID,
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

	// Do the POST request to create a transaction
	response, err := cl.Do(client.NewPostRequest(
		fmt.Sprintf(createTransactionPath, userID),
		requestPayload,
	))
	require.NoError(t, err)
	require.Equal(t, response.StatusCode(), 400, "Expected a 400 Bad Request for negative balance transaction")
}
