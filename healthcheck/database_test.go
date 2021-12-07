package healthcheck

import (
	"errors"
	"testing"
	"time"

	"github.com/cloudtrust/common-service/v2/healthcheck/mock"
	"github.com/golang/mock/gomock"
	"github.com/stretchr/testify/assert"
)

func TestDbHealthCheck(t *testing.T) {
	var mockCtrl = gomock.NewController(t)
	defer mockCtrl.Finish()

	var mockDB = mock.NewHealthDatabase(mockCtrl)

	{
		var dbChecker = newDatabaseChecker("alias", mockDB, 10*time.Second)
		mockDB.EXPECT().Ping().Return(nil)
		var res = dbChecker.CheckStatus()
		assert.NotNil(t, res.Connection)
		assert.Equal(t, "established", *res.Connection)
	}

	{
		var dbChecker = newDatabaseChecker("alias", mockDB, 10*time.Second)
		var errMsg = "Error message"
		var err = errors.New(errMsg)

		mockDB.EXPECT().Ping().Return(err).Times(1)

		var res = dbChecker.CheckStatus()
		assert.NotNil(t, res.Message)
		assert.Equal(t, errMsg, *res.Message)

		// Mock is configured to be called only once... A new call would let the test fail but it wont fail as result is cached
		res = dbChecker.CheckStatus()
		assert.Equal(t, errMsg, *res.Message)
	}
}
