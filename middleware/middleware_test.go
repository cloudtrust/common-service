package middleware

//go:generate mockgen -destination=./mock/logging.go -package=mock -mock_names=Logger=Logger github.com/go-kit/kit/log Logger
//go:generate mockgen -destination=./mock/tracing.go -package=mock -mock_names=OpentracingClient=OpentracingClient github.com/cloudtrust/common-service/tracing OpentracingClient
//go:generate mockgen -destination=./mock/metrics.go -package=mock -mock_names=Metrics=Metrics,Histogram=Histogram github.com/cloudtrust/common-service/metrics Metrics,Histogram
//go:generate mockgen -destination=./mock/idGenerator.go -package=mock -mock_names=IDGenerator=IDGenerator github.com/cloudtrust/common-service/idgenerator IDGenerator

import (
	"context"
	"math/rand"
	"strconv"
	"testing"
	"time"

	"github.com/cloudtrust/common-service/middleware/mock"
	"github.com/golang/mock/gomock"
	"github.com/stretchr/testify/assert"
)

func init() {
	rand.Seed(time.Now().UnixNano())
}

func dummyEndpoint(ctx context.Context, request interface{}) (response interface{}, err error) {
	return nil, nil
}

func TestHTTPTracingMW(t *testing.T) {
	/*var mockCtrl = gomock.NewController(t)
	defer mockCtrl.Finish()
	var mockTracer = mock.NewOpentracingClient(mockCtrl)

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
	m.ServeHTTP(w, req)*/
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
	mockLogger.EXPECT().Log("correlation_id", corrID, "took", gomock.Any()).Return(nil).Times(1)
	m(ctx, nil)

	// Without correlation ID.
	var f = func() {
		m(context.Background(), nil)
	}
	assert.Panics(t, f)
}

func TestEndpointInstrumentingMW(t *testing.T) {
	/*var mockCtrl = gomock.NewController(t)
	defer mockCtrl.Finish()
	var mockMetrics = mock.NewMetrics(mockCtrl)
	var mockHisto = mock.NewHistogram(mockCtrl)

	var histoName = "histo_name"
	var corrID = strconv.FormatUint(rand.Uint64(), 10)
	var ctx = context.WithValue(context.Background(), "correlation_id", corrID)

	mockMetrics.EXPECT().NewHistogram(histoName).Return(mockHisto).Times(1)
	mockHisto.EXPECT().With("correlation_id", corrID)
	mockHisto.EXPECT().Observe(gomock.Any())

	var m = MakeEndpointInstrumentingMW(mockMetrics, histoName)(dummyEndpoint)

	// With correlation ID.
	m(ctx, nil)

	// Without correlation ID.
	var f = func() {
		m(context.Background(), nil)
	}
	assert.Panics(t, f)*/
}

func TestEndpointTracingMW(t *testing.T) {
	/*var mockCtrl = gomock.NewController(t)
	defer mockCtrl.Finish()
	var mockTracer = mock.NewOpentracingClient(mockCtrl)
	//var mockSpan = mock.NewSpan(mockCtrl)
	//var mockSpanContext = mock.NewSpanContext(mockCtrl)

	var m = MakeEndpointTracingMW(mockTracer, "operationName")(dummyEndpoint)

	// Context with correlation ID.
	var corrID = strconv.FormatUint(rand.Uint64(), 10)
	var ctx = context.WithValue(context.Background(), "correlation_id", corrID)
	ctx = opentracing.ContextWithSpan(ctx, mockSpan)

	// With correlation ID.
	mockTracer.EXPECT().StartSpan("operationName", gomock.Any()).Return(mockSpan).Times(1)
	mockSpan.EXPECT().Context().Return(mockSpanContext).Times(1)
	mockSpan.EXPECT().Finish().Return().Times(1)
	mockSpan.EXPECT().SetTag("correlation_id", corrID).Return(mockSpan).Times(1)
	m(ctx, nil)

	// Without tracer.
	m(context.Background(), nil)

	// Stats without correlation ID.
	mockTracer.EXPECT().StartSpan("operationName", gomock.Any()).Return(mockSpan).Times(1)
	mockSpan.EXPECT().Context().Return(mockSpanContext).Times(1)
	mockSpan.EXPECT().Finish().Return().Times(1)
	var f = func() {
		m(opentracing.ContextWithSpan(context.Background(), mockSpan), nil)
	}
	assert.Panics(t, f)*/
}
