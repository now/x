package zap_test

import (
	"bytes"
	"fmt"
	"testing"

	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"

	"github.com/now/x/log"
	xzap "github.com/now/x/log/zap"
)

func TestLoggerEntry(t *testing.T) {
	t.Run("delegates to l.Zap.Info(…)", func(t *testing.T) {
		tests := []string{
			"abc",
		}
		for _, tt := range tests {
			var b buffer
			logger(&b).Entry(tt)
			if got, want := string(b.Bytes()), fmt.Sprintf("info\t%s", tt); got != want {
				t.Errorf("zap.Logger{…}.Entry(%#v) = %#v, want %#v", tt, got, want)
			}
		}
	})

	t.Run("maps log.Fields to zap.Fields", func(t *testing.T) {
		tests := []struct {
			field log.Field
			want  string
		}{
			{log.Error(fmt.Errorf("failed")), `"error": "failed"`},
			{log.Int64("ID", 1), `"ID": 1`},
			{log.Int64("neg", -1), `"neg": -1`},
			{log.Reflect("values", map[string]int{"a": 1}), `"values": {"a":1}`},
			{log.String("name", "something"), `"name": "something"`},
			{log.Stringer("string", stringer("stringed")), `"string": "stringed"`},
		}
		for _, tt := range tests {
			var b buffer
			logger(&b).Entry("abc", tt.field)
			if got, want := string(b.Bytes()), fmt.Sprintf("info\tabc\t{%s}", tt.want); got != want {
				t.Errorf("zap.Logger{…}.Entry(\"abc\", %#v) = %#v, want %#v", tt.field, got, want)
			}
		}
	})
}

func TestLoggerNamed(t *testing.T) {
	var b buffer
	logger(&b).Named("a").Entry("bc")
	if got, want := string(b.Bytes()), "info\ta\tbc"; got != want {
		t.Errorf(`zap.Logger{…}.Name("a").Entry("bc") = %#v, want %#v`, got, want)
	}
}

func TestLoggerWith(t *testing.T) {
	var b buffer
	logger(&b).With(log.Int64("ID", 1)).Entry("bc")
	if got, want := string(b.Bytes()), "info\tbc\t{\"ID\": 1}"; got != want {
		t.Errorf("zap.Logger{…}.With(log.Int64(\"ID\", 1)).Entry(\"bc\") = %#v, want %#v", got, want)
	}
}

type buffer struct {
	bytes.Buffer
}

func (b *buffer) Sync() error {
	return nil
}

func logger(b *buffer) xzap.Logger {
	return xzap.Logger{
		Zap: zap.New(
			zapcore.NewCore(
				zapcore.NewConsoleEncoder(
					zapcore.EncoderConfig{
						MessageKey:     "M",
						LevelKey:       "L",
						NameKey:        "N",
						SkipLineEnding: true,
						LineEnding:     "\n",
						EncodeLevel:    zapcore.LowercaseLevelEncoder,
						EncodeName:     zapcore.FullNameEncoder,
					},
				),
				b,
				zap.LevelEnablerFunc(func(zapcore.Level) bool {
					return true
				}),
			),
		),
	}
}

type stringer string

func (s stringer) String() string {
	return string(s)
}
