package http

import (
	"net/http"
	"time"

	log "github.com/go-kit/kit/log"
)

type accessLogHandler struct {
	logger  log.Logger
	handler http.Handler
}

// MakeAccessLogHandler creates an HTTP handler to log access logs.
func MakeAccessLogHandler(logger log.Logger, handler http.Handler) http.Handler {
	return &accessLogHandler{
		logger, handler,
	}
}

func (h accessLogHandler) ServeHTTP(w http.ResponseWriter, req *http.Request) {
	writer := &responseWrapper{
		w, http.StatusOK, 0,
	}

	url := *req.URL
	uri := req.RequestURI

	// Requests using the CONNECT method over HTTP/2.0 must use
	// the authority field (aka r.Host) to identify the target.
	// Refer: https://httpwg.github.io/specs/rfc7540.html#CONNECT
	if req.ProtoMajor == 2 && req.Method == "CONNECT" {
		uri = req.Host
	}
	if uri == "" {
		uri = url.RequestURI()
	}

	defer func(begin time.Time) {
		_ = h.logger.Log("method", req.Method, "uri", uri, "status_code", writer.Status(), "size", writer.Size(), "time", time.Since(begin))
	}(time.Now())

	h.handler.ServeHTTP(writer, req)
}

// responseLogger is wrapper of http.ResponseWriter that keeps track of its HTTP
// status code and body size
type responseWrapper struct {
	w      http.ResponseWriter
	status int
	size   int
}

func (l *responseWrapper) Header() http.Header {
	return l.w.Header()
}

func (l *responseWrapper) Write(b []byte) (int, error) {
	size, err := l.w.Write(b)
	l.size += size
	return size, err
}

func (l *responseWrapper) WriteHeader(s int) {
	l.w.WriteHeader(s)
	l.status = s
}

func (l *responseWrapper) Status() int {
	return l.status
}

func (l *responseWrapper) Size() int {
	return l.size
}

func (l *responseWrapper) Flush() {
	f, ok := l.w.(http.Flusher)
	if ok {
		f.Flush()
	}
}
