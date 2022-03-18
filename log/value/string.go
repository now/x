package value

// String log.Value that’ll write itself.
type String string

// Write w.String(string(s))
func (s String) Write(w Writer) error {
	return w.String(string(s))
}
