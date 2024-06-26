package middleware

import (
	"context"
	"errors"
	"net/http"
	"time"

	cs "github.com/cloudtrust/common-service/v2"
	"github.com/cloudtrust/common-service/v2/configuration"
	errorhandler "github.com/cloudtrust/common-service/v2/errors"
	"github.com/cloudtrust/common-service/v2/log"
	"github.com/cloudtrust/common-service/v2/metrics"
)

// MakeEndpointLoggingMW makes a logging middleware.
func MakeEndpointLoggingMW(logger log.Logger) cs.Middleware {
	return MakeEndpointLoggingInOutMW(logger)
}

// MakeEndpointLoggingInMW makes a logging middleware
func MakeEndpointLoggingInMW(logger log.Logger) cs.Middleware {
	return func(next cs.Endpoint) cs.Endpoint {
		return func(ctx context.Context, req any) (any, error) {
			logger.Debug(ctx, "req", req)
			res, err := next(ctx, req)
			logger.Debug(ctx, "res", "hidden")
			return res, err
		}
	}
}

// MakeEndpointLoggingOutMW makes a logging middleware
func MakeEndpointLoggingOutMW(logger log.Logger) cs.Middleware {
	return func(next cs.Endpoint) cs.Endpoint {
		return func(ctx context.Context, req any) (any, error) {
			logger.Debug(ctx, "req", "hidden")
			res, err := next(ctx, req)
			logger.Debug(ctx, "res", res)
			return res, err
		}
	}
}

// MakeEndpointLoggingInOutMW makes a logging middleware
func MakeEndpointLoggingInOutMW(logger log.Logger) cs.Middleware {
	return func(next cs.Endpoint) cs.Endpoint {
		return func(ctx context.Context, req any) (any, error) {
			logger.Debug(ctx, "req", req)
			res, err := next(ctx, req)
			logger.Debug(ctx, "res", res)
			return res, err
		}
	}
}

// MakeEndpointInstrumentingMW makes a middleware that measure the endpoints response time and
// send the metrics to influx DB.
func MakeEndpointInstrumentingMW(m metrics.Metrics, histoName string) cs.Middleware {
	h := m.NewHistogram(histoName)
	return func(next cs.Endpoint) cs.Endpoint {
		return func(ctx context.Context, req any) (any, error) {
			defer func(begin time.Time) {
				h.With("corr_id", ctx.Value(cs.CtContextCorrelationID).(string)).Observe(time.Since(begin).Seconds())
			}(time.Now())
			return next(ctx, req)
		}
	}
}

// IDRetriever is an interface to get an ID using an object's name
type IDRetriever interface {
	GetID(accessToken, name string) (string, error)
}

// AdminConfigurationRetriever is an interface to get an admin configuration
type AdminConfigurationRetriever interface {
	GetAdminConfiguration(ctx context.Context, realmID string) (configuration.RealmAdminConfiguration, error)
}

// EndpointAvailabilityChecker interface
type EndpointAvailabilityChecker interface {
	CheckAvailability(ctx context.Context, logger log.Logger) (context.Context, error)
	CheckAvailabilityForRealm(ctx context.Context, targetRealm string, logger log.Logger) (context.Context, error)
}

type endpointAvailabilityChecker struct {
	enabledKey    string
	idRetriever   IDRetriever
	confRetriever AdminConfigurationRetriever
}

// NewEndpointAvailabilityChecker creates a new EndpointAvailabilityChecker
func NewEndpointAvailabilityChecker(enabledKey string, idRetriever IDRetriever, confRetriever AdminConfigurationRetriever) EndpointAvailabilityChecker {
	return &endpointAvailabilityChecker{
		enabledKey:    enabledKey,
		idRetriever:   idRetriever,
		confRetriever: confRetriever,
	}
}

// CheckAvailability checks if an endpoint is enabled for the realm stored in the context
func (eac *endpointAvailabilityChecker) CheckAvailability(ctx context.Context, logger log.Logger) (context.Context, error) {
	return eac.CheckAvailabilityForRealm(ctx, ctx.Value(cs.CtContextRealm).(string), logger)
}

// CheckAvailabilityForRealm checks if an endpoint is enabled for the given realm
func (eac *endpointAvailabilityChecker) CheckAvailabilityForRealm(ctx context.Context, targetRealm string, logger log.Logger) (context.Context, error) {
	var accessToken = ctx.Value(cs.CtContextAccessToken).(string)

	// Get realm ID
	var realmID, err = eac.idRetriever.GetID(accessToken, targetRealm)
	if err != nil {
		logger.Info(ctx, "msg", "Can't get realm ID", "realm", targetRealm)
		return ctx, err
	}
	// Get admin configuration
	var conf configuration.RealmAdminConfiguration
	conf, err = eac.confRetriever.GetAdminConfiguration(ctx, realmID)
	if err != nil {
		logger.Info(ctx, "msg", "Can't get realm admin configuration", "realm", targetRealm)
		return ctx, err
	}
	if !conf.AvailableChecks[eac.enabledKey] {
		logger.Info(ctx, "msg", "Feature not enabled", "realm", targetRealm, "feat", eac.enabledKey)
		return ctx, errorhandler.CreateEndpointNotEnabled(targetRealm)
	}

	return context.WithValue(ctx, cs.CtContextRealmID, realmID), nil
}

// MakeEndpointAvailableCheckMW makes a middleware that ensure a feature is enabled at admin configuration level in the current context
func MakeEndpointAvailableCheckMW(availabilityChecker EndpointAvailabilityChecker, logger log.Logger) func(http.Handler) http.Handler {
	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, req *http.Request) {
			var ctx = req.Context()
			var realmName = ctx.Value(cs.CtContextRealm).(string)
			var err error
			ctx, err = availabilityChecker.CheckAvailabilityForRealm(ctx, realmName, logger)
			if err != nil {
				handleError(req.Context(), err, w)
				return
			}

			next.ServeHTTP(w, req.WithContext(ctx))
		})
	}
}

func handleError(ctx context.Context, err error, w http.ResponseWriter) {
	switch e := err.(type) {
	case errorhandler.Error:
		httpErrorHandler(ctx, e.Status, e, w)
	default:
		httpErrorHandler(ctx, http.StatusInternalServerError, errors.New("unexpected.error"), w)
	}
}

func httpErrorHandler(_ context.Context, statusCode int, err error, w http.ResponseWriter) {
	w.WriteHeader(statusCode)
	w.Header().Set("Content-Type", "application/json; charset=utf-8")

	_, _ = w.Write([]byte(errorhandler.GetEmitter() + "." + err.Error()))
}
