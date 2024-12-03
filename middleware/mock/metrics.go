// Code generated by MockGen. DO NOT EDIT.
// Source: github.com/cloudtrust/common-service/v2/metrics (interfaces: Metrics,Histogram)
//
// Generated by this command:
//
//	mockgen --build_flags=--mod=mod -destination=./mock/metrics.go -package=mock -mock_names=Metrics=Metrics,Histogram=Histogram github.com/cloudtrust/common-service/v2/metrics Metrics,Histogram
//

// Package mock is a generated GoMock package.
package mock

import (
	context "context"
	reflect "reflect"
	time "time"

	metrics "github.com/cloudtrust/common-service/v2/metrics"
	gomock "go.uber.org/mock/gomock"
)

// Metrics is a mock of Metrics interface.
type Metrics struct {
	ctrl     *gomock.Controller
	recorder *MetricsMockRecorder
	isgomock struct{}
}

// MetricsMockRecorder is the mock recorder for Metrics.
type MetricsMockRecorder struct {
	mock *Metrics
}

// NewMetrics creates a new mock instance.
func NewMetrics(ctrl *gomock.Controller) *Metrics {
	mock := &Metrics{ctrl: ctrl}
	mock.recorder = &MetricsMockRecorder{mock}
	return mock
}

// EXPECT returns an object that allows the caller to indicate expected use.
func (m *Metrics) EXPECT() *MetricsMockRecorder {
	return m.recorder
}

// Close mocks base method.
func (m *Metrics) Close() {
	m.ctrl.T.Helper()
	m.ctrl.Call(m, "Close")
}

// Close indicates an expected call of Close.
func (mr *MetricsMockRecorder) Close() *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "Close", reflect.TypeOf((*Metrics)(nil).Close))
}

// NewCounter mocks base method.
func (m *Metrics) NewCounter(name string) metrics.Counter {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "NewCounter", name)
	ret0, _ := ret[0].(metrics.Counter)
	return ret0
}

// NewCounter indicates an expected call of NewCounter.
func (mr *MetricsMockRecorder) NewCounter(name any) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "NewCounter", reflect.TypeOf((*Metrics)(nil).NewCounter), name)
}

// NewGauge mocks base method.
func (m *Metrics) NewGauge(name string) metrics.Gauge {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "NewGauge", name)
	ret0, _ := ret[0].(metrics.Gauge)
	return ret0
}

// NewGauge indicates an expected call of NewGauge.
func (mr *MetricsMockRecorder) NewGauge(name any) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "NewGauge", reflect.TypeOf((*Metrics)(nil).NewGauge), name)
}

// NewHistogram mocks base method.
func (m *Metrics) NewHistogram(name string) metrics.Histogram {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "NewHistogram", name)
	ret0, _ := ret[0].(metrics.Histogram)
	return ret0
}

// NewHistogram indicates an expected call of NewHistogram.
func (mr *MetricsMockRecorder) NewHistogram(name any) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "NewHistogram", reflect.TypeOf((*Metrics)(nil).NewHistogram), name)
}

// Ping mocks base method.
func (m *Metrics) Ping(timeout time.Duration) (time.Duration, string, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "Ping", timeout)
	ret0, _ := ret[0].(time.Duration)
	ret1, _ := ret[1].(string)
	ret2, _ := ret[2].(error)
	return ret0, ret1, ret2
}

// Ping indicates an expected call of Ping.
func (mr *MetricsMockRecorder) Ping(timeout any) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "Ping", reflect.TypeOf((*Metrics)(nil).Ping), timeout)
}

// Stats mocks base method.
func (m *Metrics) Stats(arg0 context.Context, name string, tags map[string]string, fields map[string]any) error {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "Stats", arg0, name, tags, fields)
	ret0, _ := ret[0].(error)
	return ret0
}

// Stats indicates an expected call of Stats.
func (mr *MetricsMockRecorder) Stats(arg0, name, tags, fields any) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "Stats", reflect.TypeOf((*Metrics)(nil).Stats), arg0, name, tags, fields)
}

// WriteLoop mocks base method.
func (m *Metrics) WriteLoop(c <-chan time.Time) {
	m.ctrl.T.Helper()
	m.ctrl.Call(m, "WriteLoop", c)
}

// WriteLoop indicates an expected call of WriteLoop.
func (mr *MetricsMockRecorder) WriteLoop(c any) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "WriteLoop", reflect.TypeOf((*Metrics)(nil).WriteLoop), c)
}

// Histogram is a mock of Histogram interface.
type Histogram struct {
	ctrl     *gomock.Controller
	recorder *HistogramMockRecorder
	isgomock struct{}
}

// HistogramMockRecorder is the mock recorder for Histogram.
type HistogramMockRecorder struct {
	mock *Histogram
}

// NewHistogram creates a new mock instance.
func NewHistogram(ctrl *gomock.Controller) *Histogram {
	mock := &Histogram{ctrl: ctrl}
	mock.recorder = &HistogramMockRecorder{mock}
	return mock
}

// EXPECT returns an object that allows the caller to indicate expected use.
func (m *Histogram) EXPECT() *HistogramMockRecorder {
	return m.recorder
}

// Observe mocks base method.
func (m *Histogram) Observe(value float64) {
	m.ctrl.T.Helper()
	m.ctrl.Call(m, "Observe", value)
}

// Observe indicates an expected call of Observe.
func (mr *HistogramMockRecorder) Observe(value any) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "Observe", reflect.TypeOf((*Histogram)(nil).Observe), value)
}

// With mocks base method.
func (m *Histogram) With(labelValues ...string) metrics.Histogram {
	m.ctrl.T.Helper()
	varargs := []any{}
	for _, a := range labelValues {
		varargs = append(varargs, a)
	}
	ret := m.ctrl.Call(m, "With", varargs...)
	ret0, _ := ret[0].(metrics.Histogram)
	return ret0
}

// With indicates an expected call of With.
func (mr *HistogramMockRecorder) With(labelValues ...any) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "With", reflect.TypeOf((*Histogram)(nil).With), labelValues...)
}
