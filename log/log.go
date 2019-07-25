package log

import (
	"fmt"

	kit_log "github.com/go-kit/kit/log"
	kit_level "github.com/go-kit/kit/log/level"
)

// Logger interface for logging with level
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
		logger: kit_log.With(logger.ToGoKitLogger(), keyvals...),
	}
}

// AllowLevel sets the log filtering according to the provided level
func AllowLevel(logger Logger, level kit_level.Option) Logger {
	return &ctLogger{
		logger: kit_level.NewFilter(logger.ToGoKitLogger(), level),
	}
}

// ConvertToLevel transform string value in level
func ConvertToLevel(strLevel string) (kit_level.Option, error) {
	switch strLevel {
	case "debug":
		return kit_level.AllowDebug(), nil
	case "info":
		return kit_level.AllowInfo(), nil
	case "warn":
		return kit_level.AllowWarn(), nil
	case "error":
		return kit_level.AllowError(), nil
	default:
		return nil, fmt.Errorf("Invalid level")
	}
}

func (l *ctLogger) Debug(keyvals ...interface{}) error {
	return kit_level.Debug(l.logger).Log(keyvals...)
}

func (l *ctLogger) Info(keyvals ...interface{}) error {
	return kit_level.Info(l.logger).Log(keyvals...)
}

func (l *ctLogger) Warn(keyvals ...interface{}) error {
	return kit_level.Warn(l.logger).Log(keyvals...)
}

func (l *ctLogger) Error(keyvals ...interface{}) error {
	return kit_level.Error(l.logger).Log(keyvals...)
}

func (l *ctLogger) ToGoKitLogger() kit_log.Logger {
	return l.logger
}
