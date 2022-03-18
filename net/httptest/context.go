// Package httptest contains extension for working with the net/httptest package.
//
// Using lets you set a http.RoundTripperFunc to use as a http.Client.Transport
// that allows you to intercept any requests being made, verify their content,
// and provide mock responses.
//
// The JSONRequest type and the NewJSONRequest() constructor allow you to more
// easily compare JSON requests that are being made, primarily by unmarshaling
// the body of the request.
package httptest

import (
	"context"
	"net/http"

	xhttp "github.com/now/x/net/http"
)

// Using is ctxʹ ≈ ctx such that http.In(ctxʹ) is a *http.Client that uses transport.
//
// This is useful for intercepting requests, verifying their content, and
// providing mock responses.
func Using(ctx context.Context, transport func(*http.Request) (*http.Response, error)) context.Context {
	return xhttp.Using(ctx, &http.Client{Transport: xhttp.RoundTripperFunc(transport)})
}
