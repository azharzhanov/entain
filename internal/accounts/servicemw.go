package accounts

import (
	"context"

	"entain/internal/domain"

	"github.com/go-kit/log"
)

// Middleware describes a service (as opposed to endpoint) middleware.
type middleware func(service domain.AccountService) domain.AccountService

// LoggingServiceMiddleware takes a logger as a dependency
// and returns a service Middleware.
func loggingServiceMiddleware(logger log.Logger) middleware {
	return func(next domain.AccountService) domain.AccountService {
		return loggingMiddleware{logger, next}
	}
}

type loggingMiddleware struct {
	logger log.Logger
	next   domain.AccountService
}

func (mw loggingMiddleware) GetAccount(
	ctx context.Context,
	userID uint64,
) (result *domain.AccountRead, err error) {
	defer func() {
		_ = mw.logger.Log("method", "GetAccount",
			"userID", userID,
			"result", result,
			"err", err,
		)
	}()
	return mw.next.GetAccount(ctx, userID)
}
