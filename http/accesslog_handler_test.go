package http

import (
	"context"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/cloudtrust/common-service/v2/http/mock"
	"github.com/golang/mock/gomock"
	"github.com/gorilla/mux"
	"github.com/stretchr/testify/assert"
)

func TestAccessLogHTTPHandler(t *testing.T) {
	var mockCtrl = gomock.NewController(t)
	defer mockCtrl.Finish()

	mockLogger := mock.NewLogger(mockCtrl)
	mockLogger.EXPECT().Log("method", "GET", "uri", "/path/to/resource", "status_code", http.StatusOK, "size", 2, "request_duration", gomock.Any())

	r := mux.NewRouter()

	r.Handle("/path/to/resource", makeHandler(func(_ context.Context, _ any) (response any, err error) {
		return GenericResponse{
			StatusCode:       http.StatusOK,
			Headers:          map[string]string{"X-Test": "here"},
			MimeContent:      nil,
			JSONableResponse: make([]int, 0),
		}, nil
	}))

	h := MakeAccessLogHandler(mockLogger, r)

	ts := httptest.NewServer(h)
	defer ts.Close()

	{
		res, err := http.Get(ts.URL + "/path/to/resource")
		assert.Nil(t, err)
		assert.Equal(t, http.StatusOK, res.StatusCode)
	}

}
