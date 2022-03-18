package value_test

import (
	"fmt"
	"testing"

	"github.com/now/x/log/value"
)

type panicStringer struct {
	message string
}

func (s *panicStringer) String() string {
	panic(s.message)
}

type stringer string

func (s stringer) String() string {
	return string(s)
}

func TestStringerWriter(t *testing.T) {
	tests := []struct {
		stringer fmt.Stringer
		want     string
		wantErr  error
	}{
		{stringer("abc"), "abc", nil},
		{nil, "<nil>", nil},
		{(*panicStringer)(nil), "<nil>", nil},
		{&panicStringer{"boom"}, "", fmt.Errorf("PANIC=boom")},
	}
	for _, tt := range tests {
		expression := fmt.Sprintf("value.Stringer{Value: %#v}.Write(…)", tt.stringer)
		var w value.BytesWriter
		if err := (value.Stringer{Value: tt.stringer}.Write(&w)); err != tt.wantErr &&
			!(err != nil && tt.wantErr != nil && err.Error() == tt.wantErr.Error()) {
			t.Errorf("%s = %v, want %v", expression, err, tt.wantErr)
		}
		if got := string(w.Bytes); got != tt.want {
			t.Errorf("%s = %#v, want %#v", expression, got, tt.want)
		}
	}
}
