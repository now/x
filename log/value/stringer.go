package value

import "fmt"

// Stringer log.Value that’ll write Value.String().
type Stringer struct {
	Value fmt.Stringer
}

// Write w.String(s.Value.String()).
//
// Any panic caused by s.Value.String() is caught.  A panic caused by s.Value
// being nil will result in w.String("<nil>") instead.  Errors on any other
// panic.
func (s Stringer) Write(w Writer) (err error) {
	// Capture panics from s.Value.String().
	defer func() {
		if rerr := recover(); rerr != nil {
			err = panicNilCheck(rerr, s.Value, w)
		}
	}()
	return w.String(s.Value.String())
}
