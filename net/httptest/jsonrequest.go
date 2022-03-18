package httptest

import (
	"net/http"

	"github.com/now/x/encoding/json"
)

// JSONRequest is a simplified http.Request with a json.Value as its Body.
//
// This mainly makes it easier to set up wanted values in tests.
type JSONRequest struct {
	Method string
	URL    string
	Header http.Header
	Body   json.Value
}

// NewJSONRequest with method, URL, header, and body based on r.
//
// The method, URL, and header are straight clones of those in r.  The body is
// json.DecodeAndClose()d into a json.Value.
func NewJSONRequest(r *http.Request) (JSONRequest, error) {
	var body json.Value
	if err := json.DecodeAndClose(r.Body, &body); err != nil {
		return JSONRequest{}, err
	}
	return JSONRequest{
		Method: r.Method,
		URL:    r.URL.String(),
		Header: r.Header.Clone(),
		Body:   body,
	}, nil
}
