package errdefs

import (
	"net/http"
)

// GetStatusCode retrieves status code from error message
func GetStatusCode(err error) int {
	if err == nil {
		return http.StatusInternalServerError
	}

	var statusCode int
	switch {
	case IsNotFound(err):
		statusCode = http.StatusNotFound
	case IsInvalidParameter(err):
		statusCode = http.StatusBadRequest
	case IsUnauthorized(err):
		statusCode = http.StatusUnauthorized
	case IsForbidden(err):
		statusCode = http.StatusForbidden
	case IsConflict(err):
		statusCode = http.StatusConflict
	case IsNotImplemented(err):
		statusCode = http.StatusNotImplemented
	default:
		if e, ok := err.(Error); ok {
			if e.Cause != nil {
				return GetStatusCode(e.Cause)
			}
		}
	}

	if statusCode == 0 {
		statusCode = http.StatusInternalServerError
	}
	return statusCode
}

// FromStatusCode creates an errdef error, based on the provided HTTP status-code
func FromStatusCode(statusCode int, err error) error {
	if err == nil {
		return err
	}

	var cause error
	if known, ok := err.(Error); ok {
		cause = known.Cause
	}

	text := err.Error()
	switch statusCode {
	case http.StatusNotFound:
		err = NotFound(text).CausedBy(cause)
	case http.StatusBadRequest:
		err = InvalidParameter(text).CausedBy(cause)
	case http.StatusUnauthorized:
		err = Unauthorized(text).CausedBy(cause)
	case http.StatusForbidden:
		err = Forbidden(text).CausedBy(cause)
	case http.StatusConflict:
		err = Conflict(text).CausedBy(cause)
	case http.StatusNotImplemented:
		err = NotImplemented(text).CausedBy(cause)
	default:
		err = Unknown(text).CausedBy(cause)
	}
	return err
}
