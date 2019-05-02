package tracing

import (
	"testing"

	"github.com/spf13/viper"
	"github.com/stretchr/testify/assert"
)

func TestCreateJaegerClientFails(t *testing.T) {
	var cfg = viper.New()
	var _, err = CreateJaegerClient(cfg, "prefix", "")

	assert.NotNil(t, err)
}

func TestCreateJaegerClientSucceeds(t *testing.T) {
	var cfg = viper.New()
	var jaeger, err = CreateJaegerClient(cfg, "prefix", "cloudtrusttester")
	defer jaeger.Close()

	assert.Nil(t, err)
}
