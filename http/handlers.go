package http

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"net/http"
	"regexp"

	cs "github.com/cloudtrust/common-service"
	"github.com/cloudtrust/common-service/security"
	"github.com/go-kit/kit/log"
	"github.com/go-kit/kit/ratelimit"
	"github.com/gorilla/mux"
	"github.com/pkg/errors"
)

// BasicDecodeRequest does not expect parameters
func BasicDecodeRequest(ctx context.Context, req *http.Request) (interface{}, error) {
	return DecodeRequest(ctx, req, map[string]string{}, map[string]string{})
}

// DecodeEventsRequest gets the HTTP parameters and body content
func DecodeRequest(_ context.Context, req *http.Request, pathParams map[string]string, queryParams map[string]string) (interface{}, error) {
	var request = map[string]string{}

	// Fetch and validate path parameter such as realm, userID, ...
	var m = mux.Vars(req)
	for key, validationRegExp := range pathParams {
		if v, ok := m[key]; ok {
			if matched, _ := regexp.Match(validationRegExp, []byte(v)); !matched {
				return nil, fmt.Errorf("Invalid path param: %s", key)
			}
			request[key] = m[key]
		}
	}

	request["scheme"] = getScheme(req)
	request["host"] = req.Host

	buf := new(bytes.Buffer)
	buf.ReadFrom(req.Body)
	// Input validation of body content should be performed once the content is unmarshalled (Endpoint layer)
	request["body"] = buf.String()

	// Fetch and validate query parameter such as email, firstName, ...
	for key, validationRegExp := range queryParams {
		if value := req.URL.Query().Get(key); value != "" {
			if matched, _ := regexp.Match(validationRegExp, []byte(value)); !matched {
				return nil, fmt.Errorf("Invalid path param: %s", key)
			}

			request[key] = value
		}
	}

	return request, nil
}

func getScheme(req *http.Request) string {
	var xForwardedProtoHeader = req.Header.Get("X-Forwarded-Proto")

	if xForwardedProtoHeader != "" {
		return xForwardedProtoHeader
	}

	if req.TLS == nil {
		return "http"
	}

	return "https"
}

// EncodeReply encodes the reply.
func EncodeReply(_ context.Context, w http.ResponseWriter, rep interface{}) error {
	if rep == nil {
		w.WriteHeader(http.StatusOK)
		return nil
	}

	w.Header().Set("Content-Type", "application/json; charset=utf-8")
	w.WriteHeader(http.StatusOK)

	var json, err = json.MarshalIndent(rep, "", " ")

	if err == nil {
		w.Write(json)
	}

	return nil
}

// ErrorHandlerNoLog calls ErrorHandler without logger
func ErrorHandlerNoLog() func(context.Context, error, http.ResponseWriter) {
	return ErrorHandler(log.NewNopLogger())
}

// ErrorHandler encodes the reply when there is an error.
func ErrorHandler(logger cs.Logger) func(context.Context, error, http.ResponseWriter) {
	return func(_ context.Context, err error, w http.ResponseWriter) {
		switch e := errors.Cause(err).(type) {
		case security.ForbiddenError:
			logger.Log("ErrorHandler", http.StatusForbidden, "msg", e.Error())
			w.WriteHeader(http.StatusForbidden)
		case Error:
			logger.Log("ErrorHandler", e.Status, "msg", e.Error())
			w.WriteHeader(e.Status)
			// You should really take care of what you are sending here : e.Message should not leak any sensitive information
			w.Write([]byte(e.Message))
		default:
			logger.Log("ErrorHandler", http.StatusInternalServerError, "msg", e.Error())
			if err == ratelimit.ErrLimited {
				w.WriteHeader(http.StatusTooManyRequests)
			} else {
				w.WriteHeader(http.StatusInternalServerError)
			}
		}
	}
}
