package metrics

import (
	"fmt"
	"time"

	"github.com/go-kit/kit/log"
	"github.com/go-kit/kit/metrics"
	gokit_influx "github.com/go-kit/kit/metrics/influx"
	metric "github.com/go-kit/kit/metrics/influx"
	influx "github.com/influxdata/influxdb/client/v2"
	"github.com/spf13/viper"
)

// Metrics client.
type Metrics interface {
	NewCounter(name string) metrics.Counter
	NewGauge(name string) metrics.Gauge
	NewHistogram(name string) metrics.Histogram
	WriteLoop(c <-chan time.Time)
	Write(bp influx.BatchPoints) error
	Ping(timeout time.Duration) (time.Duration, string, error)
	Close()
}

// Influx is the Influx client interface.
type Influx interface {
	Ping(timeout time.Duration) (time.Duration, string, error)
	Write(bp influx.BatchPoints) error
	Close() error
}

// GoKitMetrics is the interface of the go-kit metrics.
type GoKitMetrics interface {
	NewCounter(name string) *metric.Counter
	NewGauge(name string) *metric.Gauge
	NewHistogram(name string) *metric.Histogram
	WriteLoop(c <-chan time.Time, w metric.BatchPointsWriter)
}

// GetBatchPointsConfig gets the influx configuration
func GetBatchPointsConfig(v *viper.Viper, prefix string) influx.BatchPointsConfig {
	return influx.BatchPointsConfig{
		Precision:        v.GetString(prefix + "-precision"),
		Database:         v.GetString(prefix + "-database"),
		RetentionPolicy:  v.GetString(prefix + "-retention-policy"),
		WriteConsistency: v.GetString(prefix + "-write-consistency"),
	}
}

// NewMetrics returns an InfluxMetrics.
func NewMetrics(v *viper.Viper, prefix string, logger log.Logger) (Metrics, error) {
	if !v.GetBool(prefix) {
		return &NoopMetrics{}, nil
	}
	logger = log.With(logger, "unit", "influx")

	// Create Influx client
	influxHTTPConfig := influx.HTTPConfig{
		Addr:     fmt.Sprintf("http://%s", v.GetString(prefix+"-host-port")),
		Username: v.GetString(prefix + "-username"),
		Password: v.GetString(prefix + "-password"),
	}
	var influxClient, err = influx.NewHTTPClient(influxHTTPConfig)
	if err != nil {
		return nil, err
	}

	// Create gokit influx
	influxBatchPointsConfig := GetBatchPointsConfig(v, prefix)
	var gokitInflux = gokit_influx.New(
		map[string]string{},
		influxBatchPointsConfig,
		log.With(logger, "unit", "go-kit influx"),
	)

	return &influxMetrics{
		influx:            influxClient,
		metrics:           gokitInflux,
		BatchPointsConfig: influxBatchPointsConfig,
	}, nil
}

// influxMetrics sends metrics to the Influx DB.
type influxMetrics struct {
	influx            Influx
	metrics           GoKitMetrics
	BatchPointsConfig influx.BatchPointsConfig
}

// Close closes the influx client
func (m *influxMetrics) Close() {
	m.influx.Close()
}

// NewCounter returns a go-kit Counter.
func (m *influxMetrics) NewCounter(name string) metrics.Counter {
	return m.metrics.NewCounter(name)
}

// NewGauge returns a go-kit Gauge.
func (m *influxMetrics) NewGauge(name string) metrics.Gauge {
	return m.metrics.NewGauge(name)
}

// NewHistogram returns a go-kit Histogram.
func (m *influxMetrics) NewHistogram(name string) metrics.Histogram {
	return m.metrics.NewHistogram(name)
}

// Write writes the data to the Influx DB.
func (m *influxMetrics) Write(bp influx.BatchPoints) error {
	return m.influx.Write(bp)
}

// WriteLoop writes the data to the Influx DB.
func (m *influxMetrics) WriteLoop(c <-chan time.Time) {
	m.metrics.WriteLoop(c, m.influx)
}

// Ping test the connection to the Influx DB.
func (m *influxMetrics) Ping(timeout time.Duration) (time.Duration, string, error) {
	return m.influx.Ping(timeout)
}

// NoopMetrics is an Influx metrics that does nothing.
type NoopMetrics struct{}

// Close does nothing.
func (m *NoopMetrics) Close() {}

// NewCounter returns a Counter that does nothing.
func (m *NoopMetrics) NewCounter(name string) metrics.Counter { return &NoopCounter{} }

// NewGauge returns a Gauge that does nothing.
func (m *NoopMetrics) NewGauge(name string) metrics.Gauge { return &NoopGauge{} }

// NewHistogram returns an Histogram that does nothing.
func (m *NoopMetrics) NewHistogram(name string) metrics.Histogram { return &NoopHistogram{} }

// Write does nothing.
func (m *NoopMetrics) Write(bp influx.BatchPoints) error { return nil }

// WriteLoop does nothing.
func (m *NoopMetrics) WriteLoop(c <-chan time.Time) {}

// Ping does nothing.
func (m *NoopMetrics) Ping(timeout time.Duration) (time.Duration, string, error) {
	return time.Duration(0), "", nil
}

// NoopCounter is a Counter that does nothing.
type NoopCounter struct{}

// With does nothing.
func (c *NoopCounter) With(labelValues ...string) metrics.Counter { return c }

// Add does nothing.
func (c *NoopCounter) Add(delta float64) {}

// NoopGauge is a Gauge that does nothing.
type NoopGauge struct{}

// With does nothing.
func (g *NoopGauge) With(labelValues ...string) metrics.Gauge { return g }

// Set does nothing.
func (g *NoopGauge) Set(value float64) {}

// Add does nothing.
func (g *NoopGauge) Add(delta float64) {}

// NoopHistogram is an Histogram that does nothing.
type NoopHistogram struct{}

// With does nothing.
func (h *NoopHistogram) With(labelValues ...string) metrics.Histogram { return h }

// Observe does nothing.
func (h *NoopHistogram) Observe(value float64) {}
