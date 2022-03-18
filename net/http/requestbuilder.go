package http

import (
	"context"
	"fmt"
	"io"
	"net/http"
	"strings"

	"github.com/now/x/encoding/json"
)

// RequestBuilder streamlines set-up and send-off of *http.Requests.
//
// Set headers with BasicAuth(username, password) and ContentType(t).  Pass JSON
// in the request body with JSONBody(v) and get back JSON results with
// JSONResult(v).  Finally, extract your request with Build(ctx, method) or send
// off your request with Post(ctx).
type RequestBuilder struct {
	url         string
	auth        bool
	username    string
	password    string
	contentType string
	body        io.Reader
	jsonResult  json.Value
}

// NewRequestBuilder to url with no header or body.
func NewRequestBuilder(url string) *RequestBuilder {
	return &RequestBuilder{url: url}
}

func (b *RequestBuilder) GoString() string {
	s := strings.Builder{}
	fmt.Fprintf(&s, "http.NewRequestBuilder(%#v)", b.url)
	if b.auth {
		fmt.Fprintf(&s, ".BasicAuth(%#v, %#v)", b.username, b.password)
	}
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
	if b.jsonResult != nil {
		fmt.Fprintf(&s, ".JSONResult(%#v)", b.jsonResult)
	}
	return s.String()
}

// BasicAuth with username and password.
func (b *RequestBuilder) BasicAuth(username, password string) *RequestBuilder {
	b.auth, b.username, b.password = true, username, password
	return b
}

// ContentType header will be set to t.
func (b *RequestBuilder) ContentType(t string) *RequestBuilder {
	b.contentType = t
	return b
}

// Body of the request will be set to r.
func (b *RequestBuilder) Body(r io.Reader) *RequestBuilder {
	b.body = r
	return b
}

// JSONBody of the request will be set to v.
//
// This sets r.ContentType("application/json").
func (b *RequestBuilder) JSONBody(v json.Value) *RequestBuilder {
	return b.ContentType("application/json").Body(json.Encode(v))
}

// JSONResult of the response will be unmarshaled into v.
func (b *RequestBuilder) JSONResult(v json.Value) *RequestBuilder {
	b.jsonResult = v
	return b
}

// Build a *http.Request in ctx using method with URL, header, and body from b.
//
// Errors if BasicAuth(username, password) was passed a username that contains a
// colon, ‘:’, or if http.NewRequestWithContext(ctx, method, …) errors.
func (b *RequestBuilder) Build(ctx context.Context, method string) (*http.Request, error) {
	if b.auth && strings.ContainsRune(b.username, ':') {
		return nil, fmt.Errorf("httpx: BasicAuth username can’t contain colon, got “%s”", b.username)
	} else if r, err := http.NewRequestWithContext(ctx, method, b.url, b.body); err != nil {
		return nil, err
	} else {
		if b.auth {
			r.SetBasicAuth(b.username, b.password)
		}
		if b.contentType != "" {
			r.Header.Set("Content-Type", b.contentType)
		}
		return r, nil
	}
}

// Post a *http.Request r in ctx with URL, header, and body from b.
//
// The *http.Client In(ctx) is used to Do(r).  A response is unmarshaled
// into any v given to JSONResult().
//
// Errors if r.Build(ctx, …), In(ctx).Do(…), or json.Decode(…)  errors.
func (b *RequestBuilder) Post(ctx context.Context) (*http.Response, error) {
	if r, err := b.Build(ctx, http.MethodPost); err != nil {
		return nil, err
	} else if rʹ, err := In(ctx).Do(r); err != nil {
		return rʹ, err
	} else {
		if b.jsonResult != nil {
			err = json.DecodeAndClose(rʹ.Body, b.jsonResult)
		}
		return rʹ, err
	}
}
