package helpers

import (
	"context"
	"encoding/json"
	"net/http"

	"entain/internal/domain"
	errors "entain/internal/error"
)

// errorer is implemented by all concrete response types that may contain
// errors. It allows us to change the HTTP response code without needing to
// trigger an endpoint (transport-level) error.
type errorer interface {
	error() error
}

// EncodeResponse is the common method to encode all response types to the
// client. I chose to do it this way because, since we're using JSON, there's no
// reason to provide anything more specific. It's certainly possible to
// specialize on a per-response (per-method) basis.
func EncodeResponse(ctx context.Context, w http.ResponseWriter, response interface{}) error {
	if e, ok := response.(errorer); ok && e.error() != nil {
		// Not a Go kit transport error, but a business-logic error.
		// Provide those as HTTP errors.
		errors.EncodeError(ctx, e.error(), w)
		return nil
	}
	w.Header().Set("Content-Type", "application/json; charset=utf-8")
	if b, ok := response.(domain.Bodier); ok {
		// If the response implements Bodier interface then we will extract body.
		return json.NewEncoder(w).Encode(b.GetBody())
	}
	return json.NewEncoder(w).Encode(response)
}
