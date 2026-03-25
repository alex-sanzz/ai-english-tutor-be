package apperr

import (
	"errors"
	"net/http"
)

type Error struct {
	Status  int
	Code    string
	Message string
	Err     error
}

func (e *Error) Error() string {
	if e.Err != nil {
		return e.Err.Error()
	}

	return e.Message
}

func (e *Error) Unwrap() error { return e.Err }

func New(status int, message string, code string, err error) *Error {
	return &Error{
		Status:  status,
		Code:    code,
		Message: message,
		Err:     err,
	}
}

func BadRequest(code, message string, err error) *Error {
	return &Error{
		Status: http.StatusBadRequest,
		Code: code,
		Message: message,
		Err: err,
	}
}

func Unauthorized(code, message string, err error) *Error {
	return &Error{
		Status: http.StatusUnauthorized,
		Code: code,
		Message: message,
		Err: err,
	}
}

func Forbidden(code, message string, err error) *Error {
	return &Error{
		Status: http.StatusForbidden,
		Code: code,
		Message: message,
		Err: err,
	}
}

func NotFound(code, message string, err error) *Error {
	return &Error{
		Status: http.StatusNotFound,
		Code: code,
		Message: message,
		Err: err,
	}
}

func Conflict(code, message string, err error) *Error {
	return &Error{
		Status: http.StatusConflict,
		Code: code,
		Message: message,
		Err: err,
	}
}

func Unprocessable(code, message string, err error) *Error {
	return &Error{
		Status: http.StatusUnprocessableEntity,
		Code: code,
		Message: message,
		Err: err,
	}
}

func Internal(err error) *Error{
	return &Error{
		Status: http.StatusInternalServerError,
		Code: "internal_error",
		Message: "internal server error",
		Err: err,
	}
}


func From(err error) *Error {
	if err == nil {
		return nil 
	}

	var ae *Error

	// errors.As(err, &ae) means err is already as apperr.Error type
	// SO just return it as it is

	if ae != nil && errors.As(err, &ae) {
		return ae 
	}

	// if the error is not apperr.Error type, then return it as internal error
	return Internal(err)
}

