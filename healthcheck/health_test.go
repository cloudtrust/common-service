package healthcheck

import (
	"bytes"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"
	"time"

	"github.com/cloudtrust/common-service/healthcheck/mock"
	"github.com/golang/mock/gomock"
	"github.com/gorilla/mux"
	"github.com/stretchr/testify/assert"
)

func TestEmptyHealthCheck(t *testing.T) {
	var hc = NewHealthChecker("test-module")
	var res = hc.CheckStatus()
	assert.Nil(t, res.Details)
}

func TestHealthCheckHandler(t *testing.T) {
	var mockCtrl = gomock.NewController(t)
	defer mockCtrl.Finish()

	var mockDB = mock.NewHealthDatabase(mockCtrl)
	mockDB.EXPECT().Ping().Return(nil).Times(1)

	var alias1 = "alias-localhost"
	var alias2 = "alias-db"
	var hc = NewHealthChecker("http-test-module")
	hc.AddDatabase(alias1, mockDB, 15*time.Second)

	r := mux.NewRouter()
	r.Handle("/health/check", hc.MakeHandler())

	ts := httptest.NewServer(r)
	defer ts.Close()

	var healthCheckURL = ts.URL + "/health/check"

	{
		// First call: only one health checker, state is UP
		resp, statusCode, _ := httpGet(healthCheckURL)
		assert.Equal(t, http.StatusOK, statusCode)

		assert.True(t, strings.Contains(resp, alias1))
		assert.False(t, strings.Contains(resp, alias2))
	}

	hc.AddHTTPEndpoint(alias2, "http://localhost:11111/", 2*time.Second, 200, time.Duration(0))

	{
		// Second call: 2 health checkers, one state is DOWN
		resp, statusCode, _ := httpGet(healthCheckURL)
		assert.Equal(t, http.StatusServiceUnavailable, statusCode)

		assert.True(t, strings.Contains(resp, alias1))
		assert.True(t, strings.Contains(resp, alias2))
	}
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
