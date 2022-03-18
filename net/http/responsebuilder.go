package http

import (
	"fmt"
	"io"
	"net/http"
	"strings"

	xjson "github.com/now/x/encoding/json"
)

// ResponseBuilder streamlines set-up and send-off of *http.Responses.
//
// Set headers with ContentType(t).  Pass JSON in the response body with
// JSONBody(v) and get back JSON results with JSONResult(result).  Finally,
// extract your request with Build(ctx, method) or send off your request with
// Post(ctx).
type ResponseBuilder struct {
	contentType string
	body        io.ReadCloser
}

// NewResponseBuilder with no header or body.
func NewResponseBuilder() *ResponseBuilder {
	return &ResponseBuilder{}
}

func (b *ResponseBuilder) GoString() string {
	s := strings.Builder{}
	fmt.Fprintf(&s, "http.NewResponseBuilder()")
	if b.contentType == "application/json" {
		fmt.Fprintf(&s, ".JSONBody(%#v)", b.body)
	} else {
		if b.contentType != "" {
			fmt.Fprintf(&s, ".ContentType(%#v)", b.contentType)
		}
		if b.body != nil {
			fmt.Fprintf(&s, ".Body(%#v)", b.body)
		}
	}
	return s.String()
}

// ContentType header will be set to t.
func (b *ResponseBuilder) ContentType(t string) *ResponseBuilder {
	b.contentType = t
	return b
}

// Body of the response will be set to r.
func (b *ResponseBuilder) Body(r io.ReadCloser) *ResponseBuilder {
	b.body = r
	return b
}

// JSONBody of the response will be set to value.
//
// This sets r.ContentType("application/json").
func (b *ResponseBuilder) JSONBody(v xjson.Value) *ResponseBuilder {
	return b.ContentType("application/json").Body(xjson.Encode(v))
}

// Build a *http.Response using statusCode with header and body from b.
//
// The Status field will be set to http.StatusText(statusCode) and the
// StatusCode field will be set to statusCode.
func (b *ResponseBuilder) Build(statusCode int) *http.Response {
	h := http.Header{}
	if b.contentType != "" {
		h.Set("Content-Type", b.contentType)
	}
	return &http.Response{
		Status: func() string {
			if s := http.StatusText(statusCode); s == "" {
				return fmt.Sprintf("status code %d", statusCode)
			} else {
				return s
			}
		}(),
		StatusCode: statusCode,
		Header:     h,
		Body:       b.body,
	}
}

// OK a *http.Response using b.Build(http.StatusOK).
func (b *ResponseBuilder) OK() *http.Response {
	return b.Build(http.StatusOK)
}
