package errdefs

// IsNotFound returns true if the passed in error is an errNotFound
func IsNotFound(err error) bool {
	if e, ok := err.(Error); ok {
		return e.Type == ErrNotFound
	}
	return false
}

// IsInvalidParameter returns true if the passed in error is an errInvalidParameter
func IsInvalidParameter(err error) bool {
	if e, ok := err.(Error); ok {
		return e.Type == ErrInvalidParameter
	}
	return false
}

// IsUnauthorized returns true if the passed in error is an errUnauthorized
func IsUnauthorized(err error) bool {
	if e, ok := err.(Error); ok {
		return e.Type == ErrUnauthorized
	}
	return false
}

// IsForbidden returns true if the passed in error is an errForbidden
func IsForbidden(err error) bool {
	if e, ok := err.(Error); ok {
		return e.Type == ErrForbidden
	}
	return false
}

// IsConflict returns true if the passed in error is an errConflict
func IsConflict(err error) bool {
	if e, ok := err.(Error); ok {
		return e.Type == ErrConflict
	}
	return false
}

// IsNotImplemented returns true if the passed in error is an errNotImplemented
func IsNotImplemented(err error) bool {
	if e, ok := err.(Error); ok {
		return e.Type == ErrNotImplemented
	}
	return false
}

// IsUnknown returns true if the passed in error is an errUnknown
func IsUnknown(err error) bool {
	if e, ok := err.(Error); ok {
		return e.Type == ErrUnknown
	}
	return false
}
