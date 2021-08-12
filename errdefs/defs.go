package errdefs

import (
	"fmt"
)

// ErrorType represents a basic error type
type ErrorType string

// Known error types
const (
	ErrNotFound         ErrorType = "not found"
	ErrInvalidParameter ErrorType = "invalid parameter"
	ErrUnauthorized     ErrorType = "unauthorized"
	ErrForbidden        ErrorType = "forbidden"
	ErrConflict         ErrorType = "conflict"
	ErrNotImplemented   ErrorType = "not implemented"
	ErrUnknown          ErrorType = "unknown"
)

// Error is a generic error type
type Error struct {
	Text  string    `json:"text"`
	Type  ErrorType `json:"type"`
	Cause error     `json:"-"`
}

func (e Error) Error() string {
	s := e.Text
	if e.Cause != nil {
		s += " caused by " + e.Cause.Error()
	}
	return s
}

// CausedBy adds error cause
func (e Error) CausedBy(cause error) Error {
	e.Cause = cause
	return e
}

func new(err interface{}, etype ErrorType, args ...interface{}) Error {
	var text string
	text, ok := err.(string)
	if !ok {
		text = fmt.Sprintf("%v", err)
	}
	return Error{
		Text: text,
		Type: etype,
	}
}

// NotFound signals that the requested object doesn't exist
func NotFound(err interface{}, args ...interface{}) Error {
	return new(err, ErrNotFound, args...)
}

// InvalidParameter signals that the user input is invalid
func InvalidParameter(err interface{}, args ...interface{}) Error {
	return new(err, ErrInvalidParameter, args...)
}

// Unauthorized is used to signify that the user is not authorized to perform a specific action
func Unauthorized(err interface{}, args ...interface{}) Error {
	return new(err, ErrUnauthorized, args...)
}

// Forbidden signals that the requested action cannot be performed under any circumstances
func Forbidden(err interface{}, args ...interface{}) Error {
	return new(err, ErrForbidden, args...)
}

// Conflict signals that the requested action cannot be performed in current system state
func Conflict(err interface{}, args ...interface{}) Error {
	return new(err, ErrConflict, args...)
}

// NotImplemented signals that the requested action/feature is not implemented on the system as configured
func NotImplemented(err interface{}, args ...interface{}) Error {
	return new(err, ErrNotImplemented, args...)
}

// Unknown signals that the kind of error that occurred is not known
func Unknown(err interface{}, args ...interface{}) Error {
	return new(err, ErrUnknown, args...)
}
