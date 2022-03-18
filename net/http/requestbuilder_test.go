package http_test

import (
	"context"
	"fmt"
	"net/http"
	"testing"

	"github.com/google/go-cmp/cmp"

	"github.com/now/x/encoding/json"
	xhttp "github.com/now/x/net/http"
	"github.com/now/x/net/httptest"
)

func TestRequestBuilderHeaderBody(t *testing.T) {
	tests := []struct {
		builder *xhttp.RequestBuilder
		want    httptest.JSONRequest
	}{
		{
			xhttp.NewRequestBuilder("http://a.b"),
			httptest.JSONRequest{
				URL:    "http://a.b",
				Header: http.Header{},
			},
		},
		{
			xhttp.NewRequestBuilder("http://a.b").BasicAuth("u", "p"),
			httptest.JSONRequest{
				URL:    "http://a.b",
				Header: http.Header{"Authorization": []string{"Basic dTpw"}},
			},
		},
		{
			xhttp.NewRequestBuilder("http://a.b").ContentType("text/plain"),
			httptest.JSONRequest{
				URL:    "http://a.b",
				Header: http.Header{"Content-Type": []string{"text/plain"}},
			},
		},
		{
			xhttp.NewRequestBuilder("http://a.b").JSONBody([]int{1, 2, 3}),
			httptest.JSONRequest{
				URL:    "http://a.b",
				Header: http.Header{"Content-Type": []string{"application/json"}},
				Body:   json.Array{1.0, 2.0, 3.0},
			},
		},
	}
	for _, tt := range tests {
		tt.want.Method = http.MethodGet
		if r, err := tt.builder.Build(context.Background(), tt.want.Method); err != nil {
			t.Errorf("%#v = %v, want %#v", tt.builder, err, nil)
		} else if got, err := httptest.NewJSONRequest(r); err != nil {
			t.Errorf("%#v = %v, want %#v", tt.builder, err, tt.want)
		} else if diff := cmp.Diff(got, tt.want); diff != "" {
			t.Errorf("%#v diff -got +want\n%s", tt.builder, diff)
		}
	}
}

func TestRequestBuilderPost(t *testing.T) {
	want := 1
	var got int
	b := xhttp.NewRequestBuilder("https://example.com").
		BasicAuth("u", "p").
		JSONBody("abc").
		JSONResult(&got)
	expression := fmt.Sprintf("%#v.Post(…)", b)
	if r, err := b.Post(httptest.Using(context.Background(), func(r *http.Request) (*http.Response, error) {
		want := httptest.JSONRequest{
			Method: http.MethodPost,
			URL:    "https://example.com",
			Header: http.Header{
				"Authorization": []string{"Basic dTpw"},
				"Content-Type":  []string{"application/json"},
			},
			Body: "abc",
		}
		if got, err := httptest.NewJSONRequest(r); err != nil {
			t.Errorf("%s = %v, want %#v", expression, err, nil)
		} else if diff := cmp.Diff(got, want); diff != "" {
			t.Errorf("%s diff -got +want\n%s", expression, diff)
		}
		return xhttp.NewResponseBuilder().JSONBody(1).OK(), nil
	})); err != nil {
		t.Errorf("%s = %v, want %#v", expression, err, nil)
	} else if r.StatusCode != http.StatusOK {
		t.Errorf("%s.StatusCode = %d, want %d", expression, r.StatusCode, http.StatusOK)
	} else if got != want {
		t.Errorf("%s = %#v, want %#v", expression, got, want)
	}
}
