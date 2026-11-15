package domain

import "errors"

// Domain errors live in core. Adapters translate them later
// (for example, HTTP turns ErrInvalidTitle into 400).
var (
	ErrInvalidTitle      = errors.New("title must be 1 to 100 characters")
	ErrInvalidStatus     = errors.New("unknown status")
	ErrInvalidTransition = errors.New("status change not allowed")

	ErrInvalidTeamName = errors.New("team name must be 1 to 100 characters")
	ErrInvalidRole     = errors.New("unknown role")
	ErrNotMember       = errors.New("user is not a member of the team")
	ErrForbidden       = errors.New("your role does not allow this")
	ErrAlreadyMember   = errors.New("user is already a member")
	ErrLastOwner       = errors.New("a team must keep at least one owner")
	ErrTaskNotInTeam   = errors.New("the team passed in is not the task’s team")

	ErrTaskNotAssigned = errors.New("the task is not assigned")
	ErrTaskNotDueDate  = errors.New("the task is not due date")
)
