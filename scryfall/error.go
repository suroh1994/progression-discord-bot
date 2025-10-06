package scryfall

import "errors"

// ErrMoreThanOneCardFound is returned when a given name matches more than one card.
var ErrMoreThanOneCardFound = errors.New("ambiguous name given, more than one card found")
