package time_test

import (
	"testing"
	"time"

	"github.com/google/go-cmp/cmp"

	xtesting "github.com/now/x/testing"
	xtime "github.com/now/x/time"
)

func TestLoadLocation(t *testing.T) {
	t.Run("valid", func(t *testing.T) {
		tests := []string{
			"Europe/Stockholm",
		}
		for _, tt := range tests {
			if want, err := time.LoadLocation(tt); err != nil {
				t.Errorf("time.LoadLocation(%q) = %v, want time.Location(%q)", tt, err, tt)
			} else if got, err := xtime.LoadLocation(tt); err != nil {
				t.Errorf("xtime.LoadLocation(%q) = %v, want %#v", tt, err, nil)
			} else if !cmp.Equal(got, want, xtime.LocationComparer) {
				t.Errorf("xtime.LoadLocation(%q) = %v, want %v", tt, got, want)
			}
		}
	})

	t.Run("invalid", func(t *testing.T) {
		tests := []string{
			"junk",
		}
		for _, tt := range tests {
			if got, err := xtime.LoadLocation(tt); err == nil {
				t.Errorf("xtime.LoadLocation(%q) = %#v, want err", tt, got)
			}
		}
	})
}

func TestMustLoadLocation(t *testing.T) {
	t.Run("valid", func(t *testing.T) {
		tests := []string{
			"Europe/Stockholm",
		}
		for _, tt := range tests {
			if want, err := time.LoadLocation(tt); err != nil {
				t.Errorf("time.LoadLocation(%q) = %v, want time.Location(%q)", tt, err, tt)
			} else {
				var rt xtesting.Recorder
				var got *time.Location
				rt.Exec(func() {
					got = xtime.MustLoadLocation(&rt, tt)
				})
				if rt.WasFatal {
					t.Errorf("xtime.MustLoadLocation(…, %q) was fatal with %v, want %#v", tt, rt.FatalArguments, want)
					if rt.CallsToHelper != 1 {
						t.Errorf("xtime.MustLoadLocation(…, %q) didn’t register 1 test helper function, but %d", tt, rt.CallsToHelper)
					}
				} else {
					if !cmp.Equal(got, want, xtime.LocationComparer) {
						t.Errorf("xtime.MustLoadLocation(…, %q) = %v, want %v", tt, got, want)
					}
					if rt.CallsToHelper != 0 {
						t.Errorf("xtime.MustLoadLocation(…, %q) registered %d test helper functions, but didn’t fatal", tt, rt.CallsToHelper)
					}
				}
			}
		}
	})

	t.Run("invalid", func(t *testing.T) {
		tests := []string{
			"junk",
		}
		for _, tt := range tests {
			var rt xtesting.Recorder
			var got *time.Location
			rt.Exec(func() {
				got = xtime.MustLoadLocation(&rt, tt)
			})
			if !rt.WasFatal {
				t.Errorf("xtime.MustLoadLocation(…, %q) = %#v, want fatal", tt, got)
			}
			if rt.CallsToHelper != 1 {
				t.Errorf("xtime.MustLoadLocation(…, %q) didn’t register as a test helper function 1 time, but %d", tt, rt.CallsToHelper)
			}
		}
	})
}

func BenchmarkTimeLoadLocation(b *testing.B) {
	for i := 0; i < b.N; i++ {
		time.LoadLocation("Europe/Stockholm")
	}
}

func BenchmarkLoadLocation(b *testing.B) {
	for i := 0; i < b.N; i++ {
		xtime.LoadLocation("Europe/Stockholm")
	}
}
