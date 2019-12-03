package middleware

import (
	"context"
	"net/http"

	cs "github.com/cloudtrust/common-service"
	gen "github.com/cloudtrust/common-service/idgenerator"
	"github.com/cloudtrust/common-service/log"
	"github.com/cloudtrust/common-service/tracing"
)

// MakeHTTPCorrelationIDMW retrieve the correlation ID from the HTTP header 'X-Correlation-ID'.
// It there is no such header, it generates a correlation ID.
func MakeHTTPCorrelationIDMW(idGenerator gen.IDGenerator, tracer tracing.OpentracingClient, logger log.Logger, componentName, componentID string) func(http.Handler) http.Handler {
	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, req *http.Request) {
			var correlationID = req.Header.Get("X-Correlation-ID")

			if correlationID == "" {
				correlationID = idGenerator.NextID()
			}

			var ctx = context.WithValue(req.Context(), cs.CtContextCorrelationID, correlationID)

			// Set X-Correlation-ID header for future response
			w.Header().Set("X-Correlation-ID", correlationID)

			next.ServeHTTP(w, req.WithContext(ctx))
		})
	}
}
