package middleware

//go:generate mockgen -destination=./mock/logging.go -package=mock -mock_names=Logger=Logger github.com/cloudtrust/common-service/log Logger
//go:generate mockgen -destination=./mock/tracing.go -package=mock -mock_names=OpentracingClient=OpentracingClient github.com/cloudtrust/common-service/tracing OpentracingClient
//go:generate mockgen -destination=./mock/metrics.go -package=mock -mock_names=Metrics=Metrics,Histogram=Histogram github.com/cloudtrust/common-service/metrics Metrics,Histogram
//go:generate mockgen -destination=./mock/idGenerator.go -package=mock -mock_names=IDGenerator=IDGenerator github.com/cloudtrust/common-service/idgenerator IDGenerator

import (
	"context"
	"math/rand"
	"strconv"
	"testing"
	"time"

	cs "github.com/cloudtrust/common-service"
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

func TestEndpointLoggingMW(t *testing.T) {
	var mockCtrl = gomock.NewController(t)
	defer mockCtrl.Finish()
	var mockLogger = mock.NewLogger(mockCtrl)

	var m = MakeEndpointLoggingMW(mockLogger)(dummyEndpoint)

	// Context with correlation ID.
	var corrID = strconv.FormatUint(rand.Uint64(), 10)
	var ctx = context.WithValue(context.Background(), cs.CtContextCorrelationID, corrID)

	// With correlation ID.
	mockLogger.EXPECT().Log("correlation_id", corrID).Return(nil).Times(1)
	m(ctx, nil)

	// Without correlation ID.
	var f = func() {
		m(context.Background(), nil)
	}
	assert.Panics(t, f)
}

func TestEndpointInstrumentingMW(t *testing.T) {
	var mockCtrl = gomock.NewController(t)
	defer mockCtrl.Finish()
	var mockMetrics = mock.NewMetrics(mockCtrl)
	var mockHisto = mock.NewHistogram(mockCtrl)

	var histoName = "histo_name"
	var corrID = strconv.FormatUint(rand.Uint64(), 10)
	var ctx = context.WithValue(context.Background(), cs.CtContextCorrelationID, corrID)

	mockMetrics.EXPECT().NewHistogram(histoName).Return(mockHisto).Times(1)
	mockHisto.EXPECT().With("correlation_id", corrID).Return(mockHisto).Times(1)
	mockHisto.EXPECT().Observe(gomock.Any())

	var m = MakeEndpointInstrumentingMW(mockMetrics, histoName)(dummyEndpoint)

	// With correlation ID.
	m(ctx, nil)

	// Without correlation ID.
	var f = func() {
		m(context.Background(), nil)
	}
	assert.Panics(t, f)
}
