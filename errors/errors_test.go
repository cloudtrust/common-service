package errorhandler

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestHTTPResponse(t *testing.T) {
	// Coverage
	var error = CreateMissingParameterError("parameter")
	assert.Contains(t, error.Error(), error.Message)

	error = CreateBadRequestError("message")
	assert.Contains(t, error.Error(), error.Message)

	error = CreateInternalServerError("message")
	assert.Contains(t, error.Error(), error.Message)

	error = CreateInvalidQueryParameterError("message")
	assert.Contains(t, error.Error(), error.Message)

	error = CreateInvalidPathParameterError("message")
	assert.Contains(t, error.Error(), error.Message)

	error = CreateNotAllowedError("message")
	assert.Contains(t, error.Error(), error.Message)

	error = CreateNotFoundError("message")
	assert.Contains(t, error.Error(), error.Message)
}

func TestEmitter(t *testing.T) {
	var emitter = "component"
	SetEmitter(emitter)
	assert.Equal(t, GetEmitter(), emitter)
}
