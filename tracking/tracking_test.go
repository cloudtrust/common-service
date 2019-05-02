package tracking

import (
	"fmt"
	"testing"

	"github.com/spf13/viper"
	"github.com/stretchr/testify/assert"
)

func TestNoopSentry(t *testing.T) {
	var cfg = viper.New()
	var sentry, _ = NewSentry(cfg, "sentry")
	defer sentry.Close()

	// CaptureError
	assert.Zero(t, sentry.CaptureError(nil, nil))
	assert.Zero(t, sentry.CaptureError(fmt.Errorf("fail"), map[string]string{"key": "val"}))

	// URL
	assert.Zero(t, sentry.URL())
}

func TestTrueSentry(t *testing.T) {
	var cfg = viper.New()
	cfg.Set("sentry", true)
	cfg.Set("sentry-dsn", "dsn")
	var sentry, _ = NewSentry(cfg, "sentry")
	defer sentry.Close()

	assert.NotNil(t, sentry)
}
