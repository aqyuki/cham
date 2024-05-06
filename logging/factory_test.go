package logging

import (
	"reflect"
	"testing"

	"github.com/google/go-cmp/cmp"
	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"
)

func TestNewLoggerFromEnv(t *testing.T) {
	t.Parallel()

	result := NewLoggerFromEnv()
	if result == nil {
		t.Error("expect not nil, but received nil")
	}
}

func TestNewLogger(t *testing.T) {
	t.Parallel()

	result1 := NewLogger(true, "debug")
	if result1 == nil {
		t.Errorf("expect not nil, but received nil")
	}

	result2 := NewLogger(false, "debug")
	if result2 == nil {
		t.Errorf("expect not nil, but received nil")
	}
}

func TestDefaultLogger(t *testing.T) {
	t.Parallel()

	result1 := DefaultLogger()
	if result1 == nil {
		t.Errorf("expect not nil, but received nil")
	}

	result2 := DefaultLogger()
	if result2 == nil {
		t.Errorf("expect not nil, but received nil")
	}

	if !reflect.DeepEqual(result1, result2) {
		t.Errorf("expect same logger, but received different logger")
	}
}

func TestStringToZapLevel(t *testing.T) {
	t.Parallel()

	cases := []struct {
		input string
		want  zapcore.Level
	}{
		{input: "debug", want: zap.DebugLevel},
		{input: "info", want: zap.InfoLevel},
		{input: "warn", want: zap.WarnLevel},
		{input: "error", want: zap.ErrorLevel},
		{input: "dpanic", want: zap.DPanicLevel},
		{input: "panic", want: zap.PanicLevel},
		{input: "fatal", want: zap.FatalLevel},
		{input: "unknown", want: zap.InfoLevel},
	}

	for _, cs := range cases {
		cs := cs
		t.Run(cs.input, func(t *testing.T) {
			t.Parallel()

			result := stringToZapLevel(cs.input)
			if diff := cmp.Diff(cs.want, result); diff != "" {
				t.Errorf("(-want, +got)\n%s", diff)
			}
		})
	}
}
