package integrationtest

var (
	integrationTestBaseURL = "http://localhost:8080" // Default base URL for integration tests
)

// Endpoints
const (
	getAccountPath        = "/user/%d/balance"
	createTransactionPath = "/user/%d/transaction"
)
