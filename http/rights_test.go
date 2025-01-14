package http

import (
	"bytes"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/cloudtrust/common-service/v2/http/mock"
	"github.com/gorilla/mux"
	"github.com/stretchr/testify/assert"
	"go.uber.org/mock/gomock"
)

func TestMakeRightsHandler(t *testing.T) {
	var mockCtrl = gomock.NewController(t)
	defer mockCtrl.Finish()

	mockAuthManager := mock.NewAuthorizationManager(mockCtrl)

	var rights = map[string]map[string]map[string]map[string]struct{}{
		"toe_administrator": {
			"GetUsers": {
				"master": {
					"*": {},
				},
			},
		},
	}

	mockAuthManager.EXPECT().GetRightsOfCurrentUser(gomock.Any()).Return(rights)

	r := mux.NewRouter()
	r.Handle("/rights", MakeRightsHandler(mockAuthManager))

	ts := httptest.NewServer(r)
	defer ts.Close()

	res, err := http.Get(ts.URL + "/rights")
	assert.Nil(t, err)
	assert.Equal(t, http.StatusOK, res.StatusCode)

	buf := new(bytes.Buffer)
	buf.ReadFrom(res.Body)
	body := buf.String()

	var response map[string]map[string]map[string]map[string]struct{}
	err = json.Unmarshal([]byte(body), &response)
	assert.Equal(t, response, rights)
	assert.Nil(t, err)
}
