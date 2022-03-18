package time

import (
	"time"

	"github.com/google/go-cmp/cmp"
)

// Comparer of time.Times Date(), Clock(), Nanosecond(), and Location().String()s.
var Comparer = cmp.Comparer(func(a, b time.Time) bool {
	{
		aYear, aMonth, aDay := a.Date()
		bYear, bMonth, bDay := b.Date()
		if !(aYear == bYear && aMonth == bMonth && aDay == bDay) {
			return false
		}
	}

	{
		aHour, aMin, aSec := a.Clock()
		bHour, bMin, bSec := b.Clock()
		if !(aHour == bHour && aMin == bMin && aSec == bSec) {
			return false
		}
	}

	return a.Nanosecond() == b.Nanosecond() &&
		a.Location().String() == b.Location().String()
})
