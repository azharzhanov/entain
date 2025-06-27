package integrationtest

import (
	"os"
)

var (
	integrationTestBaseURL = os.Getenv("BASE_URL") // Default base URL for integration tests
)

// Endpoints
const (
	getAccountPath        = "/user/%d/balance"
	createTransactionPath = "/user/%d/transaction"
)
