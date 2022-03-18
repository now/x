package httptest_test

import (
	"fmt"
	"net/http"
	"testing"

	"github.com/google/go-cmp/cmp"

	xhttp "github.com/now/x/net/http"
	"github.com/now/x/net/httptest"
)

func TestJSONResponseFromHTTPResponse(t *testing.T) {
	b := xhttp.NewResponseBuilder().JSONBody("abc")
	expression := fmt.Sprintf("%#v.OK()", b)
	want := httptest.JSONResponse{
		StatusCode: http.StatusOK,
		Header: http.Header{
			"Content-Type": []string{"application/json"},
		},
		Body: "abc",
	}
	if got, err := httptest.NewJSONResponse(b.OK()); err != nil {
		t.Errorf("%s = %#v, want nil\n", expression, err)
	} else if diff := cmp.Diff(got, want); diff != "" {
		t.Errorf("%s diff -got +want\n%s", expression, diff)
	}
}
