package tracking

import (
	sentry "github.com/getsentry/raven-go"
	"github.com/spf13/viper"
)

// Sentry is the Sentry client interface.
type Sentry interface {
	URL() string
}

// SentryTracking interface
type SentryTracking interface {
	CaptureError(err error, tags map[string]string, interfaces ...sentry.Interface) string
	URL() string
	Close()
}

// NewSentry creates a Sentry instance
func NewSentry(v *viper.Viper, prefix string) (SentryTracking, error) {
	sentryEnabled := v.GetBool("sentry")
	if !sentryEnabled {
		return &NoopSentry{}, nil
	}
	return sentry.New(v.GetString("sentry-dsn"))
}

// NoopSentry is a Sentry client that does nothing.
type NoopSentry struct{}

// CaptureError does nothing.
func (s *NoopSentry) CaptureError(err error, tags map[string]string, interfaces ...sentry.Interface) string {
	return ""
}

// URL does nothing.
func (s *NoopSentry) URL() string { return "" }

// Close does nothing.
func (s *NoopSentry) Close() {}
