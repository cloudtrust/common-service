package errorhandler

import (
	"fmt"
	"net/http"
	"strings"
)

// Constants for error messages
const (
	MsgErrMissingParam = "missingParameter"

	MsgErrInvalidQueryParam         = "invalidQueryParameter"
	MsgErrInvalidPathParam          = "invalidPathParameter"
	MsgErrInvalidParam              = "invalidParameter"
	MsgErrOpNotPermitted            = "operationNotPermitted"
	MsgErrDisabledEndpoint          = "disabledEndpoint"
	MsgErrInvalidLength             = "invalidLength"
	MsgErrDecryptionKeyNotAvailable = "decryptionKeyNotAvailable"
	MsgErrUnknown                   = "unknownError"

	EncryptDecrypt = "encryptOrDecrypt"
	Ciphertext     = "ciphertext"
	AuthHeader     = "authorizationHeader"
	BasicToken     = "basicToken"
	BearerToken    = "bearerToken"
	Token          = "token"
	Level          = "level"
	JSONExpected   = "JSONExpected"
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

// DetailedError interface
type DetailedError interface {
	Error() string
	Status() int
	ErrorMessage() string
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

// CreateForbiddenError creates an error relative to a not allowed request
func CreateForbiddenError(messageParts ...string) Error {
	return Error{
		Status:  http.StatusForbidden,
		Message: concatTo2StaticValues(GetEmitter(), "forbidden", messageParts...),
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

// UnauthorizedError when access to the service is disallowed
type UnauthorizedError struct{}

func (e UnauthorizedError) Error() string {
	return "UnauthorizedError: Unauthorized"
}

func concatTo2StaticValues(val1 string, val2 string, values ...string) string {
	var msg = append([]string{val1}, val2)
	msg = append(msg, values...)
	return strings.Join(msg, ".")
}
