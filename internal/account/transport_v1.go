package account

import (
	"context"
	"net/http"

	"entain/internal/helpers"

	httptransport "github.com/go-kit/kit/transport/http"
	"github.com/go-kit/log"
	"github.com/gorilla/mux"
)

// RegisterRoutersV1 registers routers with version 1 endpoints.
func RegisterRoutersV1(
	r *mux.Router,
	e Endpoints,
	logger log.Logger,
) {
	options := helpers.SetupServerOptions(logger)

	// GetAccount swagger:route GET /user/{user_id}/balance users GetAccountRequest
	//
	// Returns the account of the user by id.
	//
	// Responses:
	// default: errorResponse
	//	   200: GetAccountResponse
	r.Methods("GET").Path("/user/{user_id}/balance").Handler(httptransport.NewServer(
		e.GetAccountEndpoint,
		func(ctx context.Context, r *http.Request) (interface{}, error) {
			vars := mux.Vars(r)
			userID, err := helpers.ExtractInt64Route(vars, "user_id")
			if err != nil {
				return nil, err
			}

			return GetAccountRequest{
				UserID: userID,
			}, nil
		},
		helpers.EncodeResponse,
		options...,
	))
}
