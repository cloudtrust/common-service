package middleware

import (
	"context"
	"errors"
	"net/http"
	"time"

	cs "github.com/cloudtrust/common-service"
	"github.com/cloudtrust/common-service/configuration"
	errorhandler "github.com/cloudtrust/common-service/errors"
	"github.com/cloudtrust/common-service/log"
	"github.com/cloudtrust/common-service/metrics"
)

// MakeEndpointLoggingMW makes a logging middleware.
func MakeEndpointLoggingMW(logger log.Logger) cs.Middleware {
	return func(next cs.Endpoint) cs.Endpoint {
		return func(ctx context.Context, req interface{}) (interface{}, error) {
			logger.Debug(ctx, "req", req)
			res, err := next(ctx, req)
			logger.Debug(ctx, "res", res)
			return res, err
		}
	}
}

// MakeEndpointLoggingNoInputMW makes a logging middleware.
func MakeEndpointLoggingNoInputMW(logger log.Logger) cs.Middleware {
	return func(next cs.Endpoint) cs.Endpoint {
		return func(ctx context.Context, req interface{}) (interface{}, error) {
			res, err := next(ctx, req)
			logger.Debug(ctx, "res", res)
			return res, err
		}
	}
}

// MakeEndpointLoggingNoOutputMW makes a logging middleware.
func MakeEndpointLoggingNoOutputMW(logger log.Logger) cs.Middleware {
	return func(next cs.Endpoint) cs.Endpoint {
		return func(ctx context.Context, req interface{}) (interface{}, error) {
			logger.Debug(ctx, "req", req)
			res, err := next(ctx, req)
			return res, err
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

// MakeEndpointAvailableCheckMW makes a middleware that ensure a feature is enabled at admin configuration level in the current context
func MakeEndpointAvailableCheckMW(enabledKey string, idRetriever IDRetriever, confRetriever AdminConfigurationRetriever, logger log.Logger) func(http.Handler) http.Handler {
	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, req *http.Request) {
			var ctx = req.Context()
			var realmName = ctx.Value(cs.CtContextRealm).(string)
			var accessToken = ctx.Value(cs.CtContextAccessToken).(string)
			// Get realm ID
			var realmID, err = idRetriever.GetID(accessToken, realmName)
			if err != nil {
				logger.Info(ctx, "msg", "Can't get realm ID", "realm", realmName)
				handleError(req.Context(), err, w)
				return
			}
			// Get admin configuration
			var conf configuration.RealmAdminConfiguration
			conf, err = confRetriever.GetAdminConfiguration(ctx, realmID)
			if err != nil {
				logger.Info(ctx, "msg", "Can't get realm admin configuration", "realm", realmName)
				handleError(req.Context(), err, w)
				return
			}
			if !conf.AvailableChecks[enabledKey] {
				logger.Info(ctx, "msg", "Feature not enabled", "realm", realmName, "feat", enabledKey)
				handleError(req.Context(), errorhandler.CreateEndpointNotEnabled(realmName), w)
				return
			}

			ctx = context.WithValue(ctx, cs.CtContextRealmID, realmID)
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

	w.Write([]byte(errorhandler.GetEmitter() + "." + err.Error()))
}
