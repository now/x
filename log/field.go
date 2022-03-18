package log

import (
	"fmt"

	"github.com/now/x/log/value"
)

// Field consisting of a Label and a Value.
//
// A Field is used to provide contextual meta-data to a log entry.  Fields can
// be added to a Logger.With(...Fields), so that any following log entries will
// include those fields, or added to a specific Logger.Entry(string, ...Fields).
//
// Field labels are, as the name suggests, not required to be unique.
type Field struct {
	Label string
	Value Value
}

// Error Field with Label “error” and value.Error{Err: err}.
func Error(err error) Field {
	return Field{"error", value.Error{Err: err}}
}

// Int Field with label and value.Int(i).
func Int(label string, i int) Field {
	return Field{label, value.Int(i)}
}

// Int64 Field with label and value.Int64(i).
func Int64(label string, i int64) Field {
	return Field{label, value.Int64(i)}
}

// Reflect Field with label and value.Reflect{Value: r}.
func Reflect(label string, r interface{}) Field {
	return Field{label, value.Reflect{Value: r}}
}

// String Field with label and value.String(s).
func String(label, s string) Field {
	return Field{label, value.String(s)}
}

// Stringer Field with label and value.Stringer{Value: s}.
func Stringer(label string, s fmt.Stringer) Field {
	return Field{label, value.Stringer{Value: s}}
}

// Write f to w using w.Field(f.Label, f.Value.Write).
func (f Field) Write(w value.Writer) error {
	return w.Field(f.Label, f.Value.Write)
}
