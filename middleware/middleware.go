package middleware

import (
	"context"
	"encoding/json"
	"net/http"
	"time"

	cs "github.com/cloudtrust/common-service"
	"github.com/cloudtrust/common-service/metrics"
)

// MakeEndpointLoggingMW makes a logging middleware.
func MakeEndpointLoggingMW(logger cs.Logger) cs.Middleware {
	return func(next cs.Endpoint) cs.Endpoint {
		return func(ctx context.Context, req interface{}) (interface{}, error) {
			logger.Log("correlation_id", ctx.Value(cs.CtContextCorrelationID).(string))
			return next(ctx, req)
		}
	}
}

// MakeEndpointInstrumentingMW makes a middleware that measure the endpoints response time and
// send the metrics to influx DB.
func MakeEndpointInstrumentingMW(m metrics.Metrics, histoName string) cs.Middleware {
	h := m.NewHistogram(histoName)
	return func(next cs.Endpoint) cs.Endpoint {
		return func(ctx context.Context, req interface{}) (interface{}, error) {
			defer func(begin time.Time) {
				h.With("correlation_id", ctx.Value(cs.CtContextCorrelationID).(string)).Observe(time.Since(begin).Seconds())
			}(time.Now())
			return next(ctx, req)
		}
	}
}

func httpErrorHandler(_ context.Context, statusCode int, err error, w http.ResponseWriter) {
	w.WriteHeader(statusCode)
	w.Header().Set("Content-Type", "application/json; charset=utf-8")

	var reply, _ = json.MarshalIndent(map[string]string{"error": err.Error()}, "", "  ")
	w.Write(reply)
}
