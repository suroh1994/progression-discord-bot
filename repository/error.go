package repository

import "errors"

// ErrPairingNotFound is returned when either a pairing between the given players never existed or has already been reported.
var ErrPairingNotFound = errors.New("pairing not found")
