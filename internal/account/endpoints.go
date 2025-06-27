package account

import (
	"context"

	"entain/internal/domain"

	"github.com/go-kit/kit/endpoint"
)

// Endpoints collects all the endpoints that compose an add service. It's meant to
// be used as a helper struct, to collect all the endpoints into a single
// parameter.
type Endpoints struct {
	GetAccountEndpoint endpoint.Endpoint
}

// NewEndpoints creates a new Endpoints struct with the provided service.
func NewEndpoints(service domain.AccountService) Endpoints {
	return Endpoints{
		GetAccountEndpoint: MakeGetAccountEndpoint(service),
	}
}

// MakeGetAccountEndpoint creates an endpoint for getting a user's balance.
func MakeGetAccountEndpoint(service domain.AccountService) endpoint.Endpoint {
	return func(ctx context.Context, request interface{}) (response interface{}, err error) {
		req := request.(GetAccountRequest)

		// Call the service
		resp, err := service.GetAccount(
			ctx,
			req.UserID,
		)
		if err != nil {
			return nil, err
		}

		return GetAccountResponse{
			Body: GetAccountResponseBody{
				Status: domain.ResponseStatusSuccess,
				Data:   resp,
			},
		}, nil
	}
}

// swagger:parameters GetAccountRequest
type GetAccountRequest struct {
	// The ID of the user whose balance is being requested.
	//
	// in: path
	// required: true
	UserID uint64 `json:"user_id"`
}

// swagger:response GetAccountResponse
type GetAccountResponse struct {
	// in: body
	Body GetAccountResponseBody `json:"body"`
}

// swagger:parameters GetAccountResponseBody
type GetAccountResponseBody struct {
	Status domain.ResponseStatus `json:"status"`
	Data   *domain.AccountRead   `json:"data"`
}

var (
	_ domain.Bodier = GetAccountResponse{}
)

func (r GetAccountResponse) GetBody() interface{} { return r.Body }
