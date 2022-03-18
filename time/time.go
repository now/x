// Package time contains extension for working with the time package.
//
// The Clock interface abstracts the concept of providing the current time.Time.
// This is used in three functions that set the Clock used in a context.Context
// ctx: Default(ctx), Stopped(ctx, time.Time), and Using(ctx, Clock).  This is
// useful for controlling the source of time used throughout a context.Context,
// for testing purposes or similar.  Getting the current time.Time out of ctx is
// done with In(ctx).
//
// WallIn(t, l) allows you to move a time.Time t to a *time.Location l without
// changing the actual wall time, which is semantically different from how
// t.In(l) works.
//
// LoadLocation(code) works exactly like time.LoadLocation(code), but caches the
// result, improving performance immensely.
//
// For testing, MustLoadLocation(t, code) calls LoadLocation(code) and, if that
// fails, fatals the testing.T t.  This makes testing code easier to read and
// write, as error checks can be removed.  Similarly, MustDate(t, …, code) is
// exactly the same as time.Date(…, code), but uses MustLoadLocation(t, code)
// and thus fatals the testing.T t if the location can’t be loaded.
//
// Also for testing, Comparer and LocationComparer allows for comparing
// time.Times and *time.Locations with the github.com/google/go-cmp library for
// slightly more precise results than time.Time.Equal provides.
package time

import (
	"time"

	"github.com/now/x/testing"
)

// MustDate is the time.Time at year, month, day, hour, minute, second, and nanosecond in the location of locationCode.
//
// Registers as a t.Helper().
//
// Fatals t if the location of locationCode can’t be loaded.
func MustDate(t testing.T, year int, month time.Month, day, hour, minute, second, nanosecond int, locationCode string) time.Time {
	t.Helper()
	return time.Date(year, month, day, hour, minute, second, nanosecond, MustLoadLocation(t, locationCode))
}

// WallIn is tʹ ≈ t such that tʹ.Location() is l.
//
// That is, Year(), Month(), Day(), Hour(), Minute(), Second(), and Nanosecond()
// will all be equal for t and tʹ, which isn’t necessarily true for t and
// t.In(l).
func WallIn(t time.Time, l *time.Location) time.Time {
	return time.Date(t.Year(), t.Month(), t.Day(), t.Hour(), t.Minute(), t.Second(), t.Nanosecond(), l)
}
