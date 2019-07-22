package log

import kit_level "github.com/go-kit/kit/log/level"
import kit_log "github.com/go-kit/kit/log"

type Logger interface {
	Debug(keyvals ...interface{}) error
	Info(keyvals ...interface{}) error
	Warn(keyvals ...interface{}) error
	Error(keyvals ...interface{}) error
	ToGoKitLogger() kit_log.Logger
}

type ctLogger struct {
	logger kit_log.Logger
}

// NewLeveledLogger is a wrapper around gokit logger with level
func NewLeveledLogger(l kit_log.Logger) Logger {
	return &ctLogger{
		logger: l,
	}
}

// With returns a new contextual logger with keyvals prepended to those passed
// to calls to Log. If logger is also a contextual logger created by With or
// WithPrefix, keyvals is appended to the existing context.
//
// The returned Logger replaces all value elements (odd indexes) containing a
// Valuer with their generated value for each call to its Log method.
func With(logger Logger, keyvals ...interface{}) Logger {
	return &ctLogger{
		logger: kit_log.With(logger.ToGoKitLogger(), keyvals),
	}
}

func (l *ctLogger) Debug(keyvals ...interface{}) error {
	return kit_level.Debug(l.logger).Log(keyvals)
}

func (l *ctLogger) Info(keyvals ...interface{}) error {
	return kit_level.Info(l.logger).Log(keyvals)
}

func (l *ctLogger) Warn(keyvals ...interface{}) error {
	return kit_level.Warn(l.logger).Log(keyvals)
}

func (l *ctLogger) Error(keyvals ...interface{}) error {
	return kit_level.Error(l.logger).Log(keyvals)
}

func (l *ctLogger) ToGoKitLogger() kit_log.Logger {
	return l.logger
}
