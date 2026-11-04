package domain

import "errors"

// Centralized domain-level errors translated into HTTP status codes by the handler middleware.
var (
	ErrZoneNotFound    = errors.New("requested office zone was not found")
	ErrZoneLocked      = errors.New("zone is currently locked; complete the algorithmic quest first")
	ErrQuestNotFound   = errors.New("algorithmic quest for the specified zone was not found")
	ErrInvalidSequence = errors.New("submitted code block sequence is incorrect")
	ErrHotspotNotFound = errors.New("requested hotspot was not found")
	ErrNotAnEasterEgg  = errors.New("specified hotspot is not a claimable easter egg")
	ErrInvalidInput    = errors.New("invalid request payload or parameters")
	ErrUnauthorized    = errors.New("missing or invalid student session identity")
)
