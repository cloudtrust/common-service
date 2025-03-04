package healthcheck

import (
	"fmt"
	"time"

	"gopkg.in/h2non/gentleman.v2"
	"gopkg.in/h2non/gentleman.v2/plugins/timeout"
)

type httpChecker struct {
	alias          string
	httpClient     *gentleman.Client
	expectedStatus int
	response       HealthStatus
}

func newHTTPEndpointChecker(alias string, targetURL string, timeoutDuration time.Duration, expectedStatus int, cacheDuration time.Duration, timeProvider TimeProvider) BasicChecker {
	var httpClient = gentleman.New()
	{
		httpClient = httpClient.URL(targetURL)
		httpClient = httpClient.Use(timeout.Request(timeoutDuration))
	}
	var http = "http"
	return &httpChecker{
		alias:          alias,
		httpClient:     httpClient,
		expectedStatus: expectedStatus,
		response:       HealthStatus{Name: &alias, Type: &http, CacheDuration: cacheDuration, TimeProvider: timeProvider},
	}
}

func (hc *httpChecker) CheckStatus() HealthStatus {
	if !hc.response.hasExpired() {
		return hc.response
	}

	req := hc.httpClient.Get()
	var resp *gentleman.Response
	{
		var err error

		resp, err = req.Do()
		if err != nil {
			hc.response.stateDown("Can't hit target: " + err.Error())
		} else if resp.StatusCode == hc.expectedStatus {
			hc.response.stateUp()
		} else {
			hc.response.stateDown(fmt.Sprintf("Expected status %d but received %d", hc.expectedStatus, resp.StatusCode))
		}

		hc.response.touch()
		return hc.response
	}
}
