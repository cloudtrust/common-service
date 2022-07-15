package healthcheck

import (
	"strings"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
)

func TestHTTPHealth(t *testing.T) {
	{
		var checker = newHTTPEndpointChecker("github", "http://github.com/", 10*time.Second, 200, 10*time.Second)
		var status = checker.CheckStatus()
		assert.Equal(t, *status.State, "UP")
		assert.Nil(t, status.Message)

		// Call twice uses the cached result
		checker.CheckStatus()
		assert.Equal(t, *status.State, "UP")
	}

	{
		var checker = newHTTPEndpointChecker("dummy", "http://dummy.server.elca.ch/", 10*time.Second, 200, 10*time.Second)
		var status = checker.CheckStatus()
		assert.Equal(t, *status.State, "DOWN")
		assert.NotNil(t, status.Message)
		assert.True(t, strings.HasPrefix(*status.Message, "Can't hit target"))
	}

	{
		var checker = newHTTPEndpointChecker("dummy", "http://www.elca.ch/not/found", 10*time.Second, 200, 10*time.Second)
		var status = checker.CheckStatus()
		assert.Equal(t, *status.State, "DOWN")
		assert.NotNil(t, status.Message)
		assert.True(t, strings.Contains(*status.Message, "404"))
	}
}
