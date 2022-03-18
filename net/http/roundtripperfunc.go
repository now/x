package http

import "net/http"

// RoundTripperFunc simpflifies creating a http.RoundTripper.
type RoundTripperFunc func(*http.Request) (*http.Response, error)

func (f RoundTripperFunc) RoundTrip(request *http.Request) (*http.Response, error) {
	return f(request)
}
