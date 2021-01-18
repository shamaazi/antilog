package antilog

import (
	"context"
)

type loggerKey struct{}

// AttachToContext attaches a logger to a context object
func AttachToContext(ctx context.Context, logger AntiLog) context.Context {
	return context.WithValue(ctx, loggerKey{}, logger)
}

// FromContext returns a logger from a context
func FromContext(ctx context.Context) AntiLog {
	logger := ctx.Value(loggerKey{})

	if logger == nil {
		return AntiLog{}
	}

	return logger.(AntiLog)
}
