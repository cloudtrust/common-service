package http

import (
	"bytes"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"

	"github.com/gorilla/mux"
	"github.com/stretchr/testify/assert"
)

func TestMakeVersionHandler(t *testing.T) {
	var componentName = "name"
	var componentID = "123-456-789"
	var version = "1.0.0b"
	var environment = "environment"
	var gitCommit = "common-service/master 0.99alpha"

	r := mux.NewRouter()
	r.Handle("/version", MakeVersionHandler(componentName, componentID, version, environment, gitCommit))

	ts := httptest.NewServer(r)
	defer ts.Close()

	res, err := http.Get(ts.URL + "/version")
	assert.Nil(t, err)
	assert.Equal(t, http.StatusOK, res.StatusCode)

	buf := new(bytes.Buffer)
	buf.ReadFrom(res.Body)
	body := buf.String()
	assert.True(t, strings.Contains(body, componentName))
	assert.True(t, strings.Contains(body, componentID))
	assert.True(t, strings.Contains(body, version))
	assert.True(t, strings.Contains(body, environment))
	assert.True(t, strings.Contains(body, gitCommit))
}
