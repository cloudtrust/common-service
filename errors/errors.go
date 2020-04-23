package errorhandler

import (
	"fmt"
	"net/http"
)

// Constants for error messages
const (
	MsgErrMissingParam = "missingParameter"

	MsgErrInvalidQueryParam = "invalidQueryParameter"
	MsgErrInvalidPathParam  = "invalidPathParameter"
	MsgErrInvalidParam      = "invalidParameter"
	MsgErrOpNotPermitted    = "operationNotPermitted"
	MsgErrDisabledEndpoint  = "disabledEndpoint"

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
func CreateBadRequestError(messageKey string) Error {
	return Error{
		Status:  http.StatusBadRequest,
		Message: GetEmitter() + "." + messageKey,
	}
}

// CreateNotAllowedError creates an error relative to a not allowed request
func CreateNotAllowedError(messageKey string) Error {
	return Error{
		Status:  http.StatusMethodNotAllowed,
		Message: GetEmitter() + "." + messageKey,
	}
}

// CreateNotFoundError creates an error relative to a not found request
func CreateNotFoundError(messageKey string) Error {
	return Error{
		Status:  http.StatusNotFound,
		Message: GetEmitter() + "." + messageKey,
	}
}

// CreateEndpointNotEnabled creates an error relative to an attempt to access a not enabled endpoint
func CreateEndpointNotEnabled(param string) Error {
	return Error{
		Status:  http.StatusConflict,
		Message: fmt.Sprintf("%s.%s.%s", GetEmitter(), MsgErrDisabledEndpoint, param),
	}
}
