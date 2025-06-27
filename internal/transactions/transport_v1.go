package transactions

import (
	"context"
	"encoding/json"
	"net/http"

	"entain/internal/domain"
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

	// AggregateTransaction swagger:route POST /user/{user_id}/transaction users AggregateTransactionRequest
	//
	// Aggregates a transaction for a user.
	//
	// Responses:
	// default: errorResponse
	//	   200: AggregateTransactionResponse
	r.Methods("POST").Path("/user/{user_id}/transaction").Handler(httptransport.NewServer(
		e.AggregateTransactionEndpoint,
		func(ctx context.Context, r *http.Request) (interface{}, error) {
			vars := mux.Vars(r)
			userID, err := helpers.ExtractInt64Route(vars, "user_id")
			if err != nil {
				return nil, err
			}

			var body domain.TransactionWrite
			if err := json.NewDecoder(r.Body).Decode(&body); err != nil {
				return nil, err
			}

			return AggregateTransactionRequest{
				UserID:      userID,
				SourceType:  r.Header.Get("Source-Type"),
				Transaction: &body,
			}, nil
		},
		helpers.EncodeResponse,
		options...,
	))
}
