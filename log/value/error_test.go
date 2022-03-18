package value_test

import (
	"fmt"
	"testing"

	"github.com/now/x/log/value"
)

type panicError struct {
	message string
}

func (e *panicError) Error() string {
	panic(e.message)
}

func TestErrorWrite(t *testing.T) {
	tests := []struct {
		err     error
		want    string
		wantErr error
	}{
		{fmt.Errorf("abc"), "abc", nil},
		{nil, "<nil>", nil},
		{(*panicError)(nil), "<nil>", nil},
		{&panicError{"boom"}, "", fmt.Errorf("PANIC=boom")},
	}
	for _, tt := range tests {
		expression := fmt.Sprintf("value.Error{Err: %#v}.Write(…)", tt.err)
		var w value.BytesWriter
		if err := (value.Error{Err: tt.err}.Write(&w)); err != tt.wantErr &&
			!(err != nil && tt.wantErr != nil && err.Error() == tt.wantErr.Error()) {
			t.Errorf("%s = %v, want %v", expression, err, tt.wantErr)
		}
		if got := string(w.Bytes); got != tt.want {
			t.Errorf("%s = %#v, want %#v", expression, got, tt.want)
		}
	}
}
