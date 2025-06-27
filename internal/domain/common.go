package domain

// ResponseStatus - represents the response status.
type ResponseStatus string

// Response statuses.
const (
	ResponseStatusSuccess = ResponseStatus("success")
	ResponseStatusError   = ResponseStatus("error")
)

// Bodier may be implemented by the gateway response types.
// It's not necessary for your response types to implement Bodier, but it may
// help for more sophisticated use cases.
type Bodier interface {
	GetBody() interface{}
}
