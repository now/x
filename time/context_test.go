package time_test

import (
	"context"
	"testing"
	"time"

	xtime "github.com/now/x/time"
)

func TestDefault(t *testing.T) {
	t.Run("is time.Now()", func(t *testing.T) {
		got := xtime.In(xtime.Default(context.Background()))
		want := time.Now()
		if want.Sub(got).Round(time.Second).Seconds() > 0 {
			t.Errorf("xtime.In(xtime.Default(…)) = %v, want ≈ %v", got, want)
		}
	})

	t.Run("rest of ctx remains the same", func(t *testing.T) {
		type keytype string
		key := keytype("test")
		want := 1
		if got, _ := xtime.Default(context.WithValue(context.Background(), key, want)).Value(key).(int); got != want {
			t.Errorf("xtime.Default(context.WithValue(…, %#v)) = %#v, want %#v", want, got, want)
		}
	})
}

func TestStopped(t *testing.T) {
	t.Run("is t", func(t *testing.T) {
		want := time.Now()
		ctx := xtime.Stopped(context.Background(), want)
		for i := 0; i < 3; i++ {
			if got := xtime.In(ctx); !got.Equal(want) {
				t.Fatalf("xtime.In(xtime.Stopped(…, %#v)) = %#v, want %#v", want, got, want)
			}
		}
	})

	t.Run("rest of ctx remains the same", func(t *testing.T) {
		type keytype string
		key := keytype("test")
		want := 1
		if got, _ := xtime.Stopped(context.WithValue(context.Background(), key, want), time.Now()).Value(key).(int); got != want {
			t.Errorf("xtime.Stopped(context.WithValue(…, %#v)) = %#v, want %#v", want, got, want)
		}
	})
}
