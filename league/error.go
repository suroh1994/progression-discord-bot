package league

import "errors"

// ErrPlayerAlreadyJoined is returned when a player attempts to join a league, which they are already part of.
var ErrPlayerAlreadyJoined = errors.New("player already joined league")
