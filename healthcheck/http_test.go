package healthcheck

import (
	"strings"
	"testing"
	"time"

	"github.com/cloudtrust/common-service/v2/healthcheck/mock"
	"github.com/stretchr/testify/assert"
	"go.uber.org/mock/gomock"
)

func TestHTTPHealth(t *testing.T) {
	mockCtrl := gomock.NewController(t)
	defer mockCtrl.Finish()

	mockTime := mock.NewTimeProvider(mockCtrl)
	mockTime.EXPECT().Now().Return(testTime).AnyTimes()

	{
		var checker = newHTTPEndpointChecker("github", "https://github.com/", 10*time.Second, 200, 10*time.Second, mockTime)
		var status = checker.CheckStatus()
		assert.Equal(t, *status.State, "UP")
		assert.Nil(t, status.Message)

		// Call twice uses the cached result
		checker.CheckStatus()
		assert.Equal(t, *status.State, "UP")
	}

	{
		var checker = newHTTPEndpointChecker("dummy", "https://dummy.server.elca.ch/", 10*time.Second, 200, 10*time.Second, mockTime)
		var status = checker.CheckStatus()
		assert.Equal(t, *status.State, "DOWN")
		assert.NotNil(t, status.Message)
		assert.True(t, strings.HasPrefix(*status.Message, "Can't hit target"))
	}

	{
		var checker = newHTTPEndpointChecker("dummy", "https://github.com/not/found", 10*time.Second, 200, 10*time.Second, mockTime)
		var status = checker.CheckStatus()
		assert.Equal(t, *status.State, "DOWN")
		assert.NotNil(t, status.Message)
		assert.True(t, strings.Contains(*status.Message, "404"))
	}
}
