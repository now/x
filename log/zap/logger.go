package zap

import (
	"go.uber.org/zap"

	"github.com/now/x/log"
	"github.com/now/x/log/value"
)

// ZapLogger delegates to a *zap.Logger.
//
// A log.Field{Label: l, Value: v} is mapped to a zap.Field z as follows:
//
// If v = value.Error{Err: err}, z = zap.Error(l, err).
//
// If v = value.Int64(i), z = zap.Int64(l, i).
//
// If v = value.Reflect{Value: r}, z = zap.Reflect(l, r).
//
// If v = value.String(s), z = zap.String(l, s).
//
// If v = value.Stringer{Value: s}, z = zap.Stringer(l, s).
//
// Otherwise, which shouldn’t happen, z = zap.Any(l, v).
type Logger struct {
	Zap *zap.Logger
}

// Entry delegates to l.Zap.Info(message, fields...).
func (l Logger) Entry(message string, fields ...log.Field) error {
	l.Zap.Info(message, zapFields(fields...)...)
	return nil
}

// Named is a new Logger wrapping l.Zap.Named(name).
func (l Logger) Named(name string) log.Logger {
	return Logger{l.Zap.Named(name)}
}

// With is a new Logger wrapping l.Zap.With(fields...).
func (l Logger) With(fields ...log.Field) log.Logger {
	return Logger{l.Zap.With(zapFields(fields...)...)}
}

func zapFields(fields ...log.Field) []zap.Field {
	zapFields := make([]zap.Field, len(fields))
	for i, f := range fields {
		switch v := f.Value.(type) {
		case value.Error:
			zapFields[i] = zap.Error(v.Err)
		case value.Int:
			zapFields[i] = zap.Int(f.Label, int(v))
		case value.Int64:
			zapFields[i] = zap.Int64(f.Label, int64(v))
		case value.Reflect:
			zapFields[i] = zap.Reflect(f.Label, v.Value)
		case value.String:
			zapFields[i] = zap.String(f.Label, string(v))
		case value.Stringer:
			zapFields[i] = zap.Stringer(f.Label, v.Value)
		default:
			zapFields[i] = zap.Any(f.Label, v)
		}
	}
	return zapFields
}
