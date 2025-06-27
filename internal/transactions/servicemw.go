package transactions

import (
	"context"

	"entain/internal/domain"

	"github.com/go-kit/log"
)

// Middleware describes a service (as opposed to endpoint) middleware.
type middleware func(service domain.TransactionService) domain.TransactionService

// LoggingServiceMiddleware takes a logger as a dependency
// and returns a service Middleware.
func loggingServiceMiddleware(logger log.Logger) middleware {
	return func(next domain.TransactionService) domain.TransactionService {
		return loggingMiddleware{logger, next}
	}
}

type loggingMiddleware struct {
	logger log.Logger
	next   domain.TransactionService
}

func (mw loggingMiddleware) AggregateTransaction(
	ctx context.Context,
	userID uint64,
	sourceType string,
	transaction *domain.TransactionWrite,
) (err error) {
	defer func() {
		_ = mw.logger.Log("method", "AggregateTransaction",
			"userID", userID,
			"sourceType", sourceType,
			"transaction", transaction,
			"err", err,
		)
	}()
	return mw.next.AggregateTransaction(ctx, userID, sourceType, transaction)
}
