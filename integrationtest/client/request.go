package integrationtestclient

import (
	"net/url"
)

// RequestHeader - http request headers.
type RequestHeader struct {
	SourceType  string
	ContentType string
}

// Request - http request.
type Request struct {
	method     string
	header     *RequestHeader
	requestURL string
	body       interface{}
}

func (req *Request) SetHeader(header *RequestHeader) *Request {
	req.header = header
	return req
}

func newRequest(
	method string,
	header *RequestHeader,
	path string,
	urlValues url.Values,
	body interface{},
) *Request {
	// Create the complete request URL with query parameters.
	requestURL := path
	if len(urlValues) > 0 {
		requestURL = requestURL + "?" + urlValues.Encode()
	}
	return &Request{
		method:     method,
		header:     header,
		requestURL: requestURL,
		body:       body,
	}
}

// NewGetRequest - returns new GET request.
func NewGetRequest(
	path string,
	urlValues url.Values,
) *Request {
	return newRequest(
		"GET",
		nil, // no header
		path,
		urlValues,
		nil, // no body
	)
}

// NewPostRequest - returns new POST request.
func NewPostRequest(
	path string,
	body interface{},
) *Request {
	return newRequest(
		"POST",
		nil, // no header
		path,
		nil, // no url values
		body,
	)
}

// NewPatchRequest - returns new PATCH request.
func NewPatchRequest(
	path string,
	body interface{},
) *Request {
	return newRequest(
		"PATCH",
		nil, // no header
		path,
		nil, // no url values
		body,
	)
}

// NewPutRequest - returns new PUT request.
func NewPutRequest(
	path string,
	body interface{},
) *Request {
	return newRequest(
		"PUT",
		nil, // no header
		path,
		nil, // no url values
		body,
	)
}

// NewDeleteRequest - returns new DELETE request.
func NewDeleteRequest(
	path string,
) *Request {
	return newRequest(
		"DELETE",
		nil, // no header
		path,
		nil, // no url values
		nil, // no body
	)
}
