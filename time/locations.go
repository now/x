package time

import (
	"sync"
	"time"

	"github.com/now/x/testing"
)

// LoadLocation with code.
//
// This is exactly like time.LoadLocation(code), except that results are cached.
func LoadLocation(code string) (*time.Location, error) {
	if l, found := locations.Load(code); found {
		return l.(*time.Location), nil
	} else if l, err := time.LoadLocation(code); err != nil {
		return nil, err
	} else {
		s, _ := locations.LoadOrStore(code, l)
		return s.(*time.Location), nil
	}
}

// MustLoadLocation is LoadLocation(code), but registers as a t.Helper() and fatals t if that errors.
func MustLoadLocation(t testing.T, code string) *time.Location {
	if l, err := LoadLocation(code); err != nil {
		t.Helper()
		t.Fatal(err)
		return nil
	} else {
		return l
	}
}

var locations = &sync.Map{}
