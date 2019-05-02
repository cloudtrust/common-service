package middleware

//go:generate mockgen -destination=./mock/logging.go -package=mock -mock_names=Logger=Logger github.com/go-kit/kit/log Logger
//go:generate mockgen -destination=./mock/instrumenting.go -package=mock -mock_names=Histogram=Histogram github.com/go-kit/kit/metrics Histogram
//go:generate mockgen -destination=./mock/tracing.go -package=mock -mock_names=Tracer=Tracer,Span=Span,SpanContext=SpanContext github.com/opentracing/opentracing-go Tracer,Span,SpanContext
//go:generate mockgen -destination=./mock/idGenerator.go -package=mock -mock_names=IDGenerator=IDGenerator github.com/cloudtrust/common-service/idgenerator IDGenerator

import (
	"bytes"
	"context"
	"fmt"
	"math/rand"
	"net/http"
	"net/http/httptest"
	"strconv"
	"testing"
	"time"

	"github.com/cloudtrust/common-service/middleware/mock"
	"github.com/cloudtrust/keycloak-bridge/pkg/event"
	"github.com/golang/mock/gomock"
	opentracing "github.com/opentracing/opentracing-go"
	"github.com/stretchr/testify/assert"
)

func init() {
	rand.Seed(time.Now().UnixNano())
}

func dummyEndpoint(ctx context.Context, request interface{}) (response interface{}, err error) {
	return nil, nil
}

func TestHTTPTracingMW(t *testing.T) {
	var mockCtrl = gomock.NewController(t)
	defer mockCtrl.Finish()
	var mockTracer = mock.NewTracer(mockCtrl)
	var mockSpan = mock.NewSpan(mockCtrl)
	var mockSpanContext = mock.NewSpanContext(mockCtrl)

	var m = MakeHTTPTracingMW(mockTracer, "componentName", "operationName")(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {}))

	// HTTP request.
	var req = httptest.NewRequest("POST", "http://cloudtrust.io/getusers", bytes.NewReader([]byte{}))
	var w = httptest.NewRecorder()

	// With existing tracer.
	mockTracer.EXPECT().Extract(opentracing.HTTPHeaders, gomock.Any()).Return(mockSpanContext, nil).Times(1)
	mockTracer.EXPECT().StartSpan("operationName", gomock.Any()).Return(mockSpan).Times(1)
	mockSpan.EXPECT().Finish().Return().Times(1)
	mockSpan.EXPECT().SetTag(gomock.Any(), gomock.Any()).Return(mockSpan).Times(3)
	m.ServeHTTP(w, req)

	// Without existing tracer.
	mockTracer.EXPECT().Extract(opentracing.HTTPHeaders, gomock.Any()).Return(nil, fmt.Errorf("fail")).Times(1)
	mockTracer.EXPECT().StartSpan("operationName").Return(mockSpan).Times(1)
	mockSpan.EXPECT().Finish().Return().Times(1)
	mockSpan.EXPECT().SetTag(gomock.Any(), gomock.Any()).Return(mockSpan).Times(3)
	m.ServeHTTP(w, req)
}

func TestEndpointLoggingMW(t *testing.T) {
	var mockCtrl = gomock.NewController(t)
	defer mockCtrl.Finish()
	var mockLogger = mock.NewLogger(mockCtrl)

	var m = MakeEndpointLoggingMW(mockLogger)(dummyEndpoint)

	// Context with correlation ID.
	var corrID = strconv.FormatUint(rand.Uint64(), 10)
	var ctx = context.WithValue(context.Background(), "correlation_id", corrID)

	// With correlation ID.
	var req = event.Request{
		Type:   "Event",
		Object: nil,
	}
	mockLogger.EXPECT().Log("correlation_id", corrID, "took", gomock.Any()).Return(nil).Times(1)
	m(ctx, req)

	// Without correlation ID.
	var f = func() {
		m(context.Background(), req)
	}
	assert.Panics(t, f)
}

func TestEndpointInstrumentingMW(t *testing.T) {
	var mockCtrl = gomock.NewController(t)
	defer mockCtrl.Finish()
	var mockHistogram = mock.NewHistogram(mockCtrl)

	var m = MakeEndpointInstrumentingMW(mockHistogram)(dummyEndpoint)

	// Context with correlation ID.
	var corrID = strconv.FormatUint(rand.Uint64(), 10)
	var ctx = context.WithValue(context.Background(), "correlation_id", corrID)

	// With correlation ID.
	var req = event.Request{
		Type:   "Event",
		Object: nil,
	}
	mockHistogram.EXPECT().With("correlation_id", corrID).Return(mockHistogram).Times(1)
	mockHistogram.EXPECT().Observe(gomock.Any()).Return().Times(1)
	m(ctx, req)

	// Without correlation ID.
	var f = func() {
		m(context.Background(), req)
	}
	assert.Panics(t, f)
}

func TestEndpointTracingMW(t *testing.T) {
	var mockCtrl = gomock.NewController(t)
	defer mockCtrl.Finish()
	var mockTracer = mock.NewTracer(mockCtrl)
	var mockSpan = mock.NewSpan(mockCtrl)
	var mockSpanContext = mock.NewSpanContext(mockCtrl)

	var m = MakeEndpointTracingMW(mockTracer, "operationName")(dummyEndpoint)

	// Context with correlation ID.
	var corrID = strconv.FormatUint(rand.Uint64(), 10)
	var ctx = context.WithValue(context.Background(), "correlation_id", corrID)
	ctx = opentracing.ContextWithSpan(ctx, mockSpan)

	// With correlation ID.
	var req = event.Request{
		Type: "Event",
	}
	mockTracer.EXPECT().StartSpan("operationName", gomock.Any()).Return(mockSpan).Times(1)
	mockSpan.EXPECT().Context().Return(mockSpanContext).Times(1)
	mockSpan.EXPECT().Finish().Return().Times(1)
	mockSpan.EXPECT().SetTag("correlation_id", corrID).Return(mockSpan).Times(1)
	m(ctx, req)

	// Without tracer.
	m(context.Background(), req)

	// Stats without correlation ID.
	mockTracer.EXPECT().StartSpan("operationName", gomock.Any()).Return(mockSpan).Times(1)
	mockSpan.EXPECT().Context().Return(mockSpanContext).Times(1)
	mockSpan.EXPECT().Finish().Return().Times(1)
	var f = func() {
		m(opentracing.ContextWithSpan(context.Background(), mockSpan), req)
	}
	assert.Panics(t, f)
}
