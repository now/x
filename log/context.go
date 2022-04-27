// Package log comprises a minimal structured logging framework.
//
// The package also contains a couple of loggers built upon this framework.
// Given it’s minimal, simple, and straightforward interface, additional loggers
// can readily be implemented.
//
// The origin type is Logger, which is an interface consisting of a method for
// creating a log entry, a method for naming a new Logger based on the receiver,
// and a method for providing additional meta-data in a new Logger based on the
// receiver.  Each meta-datum, represented by the Field type, is a pair of a
// string label and a value that implements the Value interface.  Values know
// how to write themselves to a value.Writer.  A value.Writer knows how to
// marshal data expressed as various types, such as int64 and string, as well as
// how to marshal Fields.  These four types, Logger, Field, Value, and
// value.Writer make up the whole of the log package’s framework.
//
// To expand a bit on the function of this framework, log entries are created
// given a string message along with any additional meta-data coming from the
// receiver and from the method call itself.  The meta-data thus provides
// structure to log entries and can be used for tasks such as filtering and
// similar.  It’s use is not governed at by the framework in any way, but can be
// used to implement log entry severity levels, additional context to a message,
// and so on.  (As such, there’s no built-in support for log entry severity
// levels, both to keep the interface minimal, but also because it’s use was
// deemed to be a poor practice.)  The type that implements the Logger interface
// along with a value.Writer that it controls, is fully in charge of how the
// message and any meta-data is marshaled.  It’s also free to add additional
// meta-data not provided explicitly by the user by invoking any of the three
// methods afforded by the Logger interface.
//
// As such, much of the specifics are left to the types that implement the
// Logger interface.  There are two implementations of the Logger interface in
// the log package that are both useful and provide some insight into how more
// Loggers can be built.  But before we get to the specifics of those two
// Loggers, the auxiliary interface of the log package needs to be introduced.
//
// The auxiliary interface of the log package consists of three sets of
// functions: one to create and/or add Loggers to a context.Context (referred to
// as “Context” hereafter); one to access the Logger stored in a Context; and
// one to interact with the Logger stored in a Context directly.  All three sets
// thus work with Contexts.  The reason behind this design is that logging is a
// cross-cutting concern, thus well suited for inclusion in a Context.  A
// function or method that receives a Context can thus easily perform logging
// without the user knowing or caring.
//
// The first set of functions of the auxiliary interface creates and/or adds
// Loggers to a Context consists of three functions, Nop(), Testing(), and
// Using().  Nop() creates and adds a Logger to a Context that does nothing.  No
// log entry will be written anywhere.  This is ideal if you want to, for
// example, disable logging for a specific sub-system.  Remember, Contexts are
// values, so passing a Context to Nop() won’t change the Context, instead
// returning a new one.  Testing() creates and adds a Logger to a Context that
// works well when testing with the testing package.  It formats meta-data in a
// way that fits well with any other output produced during testing and, as it
// uses, testing.T.Log() for displaying log entries, log entries will only be
// displayed if the test fails or if testing.Verbose is true, which means that
// logging can facilitate testing without getting in the way.  Finally, Using()
// allows for an already created Logger to be added to a Context.
//
// The second set of functions of the auxiliary interface consists of one
// function, In(), that accesses the Logger previously added to a Context for
// further interaction.
//
// The third set of functions of the auxiliary interface consists of three
// functions, namely Entry(), Named(), and With().  These functions interact
// with the Logger stored in a Context directly.  Entry() accesses the Logger
// and creates a log entry with it.  Named() returns a new Context based on a
// given Context with a new Logger resulting from naming the existing Logger in
// the the given Context.  With() returns a new Context based on a given Context
// with a new Logger resulting from providing additional meta-data to the
// existing Logger in the given Context.
package log

import (
	"context"
	"fmt"
	"strconv"
	"sync"
	"unicode/utf8"

	"github.com/now/x/log/value"
	"github.com/now/x/testing"
)

// With is ctxʹ ≈ ctx such that In(ctxʹ) is a Logger that doesn’t log anything.
func Nop(ctx context.Context) context.Context {
	return Using(ctx, noper)
}

// Testing is ctxʹ ≈ ctx such that In(ctxʹ) is a Logger for use during testing in t.
//
// Invoking Entry(message, ...fields) on this Logger will call t.Helper(), then
// write an entry to t.Log().  If an error occurs during this, "write error: ",
// followed by the error’s Error(), is logged to t.Log() and then t.Fail() is
// called.
//
// The entry is formatted as [NAME ": "] MESSAGE, where NAME is the name of the
// Logger.  MESSAGE is formatted as a value.String.
//
// This is followed by zero or more fields, beginning with those added to the
// Logger, then the given fields.  The format of each field is "\n" LABEL ": "
// VALUE.  The format of VALUE depends on its type.
//
// A value.Int64 is formatted as an integer in base ten.  Negative values are
// prefixed by a hyphen-minus, U+002D.
//
// A value.Reflect’s Value r is replaced by s = fmt.Sprintf("%+v", r) and s is
// formatted as a value.String.
//
// A value.String is formatted as is, except that line feed, U+000A, is followed
// by i spaces, U+0020, where i = len(LABEL) + 2, if there’s a LABEL, i = 0,
// otherwise, and any rune c below U+0020 except for line feed, U+000A, is
// replaced by U+2400 + codepoint(c).
func Testing(ctx context.Context, t testing.T) context.Context {
	return Using(ctx, &testingLogger{t: t})
}

// Using is ctxʹ ≈ ctx such that In(ctxʹ) = l.
func Using(ctx context.Context, l Logger) context.Context {
	return context.WithValue(ctx, Key, l)
}

// In is the Logger in ctx.
//
// In is nil if no Logger has been added to ctx with Nop(ctx), Testing(ctx,
// testing.T), or Using(ctx, Logger).
func In(ctx context.Context) Logger {
	return ctx.Value(Key).(Logger)
}

// Entry consisting of message and fields is added to In(ctx).
//
// Delegates to In(ctx).Entry(message, fields...).
//
// Errors if Logger.Entry(message, fields...) errors.
func Entry(ctx context.Context, message string, fields ...Field) error {
	return In(ctx).Entry(message, fields...)
}

// Named is ctxʹ ≈ ctx such that In(ctxʹ) = l.Named(name), where l = In(ctx).
func Named(ctx context.Context, name string) context.Context {
	return Using(ctx, In(ctx).Named(name))
}

// With is ctxʹ ≈ ctx such that In(ctxʹ) = l.With(fields...), where l = In(ctx).
func With(ctx context.Context, fields ...Field) context.Context {
	return Using(ctx, In(ctx).With(fields...))
}

// Key of Logger in context.Context.
var Key = key{}

type key struct{}

func (key) GoString() string {
	return "log.Key"
}

func (k key) String() string {
	return k.GoString()
}

var noper = nopLogger{}

type nopLogger struct{}

// Entry adds nothing.
func (nopLogger) Entry(string, ...Field) error {
	return nil
}

// Named is n.
func (n nopLogger) Named(string) Logger {
	return n
}

// With is n.
func (n nopLogger) With(...Field) Logger {
	return n
}

type testingLogger struct {
	t      testing.T
	name   string
	parent *testingLogger
	fields []Field
}

func (t *testingLogger) Entry(message string, fields ...Field) (err error) {
	t.t.Helper()

	w := testingWriters.Get().(*testingWriter)
	defer func() {
		w.b = w.b[:0]
		w.indention = 0
		w.separate = false
		testingWriters.Put(w)
		if err != nil {
			t.t.Log(fmt.Sprintf("write error: %v", err))
			t.t.Fail()
		}
	}()

	if t.entryName(w) {
		w.bytes(": ")
	}

	w.string(message)

	if err := t.entryFields(w); err != nil {
		return err
	}

	for i := range fields {
		if err := fields[i].Write(w); err != nil {
			return err
		}
	}

	t.t.Log(string(w.b))

	return nil
}

func (t *testingLogger) entryName(w *testingWriter) bool {
	if t.parent != nil && t.parent.entryName(w) {
		w.byte('.')
	}

	if len(t.name) == 0 {
		return false
	}

	w.string(t.name)
	return true
}

func (t *testingLogger) entryFields(w *testingWriter) error {
	if t.parent != nil {
		if err := t.parent.entryFields(w); err != nil {
			return err
		}
	}

	for i := range t.fields {
		if err := t.fields[i].Write(w); err != nil {
			return err
		}
	}

	return nil
}

func (t *testingLogger) Named(name string) Logger {
	if name == "" {
		return t
	}
	c := t.clone()
	c.name = name
	return c
}

func (t *testingLogger) With(fields ...Field) Logger {
	if len(fields) == 0 {
		return t
	}
	c := t.clone()
	c.fields = fields
	return c
}

func (t *testingLogger) clone() *testingLogger {
	c := *t
	c.parent = t
	return &c
}

var testingWriters = sync.Pool{
	New: func() interface{} {
		return &testingWriter{b: make([]byte, 0, 1024)}
	},
}

type testingWriter struct {
	b         []byte
	indention int
	separate  bool
}

func (w *testingWriter) Int(i int) error {
	return w.Int64(int64(i))
}

func (w *testingWriter) Int64(i int64) error {
	w.separator()
	w.b = strconv.AppendInt(w.b, i, 10)
	return nil
}

func (w *testingWriter) Reflect(r interface{}) error {
	return w.String(fmt.Sprintf("%+v", r))
}

func (w *testingWriter) String(s string) error {
	w.separator()
	w.string(s)
	return nil
}

func (w *testingWriter) Field(label string, f func(value.Writer) error) error {
	w.lineFeed()
	w.string(label)
	w.bytes(": ")
	w.separate = false
	n := len(label) + 2
	w.indention += n
	err := f(w)
	w.indention -= n
	return err
}

func (w *testingWriter) separator() {
	if w.separate {
		w.bytes(", ")
	} else {
		w.separate = true
	}
}

func (w *testingWriter) byte(c byte) {
	w.b = append(w.b, c)
}

func (w *testingWriter) string(s string) {
	var i, j int
	for ; j < len(s); j++ {
		if c := s[j]; c < 0x20 {
			if i < j {
				w.bytes(s[i:j])
			}
			i = j + 1

			if c == '\n' {
				w.lineFeed()
			} else {
				bytes := make([]byte, utf8.UTFMax)
				n := utf8.EncodeRune(bytes, rune(0x2400)+rune(c))
				w.bytes(string(bytes[:n]))
			}
		}
	}
	if i < j {
		w.bytes(s[i:j])
	}
}

func (w *testingWriter) lineFeed() {
	w.byte('\n')
	for i := 0; i < w.indention; i++ {
		w.byte(' ')
	}
}

func (w *testingWriter) bytes(s string) {
	w.b = append(w.b, s...)
}
