package error

import (
	"database/sql"
)

// ErrType - error type.
type ErrType string

// GenericError - generic error interface.
type GenericError interface {
	Error() string
	ErrorType() ErrType
}

// Supported error types.
const (
	ErrTypeUnknown            = ErrType("Unknown")
	ErrTypeInvalidArgument    = ErrType("InvalidArgument")
	ErrTypeAlreadyExist       = ErrType("AlreadyExist")
	ErrTypeNotFound           = ErrType("NotFound")
	ErrTypeFailedPrecondition = ErrType("FailedPrecondition")
	ErrTypeInternal           = ErrType("Internal")
	ErrTypeUnauthenticated    = ErrType("Unauthenticated")
	ErrTypePermissionDenied   = ErrType("PermissionDenied")
)

// ErrInvalidArgument - invalid argument
type ErrInvalidArgument struct {
	Msg string
}

func (e *ErrInvalidArgument) Error() string {
	return e.Msg
}

// ErrorType - returns "ErrTypeInvalidArgument".
func (ErrInvalidArgument) ErrorType() ErrType {
	return ErrTypeInvalidArgument
}

// ErrAlreadyExist - already exist.
type ErrAlreadyExist struct {
	Msg string
}

func (e *ErrAlreadyExist) Error() string {
	return e.Msg
}

// ErrorType - returns "ErrTypeAlreadyExist".
func (ErrAlreadyExist) ErrorType() ErrType {
	return ErrTypeAlreadyExist
}

// ErrNotFound - not found.
type ErrNotFound struct {
	Msg string
}

func (e *ErrNotFound) Error() string {
	return e.Msg
}

// ErrorType - returns "ErrTypeNotFound".
func (ErrNotFound) ErrorType() ErrType {
	return ErrTypeNotFound
}

// ErrFailedPrecondition - failed pre condition.
type ErrFailedPrecondition struct {
	Msg string
}

func (e *ErrFailedPrecondition) Error() string {
	return e.Msg
}

// ErrorType - returns "ErrTypeFailedPrecondition".
func (ErrFailedPrecondition) ErrorType() ErrType {
	return ErrTypeFailedPrecondition
}

// ErrInternal - internal error.
type ErrInternal struct {
	Msg string
}

func (e *ErrInternal) Error() string {
	return e.Msg
}

// ErrorType - returns "ErrTypeInternal".
func (ErrInternal) ErrorType() ErrType {
	return ErrTypeInternal
}

// ErrUnauthorized - unauthorized access.
type ErrUnauthorized struct {
	Msg string
}

func (e *ErrUnauthorized) Error() string {
	return e.Msg
}

// ErrorType - returns "ErrTypeUnauthenticated".
func (ErrUnauthorized) ErrorType() ErrType {
	return ErrTypeUnauthenticated
}

// ErrPermissionDenied - permission or access denied.
type ErrPermissionDenied struct {
	Msg string
}

func (e *ErrPermissionDenied) Error() string {
	return e.Msg
}

// ErrorType - returns "ErrTypePermissionDenied".
func (ErrPermissionDenied) ErrorType() ErrType {
	return ErrTypePermissionDenied
}

// NewErrInvalidArgument creates a new error.
func NewErrInvalidArgument(msg string) error {
	return &ErrInvalidArgument{msg}
}

// NewErrAlreadyExist already exist.
func NewErrAlreadyExist(msg string) error {
	return &ErrAlreadyExist{msg}
}

// NewErrNotFound not found.
func NewErrNotFound(msg string) error {
	return &ErrNotFound{msg}
}

// NewErrFailedPrecondition failed pre condition.
func NewErrFailedPrecondition(msg string) error {
	return &ErrFailedPrecondition{msg}
}

// NewErrInternal internal error.
func NewErrInternal(msg string) error {
	return &ErrInternal{msg}
}

// NewErrUnauthorized unauthorized access.
func NewErrUnauthorized(msg string) error {
	return &ErrUnauthorized{msg}
}

// NewErrPermissionDenied permission denied.
func NewErrPermissionDenied(msg string) error {
	return &ErrPermissionDenied{msg}
}

// compile time assertions for our response types implementing error interface.
var (
	_ error = &ErrInvalidArgument{}
	_ error = &ErrAlreadyExist{}
	_ error = &ErrNotFound{}
	_ error = &ErrFailedPrecondition{}
	_ error = &ErrInternal{}
	_ error = &ErrUnauthorized{}
	_ error = &ErrPermissionDenied{}
)

// compile time assertions for our response types implementing GenericError interface.
var (
	_ GenericError = &ErrInvalidArgument{}
	_ GenericError = &ErrAlreadyExist{}
	_ GenericError = &ErrNotFound{}
	_ GenericError = &ErrFailedPrecondition{}
	_ GenericError = &ErrInternal{}
	_ GenericError = &ErrUnauthorized{}
	_ GenericError = &ErrPermissionDenied{}
)

// IsErrNotFound - determines if the error is "Not Found"
func IsErrNotFound(err error) bool {
	if err == nil {
		return false
	}
	if err == sql.ErrNoRows {
		return true
	}
	genericError, ok := err.(GenericError)
	if !ok {
		return false
	}
	return genericError.ErrorType() == ErrTypeNotFound
}
