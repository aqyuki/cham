package logging

import (
	"context"

	"go.uber.org/zap"
)

// contextKey is a private type used to define context key.
type contextKey string

// loggerKey is a context key to store logger in context.
const loggerKey = contextKey("logger")


// WithLogger stores a given logger to given context.
// If context is nil, it will panic. because, this function is wrapper for context.WithValue.
func WithLogger(ctx context.Context, logger *zap.SugaredLogger) context.Context{
	return context.WithValue(ctx, loggerKey, logger)
}

// FromContext returns a logger from given context.
// If not contained logger from given context, it will return a default logger.
func FromContext(ctx context.Context) *zap.SugaredLogger{
	if logger, ok := ctx.Value(loggerKey).(*zap.SugaredLogger); ok {
		return logger
	}
	return DefaultLogger()
}