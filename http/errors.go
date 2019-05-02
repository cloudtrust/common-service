package http

import (
	"fmt"
	"net/http"
)

// Error can be returned by the API endpoints
type Error struct {
	Status  int
	Message string
}

func (e Error) Error() string {
	return fmt.Sprintf("%d %s", e.Status, e.Message)
}

// CreateMissingParameterError creates a HTTPResponse for an error relative to a missing mandatory parameter
func CreateMissingParameterError(name string) Error {
	return Error{
		Status:  http.StatusBadRequest,
		Message: fmt.Sprintf("Missing mandatory parameter %s", name),
	}
}
