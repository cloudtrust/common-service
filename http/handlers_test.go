package http

import (
	"bytes"
	"context"
	"crypto/tls"
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"net/http"
	"net/http/httptest"
	"net/url"
	"testing"

	errorhandler "github.com/cloudtrust/common-service/v2/errors"
	"github.com/cloudtrust/common-service/v2/http/mock"
	"github.com/cloudtrust/common-service/v2/security"
	"github.com/go-kit/kit/endpoint"
	"github.com/go-kit/kit/ratelimit"
	http_transport "github.com/go-kit/kit/transport/http"
	"github.com/gorilla/mux"
	"github.com/stretchr/testify/assert"
	"go.uber.org/mock/gomock"
)

func makeHandler(e endpoint.Endpoint) *http_transport.Server {
	return http_transport.NewServer(e,
		func(ctx context.Context, req *http.Request) (any, error) {
			return BasicDecodeRequest(ctx, req)
		},
		EncodeReply,
		http_transport.ServerErrorEncoder(ErrorHandlerNoLog()),
	)
}

func TestHttpGenericResponse(t *testing.T) {
	r := mux.NewRouter()
	// Test with JSON content
	r.Handle("/path/to/json", makeHandler(func(_ context.Context, _ any) (response any, err error) {
		return GenericResponse{
			StatusCode:       http.StatusNotFound,
			Headers:          map[string]string{"Location": "here"},
			MimeContent:      nil,
			JSONableResponse: make([]int, 0),
		}, nil
	}))
	// Test with MimeContent
	r.Handle("/path/to/jpeg", makeHandler(func(_ context.Context, _ any) (response any, err error) {
		mime := MimeContent{
			MimeType: "image/jpg",
			Content:  []byte("not a real jpeg"),
			Filename: "filename.jpg",
		}
		return GenericResponse{
			StatusCode:       http.StatusCreated,
			Headers:          nil,
			MimeContent:      &mime,
			JSONableResponse: nil,
		}, nil
	}))
	// Test with neither JSON content nor MimeContent
	r.Handle("/path/to/empty", makeHandler(func(_ context.Context, _ any) (response any, err error) {
		return GenericResponse{StatusCode: http.StatusNoContent}, nil
	}))

	ts := httptest.NewServer(r)
	defer ts.Close()

	t.Run("JSON test", func(t *testing.T) {
		res, err := http.Get(ts.URL + "/path/to/json")
		assert.Nil(t, err)
		assert.Equal(t, http.StatusNotFound, res.StatusCode)
		assert.Equal(t, "here", res.Header.Get("Location"))

		buf := new(bytes.Buffer)
		buf.ReadFrom(res.Body)
		assert.Equal(t, "[]", buf.String())
	})

	t.Run("MIME test", func(t *testing.T) {
		res, err := http.Get(ts.URL + "/path/to/jpeg")
		assert.Nil(t, err)
		assert.Equal(t, http.StatusCreated, res.StatusCode)

		buf := new(bytes.Buffer)
		buf.ReadFrom(res.Body)
		assert.Equal(t, "not a real jpeg", buf.String())
	})

	t.Run("Empty test", func(t *testing.T) {
		res, err := http.Get(ts.URL + "/path/to/empty")
		assert.Nil(t, err)
		assert.Equal(t, http.StatusNoContent, res.StatusCode)

		buf := new(bytes.Buffer)
		buf.ReadFrom(res.Body)
		assert.Equal(t, "", buf.String())
	})
}

type nonJsonable struct {
}

func (nj nonJsonable) MarshalJSON() ([]byte, error) {
	return nil, errors.New("Non serializable")
}

func TestHTTPManagementHandler(t *testing.T) {
	var mockCtrl = gomock.NewController(t)
	defer mockCtrl.Finish()

	var e = func(ctx context.Context, req any) (any, error) {
		var m = req.(map[string]string)
		res := fmt.Sprint(len(m))
		if v, err := m["pathParameter1"]; err {
			res = res + "-p1" + v
		}
		if v, err := m["pathParameter2"]; err {
			res = res + "-p2" + v
		}
		return res, nil
	}
	handler := http_transport.NewServer(e,
		func(ctx context.Context, req *http.Request) (any, error) {
			pathParams := map[string]string{
				"pathParameter1": "^[a-zA-Z0-9-]+$",
				"pathParameter2": "^[a-zA-Z0-9-]+$",
				"pathParameter3": "^[a-zA-Z0-9-]+$",
			}
			queryParams := map[string]string{
				"queryParameter1": "^[a-zA-Z0-9-]+$",
				"queryParameter2": "^[a-zA-Z0-9-]+$",
				"queryParameter3": "^[a-zA-Z0-9-]+$",
			}
			return DecodeRequest(ctx, req, pathParams, queryParams)
		},
		EncodeReply,
		http_transport.ServerErrorEncoder(ErrorHandlerNoLog()),
	)

	r := mux.NewRouter()
	r.Handle("/noparam", handler)
	r.Handle("/param1/{pathParameter1}", handler)
	r.Handle("/param2/{pathParameter2}", handler)
	r.Handle("/param1/{pathParameter1}/param2/{pathParameter2}", handler)

	ts := httptest.NewServer(r)
	defer ts.Close()

	// No path parameter
	checkPathParameter(t, ts.URL+"/noparam", "3")
	// Has pathParam1 in path
	checkPathParameter(t, ts.URL+"/param1/valueA", "4-p1valueA")
	// Has pathParam2 in path
	checkPathParameter(t, ts.URL+"/param2/valueB", "4-p2valueB")
	// Has both pathParam1 and pathParam2 in path
	checkPathParameter(t, ts.URL+"/param1/valueC/param2/valueD", "5-p1valueC-p2valueD")
	// Has pathParam2 invalid in path
	checkInvalidPathParameter(t, ts.URL+"/param2/valueB!")
}

func checkPathParameter(t *testing.T, url string, expected string) {
	res, err := http.Get(url)
	assert.Nil(t, err)
	assert.Equal(t, http.StatusOK, res.StatusCode)

	buf := new(bytes.Buffer)
	buf.ReadFrom(res.Body)
	assert.Equal(t, `"`+expected+`"`, buf.String())
}

func checkInvalidPathParameter(t *testing.T, url string) {
	res, err := http.Get(url)
	assert.Nil(t, err)
	assert.Equal(t, http.StatusBadRequest, res.StatusCode)
}

func genericTestDecodeRequest(ctx context.Context, tls *tls.ConnectionState, forwarded *string, rawQuery string) (map[string]string, error) {
	return genericTestDecodeRequestWithHeader(ctx, tls, forwarded, rawQuery, nil)
}

func genericTestDecodeRequestWithHeader(ctx context.Context, tls *tls.ConnectionState, forwarded *string, rawQuery string, headers map[string]string) (map[string]string, error) {
	input := "the body"
	var req http.Request
	var url url.URL
	url.RawQuery = rawQuery
	req.Host = "localhost"
	req.TLS = tls
	req.Header = make(http.Header)
	if forwarded != nil {
		req.Header.Set("Forwarded", *forwarded)
	}
	req.Body = io.NopCloser(bytes.NewBufferString(input))
	req.URL = &url
	pathParams := map[string]string{
		"pathParam1": "^[a-zA-Z0-9-]+$",
		"pathParam2": "^[a-zA-Z0-9-]+$",
	}
	queryParams := map[string]string{
		"queryParam1": "^[a-zA-Z0-9-]+$",
		"queryParam2": "^[a-zA-Z0-9-]+$",
	}

	var r any
	var err error

	if headers != nil {
		for key, value := range headers {
			req.Header.Add(key, value)
		}
		r, err = DecodeRequestWithHeaders(ctx, &req, pathParams, queryParams, []string{"Authorization"})
	} else {
		r, err = DecodeRequest(ctx, &req, pathParams, queryParams)
	}

	if err != nil {
		return nil, err
	}

	return r.(map[string]string), nil
}

func TestDecodeRequestHTTP(t *testing.T) {
	request, _ := genericTestDecodeRequest(context.Background(), nil, nil, "")

	// Minimum parameters are scheme, host and body
	assert.Equal(t, 3, len(request))
	assert.Equal(t, "localhost", request["host"])
	assert.Equal(t, "http", request["scheme"])
	assert.Equal(t, "the body", request["body"])
}

func TestDecodeRequestHTTPS(t *testing.T) {
	var reqConnState tls.ConnectionState

	request, _ := genericTestDecodeRequest(context.Background(), &reqConnState, nil, "")

	// Minimum parameters are scheme, host and body
	assert.Equal(t, 3, len(request))
	assert.Equal(t, "localhost", request["host"])
	assert.Equal(t, "https", request["scheme"])
	assert.Equal(t, "the body", request["body"])
}

func TestDecodeRequestForwardProto(t *testing.T) {
	forwardedHeader := "for=192.0.2.60;proto=http;by=203.0.113.43"

	request, _ := genericTestDecodeRequest(context.Background(), nil, &forwardedHeader, "")

	// Minimum parameters are scheme, host and body
	assert.Equal(t, 3, len(request))
	assert.Equal(t, "localhost", request["host"])
	assert.Equal(t, "http", request["scheme"])
	assert.Equal(t, "the body", request["body"])
}

func TestDecodeRequestQueryParams(t *testing.T) {
	value := "valueX"
	request, _ := genericTestDecodeRequest(context.Background(), nil, nil, "queryParam2="+value)

	// Minimum parameters are scheme, host and body
	assert.Equal(t, 4, len(request))
	assert.Equal(t, value, request["queryParam2"])
}

func TestDecodeRequestInvalidQueryParams(t *testing.T) {
	value := "valueX!"
	_, err := genericTestDecodeRequest(context.Background(), nil, nil, "queryParam2="+value)

	assert.NotNil(t, err)
}

func TestDecodeRequestWithHeaders(t *testing.T) {
	var auth = "Basic ABC="
	var headers = map[string]string{"Content-Type": "text/plain", "Authorization": auth}
	request, err := genericTestDecodeRequestWithHeader(context.Background(), nil, nil, "", headers)

	assert.Nil(t, err)
	assert.NotContains(t, request, "Content-Type")
	assert.Contains(t, request, "Authorization")
	assert.Equal(t, auth, request["Authorization"])
}

func TestEncodeReply(t *testing.T) {
	var mockCtrl = gomock.NewController(t)
	defer mockCtrl.Finish()

	mockRespWriter := mock.NewResponseWriter(mockCtrl)

	t.Run("Nil response", func(t *testing.T) {
		mockRespWriter.EXPECT().WriteHeader(http.StatusOK)
		assert.Nil(t, EncodeReply(context.Background(), mockRespWriter, nil))
	})
	t.Run("NoContent response", func(t *testing.T) {
		mockRespWriter.EXPECT().WriteHeader(http.StatusNoContent)
		assert.Nil(t, EncodeReply(context.Background(), mockRespWriter, StatusNoContent{}))
	})
	t.Run("Created without location response", func(t *testing.T) {
		mockRespWriter.EXPECT().WriteHeader(http.StatusCreated)
		assert.Nil(t, EncodeReply(context.Background(), mockRespWriter, StatusCreated{}))
	})
	t.Run("Created with location response", func(t *testing.T) {
		var location = "http://localhost/path/to/new/resource"
		headers := make(http.Header)
		gomock.InOrder(
			mockRespWriter.EXPECT().Header().Return(headers),
			mockRespWriter.EXPECT().WriteHeader(http.StatusCreated),
		)
		assert.Nil(t, EncodeReply(context.Background(), mockRespWriter, StatusCreated{Location: location}))
		assert.Len(t, headers, 1)
		assert.Equal(t, location, headers["Location"][0])
	})
	t.Run("JSON serializable response", func(t *testing.T) {
		headers := make(http.Header)
		resp := map[string]string{"key1": "value1", "key2": "value2"}
		json, _ := json.MarshalIndent(resp, "", " ")
		gomock.InOrder(
			mockRespWriter.EXPECT().Header().Return(headers),
			mockRespWriter.EXPECT().WriteHeader(http.StatusOK),
			mockRespWriter.EXPECT().Write(json),
		)
		assert.Nil(t, EncodeReply(context.Background(), mockRespWriter, resp))
		assert.Len(t, headers, 1)
		assert.Equal(t, "application/json; charset=utf-8", headers["Content-Type"][0])
	})
	t.Run("Non-JSON serializable response", func(t *testing.T) {
		headers := make(http.Header)
		var resp nonJsonable

		mockRespWriter.EXPECT().WriteHeader(http.StatusOK)
		mockRespWriter.EXPECT().Header().Return(headers)

		assert.Nil(t, EncodeReply(context.Background(), mockRespWriter, resp))
		assert.Equal(t, 1, len(headers))
	})
}

func TestErrorHandler(t *testing.T) {
	var mockCtrl = gomock.NewController(t)
	defer mockCtrl.Finish()

	mockRespWriter := mock.NewResponseWriter(mockCtrl)

	t.Run("ForbiddenError", func(t *testing.T) {
		mockRespWriter.EXPECT().WriteHeader(http.StatusForbidden)
		mockRespWriter.EXPECT().Write(gomock.Any())
		ErrorHandlerNoLog()(context.Background(), security.ForbiddenError{}, mockRespWriter)
	})
	t.Run("HTTPError", func(t *testing.T) {
		err := errorhandler.Error{
			Status:  123,
			Message: "abc",
		}
		mockRespWriter.EXPECT().WriteHeader(err.Status)
		mockRespWriter.EXPECT().Write([]byte(err.Message))
		ErrorHandlerNoLog()(context.Background(), err, mockRespWriter)
	})
	t.Run("DetailedError", func(t *testing.T) {
		var mockError = mock.NewDetailedError(mockCtrl)
		var status = 403
		var message = "error.message"
		mockError.EXPECT().Status().Return(status)
		mockError.EXPECT().ErrorMessage().Return(message)
		mockRespWriter.EXPECT().WriteHeader(status)
		mockRespWriter.EXPECT().Write([]byte(message))
		ErrorHandlerNoLog()(context.Background(), mockError, mockRespWriter)
	})
	t.Run("ratelimit.ErrLimited", func(t *testing.T) {
		mockRespWriter.EXPECT().WriteHeader(http.StatusTooManyRequests)
		ErrorHandlerNoLog()(context.Background(), ratelimit.ErrLimited, mockRespWriter)
	})
	t.Run("Internal server error", func(t *testing.T) {
		message := "500"
		mockRespWriter.EXPECT().WriteHeader(http.StatusInternalServerError)
		ErrorHandlerNoLog()(context.Background(), errors.New(message), mockRespWriter)
	})
}
