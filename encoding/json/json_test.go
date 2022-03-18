package json_test

import (
	"fmt"
	"io"
	"strings"
	"testing"

	"github.com/google/go-cmp/cmp"

	"github.com/now/x/encoding/json"
)

func TestDecodeAndClose(t *testing.T) {
	expression := "json.DecodeAndClose(…)"

	t.Run("valid input generates wanted output", func(t *testing.T) {
		tests := []struct {
			r    io.ReadCloser
			want float64
		}{
			{io.NopCloser(strings.NewReader("1")), 1},
			{nil, 0},
		}
		for _, tt := range tests {
			var got float64
			if err := json.DecodeAndClose(tt.r, &got); err != nil {
				t.Errorf("%s = %v, want %#v", expression, err, tt.want)
			} else if got != tt.want {
				t.Errorf("%s = %#v, want %#v", expression, got, tt.want)
			}
		}
	})

	t.Run("errors if read errors", func(t *testing.T) {
		var got json.Value
		if err := json.DecodeAndClose(io.NopCloser(failingReader{}), &got); err == nil {
			t.Errorf("%s = %#v, want err", expression, got)
		}
	})

	t.Run("errors on invalid input", func(t *testing.T) {
		var got json.Value
		if err := json.DecodeAndClose(io.NopCloser(strings.NewReader("")), &got); err == nil {
			t.Errorf("%s = %#v, want err", expression, got)
		}
	})

	t.Run("ignores error on close", func(t *testing.T) {
		want := 1.0
		var got float64
		if err := json.DecodeAndClose(failingCloser{strings.NewReader(fmt.Sprint(want))}, &got); err != nil {
			t.Errorf("%s = %v, want %#v", expression, err, want)
		} else if got != want {
			t.Errorf("%s = %#v, want %#v", expression, got, want)
		}
	})

	t.Run("decodes to expected types", func(t *testing.T) {
		want := json.Object{"a": json.Array{1.0, 2.0, 3.0}}
		var got json.Object
		if err := json.DecodeAndClose(io.NopCloser(strings.NewReader(`{"a": [1, 2, 3]}`)), &got); err != nil {
			t.Error(err)
		} else if diff := cmp.Diff(got, want); diff != "" {
			t.Errorf("json.DecodeAndClose(io.NopCloser(strings.NewReader(`{\"a\": [1, 2, 3]}`)), …) diff -got +want\n%s", diff)
		}
	})
}

func TestEncode(t *testing.T) {
	t.Run("valid input generates wanted output", func(t *testing.T) {
		v := json.Value(1)
		r := json.Encode(v)
		defer r.Close()
		if b, err := io.ReadAll(r); err != nil {
			t.Error(err)
		} else if got, want := string(b), "1\n"; got != want {
			t.Errorf("json.Encode(%#v) = %#v, want %#v", v, got, want)
		}
	})

	t.Run("errors on invalid input (circular reference)", func(t *testing.T) {
		v := json.Object{}
		v["v"] = v
		r := json.Encode(v)
		defer r.Close()
		if b, err := io.ReadAll(r); err == nil {
			t.Errorf("json.Encode(%#v) = %q, want err", v, string(b))
		}
	})
}

type failingReader struct{}

func (failingReader) Read([]byte) (int, error) {
	return 0, fmt.Errorf("can’t read")
}

type failingCloser struct {
	io.Reader
}

func (failingCloser) Close() error {
	return fmt.Errorf("can’t close")
}
