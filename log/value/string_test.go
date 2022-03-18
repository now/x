package value_test

import (
	"fmt"
	"testing"

	"github.com/now/x/log/value"
)

func TestStringWrite(t *testing.T) {
	tests := []struct {
		s    string
		want string
	}{
		{"abc", "abc"},
	}
	for _, tt := range tests {
		expression := fmt.Sprintf("value.String(%#v).Write(…)", tt.s)
		var w value.BytesWriter
		if err := value.String(tt.s).Write(&w); err != nil {
			t.Errorf("%s = %v, want %#v", expression, err, nil)
		} else if got := string(w.Bytes); got != tt.want {
			t.Errorf("%s = %#v, want %#v", expression, got, tt.want)
		}
	}
}
