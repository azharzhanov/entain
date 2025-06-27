package integrationtestclient

import (
	"encoding/json"
	"fmt"
)

// Response - http response.
type Response struct {
	status int
	body   []byte
}

func (r *Response) Read(v any) error {
	err := json.Unmarshal(r.body, v)
	if err != nil {
		return fmt.Errorf("error unmarshaling JSON response: %v", err)
	}
	return nil
}

func (r *Response) ReadOptionalWrapped(v any) error {
	type wrappedType struct {
		Body json.RawMessage `json:"body"`
	}

	var wrapped wrappedType
	err := json.Unmarshal(r.body, &wrapped)
	if err != nil {
		return fmt.Errorf("error unmarshaling wrapped JSON response: %v", err)
	}

	if len(wrapped.Body) > 0 {
		err := json.Unmarshal(wrapped.Body, v)
		if err != nil {
			return fmt.Errorf("error unmarshaling nested JSON response: %v", err)
		}

		return nil
	}

	return r.Read(v)
}

func (r *Response) IsOK() bool {
	return r.status == 200
}

func (r *Response) StatusCode() int {
	return r.status
}

func (r *Response) BodyString() string {
	return string(r.body)
}
