package transactions

import (
	"context"

	"entain/internal/domain"

	"github.com/go-kit/kit/endpoint"
)

// Endpoints collects all the endpoints that compose an add service. It's meant to
// be used as a helper struct, to collect all the endpoints into a single
// parameter.
type Endpoints struct {
	AggregateTransactionEndpoint endpoint.Endpoint
}

// NewEndpoints creates a new Endpoints struct with the provided service.
func NewEndpoints(service domain.TransactionService) Endpoints {
	return Endpoints{
		AggregateTransactionEndpoint: MakeAggregateTransactionEndpoint(service),
	}
}

// MakeAggregateTransactionEndpoint Impl.
func MakeAggregateTransactionEndpoint(service domain.TransactionService) endpoint.Endpoint {
	return func(ctx context.Context, request interface{}) (response interface{}, err error) {
		req := request.(AggregateTransactionRequest)

		// Call the service
		err = service.AggregateTransaction(
			ctx,
			req.UserID,
			req.SourceType,
			req.Transaction,
		)
		if err != nil {
			return nil, err
		}

		return AggregateTransactionResponse{
			Body: AggregateTransactionResponseBody{
				Status: domain.ResponseStatusSuccess,
			},
		}, nil
	}
}

// swagger:parameters AggregateTransactionRequest
type AggregateTransactionRequest struct {
	// UserID is the ID of the user for whom the transaction is being created.
	//
	// in:path
	// required: true
	UserID uint64 `json:"user_id"`

	// SourceType indicates the source of the transaction, such as "game", "server", "payment".
	//
	// in:header
	// required: true
	SourceType string `json:"source_type"`

	// Transaction contains the details of the transaction to be created.
	//
	// in:body
	// required: true
	Transaction *domain.TransactionWrite `json:"body"`
}

// swagger:response AggregateTransactionResponse
type AggregateTransactionResponse struct {
	// in:body
	Body AggregateTransactionResponseBody `json:"body"`
}

// swagger:parameters AggregateTransactionResponseBody
type AggregateTransactionResponseBody struct {
	Status domain.ResponseStatus `json:"status"`
}

var (
	_ domain.Bodier = AggregateTransactionResponse{}
)

// GetBody returns the body of the response.
func (r AggregateTransactionResponse) GetBody() interface{} { return r.Body }
