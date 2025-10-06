package league

import "errors"

// ErrPlayerAlreadyJoined is returned when a player attempts to join a league, which they are already part of.
var ErrMatchAlreadyReported = errors.New("match has already been reported")
var ErrInvalidMatchResult = errors.New("invalid match result")
