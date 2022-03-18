package http_test

import (
	"context"
	"net/http"
	"testing"

	xhttp "github.com/now/x/net/http"
)

func TestUsing(t *testing.T) {
	type keytype string
	key := keytype("test")
	v := 1
	want := &http.Client{}
	ctx := xhttp.Using(context.WithValue(context.Background(), key, v), want)
	if got := xhttp.In(ctx); got != want {
		t.Errorf("xhttp.In(xhttp.Using(…, %#v)) = %#v, want %#v", want, got, want)
	} else if got, _ := ctx.Value(key).(int); got != v {
		t.Errorf("xhttp.Using(context.WithValue(…, %#v, %#v), …).Value(%#v) = %#v, want %#v", key, v, key, got, v)
	}
}

func TestIn(t *testing.T) {
	want := http.DefaultClient
	if got := xhttp.In(context.Background()); got != want {
		t.Errorf("xhttp.In(…) = %#v, want %#v", got, want)
	}
}
