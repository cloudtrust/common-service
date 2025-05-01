package log

import (
	"context"

	kit_log "github.com/go-kit/log"
)

type nopLogger struct{}

// NewNopLogger returns a logger that doesn't do anything.
func NewNopLogger() Logger { return &nopLogger{} }

func (*nopLogger) Debug(context.Context, ...any) {
	// Nothing to do in NopLogger
}

func (*nopLogger) Info(context.Context, ...any) {
	// Nothing to do in NopLogger
}

func (*nopLogger) Warn(context.Context, ...any) {
	// Nothing to do in NopLogger
}

func (*nopLogger) Error(context.Context, ...any) {
	// Nothing to do in NopLogger
}

func (*nopLogger) ToGoKitLogger() kit_log.Logger {
	return kit_log.NewNopLogger()
}
