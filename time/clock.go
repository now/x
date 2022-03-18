package time

import "time"

// Clock lets you know the current time.Time.
type Clock interface {
	// Now is the current time.Time.
	Now() time.Time
}
