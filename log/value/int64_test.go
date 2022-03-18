package value_test

import (
	"fmt"
	"testing"

	"github.com/now/x/log/value"
)

func TestInt64Write(t *testing.T) {
	tests := []struct {
		i    int64
		want string
	}{
		{1, "1"},
	}
	for _, tt := range tests {
		expression := fmt.Sprintf("value.Int64(%#v).Write(…)", tt.i)
		var w value.BytesWriter
		if err := value.Int64(tt.i).Write(&w); err != nil {
			t.Errorf("%s = %v, want %#v", expression, err, nil)
		} else if got := string(w.Bytes); got != tt.want {
			t.Errorf("%s = %#v, want %#v", expression, got, tt.want)
		}
	}
}
