package healthcheck

import (
	"context"
	"time"

	"github.com/cloudtrust/common-service/v2/events"
	log "github.com/cloudtrust/common-service/v2/log"
)

type auditEventsReporterChecker struct {
	alias          string
	reporter       events.AuditEventsReporterModule
	timeout        time.Duration
	response       HealthStatus
	logger         log.Logger
	failureCounter int
}

func newAuditEventsReporterChecker(alias string, reporter events.AuditEventsReporterModule, timeout time.Duration, cacheDuration time.Duration, logger log.Logger, timeProvider TimeProvider) BasicChecker {
	healthStatusType := "auditEventreporter"
	response := HealthStatus{Name: &alias, Type: &healthStatusType, CacheDuration: cacheDuration, TimeProvider: timeProvider}
	response.connection("init")
	response.stateUp()
	return &auditEventsReporterChecker{
		alias:          alias,
		reporter:       reporter,
		timeout:        timeout,
		response:       response,
		logger:         logger,
		failureCounter: 0,
	}
}

func (a *auditEventsReporterChecker) CheckStatus() HealthStatus {
	if a.response.hasExpired() {
		go a.updateStatus()
	}

	return a.response
}

func (a *auditEventsReporterChecker) updateStatus() {
	finished := make(chan bool)
	go func() {
		event := events.NewEvent("healthcheck", "", "master", "health_checker", "health_checker", "master", map[string]string{})
		a.reporter.ReportEvent(context.Background(), event)
		finished <- true
	}()

	select {
	case <-finished:
		a.response.connection("established")
		a.response.stateUp()
		a.failureCounter = 0
	case <-time.After(a.timeout):
		a.response.stateDown("Events reporter timeout")
		a.failureCounter++
		a.logger.Error(context.Background(), "msg", "Audit Events Reporter timeout to produce events", "timeout", a.timeout, "failureCounter", a.failureCounter)
	}

	a.response.touch()
}
