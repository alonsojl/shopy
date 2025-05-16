package errorx

import (
	"fmt"
)

type (
	// Code defines supported error codes.
	Code uint

	// Error interface represents an wrap error.
	Error interface {
		Error() string
		Message() string
		Wrap(err error) Error
		Unwrap() error
		Code() Code
	}

	// ObjectError is the default wrap error
	// that implements the Error interface.
	ObjectError struct {
		err  error
		msg  string
		code Code
	}
)

// WrapErrorf returns a wrapped error.
func WrapErrorf(err error, code Code, format string, a ...any) Error {
	return &ObjectError{
		err:  err,
		code: code,
		msg:  fmt.Sprintf(format, a...),
	}
}

// NewErrorf instantiates a new error.
func NewErrorf(code Code, format string, a ...any) Error {
	return WrapErrorf(nil, code, format, a...)
}

// Error returns the message, when wrapping errors the wrapped error is returned.
func (e *ObjectError) Error() string {
	if e.err != nil {
		return fmt.Sprintf("%s: %v", e.msg, e.err)
	}
	return e.msg
}

// Message return the error's message.
func (e *ObjectError) Message() string {
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
func (e *ObjectError) Code() Code {
	return e.code
}
