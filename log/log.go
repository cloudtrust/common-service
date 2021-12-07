package log

import (
	"context"
	"errors"

	cs "github.com/cloudtrust/common-service/v2"
	errorhandler "github.com/cloudtrust/common-service/v2/errors"
	kit_log "github.com/go-kit/kit/log"
	kit_level "github.com/go-kit/kit/log/level"
)

// Logger interface for logging with level
type Logger interface {
	Debug(ctx context.Context, keyvals ...interface{})
	Info(ctx context.Context, keyvals ...interface{})
	Warn(ctx context.Context, keyvals ...interface{})
	Error(ctx context.Context, keyvals ...interface{})
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

var levels = map[string]kit_level.Option{
	"debug": kit_level.AllowDebug(),
	"info":  kit_level.AllowInfo(),
	"warn":  kit_level.AllowWarn(),
	"error": kit_level.AllowError(),
}

// ConvertToLevel transform string value in level
func ConvertToLevel(strLevel string) (kit_level.Option, error) {
	var level, ok = levels[strLevel]

	if !ok {
		return nil, errors.New(errorhandler.MsgErrInvalidParam + errorhandler.Level)
	}

	return level, nil
}

func (l *ctLogger) Debug(ctx context.Context, keyvals ...interface{}) {
	keyvals = append(keyvals, extractInfoFromContext(ctx)...)
	_ = kit_level.Debug(l.logger).Log(keyvals...)
}

func (l *ctLogger) Info(ctx context.Context, keyvals ...interface{}) {
	keyvals = append(keyvals, extractInfoFromContext(ctx)...)
	_ = kit_level.Info(l.logger).Log(keyvals...)
}

func (l *ctLogger) Warn(ctx context.Context, keyvals ...interface{}) {
	keyvals = append(keyvals, extractInfoFromContext(ctx)...)
	_ = kit_level.Warn(l.logger).Log(keyvals...)
}

func (l *ctLogger) Error(ctx context.Context, keyvals ...interface{}) {
	keyvals = append(keyvals, extractInfoFromContext(ctx)...)
	_ = kit_level.Error(l.logger).Log(keyvals...)
}

func (l *ctLogger) ToGoKitLogger() kit_log.Logger {
	return l.logger
}

var (
	contextInfo = map[string]cs.CtContext{
		"user_id":  cs.CtContextUserID,
		"realm_id": cs.CtContextRealmID,
		"corr_id":  cs.CtContextCorrelationID,
	}
)

func extractInfoFromContext(ctx context.Context) []interface{} {
	var keyvals = []interface{}{}

	if ctx == nil {
		return keyvals
	}

	for k, v := range contextInfo {
		var value = ctx.Value(v)
		if value != nil {
			keyvals = append(keyvals, k, value.(string))
		}
	}

	return keyvals
}
