package http

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestHTTPResponse(t *testing.T) {
	// Coverage
	var error = CreateMissingParameterError("parameter")
	assert.Contains(t, error.Error(), error.Message)
}
