// Package http contains extensions for working with the net/http package.
//
// There are three parts to this package: dealing with an http.Client in a
// context.Context, creating http.Requests with RequestBuilder, and creating
// http.Responses with ResponseBuilder.
//
// Adding and accessing a http.Client in a context.Context is useful in
// production code for handling cookies and timeouts in a controlled manner.
// (The http.Client type is safe to share between many goroutines, so it’s a
// good idea to only have one.)  It’s also useful in testing in that
// http.Client.Transport can be overridden to intercept http.Requests being
// made, verify that they’re the wanted ones, and providing mock http.Responses.
//
// Creating http.Requests using RequestBuilder is easy, as it can help in
// setting headers (BasicAuth, ContentType), marshaling and unmarshaling JSON
// request and response bodies (JSONBody/JSONResult), and building (Build) or
// sending the request (Post).
//
// Creating http.Responses using ResponseBuilder is equally easy, as it can help
// in setting headers (ContentType),marshaling JSON bodies (JSONBody), and
// setting statuses (Build, OK).
package http

import (
	"context"
	"net/http"
)

// Using is ctxʹ ≈ ctx such that In(ctxʹ) is c.
func Using(ctx context.Context, c *http.Client) context.Context {
	return context.WithValue(ctx, Key, c)
}

// In is the http.Client c previously added to ctx with Using(ctx, c).
//
// In is http.DefaultClient if none has previously been added to ctx with
// Using(ctx, …).
//
// TODO(now) Remove default?
func In(ctx context.Context) *http.Client {
	if client, _ := ctx.Value(Key).(*http.Client); client == nil {
		return http.DefaultClient
	} else {
		return client
	}
}

// Key of http.Client in context.Context.
const Key key = "HttpClient"

type key string
