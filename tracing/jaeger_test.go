package tracing

//go:generate mockgen --build_flags=--mod=mod -destination=./mock/configuration.go -package=mock -mock_names=Configuration=Configuration github.com/cloudtrust/common-service/v2 Configuration
//go:generate mockgen --build_flags=--mod=mod -destination=./mock/tracing.go -package=mock -mock_names=Tracer=Tracer,Span=Span,SpanContext=SpanContext github.com/opentracing/opentracing-go Tracer,Span,SpanContext

import (
	"bytes"
	"context"
	"errors"
	"math/rand"
	"net/http"
	"net/http/httptest"
	"strconv"
	"testing"
	"time"

	cs "github.com/cloudtrust/common-service/v2"
	"github.com/cloudtrust/common-service/v2/tracing/mock"
	"github.com/golang/mock/gomock"
	opentracing "github.com/opentracing/opentracing-go"
	"github.com/stretchr/testify/assert"
)

func TestCreateNoopJaegerClient(t *testing.T) {
	var mockCtrl = gomock.NewController(t)
	defer mockCtrl.Finish()
	var mockConf = mock.NewConfiguration(mockCtrl)

	var prefix = "prefix"
	var expected = "expected"

	mockConf.EXPECT().GetBool(prefix).Return(false).Times(1)

	var jaeger, err = CreateJaegerClient(mockConf, prefix, "")
	assert.Nil(t, err)

	jaeger.TryStartSpanWithTag(context.TODO(), "name", "tagName", "tagValue")
	jaeger.MakeHTTPTracingMW("component", "operation")
	jaeger.Close()

	var e = jaeger.MakeEndpointTracingMW("name")(
		func(_ context.Context, _ interface{}) (interface{}, error) {
			return expected, nil
		},
	)
	result, _ := e(context.TODO(), nil)
	assert.Equal(t, expected, result)
}

func TestCreateJaegerClientSucceeds(t *testing.T) {
	var mockCtrl = gomock.NewController(t)
	defer mockCtrl.Finish()
	var mockConf = mock.NewConfiguration(mockCtrl)

	var prefix = "prefix"

	mockConf.EXPECT().GetBool(prefix).Return(true).Times(1)
	mockConf.EXPECT().GetString(prefix + "-sampler-type").Return("").Times(1)
	mockConf.EXPECT().GetFloat64(prefix + "-sampler-param").Return(0.0).Times(1)
	mockConf.EXPECT().GetString(prefix + "-sampler-host-port").Return("").Times(1)
	mockConf.EXPECT().GetBool(prefix + "-reporter-logspan").Return(false).Times(1)
	mockConf.EXPECT().GetDuration(prefix + "-write-interval").Return(time.Duration(0)).Times(1)

	var jaeger, err = CreateJaegerClient(mockConf, "prefix", "cloudtrusttester")
	defer jaeger.Close()

	assert.Nil(t, err)

	var initialContext = context.TODO()
	ctx, f := jaeger.TryStartSpanWithTag(initialContext, "op", "tag", "val")
	assert.Equal(t, initialContext, ctx)
	assert.Nil(t, f)
}

func TestHTTPTracingMW(t *testing.T) {
	var mockCtrl = gomock.NewController(t)
	defer mockCtrl.Finish()
	var mockTracer = mock.NewTracer(mockCtrl)
	var mockSpan = mock.NewSpan(mockCtrl)
	var mockSpanContext = mock.NewSpanContext(mockCtrl)

	var tracer = &internalOpentracingClient{
		Tracer: mockTracer,
		closer: nil,
	}

	var m = tracer.MakeHTTPTracingMW("componentName", "operationName")(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {}))

	// HTTP request.
	var req = httptest.NewRequest("POST", "http://cloudtrust.io/getusers", bytes.NewReader([]byte{}))
	var w = httptest.NewRecorder()

	t.Run("With existing tracer", func(t *testing.T) {
		mockTracer.EXPECT().Extract(opentracing.HTTPHeaders, gomock.Any()).Return(mockSpanContext, nil).Times(1)
		mockTracer.EXPECT().StartSpan("operationName", gomock.Any()).Return(mockSpan).Times(1)
		mockSpan.EXPECT().Finish().Return().Times(1)
		mockSpan.EXPECT().SetTag(gomock.Any(), gomock.Any()).Return(mockSpan).Times(3)
		mockSpan.EXPECT().Tracer().Times(1)
		m.ServeHTTP(w, req)
	})

	t.Run("Without existing tracer", func(t *testing.T) {
		mockTracer.EXPECT().Extract(opentracing.HTTPHeaders, gomock.Any()).Return(nil, errors.New("fail")).Times(1)
		mockTracer.EXPECT().StartSpan("operationName").Return(mockSpan).Times(1)
		mockSpan.EXPECT().Finish().Return().Times(1)
		mockSpan.EXPECT().SetTag(gomock.Any(), gomock.Any()).Return(mockSpan).Times(3)
		mockSpan.EXPECT().Tracer().Times(1)
		m.ServeHTTP(w, req)
	})
}

func dummyEndpoint(ctx context.Context, request interface{}) (response interface{}, err error) {
	return nil, nil
}

func TestEndpointTracingMW(t *testing.T) {
	var mockCtrl = gomock.NewController(t)
	defer mockCtrl.Finish()
	var mockTracer = mock.NewTracer(mockCtrl)
	var mockSpan = mock.NewSpan(mockCtrl)
	var mockSpanContext = mock.NewSpanContext(mockCtrl)

	var tracer = &internalOpentracingClient{
		Tracer: mockTracer,
		closer: nil,
	}

	var m = tracer.MakeEndpointTracingMW("operationName")(dummyEndpoint)

	var corrID = strconv.FormatUint(rand.Uint64(), 10)
	var ctx = context.WithValue(context.Background(), cs.CtContextCorrelationID, corrID)

	t.Run("Context with correlation ID", func(t *testing.T) {
		mockSpan.EXPECT().Tracer().Times(1)
		ctx = opentracing.ContextWithSpan(ctx, mockSpan)
	})

	t.Run("With correlation ID", func(t *testing.T) {
		mockTracer.EXPECT().StartSpan("operationName", gomock.Any()).Return(mockSpan).Times(1)
		mockSpan.EXPECT().Context().Return(mockSpanContext).Times(1)
		mockSpan.EXPECT().Finish().Return().Times(1)
		mockSpan.EXPECT().SetTag("correlation_id", corrID).Return(mockSpan).Times(1)
		mockSpan.EXPECT().Tracer().Times(1)
		m(ctx, nil)
	})

	t.Run("Without tracer", func(t *testing.T) {
		m(context.Background(), nil)
	})

	t.Run("Stats without correlation ID", func(t *testing.T) {
		mockTracer.EXPECT().StartSpan("operationName", gomock.Any()).Return(mockSpan).Times(1)
		mockSpan.EXPECT().Context().Return(mockSpanContext).Times(1)
		mockSpan.EXPECT().Finish().Return().Times(1)
		mockSpan.EXPECT().Tracer().Times(1)
		var f = func() {
			m(opentracing.ContextWithSpan(context.Background(), mockSpan), nil)
		}
		assert.Panics(t, f)
	})
}
