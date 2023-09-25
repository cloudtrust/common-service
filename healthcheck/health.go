package healthcheck

import (
	"context"
	"net/http"
	"time"

	"github.com/cloudtrust/common-service/v2/events"
	commonhttp "github.com/cloudtrust/common-service/v2/http"
	log "github.com/cloudtrust/common-service/v2/log"
	"github.com/go-kit/kit/ratelimit"
)

// HealthChecker is a tool used to perform health checks
type HealthChecker interface {
	CheckStatus() HealthResponse
	AddHealthChecker(name string, checker BasicChecker)
	AddHTTPEndpoint(name string, targetURL string, timeoutDuration time.Duration, expectedStatus int, cacheDuration time.Duration)
	AddHTTPEndpoints(endpoints map[string]string, timeoutDuration time.Duration, expectedStatus int, cacheDuration time.Duration)
	AddDatabase(name string, db HealthDatabase, cacheDuration time.Duration)
	AddAuditEventsReporterModule(name string, reporter events.AuditEventsReporterModule, timeout time.Duration, cacheDuration time.Duration)
	MakeHandler(rateLimit ratelimit.Allower) http.HandlerFunc
}

type healthchecker struct {
	name     string
	checkers []BasicChecker
	logger   log.Logger
}

// BasicChecker is a basic health check processor
type BasicChecker interface {
	CheckStatus() HealthStatus
}

// HealthResponse is the full detailed response to an health check request
type HealthResponse struct {
	Name    string         `json:"name"`
	State   string         `json:"state"`
	Details []HealthStatus `json:"details,omitempty"`
	Healthy bool           `json:"-"`
}

// HealthStatus is the response to an health check of a dependency
type HealthStatus struct {
	Name          *string       `json:"name,omitempty"`
	Type          *string       `json:"type,omitempty"`
	State         *string       `json:"state,omitempty"`
	Message       *string       `json:"message,omitempty"`
	Connection    *string       `json:"connection,omitempty"`
	ValideUntil   time.Time     `json:"-"`
	CacheDuration time.Duration `json:"-"`
}

func (hs *HealthStatus) hasExpired() bool {
	return time.Now().After(hs.ValideUntil)
}

func (hs *HealthStatus) touch() {
	hs.ValideUntil = time.Now().Add(hs.CacheDuration)
}

func (hs *HealthStatus) connection(status string) {
	hs.Connection = &status
	hs.Message = nil
}

func (hs *HealthStatus) stateDown(message string) {
	var state = "DOWN"
	hs.Message = &message
	hs.Connection = nil
	hs.State = &state
}

func (hs *HealthStatus) stateUp() {
	var state = "UP"
	hs.State = &state
	hs.Message = nil
}

// NewHealthChecker Creates a new health checker
func NewHealthChecker(name string, logger log.Logger) HealthChecker {
	return &healthchecker{
		name:   name,
		logger: log.With(logger, "healthchecker", name),
	}
}

func (hc *healthchecker) CheckStatus() HealthResponse {
	var res = HealthResponse{
		Name:    hc.name,
		State:   "UP",
		Healthy: true,
	}
	// Channel to collect dependencies health status
	results := make(chan HealthStatus, len(hc.checkers))
	for _, checker := range hc.checkers {
		go hc.execHealthChecker(checker, results)
	}
	for range hc.checkers {
		status := <-results
		if *status.State == "DOWN" {
			res.Healthy = false
			res.State = "DOWN"
			hc.logger.Info(context.Background(), "msg", *status.Message, "processor", status.Name, "status", *status.State)
		}
		res.Details = append(res.Details, status)
	}
	return res
}

func (hc *healthchecker) execHealthChecker(checker BasicChecker, results chan<- HealthStatus) {
	results <- checker.CheckStatus()
}

func (hc *healthchecker) AddHealthChecker(name string, checker BasicChecker) {
	hc.checkers = append(hc.checkers, checker)
}

func (hc *healthchecker) AddHTTPEndpoint(name string, targetURL string, timeoutDuration time.Duration, expectedStatus int, cacheDuration time.Duration) {
	hc.logger.Info(context.Background(), "msg", "Adding HTTP endpoint", "processor", name, "url", targetURL)
	hc.AddHealthChecker(name, newHTTPEndpointChecker(name, targetURL, timeoutDuration, expectedStatus, cacheDuration))
}

func (hc *healthchecker) AddHTTPEndpoints(endpoints map[string]string, timeoutDuration time.Duration, expectedStatus int, cacheDuration time.Duration) {
	for key, value := range endpoints {
		hc.AddHTTPEndpoint(key, value, timeoutDuration, expectedStatus, cacheDuration)
	}
}

func (hc *healthchecker) AddDatabase(name string, db HealthDatabase, cacheDuration time.Duration) {
	hc.logger.Info(context.Background(), "msg", "Adding database", "processor", name)
	hc.AddHealthChecker(name, newDatabaseChecker(name, db, cacheDuration))
}

func (hc *healthchecker) AddAuditEventsReporterModule(name string, reporter events.AuditEventsReporterModule, timeout time.Duration, cacheDuration time.Duration) {
	hc.logger.Info(context.Background(), "msg", "Adding audit event reporter module", "processor", name)
	hc.AddHealthChecker(name, newAuditEventsReporterChecker(name, reporter, timeout, cacheDuration, hc.logger))
}

// MakeHandler makes a HTTP handler that returns health check information
func (hc *healthchecker) MakeHandler(rateLimit ratelimit.Allower) http.HandlerFunc {
	var ctx = context.Background()
	var responseProvider = ratelimit.NewErroringLimiter(rateLimit)(hc.getCurrentStatus)
	var errorHandler = commonhttp.ErrorHandlerNoLog()

	return http.HandlerFunc(func(w http.ResponseWriter, _ *http.Request) {
		var response, err = responseProvider(ctx, nil)
		if err != nil {
			errorHandler(ctx, err, w)
		} else {
			_ = commonhttp.EncodeReply(ctx, w, response)
		}
	})
}

func (hc *healthchecker) getCurrentStatus(ctx context.Context, _ interface{}) (interface{}, error) {
	var status = hc.CheckStatus()
	var response = commonhttp.GenericResponse{
		StatusCode:       http.StatusOK,
		JSONableResponse: status,
	}
	if !status.Healthy {
		response.StatusCode = http.StatusServiceUnavailable
		hc.logger.Warn(ctx, "msg", "Health check issue detected", "http_status", response.StatusCode, "health", status)
	}
	return response, nil
}
