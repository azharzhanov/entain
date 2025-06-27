package account

import (
	"context"
	"fmt"

	errors "entain/internal/error"

	"entain/internal/domain"
)

type service struct {
	repository domain.AccountRepository
}

// NewService creates a new user service with the provided repository.
func NewService(repository domain.AccountRepository) domain.AccountService {
	var service domain.AccountService
	{
		service = newBasicService(repository)
	}
	return service
}

// newBasicService returns a naive, stateless implementation of AccountService.
func newBasicService(
	repository domain.AccountRepository,
) domain.AccountService {
	return &service{
		repository: repository,
	}
}

func (s *service) GetAccount(
	ctx context.Context,
	userID uint64,
) (*domain.AccountRead, error) {
	// Validate inputs
	err := s.validateGetAccount(ctx, userID)
	if err != nil {
		return nil, err
	}

	return s.repository.GetAccount(ctx, userID)
}

func (s *service) validateGetAccount(
	ctx context.Context,
	userID uint64,
) error {
	// Validate inputs
	if userID <= 0 {
		return errors.NewErrInvalidArgument("user id required")
	}

	// Check if user exists
	isExist, err := s.repository.IsAccountExists(ctx, userID)
	if err != nil {
		return err
	}
	if !isExist {
		return errors.NewErrNotFound(
			fmt.Sprintf("user with id %d not found", userID),
		)
	}

	return nil
}
