package tracking

import (
	cs "github.com/cloudtrust/common-service"
	sentry "github.com/getsentry/raven-go"
)

// Sentry is the Sentry client interface.
type Sentry interface {
	URL() string
}

// SentryTracking interface
type SentryTracking interface {
	CaptureError(err error, tags map[string]string) string
	URL() string
	Close()
}

type internalSentry struct {
	sentry *sentry.Client
}

// NewSentry creates a Sentry instance
func NewSentry(v cs.Configuration, prefix string) (SentryTracking, error) {
	sentryEnabled := v.GetBool("sentry")
	if !sentryEnabled {
		return &NoopSentry{}, nil
	}
	sentry, err := sentry.New(v.GetString("sentry-dsn"))
	return &internalSentry{
		sentry: sentry,
	}, err
}

func (s *internalSentry) CaptureError(err error, tags map[string]string) string {
	return s.sentry.CaptureError(err, tags)
}

func (s *internalSentry) URL() string {
	return s.sentry.URL()
}

func (s *internalSentry) Close() {
	s.sentry.Close()
}

// NoopSentry is a Sentry client that does nothing.
type NoopSentry struct{}

// CaptureError does nothing.
func (s *NoopSentry) CaptureError(err error, tags map[string]string) string {
	return ""
}

// URL does nothing.
func (s *NoopSentry) URL() string { return "" }

// Close does nothing.
func (s *NoopSentry) Close() {}
