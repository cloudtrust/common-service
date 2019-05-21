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

// CreateMissingParameterError creates an error relative to a missing mandatory parameter
func CreateMissingParameterError(name string) Error {
	return Error{
		Status:  http.StatusBadRequest,
		Message: fmt.Sprintf("Missing mandatory parameter %s", name),
	}
}

// CreateInvalidQueryParameterError creates an error relative to a invalid query parameter
func CreateInvalidQueryParameterError(paramName string) Error {
	return Error{
		Status:  http.StatusBadRequest,
		Message: fmt.Sprintf("Invalid query parameter %s", paramName),
	}
}

// CreateInvalidPathParameterError creates an error relative to a invalid path parameter
func CreateInvalidPathParameterError(paramName string) Error {
	return Error{
		Status:  http.StatusBadRequest,
		Message: fmt.Sprintf("Invalid path parameter %s", paramName),
	}
}

// CreateInvalidParameterError creates an error relative to a invalid parameter
func CreateBadRequestError(publicMessage string) Error {
	return Error{
		Status:  http.StatusBadRequest,
		Message: publicMessage,
	}
}