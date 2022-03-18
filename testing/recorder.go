package testing

import "fmt"

// Recorder of mutations of a T.
//
// Any test code that may call r.Fatal should be wrapped in an r.Exec.
type Recorder struct {
	Failed         bool            // True if Fail() was called.
	WasFatal       bool            // True if Fatal() was called.
	FatalArguments []interface{}   // Arguments that Fatal() was called with, if WasFatal.
	CallsToHelper  int             // Number of times Helper() was called.
	Logs           [][]interface{} // Arguments that Log() was called with.
}

// Fail sets r.Failed to true.
func (r *Recorder) Fail() {
	r.Failed = true
}

// Fatal sets r.WasFatal to true, sets FatalArguments to arguments, then panics.
func (r *Recorder) Fatal(arguments ...interface{}) {
	r.WasFatal = true
	r.FatalArguments = arguments
	panic(fmt.Sprintf("%v", arguments))
}

// Helper adds 1 to CallsToHelper.
func (r *Recorder) Helper() {
	r.CallsToHelper++
}

// Log adds arguments to r.Logs.
func (r *Recorder) Log(arguments ...interface{}) {
	r.Logs = append(r.Logs, arguments)
}

// Exec possiblyFatal, recovering a panic by Fatal.
func (r *Recorder) Exec(possiblyFatal func()) {
	defer func() {
		if a := recover(); a != nil {
			if !r.WasFatal {
				panic(a)
			}
		}
	}()
	possiblyFatal()
}
