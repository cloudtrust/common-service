package middleware

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

	errorhandler "github.com/cloudtrust/common-service/errors"

	cs "github.com/cloudtrust/common-service"
	"github.com/cloudtrust/common-service/configuration"
	"github.com/cloudtrust/common-service/log"
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
	var req = "req"
	mockLogger.EXPECT().Debug(gomock.Any(), "req", req).Return(nil).Times(1)
	mockLogger.EXPECT().Debug(gomock.Any(), "res", nil).Return(nil).Times(1)
	m(ctx, req)

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
	mockHisto.EXPECT().With("corr_id", corrID).Return(mockHisto).Times(1)
	mockHisto.EXPECT().Observe(gomock.Any())

	var m = MakeEndpointInstrumentingMW(mockMetrics, histoName)(dummyEndpoint)

	t.Run("With correlation ID", func(t *testing.T) {
		m(ctx, nil)
	})
	t.Run("Without correlation ID", func(t *testing.T) {
		var f = func() {
			m(context.Background(), nil)
		}
		assert.Panics(t, f)
	})
}

func TestEndpointAvailableMW(t *testing.T) {
	var mockCtrl = gomock.NewController(t)
	defer mockCtrl.Finish()

	var mockIDRetriever = mock.NewIDRetriever(mockCtrl)
	var mockConfRetriever = mock.NewAdminConfigurationRetriever(mockCtrl)
	var logger = log.NewNopLogger()

	var accessToken = "jwtaccesstoken"
	var realmName = "myrealm"
	var realmID = "abcdefgh-1234-5678"
	var ctx = context.TODO()
	var expectedError = errors.New("the error")
	var feature = configuration.CheckKeyPhysical
	var adminConfig configuration.RealmAdminConfiguration
	var dummyHandler = http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {})
	var m = MakeEndpointAvailableCheckMW(feature, mockIDRetriever, mockConfRetriever, logger)(dummyHandler)

	// HTTP request
	var req = httptest.NewRequest("POST", "http://cloudtrust.io/api", bytes.NewReader([]byte{}))
	ctx = context.WithValue(ctx, cs.CtContextAccessToken, accessToken)
	ctx = context.WithValue(ctx, cs.CtContextRealm, realmName)
	req = req.WithContext(ctx)

	t.Run("Can't get realmID - unexpected error", func(t *testing.T) {
		var w = httptest.NewRecorder()
		mockIDRetriever.EXPECT().GetID(accessToken, realmName).Return(realmID, expectedError)
		m.ServeHTTP(w, req)
		var result = w.Result()
		assert.Equal(t, http.StatusInternalServerError, result.StatusCode)
	})
	t.Run("Can't get realmID - unauthorized", func(t *testing.T) {
		var w = httptest.NewRecorder()
		mockIDRetriever.EXPECT().GetID(accessToken, realmName).Return(realmID, errorhandler.Error{Status: http.StatusForbidden})
		m.ServeHTTP(w, req)
		var result = w.Result()
		assert.Equal(t, http.StatusForbidden, result.StatusCode)
	})
	t.Run("Get admin configuration fails", func(t *testing.T) {
		var w = httptest.NewRecorder()
		mockIDRetriever.EXPECT().GetID(accessToken, realmName).Return(realmID, nil)
		mockConfRetriever.EXPECT().GetAdminConfiguration(ctx, realmID).Return(adminConfig, expectedError)
		m.ServeHTTP(w, req)
		var result = w.Result()
		assert.Equal(t, http.StatusInternalServerError, result.StatusCode)
	})
	t.Run("Feature is not enabled", func(t *testing.T) {
		var w = httptest.NewRecorder()
		mockIDRetriever.EXPECT().GetID(accessToken, realmName).Return(realmID, nil)
		mockConfRetriever.EXPECT().GetAdminConfiguration(ctx, realmID).Return(adminConfig, nil)
		m.ServeHTTP(w, req)
		var result = w.Result()
		assert.Equal(t, http.StatusConflict, result.StatusCode)
	})
	t.Run("Feature is enabled", func(t *testing.T) {
		var w = httptest.NewRecorder()
		adminConfig.AvailableChecks = map[string]bool{feature: true}
		mockIDRetriever.EXPECT().GetID(accessToken, realmName).Return(realmID, nil)
		mockConfRetriever.EXPECT().GetAdminConfiguration(ctx, realmID).Return(adminConfig, nil)
		m.ServeHTTP(w, req)
		var result = w.Result()
		assert.Equal(t, http.StatusOK, result.StatusCode)
	})
}
