package http

import (
	"bytes"
	"context"
	"encoding/json"
	"net/http"
	"regexp"
	"strings"

	errorhandler "github.com/cloudtrust/common-service/v2/errors"
	"github.com/cloudtrust/common-service/v2/log"
	"github.com/cloudtrust/common-service/v2/security"
	"github.com/go-kit/kit/ratelimit"
	"github.com/gorilla/mux"
	"github.com/pkg/errors"
)

// StatusCreated let handlers return a HTTP Created (201) with the given location in headers
type StatusCreated struct {
	Location string
}

// StatusNoContent let handlers return a HTTP NoContent (204)
type StatusNoContent struct {
}

// MimeContent defines a mime content for HTTP responses
type MimeContent struct {
	Filename string
	MimeType string
	Content  []byte
}

// GenericResponse for HTTP requests
type GenericResponse struct {
	StatusCode       int
	Headers          map[string]string
	MimeContent      *MimeContent
	JSONableResponse any
}

var protoRegex = regexp.MustCompile(`(?i)(?:proto=)(https|http)`)

// WriteResponse writes a response for a mime content type
func (r *GenericResponse) WriteResponse(w http.ResponseWriter) {
	if r.Headers == nil {
		r.Headers = make(map[string]string)
	}
	// Headers
	if r.MimeContent != nil {
		r.Headers["Content-Type"] = r.MimeContent.MimeType
		if len(r.MimeContent.Filename) > 0 {
			// Does not support UTF-8 or spaces in filename
			r.Headers["Content-Disposition"] = "attachment; filename=" + r.MimeContent.Filename
		}
	}
	for k, v := range r.Headers {
		w.Header().Set(k, v)
	}

	// Body
	if r.MimeContent != nil {
		w.WriteHeader(r.StatusCode)
		_, _ = w.Write(r.MimeContent.Content)
	} else if r.JSONableResponse != nil {
		writeJSON(r.JSONableResponse, w, r.StatusCode)
	} else {
		w.WriteHeader(r.StatusCode)
	}
}

func writeJSON(jsonableResponse any, w http.ResponseWriter, statusCode int) {
	w.Header().Set("Content-Type", "application/json; charset=utf-8")
	w.WriteHeader(statusCode)

	var json, err = json.MarshalIndent(jsonableResponse, "", " ")

	if err == nil {
		_, _ = w.Write(json)
	}
}

// BasicDecodeRequest does not expect parameters
func BasicDecodeRequest(ctx context.Context, req *http.Request) (any, error) {
	return DecodeRequestWithHeaders(ctx, req, map[string]string{}, map[string]string{}, nil)
}

// DecodeRequest gets the HTTP parameters and body content
func DecodeRequest(ctx context.Context, req *http.Request, pathParams map[string]string, queryParams map[string]string) (any, error) {
	return DecodeRequestWithHeaders(ctx, req, pathParams, queryParams, nil)
}

func getBodyAsString(req *http.Request) string {
	buf := new(bytes.Buffer)
	var reader = req.Body
	if req.GetBody != nil {
		// GetBody is not nil: may be req.Body has already been used and can't be re-used
		reader, _ = req.GetBody()
	}
	_, _ = buf.ReadFrom(reader)
	return buf.String()
}

// DecodeRequestWithHeaders gets the HTTP parameters, headers and body content
func DecodeRequestWithHeaders(_ context.Context, req *http.Request, pathParams map[string]string, queryParams map[string]string, headers []string) (any, error) {
	var request = map[string]string{}

	// Fetch and validate path parameter such as realm, userID, ...
	var m = mux.Vars(req)
	for key, validationRegExp := range pathParams {
		if v, ok := m[key]; ok {
			if matched, _ := regexp.Match(validationRegExp, []byte(v)); !matched {
				return nil, errorhandler.CreateInvalidPathParameterError(key)
			}
			request[key] = m[key]
		}
	}

	request["scheme"] = getScheme(req)
	request["host"] = req.Host

	// Input validation of body content should be performed once the content is unmarshalled (Endpoint layer)
	request["body"] = getBodyAsString(req)

	// Fetch and validate query parameter such as email, firstName, ...
	for key, validationRegExp := range queryParams {
		if value := req.URL.Query().Get(key); value != "" {
			if matched, _ := regexp.Match(validationRegExp, []byte(value)); !matched {
				return nil, errorhandler.CreateInvalidQueryParameterError(key)
			}

			request[key] = value
		}
	}

	for _, headerKey := range headers {
		request[headerKey] = req.Header.Get(headerKey)
	}

	return request, nil
}

func getScheme(req *http.Request) string {
	var forwardedHeader = req.Header.Get("Forwarded")

	if forwardedHeader != "" {
		// match should contain at least two elements if the protocol was specified in the Forwarded header.
		// The first match (match[0]) will always be the 'proto=' capture, which we ignore.
		// In the case of multiple proto parameters we only extract the first.
		if match := protoRegex.FindStringSubmatch(forwardedHeader); len(match) > 1 {
			return strings.ToLower(match[1])
		}
	}

	if req.TLS == nil {
		return "http"
	}

	return "https"
}

// EncodeReply encodes the reply.
func EncodeReply(_ context.Context, w http.ResponseWriter, rep any) error {
	if rep == nil {
		w.WriteHeader(http.StatusOK)
		return nil
	}

	switch e := rep.(type) {
	case GenericResponse:
		e.WriteResponse(w)
	case StatusCreated:
		if e.Location != "" {
			w.Header().Set("Location", e.Location)
		}
		w.WriteHeader(http.StatusCreated)
	case StatusNoContent:
		w.WriteHeader(http.StatusNoContent)
	default:
		writeJSON(rep, w, http.StatusOK)
	}

	return nil
}

// ErrorHandlerNoLog calls ErrorHandler without logger
func ErrorHandlerNoLog() func(context.Context, error, http.ResponseWriter) {
	return ErrorHandler(log.NewNopLogger())
}

// ErrorHandler encodes the reply when there is an error.
func ErrorHandler(logger log.Logger) func(context.Context, error, http.ResponseWriter) {
	return func(ctx context.Context, err error, w http.ResponseWriter) {
		switch e := errors.Cause(err).(type) {
		case security.ForbiddenError:
			w.WriteHeader(http.StatusForbidden)
			_, _ = w.Write([]byte(errorhandler.GetEmitter() + "." + errorhandler.MsgErrOpNotPermitted))
		case errorhandler.Error:
			w.WriteHeader(e.Status)
			// You should really take care of what you are sending here : e.Message should not leak any sensitive information
			_, _ = w.Write([]byte(e.Message))
		case errorhandler.DetailedError:
			w.WriteHeader(e.Status())
			_, _ = w.Write([]byte(e.ErrorMessage()))
		case errorhandler.UnauthorizedError:
			w.WriteHeader(http.StatusUnauthorized)
			_, _ = w.Write([]byte(e.Error()))
		default:
			logger.Error(ctx, "errorHandler", http.StatusInternalServerError, "msg", e.Error())
			if err == ratelimit.ErrLimited {
				w.WriteHeader(http.StatusTooManyRequests)
			} else {
				w.WriteHeader(http.StatusInternalServerError)
			}
		}
	}
}
