package time

import (
	"context"
	"time"
)

// Default is ctxʹ ≈ ctx such that In(ctxʹ) is time.Now().
func Default(ctx context.Context) context.Context {
	return Using(ctx, ClockFunc(time.Now))
}

// Stopped is ctxʹ ≈ ctx such that In(ctxʹ) is t.
//
// This is useful for testing, where time is often necessary to control.
func Stopped(ctx context.Context, t time.Time) context.Context {
	return Using(ctx, ClockFunc(func() time.Time { return t }))
}

// Using is ctxʹ ≈ ctx such that In(ctxʹ) is c.Now().
func Using(ctx context.Context, c Clock) context.Context {
	return context.WithValue(ctx, Key, c)
}

// In is c.Now(), where c is the Clock previously added to ctx.
//
// The Clock c can be added to ctx with Default(ctx), Stopped(ctx), or
// Using(ctx, c).
func In(ctx context.Context) time.Time {
	return ctx.Value(Key).(Clock).Now()
}

// Key of Clock in context.Context.
const Key key = "Clock"

type key string
