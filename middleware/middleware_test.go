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

	errorhandler "github.com/cloudtrust/common-service/v2/errors"

	cs "github.com/cloudtrust/common-service/v2"
	"github.com/cloudtrust/common-service/v2/configuration"
	"github.com/cloudtrust/common-service/v2/log"
	"github.com/cloudtrust/common-service/v2/middleware/mock"
	"github.com/stretchr/testify/assert"
	"go.uber.org/mock/gomock"
)

func dummyEndpoint(ctx context.Context, request any) (response any, err error) {
	return nil, nil
}

func TestEndpointLoggingMW(t *testing.T) {
	var mockCtrl = gomock.NewController(t)
	defer mockCtrl.Finish()
	var mockLogger = mock.NewLogger(mockCtrl)

	var mInOut = MakeEndpointLoggingMW(mockLogger)(dummyEndpoint)
	var mIn = MakeEndpointLoggingInMW(mockLogger)(dummyEndpoint)
	var mOut = MakeEndpointLoggingOutMW(mockLogger)(dummyEndpoint)

	// Context with correlation ID.
	var corrID = strconv.FormatUint(rand.Uint64(), 10)
	var ctx = context.WithValue(context.Background(), cs.CtContextCorrelationID, corrID)

	t.Run("In/Out logging - With correlation ID", func(t *testing.T) {
		var req = "req"
		mockLogger.EXPECT().Debug(gomock.Any(), "req", req).Times(1)
		mockLogger.EXPECT().Debug(gomock.Any(), "res", nil).Times(1)
		mInOut(ctx, req)
	})
	t.Run("In logging - With correlation ID", func(t *testing.T) {
		var req = "req"
		mockLogger.EXPECT().Debug(gomock.Any(), "req", req).Times(1)
		mockLogger.EXPECT().Debug(gomock.Any(), "res", "hidden").Times(1)
		mIn(ctx, req)
	})
	t.Run("Out logging - With correlation ID", func(t *testing.T) {
		var req = "req"
		mockLogger.EXPECT().Debug(gomock.Any(), "req", "hidden").Times(1)
		mockLogger.EXPECT().Debug(gomock.Any(), "res", nil).Times(1)
		mOut(ctx, req)
	})

	t.Run("In/Out logging - Without correlation ID", func(t *testing.T) {
		var f = func() {
			mInOut(context.Background(), nil)
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
	var availabilityChecker = NewEndpointAvailabilityChecker(feature, mockIDRetriever, mockConfRetriever)
	var m = MakeEndpointAvailableCheckMW(availabilityChecker, logger)(dummyHandler)

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
