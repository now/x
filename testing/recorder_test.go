package testing_test

import (
	"fmt"
	"testing"

	"github.com/google/go-cmp/cmp"

	xtesting "github.com/now/x/testing"
)

func TestRecorderFailed(t *testing.T) {
	var r xtesting.Recorder
	r.Fail()
	if !r.Failed {
		t.Errorf("r.Failed = false, want true")
	}
}

func TestRecorderFatal(t *testing.T) {
	var r xtesting.Recorder
	err := fmt.Errorf("failed")
	r.Exec(func() {
		r.Fatal(err)
	})
	if !r.WasFatal {
		t.Error("r.WasFatal = false, want true after rt.Fatal(…)")
	} else if want := []interface{}{err}; !cmp.Equal(r.FatalArguments, want, cmp.Comparer(func(a, b error) bool {
		return a.Error() == b.Error()
	})) {
		t.Errorf("r.FatalArguments = %#v, want %#v", r.FatalArguments, want)
	}
}

func TestRecorderHelper(t *testing.T) {
	var r xtesting.Recorder
	r.Helper()
	r.Helper()
	if r.CallsToHelper != 2 {
		t.Errorf("r.CallsToHelper = %d, want 2", r.CallsToHelper)
	}
}

func TestRecorderLog(t *testing.T) {
	var r xtesting.Recorder
	r.Log("a", "b", "c")
	r.Log(1, 2, 3)
	if diff := cmp.Diff(r.Logs, [][]interface{}{
		{"a", "b", "c"},
		{1, 2, 3},
	}); diff != "" {
		t.Errorf("r.Logs diff -got +want\n%s", diff)
	}
}

func TestRecorderDo(t *testing.T) {
	var r xtesting.Recorder
	defer func() {
		if a := recover(); a == nil {
			t.Errorf("rt.Do(…) caught panic, want panic to propagate")
		}
	}()
	r.Exec(func() {
		panic("not our panic")
	})
}
