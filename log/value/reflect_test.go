package value_test

import (
	"fmt"
	"testing"

	"github.com/now/x/log/value"
)

func TestReflectWrite(t *testing.T) {
	tests := []struct {
		v    interface{}
		want string
	}{
		{1, "1"},
		{"abc", "abc"},
		{map[string]int{"a": 1}, "map[a:1]"},
	}
	for _, tt := range tests {
		expression := fmt.Sprintf("value.Reflect{Value: %#v}.Write(…)", tt.v)
		var w value.BytesWriter
		if err := (value.Reflect{Value: tt.v}.Write(&w)); err != nil {
			t.Errorf("%s = %v, want %#v", expression, err, nil)
		} else if got := string(w.Bytes); got != tt.want {
			t.Errorf("%s = %#v, want %#v", expression, got, tt.want)
		}
	}
}
