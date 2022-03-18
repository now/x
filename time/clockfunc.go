package time

import "time"

// ClockFunc for functions that implement Clock.
type ClockFunc func() time.Time

// Now is f().
func (f ClockFunc) Now() time.Time {
	return f()
}
