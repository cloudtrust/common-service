package http

//go:generate mockgen -destination=./mock/responsewriter.go -package=mock -mock_names=ResponseWriter=ResponseWriter net/http ResponseWriter

import (
	"bytes"
	"context"
	"crypto/tls"
	"encoding/json"
	"errors"
	"fmt"
	"io/ioutil"
	"net/http"
	"net/http/httptest"
	"net/url"
	"testing"

	"github.com/cloudtrust/common-service/http/mock"
	"github.com/go-kit/kit/endpoint"
	"github.com/go-kit/kit/ratelimit"
	http_transport "github.com/go-kit/kit/transport/http"
	"github.com/golang/mock/gomock"
	"github.com/gorilla/mux"
	"github.com/stretchr/testify/assert"
)

type nonJsonable struct {
}

func (nj nonJsonable) MarshalJSON() ([]byte, error) {
	return nil, errors.New("Non serializable")
}

func TestHTTPManagementHandler(t *testing.T) {
	var mockCtrl = gomock.NewController(t)
	defer mockCtrl.Finish()

	var e endpoint.Endpoint
	e = func(ctx context.Context, req interface{}) (interface{}, error) {
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
		func(ctx context.Context, req *http.Request) (interface{}, error) {
			pathParams := []string{"pathParameter1", "pathParameter2", "pathParameter3"}
			queryParams := []string{"queryParameter1", "queryParameter2", "queryParameter3"}
			return DecodeRequest(ctx, req, pathParams, queryParams)
		},
		EncodeReply,
		http_transport.ServerErrorEncoder(ErrorHandler),
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
}

func checkPathParameter(t *testing.T, url string, expected string) {
	res, err := http.Get(url)
	assert.Nil(t, err)
	assert.Equal(t, http.StatusOK, res.StatusCode)

	buf := new(bytes.Buffer)
	buf.ReadFrom(res.Body)
	assert.Equal(t, `"`+expected+`"`, buf.String())
}

func genericTestDecodeRequest(ctx context.Context, tls *tls.ConnectionState, xFwdProto *string, rawQuery string) map[string]string {
	input := "the body"
	var req http.Request
	var url url.URL
	url.RawQuery = rawQuery
	req.Host = "localhost"
	req.TLS = tls
	req.Header = make(http.Header)
	if xFwdProto != nil {
		req.Header.Set("X-Forwarded-Proto", *xFwdProto)
	}
	req.Body = ioutil.NopCloser(bytes.NewBufferString(input))
	req.URL = &url
	pathParams := []string{"pathParam1", "pathParam2"}
	queryParams := []string{"queryParam1", "queryParam2"}
	r, _ := DecodeRequest(ctx, &req, pathParams, queryParams)
	return r.(map[string]string)
}

func TestDecodeRequestHTTP(t *testing.T) {
	request := genericTestDecodeRequest(context.Background(), nil, nil, "")

	// Minimum parameters are scheme, host and body
	assert.Equal(t, 3, len(request))
	assert.Equal(t, "localhost", request["host"])
	assert.Equal(t, "http", request["scheme"])
	assert.Equal(t, "the body", request["body"])
}

func TestDecodeRequestHTTPS(t *testing.T) {
	var reqConnState tls.ConnectionState

	request := genericTestDecodeRequest(context.Background(), &reqConnState, nil, "")

	// Minimum parameters are scheme, host and body
	assert.Equal(t, 3, len(request))
	assert.Equal(t, "localhost", request["host"])
	assert.Equal(t, "https", request["scheme"])
	assert.Equal(t, "the body", request["body"])
}

func TestDecodeRequestForwardProto(t *testing.T) {
	proto := "ftp"

	request := genericTestDecodeRequest(context.Background(), nil, &proto, "")

	// Minimum parameters are scheme, host and body
	assert.Equal(t, 3, len(request))
	assert.Equal(t, "localhost", request["host"])
	assert.Equal(t, proto, request["scheme"])
	assert.Equal(t, "the body", request["body"])
}

func TestDecodeRequestQueryParams(t *testing.T) {
	value := "valueX"
	request := genericTestDecodeRequest(context.Background(), nil, nil, "queryParam2="+value)

	// Minimum parameters are scheme, host and body
	assert.Equal(t, 4, len(request))
	assert.Equal(t, value, request["queryParam2"])
}

func TestEncodeReplyNilResponse(t *testing.T) {
	var mockCtrl = gomock.NewController(t)
	defer mockCtrl.Finish()

	mockRespWriter := mock.NewResponseWriter(mockCtrl)

	// Nil response
	{
		mockRespWriter.EXPECT().WriteHeader(http.StatusOK).Times(1)
		mockRespWriter.EXPECT().Header().Times(0)
		assert.Nil(t, EncodeReply(context.Background(), mockRespWriter, nil))
	}
}

func TestEncodeReplyJsonableResponse(t *testing.T) {
	var mockCtrl = gomock.NewController(t)
	defer mockCtrl.Finish()

	mockRespWriter := mock.NewResponseWriter(mockCtrl)

	// JSON serializable response
	{
		headers := make(http.Header)
		resp := map[string]string{"key1": "value1", "key2": "value2"}
		json, _ := json.MarshalIndent(resp, "", " ")

		mockRespWriter.EXPECT().WriteHeader(http.StatusOK).Times(1)
		mockRespWriter.EXPECT().Header().Return(headers).Times(1)
		mockRespWriter.EXPECT().Write(json).Times(1)

		assert.Nil(t, EncodeReply(context.Background(), mockRespWriter, resp))
		assert.Equal(t, 1, len(headers))
	}
}

func TestEncodeReplyNonJsonableResponse(t *testing.T) {
	var mockCtrl = gomock.NewController(t)
	defer mockCtrl.Finish()

	mockRespWriter := mock.NewResponseWriter(mockCtrl)

	// Non-JSON serializable response
	{
		headers := make(http.Header)
		var resp nonJsonable

		mockRespWriter.EXPECT().WriteHeader(http.StatusOK).Times(1)
		mockRespWriter.EXPECT().Header().Return(headers).Times(1)
		mockRespWriter.EXPECT().Write(gomock.Any()).Times(0)

		assert.Nil(t, EncodeReply(context.Background(), mockRespWriter, resp))
		assert.Equal(t, 1, len(headers))
	}
}

func TestErrorHandler(t *testing.T) {
	var mockCtrl = gomock.NewController(t)
	defer mockCtrl.Finish()

	mockRespWriter := mock.NewResponseWriter(mockCtrl)

	// HTTPError
	{
		err := Error{
			Status:  123,
			Message: "abc",
		}
		mockRespWriter.EXPECT().WriteHeader(err.Status).Times(1)
		mockRespWriter.EXPECT().Write([]byte(err.Message)).Times(1)
		ErrorHandler(context.Background(), err, mockRespWriter)
	}

	// ratelimit.ErrLimited
	{
		mockRespWriter.EXPECT().WriteHeader(http.StatusTooManyRequests).Times(1)
		mockRespWriter.EXPECT().Write(gomock.Any()).Times(0)
		ErrorHandler(context.Background(), ratelimit.ErrLimited, mockRespWriter)
	}

	// Internal server error
	{
		message := "500"
		mockRespWriter.EXPECT().WriteHeader(http.StatusInternalServerError).Times(1)
		mockRespWriter.EXPECT().Write([]byte(message)).Times(1)
		ErrorHandler(context.Background(), errors.New(message), mockRespWriter)
	}
}
