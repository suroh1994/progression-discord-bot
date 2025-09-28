package repository

import "errors"

// ErrPairingNotFound is returned when either a pairing between the given players never existed or has already been reported.
var ErrPairingNotFound = errors.New("pairing not found")

// ErrPlayerNotFound is returned when a given userID cannot be matched to a valid player.
var ErrPlayerNotFound = errors.New("player not found")

// ErrNoActiveLeague is returned when no league is currently active.
var ErrNoActiveLeague = errors.New("no active league")

// ErrLeagueAlreadyOngoing is returned when a league cannot be started because another is already ongoing.
var ErrLeagueAlreadyOngoing = errors.New("another league is already ongoing")
