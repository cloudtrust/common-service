package log

import (
	"context"

	kit_log "github.com/go-kit/kit/log"
)

type nopLogger struct{}

// NewNopLogger returns a logger that doesn't do anything.
func NewNopLogger() Logger { return &nopLogger{} }

func (*nopLogger) Debug(context.Context, ...interface{}) error { return nil }
func (*nopLogger) Info(context.Context, ...interface{}) error  { return nil }
func (*nopLogger) Warn(context.Context, ...interface{}) error  { return nil }
func (*nopLogger) Error(context.Context, ...interface{}) error { return nil }
func (*nopLogger) ToGoKitLogger() kit_log.Logger               { return kit_log.NewNopLogger() }
