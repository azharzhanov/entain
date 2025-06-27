package integrationtestclient

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
)

// Client - http client.
type Client struct {
	baseURL       string
	defaultHeader *RequestHeader
}

// NewClient - returns a new client.
func NewClient(
	baseURL string,
	defaultHeader *RequestHeader,
) *Client {
	return &Client{
		baseURL:       baseURL,
		defaultHeader: defaultHeader,
	}
}

func (client *Client) Do(req *Request) (*Response, error) {

	// Convert request body to JSON.
	var bodyReader io.Reader = nil
	if req.body != nil {
		bodyJSON, err := json.Marshal(req.body)
		if err != nil {
			return nil, fmt.Errorf("error marshaling JSON: %v", err)
		}
		bodyReader = bytes.NewBuffer(bodyJSON)
	}

	// Create the http request with the payload.
	httpReq, err := http.NewRequest(
		req.method,
		client.baseURL+req.requestURL,
		bodyReader,
	)
	if err != nil {
		return nil, fmt.Errorf("error creating request: %v", err)
	}

	// Set the content type.
	httpReq.Header.Set("Content-Type", "application/json")

	// Apply default header.
	if req.header == nil {
		req.SetHeader(client.defaultHeader)
	}

	if req.header != nil {
		// Set source type header if present.
		if req.header.SourceType != "" {
			httpReq.Header.Set("Source-Type", req.header.SourceType)
		}
	}

	// Execute the HTTP request.
	httpResp, err := http.DefaultClient.Do(httpReq)
	if err != nil {
		return nil, fmt.Errorf("error executing request: %v", err)
	}
	defer httpResp.Body.Close()

	// Read the response body.
	body, err := io.ReadAll(httpResp.Body)
	if err != nil {
		return nil, fmt.Errorf("error reading response body: %v", err)
	}

	return &Response{
		status: httpResp.StatusCode,
		body:   body,
	}, nil
}

func (client *Client) SetHeader(header *RequestHeader) *Client {
	client.defaultHeader = header
	return client
}
