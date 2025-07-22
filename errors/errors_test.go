package errorhandler

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestHTTPResponse(t *testing.T) {
	t.Run("Missing parameter error", func(t *testing.T) {
		var error = CreateMissingParameterError("parameter")
		assert.Contains(t, error.Error(), error.Message)
	})
	t.Run("Bad request error", func(t *testing.T) {
		var error = CreateBadRequestError("message")
		assert.Contains(t, error.Error(), error.Message)
	})
	t.Run("Internal server error", func(t *testing.T) {
		var error = CreateInternalServerError("message")
		assert.Contains(t, error.Error(), error.Message)
	})
	t.Run("Invalid query parameter error", func(t *testing.T) {
		var error = CreateInvalidQueryParameterError("message")
		assert.Contains(t, error.Error(), error.Message)
	})
	t.Run("Invalid path parameter error", func(t *testing.T) {
		var error = CreateInvalidPathParameterError("message")
		assert.Contains(t, error.Error(), error.Message)
	})
	t.Run("Not allowed error", func(t *testing.T) {
		var error = CreateNotAllowedError("message")
		assert.Contains(t, error.Error(), error.Message)
	})
	t.Run("Not found error", func(t *testing.T) {
		var error = CreateNotFoundError("message")
		assert.Contains(t, error.Error(), error.Message)
	})
	t.Run("Forbidden error", func(t *testing.T) {
		var error = CreateForbiddenError("param1", "param2")
		assert.Contains(t, error.Error(), "forbidden.param1.param2")
	})
	t.Run("endpoint not enabled", func(t *testing.T) {
		var error = CreateEndpointNotEnabled("param1")
		assert.Contains(t, error.Error(), "disabledEndpoint.param1")
	})
}

func TestEmitter(t *testing.T) {
	var emitter = "component"
	SetEmitter(emitter)
	assert.Equal(t, GetEmitter(), emitter)
}
