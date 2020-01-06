package middleware

import (
	"bytes"
	"math/rand"
	"net/http"
	"net/http/httptest"
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

func TestHTTPCorrelationIDMW(t *testing.T) {
	var mockCtrl = gomock.NewController(t)
	defer mockCtrl.Finish()
	var mockLogger = mock.NewLogger(mockCtrl)
	var mockTracer = mock.NewOpentracingClient(mockCtrl)
	var mockIDGenerator = mock.NewIDGenerator(mockCtrl)

	var (
		componentName = "component"
		componentID   = strconv.FormatUint(rand.Uint64(), 10)
		corrID        = strconv.FormatUint(rand.Uint64(), 10)
	)

	// With header 'X-Correlation-ID'
	{
		var mockHandler = http.HandlerFunc(func(w http.ResponseWriter, req *http.Request) {
			var id = req.Context().Value(cs.CtContextCorrelationID).(string)
			assert.Equal(t, corrID, id)
		})

		// HTTP request with valid correlation ID
		var m = MakeHTTPCorrelationIDMW(mockIDGenerator, mockTracer, mockLogger, componentName, componentID)(mockHandler)

		var req = httptest.NewRequest("GET", "http://cloudtrust.io/getusers", bytes.NewReader([]byte{}))
		req.Header.Add("X-Correlation-ID", corrID)
		var w = httptest.NewRecorder()

		m.ServeHTTP(w, req)
	}

	// HTTP request with invalid correlation ID
	{
		var mockHandler = http.HandlerFunc(func(_ http.ResponseWriter, _ *http.Request) {
			assert.Fail(t, "should not be executed")
		})
		var m = MakeHTTPCorrelationIDMW(mockIDGenerator, mockTracer, mockLogger, componentName, componentID)(mockHandler)

		var req = httptest.NewRequest("GET", "http://cloudtrust.io/getusers", bytes.NewReader([]byte{}))
		req.Header.Add("X-Correlation-ID", "<$+invalid+$>")
		var w = httptest.NewRecorder()

		m.ServeHTTP(w, req)
	}

	// Without header 'X-Correlation-ID', so there is a call to IDGenerator.
	{
		var mockCorrID = "keycloak_bridge-123456789-12645316163-45641615174715"
		var mockHandler = http.HandlerFunc(func(w http.ResponseWriter, req *http.Request) {
			var id = req.Context().Value(cs.CtContextCorrelationID).(string)
			assert.Equal(t, mockCorrID, id)
		})
		var m = MakeHTTPCorrelationIDMW(mockIDGenerator, mockTracer, mockLogger, componentName, componentID)(mockHandler)

		// HTTP request.
		var req = httptest.NewRequest("GET", "http://cloudtrust.io/getusers", bytes.NewReader([]byte{}))
		var w = httptest.NewRecorder()

		mockIDGenerator.EXPECT().NextID().Return(mockCorrID).Times(1)
		m.ServeHTTP(w, req)
	}
}
