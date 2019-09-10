package errorhandler

import (
	"fmt"
	"net/http"
)

const (
	MsgErrMissingParam = "missingParameter"

	MsgErrInvalidQueryParam = "invalidQueryParameter"
	MsgErrInvalidPathParam  = "invalidPathParameter"
	MsgErrInvalidParam      = "invalidParameter"
	MsgErrOpNotPermitted    = "operationNotPermitted"

	AuthHeader   = "authorizationHeader"
	BasicToken   = "basicToken"
	BearerToken  = "bearerToken"
	Token        = "token"
	Level        = "level"
	JSONExpected = "JSONExpected"
)

var emitter string

// SetEmitter is called by the components that call the methods of the common service
func SetEmitter(e string) {
	emitter = e
}

// GetEmitter allows to see who called the common service
func GetEmitter() string {
	return emitter
}

// Error can be returned by the API endpoints
type Error struct {
	Status  int
	Message string
}

func (e Error) Error() string {
	return fmt.Sprintf("%d %s", e.Status, e.Message)
}

// CreateInternalServerError creates an error relative to an internal server error
func CreateInternalServerError(message string) Error {
	return Error{
		Status:  http.StatusInternalServerError,
		Message: GetEmitter() + "." + message,
	}
}

// CreateMissingParameterError creates an error relative to a missing mandatory parameter
func CreateMissingParameterError(name string) Error {
	return Error{
		Status:  http.StatusBadRequest,
		Message: fmt.Sprintf("%s.%s.%s", GetEmitter(), MsgErrMissingParam, name),
	}
}

// CreateInvalidQueryParameterError creates an error relative to a invalid query parameter
func CreateInvalidQueryParameterError(paramName string) Error {
	return Error{
		Status:  http.StatusBadRequest,
		Message: fmt.Sprintf("%s.%s.%s", GetEmitter(), MsgErrInvalidQueryParam, paramName),
	}
}

// CreateInvalidPathParameterError creates an error relative to a invalid path parameter
func CreateInvalidPathParameterError(paramName string) Error {
	return Error{
		Status:  http.StatusBadRequest,
		Message: fmt.Sprintf("%s.%s.%s", GetEmitter(), MsgErrInvalidPathParam, paramName),
	}
}

// CreateBadRequestError creates an error relative to a bad request
func CreateBadRequestError(publicMessage string) Error {
	return Error{
		Status:  http.StatusBadRequest,
		Message: GetEmitter() + "." + publicMessage,
	}
}

// CreateNotAllowedError creates an error relative to a not allowed request
func CreateNotAllowedError(publicMessage string) Error {
	return Error{
		Status:  http.StatusMethodNotAllowed,
		Message: GetEmitter() + "." + publicMessage,
	}
}

// CreateNotFoundError creates an error relative to a not found request
func CreateNotFoundError(publicMessage string) Error {
	return Error{
		Status:  http.StatusNotFound,
		Message: GetEmitter() + "." + publicMessage,
	}
}