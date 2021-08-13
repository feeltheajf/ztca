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

func new(err interface{}, etype ErrorType, a ...interface{}) Error {
	var text string
	switch t := err.(type) {
	case string:
		text = fmt.Sprintf(t, a...)
	default:
		text = fmt.Sprintf("%v", err)
	}
	return Error{
		Text: text,
		Type: etype,
	}
}

// NotFound signals that the requested object doesn't exist
func NotFound(err interface{}, a ...interface{}) Error {
	return new(err, ErrNotFound, a...)
}

// InvalidParameter signals that the user input is invalid
func InvalidParameter(err interface{}, a ...interface{}) Error {
	return new(err, ErrInvalidParameter, a...)
}

// Unauthorized is used to signify that the user is not authorized to perform a specific action
func Unauthorized(err interface{}, a ...interface{}) Error {
	return new(err, ErrUnauthorized, a...)
}

// Forbidden signals that the requested action cannot be performed under any circumstances
func Forbidden(err interface{}, a ...interface{}) Error {
	return new(err, ErrForbidden, a...)
}

// Conflict signals that the requested action cannot be performed in current system state
func Conflict(err interface{}, a ...interface{}) Error {
	return new(err, ErrConflict, a...)
}

// NotImplemented signals that the requested action/feature is not implemented on the system as configured
func NotImplemented(err interface{}, a ...interface{}) Error {
	return new(err, ErrNotImplemented, a...)
}

// Unknown signals that the kind of error that occurred is not known
func Unknown(err interface{}, a ...interface{}) Error {
	return new(err, ErrUnknown, a...)
}
