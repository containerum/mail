package errors

import (
	"fmt"
	"net/http"
)

type Error struct {
	Text string `json:"error"`
	Code int    `json:"code,omitempty"`
}

func (e *Error) Error() string {
	if e.Code == 0 {
		return e.Text
	}
	return fmt.Sprintf("description: %s, code: %d", e.Text, e.Code)
}

func New(text string) *Error {
	return &Error{
		Text: text,
	}
}

func NewWithCode(text string, code int) *Error {
	return &Error{
		Text: text,
		Code: code,
	}
}

func Format(format string, data ...interface{}) *Error {
	return &Error{
		Text: fmt.Sprintf(format, data...),
	}
}

func FormatWithCode(code int, format string, data ...interface{}) *Error {
	return &Error{
		Text: fmt.Sprintf(format, data...),
		Code: code,
	}
}

func ErrorWithHTTPStatus(err error) (int, []*Error) {
	switch err.(type) {
	case *AccessDeniedError:
		return http.StatusForbidden, []*Error{err.(*AccessDeniedError).Err}
	case *NotFoundError:
		return http.StatusNotFound, []*Error{err.(*NotFoundError).Err}
	case *BadRequestError:
		return http.StatusBadRequest, []*Error{err.(*BadRequestError).Err}
	case *AlreadyExistsError:
		return http.StatusConflict, []*Error{err.(*AlreadyExistsError).Err}
	case *InternalError:
		return http.StatusInternalServerError, []*Error{err.(*InternalError).Err}
	case *WebAPIError:
		return err.(*WebAPIError).StatusCode, []*Error{err.(*WebAPIError).Err}
	default:
		//TODO Do something with grpc errors
		/*if grpcStatus, ok := status.FromError(err); ok {
			if httpStatus, hasStatus := grpcutils.GRPCToHTTPCode[grpcStatus.Code()]; hasStatus {
				return httpStatus, []*Error{New(grpcStatus.Message())}
			}
			return http.StatusInternalServerError, []*Error{New(grpcStatus.Err().Error())}
		}*/
		return http.StatusInternalServerError, []*Error{New(err.Error())}
	}
}

// InternalError describes server errors which should not be exposed to client explicitly.
type InternalError struct {
	Err *Error
}

func (e *InternalError) Error() string {
	return e.Err.Error()
}

// AccessDeniedError describes error if client has no access to resource, method, etc.
type AccessDeniedError struct {
	Err *Error
}

func (e *AccessDeniedError) Error() string {
	return e.Err.Error()
}

// NotFoundError describes error returned if requested resource was not found
type NotFoundError struct {
	Err *Error
}

func (e *NotFoundError) Error() string {
	return e.Err.Error()
}

// BadRequestError describes error returned if request was malformed.
type BadRequestError struct {
	Err *Error
}

func (e *BadRequestError) Error() string {
	return e.Err.Error()
}

// AlreadyExistsError describes error returned if client attempts to create resource or register with username which already exists.
type AlreadyExistsError struct {
	Err *Error
}

func (e *AlreadyExistsError) Error() string {
	return e.Err.Error()
}

// WebAPIError describes error returned from web-api service.
type WebAPIError struct {
	Err        *Error
	StatusCode int
}

func (e *WebAPIError) Error() string {
	return e.Err.Error()
}
