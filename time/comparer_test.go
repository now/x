package time_test

import (
	"testing"
	"time"

	"github.com/google/go-cmp/cmp"

	xtime "github.com/now/x/time"
)

func TestComparer(t *testing.T) {
	tests := []time.Time{
		xtime.MustDate(t, 2022, time.February, 27, 20, 0, 12, 2, "Europe/Stockholm"),
	}
	for _, tt := range tests {
		if !cmp.Equal(tt, tt, xtime.Comparer) {
			t.Errorf("cmp.Equal(%#v, %#v, timex.Comparer) = false, want true", tt, tt)
		}
		if stockholm, err := time.LoadLocation(tt.Location().String()); err != nil {
			t.Errorf("time.LoadLocation(%q) = %v, want %v", tt.Location().String(), err, tt.Location())
		} else if tu := xtime.WallIn(tt, stockholm); !cmp.Equal(tt, tu, xtime.Comparer) {
			t.Errorf("cmp.Equal(%#v, %#v, xtime.Comparer) = false, want true", tt, tu)
		}
		if tu := tt.AddDate(1, 0, 0); cmp.Equal(tt, tu, xtime.Comparer) {
			t.Errorf("cmp.Equal(%#v, %#v, xtime.Comparer) = true, want false", tt, tu)
		}
		if tu := tt.AddDate(0, 1, 0); cmp.Equal(tt, tu, xtime.Comparer) {
			t.Errorf("cmp.Equal(%#v, %#v, xtime.Comparer) = true, want false", tt, tu)
		}
		if tu := tt.AddDate(0, 0, 1); cmp.Equal(tt, tu, xtime.Comparer) {
			t.Errorf("cmp.Equal(%#v, %#v, xtime.Comparer) = true, want false", tt, tu)
		}
		if tu := tt.Add(1 * time.Hour); cmp.Equal(tt, tu, xtime.Comparer) {
			t.Errorf("cmp.Equal(%#v, %#v, xtime.Comparer) = true, want false", tt, tu)
		}
		if tu := tt.Add(1 * time.Minute); cmp.Equal(tt, tu, xtime.Comparer) {
			t.Errorf("cmp.Equal(%#v, %#v, xtime.Comparer) = true, want false", tt, tu)
		}
		if tu := tt.Add(1 * time.Second); cmp.Equal(tt, tu, xtime.Comparer) {
			t.Errorf("cmp.Equal(%#v, %#v, xtime.Comparer) = true, want false", tt, tu)
		}
		if tu := tt.Add(1 * time.Nanosecond); cmp.Equal(tt, tu, xtime.Comparer) {
			t.Errorf("cmp.Equal(%#v, %#v, xtime.Comparer) = true, want false", tt, tu)
		}
		if tu := xtime.WallIn(tt, xtime.MustLoadLocation(t, "Europe/Helsinki")); cmp.Equal(tt, tu, xtime.Comparer) {
			t.Errorf("cmp.Equal(%#v, %#v, xtime.Comparer) = true, want false", tt, tu)
		}
	}
}
