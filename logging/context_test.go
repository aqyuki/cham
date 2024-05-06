package logging

import (
	"context"
	"reflect"
	"testing"
)

func TestContext(t *testing.T) {
	t.Parallel()

	logger := FromContext(context.Background())
	if logger == nil {
		t.Fatal("expect logger, but received nil")
	}

	ctx := WithLogger(context.Background(), logger)
	if ctx == nil {
		t.Fatal("expect context, but received nil")
	}

	logger2 := FromContext(ctx)
	if logger2 == nil {
		t.Fatal("expect logger, but received nil")
	}

	if !reflect.DeepEqual(logger, logger2) {
		t.Error("expect same logger, but received different logger")
	}
}
