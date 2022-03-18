package httptest_test

import (
	"context"
	"fmt"
	"net/http"
	"testing"

	"github.com/google/go-cmp/cmp"

	xhttp "github.com/now/x/net/http"
	"github.com/now/x/net/httptest"
)

func TestJSONRequestFromHTTPRequest(t *testing.T) {
	want := 1
	var got int
	b := xhttp.NewRequestBuilder("https://example.com").JSONBody("abc").JSONResult(&got)
	expression := fmt.Sprintf("%#v.Post(…)", b)
	if r, err := b.Post(httptest.Using(context.Background(), func(r *http.Request) (*http.Response, error) {
		want := httptest.JSONRequest{
			Method: http.MethodPost,
			URL:    "https://example.com",
			Header: http.Header{
				"Content-Type": []string{"application/json"},
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
