package league

import "errors"

// ErrPlayerAlreadyJoined is returned when a player attempts to join a league, which they are already part of.
var ErrPlayerAlreadyJoined = errors.New("player already joined the league")

// ErrMatchAlreadyReported is returned when a player attempts to report a match result for a match with an existing result
var ErrMatchAlreadyReported = errors.New("match has already been reported")

// ErrInvalidMatchResult is returned when a player attempts to report a match result, which is not a valid outcome. Currently, this is only 0/0/0.
var ErrInvalidMatchResult = errors.New("invalid match result")

// ErrPlayerAlreadyDropped is returned when a player attempts to drop from a league, which they have already dropped from.
var ErrPlayerAlreadyDropped = errors.New("player already dropped from the league")
