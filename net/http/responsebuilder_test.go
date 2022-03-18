package http_test

import (
	"net/http"
	"testing"

	"github.com/google/go-cmp/cmp"

	"github.com/now/x/encoding/json"
	xhttp "github.com/now/x/net/http"
	"github.com/now/x/net/httptest"
)

func TestResponseBuilderHeaderBody(t *testing.T) {
	tests := []struct {
		builder *xhttp.ResponseBuilder
		want    httptest.JSONResponse
	}{
		{
			xhttp.NewResponseBuilder(),
			httptest.JSONResponse{
				Header: http.Header{},
			},
		},
		{
			xhttp.NewResponseBuilder().ContentType("text/plain"),
			httptest.JSONResponse{
				Header: http.Header{"Content-Type": []string{"text/plain"}},
			},
		},
		{
			xhttp.NewResponseBuilder().JSONBody([]int{1, 2, 3}),
			httptest.JSONResponse{
				Header: http.Header{"Content-Type": []string{"application/json"}},
				Body:   json.Array{1.0, 2.0, 3.0},
			},
		},
	}
	for _, tt := range tests {
		tt.want.StatusCode = http.StatusOK
		if got, err := httptest.NewJSONResponse(tt.builder.OK()); err != nil {
			t.Errorf("%#v = %v, want %#v", tt.builder, err, nil)
		} else if diff := cmp.Diff(got, tt.want); diff != "" {
			t.Errorf("%#v diff -got +want\n%s", tt.builder, diff)
		}
	}
}
