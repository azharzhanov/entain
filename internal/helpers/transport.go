package helpers

import (
	errors "entain/internal/error"

	"github.com/go-kit/kit/auth/jwt"
	"github.com/go-kit/kit/transport"
	kithttp "github.com/go-kit/kit/transport/http"
	"github.com/go-kit/log"
)

// SetupServerOptions - basic setup for all servers.
func SetupServerOptions(logger log.Logger) []kithttp.ServerOption {
	options := []kithttp.ServerOption{
		kithttp.ServerErrorHandler(transport.NewLogErrorHandler(logger)),
		kithttp.ServerErrorEncoder(errors.EncodeError),
		kithttp.ServerBefore(jwt.HTTPToContext()),
	}
	return options
}
