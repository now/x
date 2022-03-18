package value

// Reflect log.Value that’ll write Value.
type Reflect struct {
	Value interface{}
}

// Write w.Reflect(r.Value).
func (r Reflect) Write(w Writer) error {
	return w.Reflect(r.Value)
}
