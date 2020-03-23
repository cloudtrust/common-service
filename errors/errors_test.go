package errorhandler

import (
	"errors"
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

func TestCommonErrorOrDefault(t *testing.T) {
	var defErr = errors.New("default error")
	t.Run("First param is not an errorhandler.Error", func(t *testing.T) {
		var err = errors.New("any error")
		assert.Equal(t, defErr, CommonErrorOrDefault(err, defErr))
	})
	t.Run("First param is an errorhandler.Error", func(t *testing.T) {
		var err = CreateNotFoundError("xxx")
		assert.Equal(t, err, CommonErrorOrDefault(err, defErr))
	})
}
