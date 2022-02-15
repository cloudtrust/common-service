package metrics

import (
	"context"
	"fmt"
	"time"

	cs "github.com/cloudtrust/common-service/v2"
	"github.com/cloudtrust/common-service/v2/log"
	"github.com/go-kit/kit/metrics"
	metric "github.com/go-kit/kit/metrics/influx"
	influx "github.com/influxdata/influxdb1-client/v2"
)

// Counter interface for go-kit/Counter
type Counter interface {
	With(labelValues ...string) metrics.Counter
	Add(delta float64)
}

// Gauge interface for go-kit/Gauge
type Gauge interface {
	With(labelValues ...string) metrics.Gauge
	Set(value float64)
	Add(delta float64)
}

// Histogram interface for go-kit/Histogram
type Histogram interface {
	With(labelValues ...string) Histogram
	Observe(value float64)
}

// Metrics client.
type Metrics interface {
	NewCounter(name string) Counter
	NewGauge(name string) Gauge
	NewHistogram(name string) Histogram
	WriteLoop(c <-chan time.Time)
	Stats(_ context.Context, name string, tags map[string]string, fields map[string]interface{}) error
	Ping(timeout time.Duration) (time.Duration, string, error)
	TrackFunc(ctx context.Context, d time.Duration, name string, tags func() map[string]string, fields func() map[string]interface{})
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
	WriteLoop(ctx context.Context, c <-chan time.Time, w metric.BatchPointsWriter)
}

// GetBatchPointsConfig gets the influx configuration
func GetBatchPointsConfig(v cs.Configuration, prefix string) influx.BatchPointsConfig {
	return influx.BatchPointsConfig{
		Precision:        v.GetString(prefix + "-precision"),
		Database:         v.GetString(prefix + "-database"),
		RetentionPolicy:  v.GetString(prefix + "-retention-policy"),
		WriteConsistency: v.GetString(prefix + "-write-consistency"),
	}
}

// NewMetrics returns an InfluxMetrics.
// For its configuration, parameters are built with the given prefix, then a dash symbol, then one of these suffixes:
// host-port, username, password, precision, database, retention-policy, write-consistency
// If a parameter exists only named with the given prefix and if its value if false, the InfluxMetrics
// will be a inactive one (Noop)
func NewMetrics(v cs.Configuration, prefix string, logger log.Logger) (Metrics, error) {
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
	var gokitInflux = metric.New(
		map[string]string{},
		influxBatchPointsConfig,
		log.With(logger, "unit", "go-kit influx").ToGoKitLogger(),
	)

	return &influxMetrics{
		influx:            influxClient,
		metrics:           gokitInflux,
		BatchPointsConfig: influxBatchPointsConfig,
		logger:            logger,
	}, nil
}

// influxMetrics sends metrics to the Influx DB.
type influxMetrics struct {
	influx            Influx
	metrics           GoKitMetrics
	BatchPointsConfig influx.BatchPointsConfig
	logger            log.Logger
}

// Close closes the influx client
func (m *influxMetrics) Close() {
	_ = m.influx.Close()
}

// NewCounter returns a go-kit Counter.
func (m *influxMetrics) NewCounter(name string) Counter {
	return m.metrics.NewCounter(name)
}

// NewGauge returns a go-kit Gauge.
func (m *influxMetrics) NewGauge(name string) Gauge {
	return m.metrics.NewGauge(name)
}

type ctHistogram struct {
	Histogram metrics.Histogram
}

func (h *ctHistogram) With(labelValues ...string) Histogram {
	var histo = h.Histogram.With(labelValues...)
	return &ctHistogram{
		Histogram: histo,
	}
}

func (h *ctHistogram) Observe(value float64) {
	h.Histogram.Observe(value)
}

// NewHistogram returns a go-kit Histogram.
func (m *influxMetrics) NewHistogram(name string) Histogram {
	var histo metrics.Histogram
	histo = m.metrics.NewHistogram(name)
	return &ctHistogram{
		Histogram: histo,
	}
}

// Write writes the data to the Influx DB.
func (m *influxMetrics) Write(bp influx.BatchPoints) error {
	return m.influx.Write(bp)
}

// WriteLoop writes the data to the Influx DB.
func (m *influxMetrics) WriteLoop(c <-chan time.Time) {
	m.metrics.WriteLoop(context.Background(), c, m.influx)
}

func (m *influxMetrics) Stats(_ context.Context, name string, tags map[string]string, fields map[string]interface{}) error {
	// Create a new point batch
	var batchPoints influx.BatchPoints
	{
		var err error
		batchPoints, err = influx.NewBatchPoints(m.BatchPointsConfig)
		if err != nil {
			return err
		}
	}

	var point *influx.Point
	{
		var err error
		point, err = influx.NewPoint(name, tags, fields, time.Now())
		if err != nil {
			return err
		}
		batchPoints.AddPoint(point)
	}

	// Write the batch
	var err = m.influx.Write(batchPoints)
	if err != nil {
		return err
	}

	return nil
}

func (m *influxMetrics) TrackFunc(ctx context.Context, d time.Duration, name string, tags func() map[string]string, fields func() map[string]interface{}) {
	go func() {
		for {
			select {
			case <-time.After(d):
				err := m.Stats(ctx, name, tags(), fields())
				m.logger.Warn(ctx, "error", err)
			}
		}
	}()
}

// Ping test the connection to the Influx DB.
func (m *influxMetrics) Ping(timeout time.Duration) (time.Duration, string, error) {
	return m.influx.Ping(timeout)
}

// NoopMetrics is an Influx metrics that does nothing.
type NoopMetrics struct{}

// Close does nothing.
func (m *NoopMetrics) Close() {
	// No operation
}

// NewCounter returns a Counter that does nothing.
func (m *NoopMetrics) NewCounter(name string) Counter { return &NoopCounter{} }

// NewGauge returns a Gauge that does nothing.
func (m *NoopMetrics) NewGauge(name string) Gauge { return &NoopGauge{} }

// NewHistogram returns an Histogram that does nothing.
func (m *NoopMetrics) NewHistogram(name string) Histogram { return &NoopHistogram{} }

// Stats does nothing
func (m *NoopMetrics) Stats(_ context.Context, name string, tags map[string]string, fields map[string]interface{}) error {
	return nil
}

// Write does nothing.
//func (m *NoopMetrics) Write(bp influx.BatchPoints) error { return nil }

// WriteLoop does nothing.
func (m *NoopMetrics) WriteLoop(c <-chan time.Time) {
	// No operation
}

// Ping does nothing.
func (m *NoopMetrics) Ping(timeout time.Duration) (time.Duration, string, error) {
	return time.Duration(0), "", nil
}

// TrackFunc does nothing
func (m *NoopMetrics) TrackFunc(ctx context.Context, d time.Duration, name string, tags func() map[string]string, fields func() map[string]interface{}) {
	//Nothing to do
}

// NoopCounter is a Counter that does nothing.
type NoopCounter struct{}

// With does nothing.
func (c *NoopCounter) With(labelValues ...string) metrics.Counter { return c }

// Add does nothing.
func (c *NoopCounter) Add(delta float64) {
	// No operation
}

// NoopGauge is a Gauge that does nothing.
type NoopGauge struct{}

// With does nothing.
func (g *NoopGauge) With(labelValues ...string) metrics.Gauge { return g }

// Set does nothing.
func (g *NoopGauge) Set(value float64) {
	// No operation
}

// Add does nothing.
func (g *NoopGauge) Add(delta float64) {
	// No operation
}

// NoopHistogram is an Histogram that does nothing.
type NoopHistogram struct{}

// With does nothing.
func (h *NoopHistogram) With(labelValues ...string) Histogram { return h }

// Observe does nothing.
func (h *NoopHistogram) Observe(value float64) {
	// No operation
}
