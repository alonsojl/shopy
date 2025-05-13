package errorx

import (
	"fmt"
)

type (
	// ErrorCode defines supported error codes.
	ErrorCode uint

	// Error interface represents an wrap error.
	Error interface {
		Error() string
		Wrap(err error) Error
		Unwrap() error
		Code() ErrorCode
	}

	// ObjectError is the default wrap error
	// that implements the Error interface.
	ObjectError struct {
		err  error
		msg  string
		code ErrorCode
	}
)

// WrapErrorf returns a wrapped error.
func WrapErrorf(err error, code ErrorCode, format string, a ...interface{}) Error {
	return &ObjectError{
		err:  err,
		code: code,
		msg:  fmt.Sprintf(format, a...),
	}
}

// NewErrorf instantiates a new error.
func NewErrorf(code ErrorCode, format string, a ...interface{}) Error {
	return WrapErrorf(nil, code, format, a...)
}

// Error returns the message, when wrapping errors the wrapped error is returned.
func (e *ObjectError) Error() string {
	return e.msg
}

// SetError set the error's origin.
func (e *ObjectError) Wrap(err error) Error {
	e.err = err
	return e
}

// Unwrap returns the wrapped error, if any.
func (e *ObjectError) Unwrap() error {
	return e.err
}

// Code returns the code representing this error.
func (e *ObjectError) Code() ErrorCode {
	return e.code
}
