package log_test

import (
	"fmt"
	"testing"
	"time"

	"github.com/now/x/log"
	"github.com/now/x/log/value"
)

func TestError(t *testing.T) {
	tests := []struct {
		err  error
		want string
	}{
		{fmt.Errorf("failed"), "error: failed"},
	}
	for _, tt := range tests {
		expression := fmt.Sprintf("log.Error(%#v).Write(…)", tt.err)
		var w value.BytesWriter
		if err := log.Error(tt.err).Write(&w); err != nil {
			t.Errorf("%s = %v, want %#v", expression, err, nil)
		} else if got := string(w.Bytes); got != tt.want {
			t.Errorf("%s = %#v, want %#v", expression, got, tt.want)
		}
	}
}

func TestInt(t *testing.T) {
	tests := []struct {
		label string
		i     int
		want  string
	}{
		{"ID", 1, "ID: 1"},
	}
	for _, tt := range tests {
		expression := fmt.Sprintf("log.Int(%#v, %#v).Write(…)", tt.label, tt.i)
		var w value.BytesWriter
		if err := log.Int(tt.label, tt.i).Write(&w); err != nil {
			t.Errorf("%s = %v, want %#v", expression, err, nil)
		} else if got := string(w.Bytes); got != tt.want {
			t.Errorf("%s = %#v, want %#v", expression, got, tt.want)
		}
	}
}

func TestInt64(t *testing.T) {
	tests := []struct {
		label string
		i     int64
		want  string
	}{
		{"ID", 1, "ID: 1"},
	}
	for _, tt := range tests {
		expression := fmt.Sprintf("log.Int64(%#v, %#v).Write(…)", tt.label, tt.i)
		var w value.BytesWriter
		if err := log.Int64(tt.label, tt.i).Write(&w); err != nil {
			t.Errorf("%s = %v, want %#v", expression, err, nil)
		} else if got := string(w.Bytes); got != tt.want {
			t.Errorf("%s = %#v, want %#v", expression, got, tt.want)
		}
	}
}

func TestReflect(t *testing.T) {
	tests := []struct {
		label string
		r     interface{}
		want  string
	}{
		{"structure", struct{ A, B int }{1, 2}, "structure: {A:1 B:2}"},
	}
	for _, tt := range tests {
		expression := fmt.Sprintf("log.Reflect(%#v, %#v).Write(…)", tt.label, tt.r)
		var w value.BytesWriter
		if err := log.Reflect(tt.label, tt.r).Write(&w); err != nil {
			t.Errorf("%s = %v, want %#v", expression, err, nil)
		} else if got := string(w.Bytes); got != tt.want {
			t.Errorf("%s = %#v, want %#v", expression, got, tt.want)
		}
	}
}

func TestString(t *testing.T) {
	tests := []struct {
		label string
		s     string
		want  string
	}{
		{"name", "abc", "name: abc"},
	}
	for _, tt := range tests {
		expression := fmt.Sprintf("log.String(%#v, %#v).Write(…)", tt.label, tt.s)
		var w value.BytesWriter
		if err := log.String(tt.label, tt.s).Write(&w); err != nil {
			t.Errorf("%s = %v, want %#v", expression, err, nil)
		} else if got := string(w.Bytes); got != tt.want {
			t.Errorf("%s = %#v, want %#v", expression, got, tt.want)
		}
	}
}

func TestStringer(t *testing.T) {
	tests := []struct {
		label string
		s     fmt.Stringer
		want  string
	}{
		{"time", time.Date(2022, time.March, 9, 19, 51, 0, 0, time.UTC), "time: 2022-03-09 19:51:00 +0000 UTC"},
	}
	for _, tt := range tests {
		expression := fmt.Sprintf("log.Stringer(%#v, %#v).Write(…)", tt.label, tt.s)
		var w value.BytesWriter
		if err := log.Stringer(tt.label, tt.s).Write(&w); err != nil {
			t.Errorf("%s = %v, want %#v", expression, err, nil)
		} else if got := string(w.Bytes); got != tt.want {
			t.Errorf("%s = %#v, want %#v", expression, got, tt.want)
		}
	}
}
