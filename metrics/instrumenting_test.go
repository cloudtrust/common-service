package metrics

import (
	"testing"

	influx "github.com/influxdata/influxdb/client/v2"

	"github.com/go-kit/kit/log"
	"github.com/spf13/viper"
	"github.com/stretchr/testify/assert"
)

func TestNoopInfluxClient(t *testing.T) {
	var cfg = viper.New()
	var noop, err = NewMetrics(cfg, "noop", nil)
	assert.Nil(t, err)
	defer noop.Close()

	// Coverage
	counter := noop.NewCounter("name")
	counter.Add(1.0)
	counter.With("value")

	gauge := noop.NewGauge("name")
	gauge.Add(1.0)
	gauge.Set(1.0)
	gauge.With("value")

	histo := noop.NewHistogram("name")
	histo.With("value")
	histo.Observe(1.0)

	var bp influx.BatchPoints
	noop.Write(bp)
	noop.Ping(1)
}

func TestInvalidConfigurationInfluxClient(t *testing.T) {
	var cfg = viper.New()
	cfg.Set("name", true)
	cfg.Set("name-host-port", "influx.io#%")
	var _, err = NewMetrics(cfg, "name", log.NewNopLogger())
	assert.NotNil(t, err)
}

func TestTrueInfluxClient(t *testing.T) {
	var cfg = viper.New()
	cfg.Set("name", true)
	cfg.Set("name-host-port", "influx.io")
	var influx, err = NewMetrics(cfg, "name", log.NewNopLogger())
	assert.Nil(t, err)
	assert.NotNil(t, influx)

	influx.NewCounter("name")
	influx.NewGauge("name")
	influx.NewHistogram("name")
	influx.Ping(1)
}
