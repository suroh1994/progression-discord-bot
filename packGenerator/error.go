package packGenerator

import "errors"

// ErrSetNotFound is returned when a given set code does not match an existing set code.
var ErrSetNotFound = errors.New("set not found")
