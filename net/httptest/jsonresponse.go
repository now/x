package httptest

import (
	"net/http"

	"github.com/now/x/encoding/json"
)

// JSONResponse is a simplified http.Response with a json.Value as its Body.
//
// This mainly makes it easier to set up wanted values in tests.
type JSONResponse struct {
	StatusCode int
	Header     http.Header
	Body       json.Value
}

// NewJSONResponse with status, header, and body based on r.
//
// The status and header are straight clones of those in r.  The body is
// json.DecodeAndClose()d into a json.Value.
func NewJSONResponse(r *http.Response) (JSONResponse, error) {
	var body json.Value
	if err := json.DecodeAndClose(r.Body, &body); err != nil {
		return JSONResponse{}, err
	}
	return JSONResponse{
		StatusCode: r.StatusCode,
		Header:     r.Header.Clone(),
		Body:       body,
	}, nil
}
