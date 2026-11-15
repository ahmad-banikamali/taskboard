package domain

import "fmt"

type Status string

const (
	StatusTodo       Status = "TODO"
	StatusInProgress Status = "IN_PROGRESS"
	StatusDone       Status = "DONE"
)

// ParseStatus turns text (for example from HTTP) into a Status.
func ParseStatus(raw string) (Status, error) {
	switch s := Status(raw); s {
	case StatusTodo, StatusInProgress, StatusDone:
		return s, nil
	}
	return "", ErrInvalidStatus
}

// CanMoveTo says if the move is allowed.
// Allowed: TODO -> IN_PROGRESS, IN_PROGRESS -> DONE, IN_PROGRESS -> TODO.
// Everything else is not allowed, including staying in the same status.
func (s Status) CanMoveTo(next Status) bool {
	switch s {
	case StatusTodo:
		return next == StatusInProgress
	case StatusInProgress:
		return next == StatusDone || next == StatusTodo
	default: // DONE is final, and unknown values can't move
		return false
	}
}

// MoveTo returns the new status, or an error if the move is not allowed.
func (s Status) MoveTo(next Status) (Status, error) {
	if !s.CanMoveTo(next) {
		return s, fmt.Errorf("%w: %s -> %s", ErrInvalidTransition, s, next)
	}
	return next, nil
}
