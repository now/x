package value

import (
	"fmt"
	"reflect"
)

// Writer of typed data.
//
// This and the log.Value interface comprises a double dispatch for log.Fields
// to write themselves in a typed manner.
type Writer interface {
	Int64(int64) error         // Int64 writes an int64.
	Reflect(interface{}) error // Reflect writes any value using reflection.
	String(string) error       // String writes a string.

	// Field writes a label and then calls a function that’ll write a value.
	//
	// This is intended to be used by log.Field.  Taking a function argument
	// allows this method to do any necessary work before and/or after of any
	// writes that the function will perform.
	Field(string, func(Writer) error) error
}

func panicNilCheck(err interface{}, a interface{}, w Writer) error {
	if v := reflect.ValueOf(a); a == nil || v.Kind() == reflect.Ptr && v.IsNil() {
		return w.String("<nil>")
	}
	return fmt.Errorf("PANIC=%v", err)
}
