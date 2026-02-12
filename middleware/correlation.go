package middleware

import (
	"context"
	"net/http"
	"regexp"

	cs "github.com/cloudtrust/common-service/v2"
	errorhandler "github.com/cloudtrust/common-service/v2/errors"
	gen "github.com/cloudtrust/common-service/v2/idgenerator"
)

const (
	regExpCorrelationID = `^[\w\d_#@-]{1,100}$`

	hdrCorrelationID = "X-Correlation-ID"
)

// MakeHTTPCorrelationIDMW retrieve the correlation ID from the HTTP header 'X-Correlation-ID'.
// It there is no such header, it generates a correlation ID.
func MakeHTTPCorrelationIDMW(idGenerator gen.IDGenerator, componentName, componentID string) func(http.Handler) http.Handler {
	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, req *http.Request) {
			var correlationID = req.Header.Get(hdrCorrelationID)

			if correlationID == "" {
				correlationID = idGenerator.NextID()
			} else if match, _ := regexp.MatchString(regExpCorrelationID, correlationID); !match {
				httpErrorHandler(req.Context(), http.StatusBadRequest, errorhandler.CreateInvalidQueryParameterError(hdrCorrelationID), w)
				return
			}

			var ctx = context.WithValue(req.Context(), cs.CtContextCorrelationID, correlationID)

			// Set X-Correlation-ID header for future response
			w.Header().Set(hdrCorrelationID, correlationID)

			next.ServeHTTP(w, req.WithContext(ctx))
		})
	}
}
