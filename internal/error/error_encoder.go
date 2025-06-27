package error

import (
	"context"
	"encoding/json"
	"net/http"

	"entain/internal/domain"
)

// EncodeError - basic error encoder that encodes any error to the common gateway error.
func EncodeError(_ context.Context, err error, w http.ResponseWriter) {
	// Cannot encode nil error.
	// For these cases, we just do nothing.
	if err == nil {
		return
	}

	// Trying to get http error code and message.
	httpCode := HTTPErrorEncoder(err)
	httpMsg := err.Error()

	// Making a common gateway error response.
	errorResp := ErrResponse{
		Body: struct {
			Status domain.ResponseStatus `json:"status"`
			Data   interface{}           `json:"data"`
			Error  *DefaultErrResponse   `json:"error"`
		}{
			Status: domain.ResponseStatusError,
			Data:   nil,
			Error: &DefaultErrResponse{
				Code:    httpCode,
				Message: httpMsg,
			},
		},
	}

	w.Header().Set("Content-Type", "application/json; charset=utf-8")
	w.WriteHeader(httpCode)

	_ = json.NewEncoder(w).Encode(errorResp.GetBody())
}

// HTTPErrorEncoder - returns mapped http code based on the provided error.
func HTTPErrorEncoder(err error) int {
	switch err.(type) {
	case *ErrInvalidArgument:
		return http.StatusBadRequest
	case *ErrAlreadyExist:
		return http.StatusConflict
	case *ErrNotFound:
		return http.StatusNotFound
	case *ErrFailedPrecondition:
		return http.StatusBadRequest
	case *ErrInternal:
		return http.StatusInternalServerError
	case *ErrUnauthorized:
		return http.StatusUnauthorized
	case *ErrPermissionDenied:
		return http.StatusForbidden
	default:
		return http.StatusInternalServerError
	}
}
