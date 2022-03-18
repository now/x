package value

import (
	"fmt"
	"strconv"
)

type BytesWriter struct {
	Bytes    []byte
	separate bool
}

func (w *BytesWriter) Int(i int) error {
	return w.Int64(int64(i))
}

func (w *BytesWriter) Int64(i int64) error {
	w.Bytes = strconv.AppendInt(w.Bytes, i, 10)
	return nil
}

func (w *BytesWriter) Reflect(r interface{}) error {
	return w.String(fmt.Sprintf("%+v", r))
}

func (w *BytesWriter) String(s string) error {
	w.separator()
	w.bytes([]byte(s))
	return nil
}

func (w *BytesWriter) Field(key string, f func(Writer) error) error {
	if w.separate {
		w.bytes([]byte("; "))
	}
	w.bytes([]byte(key))
	w.bytes([]byte(": "))
	w.separate = false
	return f(w)
}

func (w *BytesWriter) separator() {
	if w.separate {
		w.bytes([]byte(", "))
	} else {
		w.separate = true
	}
}

func (w *BytesWriter) bytes(b []byte) {
	w.Bytes = append(w.Bytes, b...)
}
