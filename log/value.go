package log

import "github.com/now/x/log/value"

// Value that can write itself to a value.Writer.
type Value interface {
	// Write receiver to the given value.Writer.
	//
	// The receiver may use any of the methods provided by the value.Writer
	// interface to provide a marshaling of itself.
	Write(value.Writer) error
}
