// Code generated by MockGen. DO NOT EDIT.
// Source: github.com/opentracing/opentracing-go (interfaces: Tracer,Span,SpanContext)
//
// Generated by this command:
//
//	mockgen --build_flags=--mod=mod -destination=./mock/tracing.go -package=mock -mock_names=Tracer=Tracer,Span=Span,SpanContext=SpanContext github.com/opentracing/opentracing-go Tracer,Span,SpanContext
//

// Package mock is a generated GoMock package.
package mock

import (
	reflect "reflect"

	opentracing "github.com/opentracing/opentracing-go"
	log "github.com/opentracing/opentracing-go/log"
	gomock "go.uber.org/mock/gomock"
)

// Tracer is a mock of Tracer interface.
type Tracer struct {
	ctrl     *gomock.Controller
	recorder *TracerMockRecorder
	isgomock struct{}
}

// TracerMockRecorder is the mock recorder for Tracer.
type TracerMockRecorder struct {
	mock *Tracer
}

// NewTracer creates a new mock instance.
func NewTracer(ctrl *gomock.Controller) *Tracer {
	mock := &Tracer{ctrl: ctrl}
	mock.recorder = &TracerMockRecorder{mock}
	return mock
}

// EXPECT returns an object that allows the caller to indicate expected use.
func (m *Tracer) EXPECT() *TracerMockRecorder {
	return m.recorder
}

// Extract mocks base method.
func (m *Tracer) Extract(format, carrier any) (opentracing.SpanContext, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "Extract", format, carrier)
	ret0, _ := ret[0].(opentracing.SpanContext)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// Extract indicates an expected call of Extract.
func (mr *TracerMockRecorder) Extract(format, carrier any) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "Extract", reflect.TypeOf((*Tracer)(nil).Extract), format, carrier)
}

// Inject mocks base method.
func (m *Tracer) Inject(sm opentracing.SpanContext, format, carrier any) error {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "Inject", sm, format, carrier)
	ret0, _ := ret[0].(error)
	return ret0
}

// Inject indicates an expected call of Inject.
func (mr *TracerMockRecorder) Inject(sm, format, carrier any) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "Inject", reflect.TypeOf((*Tracer)(nil).Inject), sm, format, carrier)
}

// StartSpan mocks base method.
func (m *Tracer) StartSpan(operationName string, opts ...opentracing.StartSpanOption) opentracing.Span {
	m.ctrl.T.Helper()
	varargs := []any{operationName}
	for _, a := range opts {
		varargs = append(varargs, a)
	}
	ret := m.ctrl.Call(m, "StartSpan", varargs...)
	ret0, _ := ret[0].(opentracing.Span)
	return ret0
}

// StartSpan indicates an expected call of StartSpan.
func (mr *TracerMockRecorder) StartSpan(operationName any, opts ...any) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	varargs := append([]any{operationName}, opts...)
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "StartSpan", reflect.TypeOf((*Tracer)(nil).StartSpan), varargs...)
}

// Span is a mock of Span interface.
type Span struct {
	ctrl     *gomock.Controller
	recorder *SpanMockRecorder
	isgomock struct{}
}

// SpanMockRecorder is the mock recorder for Span.
type SpanMockRecorder struct {
	mock *Span
}

// NewSpan creates a new mock instance.
func NewSpan(ctrl *gomock.Controller) *Span {
	mock := &Span{ctrl: ctrl}
	mock.recorder = &SpanMockRecorder{mock}
	return mock
}

// EXPECT returns an object that allows the caller to indicate expected use.
func (m *Span) EXPECT() *SpanMockRecorder {
	return m.recorder
}

// BaggageItem mocks base method.
func (m *Span) BaggageItem(restrictedKey string) string {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "BaggageItem", restrictedKey)
	ret0, _ := ret[0].(string)
	return ret0
}

// BaggageItem indicates an expected call of BaggageItem.
func (mr *SpanMockRecorder) BaggageItem(restrictedKey any) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "BaggageItem", reflect.TypeOf((*Span)(nil).BaggageItem), restrictedKey)
}

// Context mocks base method.
func (m *Span) Context() opentracing.SpanContext {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "Context")
	ret0, _ := ret[0].(opentracing.SpanContext)
	return ret0
}

// Context indicates an expected call of Context.
func (mr *SpanMockRecorder) Context() *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "Context", reflect.TypeOf((*Span)(nil).Context))
}

// Finish mocks base method.
func (m *Span) Finish() {
	m.ctrl.T.Helper()
	m.ctrl.Call(m, "Finish")
}

// Finish indicates an expected call of Finish.
func (mr *SpanMockRecorder) Finish() *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "Finish", reflect.TypeOf((*Span)(nil).Finish))
}

// FinishWithOptions mocks base method.
func (m *Span) FinishWithOptions(opts opentracing.FinishOptions) {
	m.ctrl.T.Helper()
	m.ctrl.Call(m, "FinishWithOptions", opts)
}

// FinishWithOptions indicates an expected call of FinishWithOptions.
func (mr *SpanMockRecorder) FinishWithOptions(opts any) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "FinishWithOptions", reflect.TypeOf((*Span)(nil).FinishWithOptions), opts)
}

// Log mocks base method.
func (m *Span) Log(data opentracing.LogData) {
	m.ctrl.T.Helper()
	m.ctrl.Call(m, "Log", data)
}

// Log indicates an expected call of Log.
func (mr *SpanMockRecorder) Log(data any) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "Log", reflect.TypeOf((*Span)(nil).Log), data)
}

// LogEvent mocks base method.
func (m *Span) LogEvent(event string) {
	m.ctrl.T.Helper()
	m.ctrl.Call(m, "LogEvent", event)
}

// LogEvent indicates an expected call of LogEvent.
func (mr *SpanMockRecorder) LogEvent(event any) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "LogEvent", reflect.TypeOf((*Span)(nil).LogEvent), event)
}

// LogEventWithPayload mocks base method.
func (m *Span) LogEventWithPayload(event string, payload any) {
	m.ctrl.T.Helper()
	m.ctrl.Call(m, "LogEventWithPayload", event, payload)
}

// LogEventWithPayload indicates an expected call of LogEventWithPayload.
func (mr *SpanMockRecorder) LogEventWithPayload(event, payload any) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "LogEventWithPayload", reflect.TypeOf((*Span)(nil).LogEventWithPayload), event, payload)
}

// LogFields mocks base method.
func (m *Span) LogFields(fields ...log.Field) {
	m.ctrl.T.Helper()
	varargs := []any{}
	for _, a := range fields {
		varargs = append(varargs, a)
	}
	m.ctrl.Call(m, "LogFields", varargs...)
}

// LogFields indicates an expected call of LogFields.
func (mr *SpanMockRecorder) LogFields(fields ...any) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "LogFields", reflect.TypeOf((*Span)(nil).LogFields), fields...)
}

// LogKV mocks base method.
func (m *Span) LogKV(alternatingKeyValues ...any) {
	m.ctrl.T.Helper()
	varargs := []any{}
	for _, a := range alternatingKeyValues {
		varargs = append(varargs, a)
	}
	m.ctrl.Call(m, "LogKV", varargs...)
}

// LogKV indicates an expected call of LogKV.
func (mr *SpanMockRecorder) LogKV(alternatingKeyValues ...any) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "LogKV", reflect.TypeOf((*Span)(nil).LogKV), alternatingKeyValues...)
}

// SetBaggageItem mocks base method.
func (m *Span) SetBaggageItem(restrictedKey, value string) opentracing.Span {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "SetBaggageItem", restrictedKey, value)
	ret0, _ := ret[0].(opentracing.Span)
	return ret0
}

// SetBaggageItem indicates an expected call of SetBaggageItem.
func (mr *SpanMockRecorder) SetBaggageItem(restrictedKey, value any) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "SetBaggageItem", reflect.TypeOf((*Span)(nil).SetBaggageItem), restrictedKey, value)
}

// SetOperationName mocks base method.
func (m *Span) SetOperationName(operationName string) opentracing.Span {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "SetOperationName", operationName)
	ret0, _ := ret[0].(opentracing.Span)
	return ret0
}

// SetOperationName indicates an expected call of SetOperationName.
func (mr *SpanMockRecorder) SetOperationName(operationName any) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "SetOperationName", reflect.TypeOf((*Span)(nil).SetOperationName), operationName)
}

// SetTag mocks base method.
func (m *Span) SetTag(key string, value any) opentracing.Span {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "SetTag", key, value)
	ret0, _ := ret[0].(opentracing.Span)
	return ret0
}

// SetTag indicates an expected call of SetTag.
func (mr *SpanMockRecorder) SetTag(key, value any) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "SetTag", reflect.TypeOf((*Span)(nil).SetTag), key, value)
}

// Tracer mocks base method.
func (m *Span) Tracer() opentracing.Tracer {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "Tracer")
	ret0, _ := ret[0].(opentracing.Tracer)
	return ret0
}

// Tracer indicates an expected call of Tracer.
func (mr *SpanMockRecorder) Tracer() *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "Tracer", reflect.TypeOf((*Span)(nil).Tracer))
}

// SpanContext is a mock of SpanContext interface.
type SpanContext struct {
	ctrl     *gomock.Controller
	recorder *SpanContextMockRecorder
	isgomock struct{}
}

// SpanContextMockRecorder is the mock recorder for SpanContext.
type SpanContextMockRecorder struct {
	mock *SpanContext
}

// NewSpanContext creates a new mock instance.
func NewSpanContext(ctrl *gomock.Controller) *SpanContext {
	mock := &SpanContext{ctrl: ctrl}
	mock.recorder = &SpanContextMockRecorder{mock}
	return mock
}

// EXPECT returns an object that allows the caller to indicate expected use.
func (m *SpanContext) EXPECT() *SpanContextMockRecorder {
	return m.recorder
}

// ForeachBaggageItem mocks base method.
func (m *SpanContext) ForeachBaggageItem(handler func(string, string) bool) {
	m.ctrl.T.Helper()
	m.ctrl.Call(m, "ForeachBaggageItem", handler)
}

// ForeachBaggageItem indicates an expected call of ForeachBaggageItem.
func (mr *SpanContextMockRecorder) ForeachBaggageItem(handler any) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "ForeachBaggageItem", reflect.TypeOf((*SpanContext)(nil).ForeachBaggageItem), handler)
}
