package time_test

import (
	"testing"
	"time"

	"github.com/google/go-cmp/cmp"

	xtime "github.com/now/x/time"
)

func TestLocationComparer(t *testing.T) {
	tests := []string{
		"Europe/Stockholm",
	}
	for _, tt := range tests {
		if l1, err := time.LoadLocation(tt); err != nil {
			t.Errorf("time.LoadLocation(%q) = %v, want time.Location(%q)", tt, err, tt)
		} else if l2, err := time.LoadLocation(tt); err != nil {
			t.Errorf("time.LoadLocation(%q) = %v, want time.Location(%q)", tt, err, tt)
		} else if l1.String() != l2.String() {
			t.Errorf("time.LoadLocation(%q).String() == time.LoadLocation(%q).String() = false, want true", tt, tt)
		} else if l1 == l2 {
			t.Errorf("time.LoadLocation(%q) == time.LoadLocation(%q) = true, want false", tt, tt)
		} else if !cmp.Equal(l1, l2, xtime.LocationComparer) {
			t.Errorf("cmp.Equal(%v, %v, xtime.LocationComparer) = false, want true", l1, l2)
		}
	}
}
