package value

// Int log.Value that’ll write itself.
type Int int

// Write w.Int(int(i)).
func (i Int) Write(w Writer) error {
	return w.Int(int(i))
}
