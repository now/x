package value

// Error log.Value that’ll write Err.Error().
type Error struct {
	Err error
}

// Write w.String(e.Err.Error()).
//
// Any panic caused by e.Err.Error() is caught.  A panic caused by e.Err being
// nil will result in w.String("<nil>") instead.  Errors on any other panic.
func (e Error) Write(w Writer) (err error) {
	// Capture panics from e.Err.Error().
	defer func() {
		if rerr := recover(); rerr != nil {
			err = panicNilCheck(rerr, e.Err, w)
		}
	}()
	w.String(e.Err.Error())
	return
}
