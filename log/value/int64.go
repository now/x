package value

// Int64 log.Value that’ll write itself.
type Int64 int64

// Write w.Int64(int64(i)).
func (i Int64) Write(w Writer) error {
	return w.Int64(int64(i))
}
