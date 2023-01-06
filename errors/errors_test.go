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

type myDetailedError struct{}

func (m myDetailedError) Error() string {
	return "200 component.my_message"
}

func (m myDetailedError) Status() int {
	return 200
}

func (m myDetailedError) ErrorMessage() string {
	return "component.my_message"
}

func TestIsDetailedError(t *testing.T) {
	var err = myDetailedError{}
	assert.False(t, IsDetailedError(nil, 200, "my_message"))
	assert.False(t, IsDetailedError(err, 204, "another_message"))
	assert.False(t, IsDetailedError(err, 204, "my_message"))
	assert.False(t, IsDetailedError(err, 200, "another_message"))
	assert.True(t, IsDetailedError(err, 200, "my_message"))
}
