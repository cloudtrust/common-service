package healthcheck

import (
	"bytes"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"
	"time"

	"github.com/cloudtrust/common-service/v2/healthcheck/mock"
	"github.com/cloudtrust/common-service/v2/log"
	"github.com/golang/mock/gomock"
	"github.com/gorilla/mux"
	"github.com/stretchr/testify/assert"
)

func TestEmptyHealthCheck(t *testing.T) {
	var hc = NewHealthChecker("test-module", log.NewNopLogger())
	var res = hc.CheckStatus()
	assert.Nil(t, res.Details)
}

type testAllower struct {
	allow bool
}

func (ta testAllower) Allow() bool {
	return ta.allow
}

func TestHealthCheckHandler(t *testing.T) {
	var mockCtrl = gomock.NewController(t)
	defer mockCtrl.Finish()

	var mockDB = mock.NewHealthDatabase(mockCtrl)
	mockDB.EXPECT().Ping().Return(nil).Times(1)

	var alias1 = "alias-localhost"
	var alias2 = "alias-db"
	var hc = NewHealthChecker("http-test-module", log.NewNopLogger())
	hc.AddDatabase(alias1, mockDB, 15*time.Second)

	var allower = testAllower{allow: true}
	var disallower = testAllower{allow: false}

	r := mux.NewRouter()
	r.Handle("/health/check", hc.MakeHandler(allower))
	r.Handle("/health/check/error", hc.MakeHandler(disallower))

	ts := httptest.NewServer(r)
	defer ts.Close()

	var healthCheckURL = ts.URL + "/health/check"
	var healthCheckErrorURL = ts.URL + "/health/check/error"

	t.Run("RateLimit is blocking call to endpoint", func(t *testing.T) {
		_, statusCode, _ := httpGet(healthCheckErrorURL)

		assert.Equal(t, http.StatusTooManyRequests, statusCode)
	})

	t.Run("State is UP", func(t *testing.T) {
		// First call: only one health checker, state is UP
		resp, statusCode, _ := httpGet(healthCheckURL)
		assert.Equal(t, http.StatusOK, statusCode)

		assert.True(t, strings.Contains(resp, alias1))
		assert.False(t, strings.Contains(resp, alias2))
	})

	hc.AddHTTPEndpoint(alias2, "http://localhost:11111/", 2*time.Second, 200, time.Duration(0))

	t.Run("State is DOWN", func(t *testing.T) {
		// Second call: 2 health checkers, one state is DOWN
		resp, statusCode, _ := httpGet(healthCheckURL)
		assert.Equal(t, http.StatusServiceUnavailable, statusCode)

		assert.True(t, strings.Contains(resp, alias1))
		assert.True(t, strings.Contains(resp, alias2))
	})
}

func httpGet(targetURL string) (string, int, error) {
	res, err := http.Get(targetURL)
	if err != nil {
		return "", 0, err
	}

	buf := new(bytes.Buffer)
	buf.ReadFrom(res.Body)

	return buf.String(), res.StatusCode, nil
}

func TestHealthStatusCache(t *testing.T) {
	var hs = HealthStatus{CacheDuration: 50 * time.Millisecond}
	assert.True(t, hs.hasExpired())

	hs.touch()
	assert.False(t, hs.hasExpired())

	time.Sleep(2 * time.Millisecond)
	assert.False(t, hs.hasExpired())

	time.Sleep(50 * time.Millisecond)
	assert.True(t, hs.hasExpired())
}
