package time

import (
	"time"

	"github.com/google/go-cmp/cmp"
)

// LocationComparer of *time.Locationsâ€™ String()s.
var LocationComparer = cmp.Comparer(func(a, b *time.Location) bool {
	return a.String() == b.String()
})
