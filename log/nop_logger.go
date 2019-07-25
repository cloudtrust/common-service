package log

import kit_log "github.com/go-kit/kit/log"

type nopLogger struct{}

// NewNopLogger returns a logger that doesn't do anything.
func NewNopLogger() Logger { return &nopLogger{} }

func (*nopLogger) Debug(...interface{}) error    { return nil }
func (*nopLogger) Info(...interface{}) error     { return nil }
func (*nopLogger) Warn(...interface{}) error     { return nil }
func (*nopLogger) Error(...interface{}) error    { return nil }
func (*nopLogger) ToGoKitLogger() kit_log.Logger { return kit_log.NewNopLogger() }
