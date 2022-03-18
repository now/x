package log_test

import (
	"context"
	"fmt"
	"testing"

	"github.com/google/go-cmp/cmp"

	"github.com/now/x/log"
	xtesting "github.com/now/x/testing"
)

func TestEntry(t *testing.T) {
	var r xtesting.Recorder
	log.Entry(log.Testing(context.Background(), &r), "abc")
	if diff := cmp.Diff(r.Logs, [][]interface{}{{"abc"}}); diff != "" {
		t.Errorf("log.Entry(…, \"abc\") diff -got +want\n%s", diff)
	}
}

func TestNamed(t *testing.T) {
	var r xtesting.Recorder
	log.Entry(log.Named(log.Testing(context.Background(), &r), "a"), "bc")
	if diff := cmp.Diff(r.Logs, [][]interface{}{{"a: bc"}}); diff != "" {
		t.Errorf("log.Entry(log.Named(…, \"a\"), \"bc\") diff -got +want\n%s", diff)
	}
}

func TestWith(t *testing.T) {
	var r xtesting.Recorder
	log.Entry(log.With(log.Testing(context.Background(), &r), log.Int64("a", 1)), "bc")
	if diff := cmp.Diff(r.Logs, [][]interface{}{{"bc\na: 1"}}); diff != "" {
		t.Errorf("log.Entry(log.With(log.Testing(…), log.Int64(\"a\", 1)), \"bc\") diff -got +want\n%s", diff)
	}
}

func TestNop(t *testing.T) {
	if err := log.Entry(log.Nop(context.Background()), "abc"); err != nil {
		t.Errorf("log.Entry(log.Nop(…), \"abc\") = %v, want nil", err)
	}

	if err := log.In(log.Nop(context.Background())).Named("a").Entry("bc"); err != nil {
		t.Errorf("log.In(log.Nop(…)).Named(\"a\").Entry(\"bc\") = %v, want nil", err)
	}

	if err := log.In(log.Nop(context.Background())).With(log.Int64("a", 1)).Entry("bc"); err != nil {
		t.Errorf("log.In(log.Nop(…)).Named(\"a\").Entry(\"bc\") = %v, want nil", err)
	}
}

func TestTesting(t *testing.T) {
	t.Run("calls t.Helper() and writes to t.Log()", func(t *testing.T) {
		var r xtesting.Recorder
		log.Entry(log.Testing(context.Background(), &r), "abc")
		if r.CallsToHelper != 1 {
			t.Errorf("log.Entry(log.Testing(…), …) didn’t call t.Helper() once")
		}
		if diff := cmp.Diff(r.Logs, [][]interface{}{
			{"abc"},
		}); diff != "" {
			t.Errorf("log.Entry(log.Testing(…), …) diff -got +want\n%s", diff)
		}
	})

	t.Run("t.Fail()s on errors", func(t *testing.T) {
		var r xtesting.Recorder
		log.Entry(log.Testing(context.Background(), &r), "abc", log.Stringer("s", stringerPanicker{}))
		if !r.Failed {
			t.Errorf("log.Entry(log.Testing(…), …) didn’t call t.Fail() on error")
		} else if diff := cmp.Diff(r.Logs, [][]interface{}{
			{"write error: PANIC=oh, no!"},
		}); diff != "" {
			t.Errorf("log.Entry(log.Testing(…), …) diff -got +want\n%s", diff)
		}
	})

	t.Run("t.Fail()s on errors in own fields", func(t *testing.T) {
		var r xtesting.Recorder
		log.In(log.Testing(context.Background(), &r)).With(log.Stringer("s", stringerPanicker{})).Entry("abc")
		if !r.Failed {
			t.Errorf("log.Entry(log.Testing(…), …) didn’t call t.Fail() on error")
		} else if diff := cmp.Diff(r.Logs, [][]interface{}{
			{"write error: PANIC=oh, no!"},
		}); diff != "" {
			t.Errorf("log.Entry(log.Testing(…), …) diff -got +want\n%s", diff)
		}
	})

	t.Run("t.Fail()s on errors in parent’s fields", func(t *testing.T) {
		var r xtesting.Recorder
		log.In(log.Testing(context.Background(), &r)).With(log.Stringer("s", stringerPanicker{})).With(log.Int64("a", 1)).Entry("bc")
		if !r.Failed {
			t.Errorf("log.Entry(log.Testing(…), …) didn’t call t.Fail() on error")
		} else if diff := cmp.Diff(r.Logs, [][]interface{}{
			{"write error: PANIC=oh, no!"},
		}); diff != "" {
			t.Errorf("log.Entry(log.Testing(…), …) diff -got +want\n%s", diff)
		}
	})

	t.Run("formats entries", func(t *testing.T) {
		tests := []struct {
			field log.Field
			want  string
		}{
			{log.Int64("ID", 1), "ID: 1"},
			{log.Int64("neg", -1), "neg: -1"},
			{log.Reflect("values", map[string]int{"a": 1}), "values: map[a:1]"},
			{log.String("name", "something"), "name: something"},
			{log.String("lines", "a\nb"), "lines: a\n       b"},
			{
				log.String(
					"control",
					"\x00\x01\x02\x03\x04\x05\x06\a\b\t\v\f\r\x0E\x0F\x10\x11\x12\x13\x14\x15\x16\x17\x18\x19\x1A\x1B\x1C\x1D\x1E\x1F",
				),
				"control: \u2400\u2401\u2402\u2403\u2404\u2405\u2406\u2407\u2408\u2409\u240b\u240c\u240d\u240e\u240f\u2410\u2411\u2412\u2413\u2414\u2415\u2416\u2417\u2418\u2419\u241a\u241b\u241c\u241d\u241e\u241f",
			},
		}
		for _, tt := range tests {
			var r xtesting.Recorder
			message := "abc"
			log.Entry(log.Testing(context.Background(), &r), message, tt.field)
			if diff := cmp.Diff(r.Logs, [][]interface{}{
				{fmt.Sprintf("%s\n%s", message, tt.want)},
			}); diff != "" {
				t.Errorf("log.Entry(…, %#v, %#v) diff -got +want\n%s", message, tt.field, diff)
			}
		}

		t.Run("one level of naming", func(t *testing.T) {
			tests := []struct {
				name, want string
			}{
				{"", "bc"},
				{"a", "a: bc"},
			}
			for _, tt := range tests {
				var r xtesting.Recorder
				log.In(log.Testing(context.Background(), &r)).Named(tt.name).Entry("bc")
				if diff := cmp.Diff(r.Logs, [][]interface{}{{tt.want}}); diff != "" {
					t.Errorf("log.In(log.Testing(…)).Named(%#v).Entry(\"bc\") diff -got +want\n%s", tt.name, diff)
				}
			}
		})

		t.Run("multiple levels of naming", func(t *testing.T) {
			var r xtesting.Recorder
			log.In(log.Testing(context.Background(), &r)).Named("a").Named("b").Entry("c")
			if diff := cmp.Diff(r.Logs, [][]interface{}{{"a.b: c"}}); diff != "" {
				t.Errorf("log.In(log.Testing(…)).Named(\"a\").Named(\"b\").Entry(\"c\") diff -got +want\n%s", diff)
			}
		})

		t.Run("with empty", func(t *testing.T) {
			var r xtesting.Recorder
			log.In(log.Testing(context.Background(), &r)).With().Entry("abc")
			if diff := cmp.Diff(r.Logs, [][]interface{}{{"abc"}}); diff != "" {
				t.Errorf("log.In(log.Testing(…)).With().Entry(\"abc\") diff -got +want\n%s", diff)
			}
		})

		t.Run("with fields from Logger and it’s parents", func(t *testing.T) {
			var r xtesting.Recorder
			log.In(log.Testing(context.Background(), &r)).With(log.Int64("a", 1)).With(log.Int64("b", 2)).Entry("c", log.Int64("d", 3))
			if diff := cmp.Diff(r.Logs, [][]interface{}{{"c\na: 1\nb: 2\nd: 3"}}); diff != "" {
				t.Errorf("log.In(log.Testing(…)).With(log.Int64(\"a\", 1)).With(log.Int64(\"b\", 2)).Entry(\"c\", log.Int64(\"d\", 3)) diff -got +want\n%s", diff)
			}
		})
	})
}

type stringerPanicker struct{}

func (stringerPanicker) String() string {
	panic("oh, no!")
}
