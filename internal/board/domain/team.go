package domain

import (
	"strings"
	"time"
	"unicode/utf8"
)

type TeamID string
type UserID string

const MaxTeamNameLength = 100

// Team is the aggregate root. The members map is private on purpose.
type Team struct {
	id        TeamID
	name      string
	members   map[UserID]Role
	createdAt time.Time
}

// NewTeam creates a team. The creator becomes its first owner.
// now is passed in (not time.Now()) so tests control the time.
func NewTeam(id TeamID, name string, creator UserID, now time.Time) (*Team, error) {
	name = strings.TrimSpace(name)
	if n := utf8.RuneCountInString(name); n < 1 || n > MaxTeamNameLength {
		return nil, ErrInvalidTeamName
	}
	return &Team{
		id:        id,
		name:      name,
		members:   map[UserID]Role{creator: RoleOwner},
		createdAt: now.UTC(),
	}, nil
}

func (t *Team) ID() TeamID           { return t.id }
func (t *Team) Name() string         { return t.name }
func (t *Team) CreatedAt() time.Time { return t.createdAt }

func (t *Team) RoleOf(u UserID) (Role, bool) {
	r, ok := t.members[u]
	return r, ok
}

// Members returns a COPY, so callers cannot change the real map.
func (t *Team) Members() map[UserID]Role {
	out := make(map[UserID]Role, len(t.members))
	for u, r := range t.members {
		out[u] = r
	}
	return out
}

// requireOwner: the actor must be a member, and an owner.
func (t *Team) requireOwner(actor UserID) error {
	role, ok := t.members[actor]
	if !ok {
		return ErrNotMember
	}
	if !role.CanManageMembers() {
		return ErrForbidden
	}
	return nil
}

func (t *Team) countOwners() int {
	n := 0
	for _, r := range t.members {
		if r == RoleOwner {
			n++
		}
	}
	return n
}

// Every method below checks first and changes last,
// so a failed call leaves the team exactly as it was.

func (t *Team) AddMember(actor, user UserID, role Role) error {
	if err := t.requireOwner(actor); err != nil {
		return err
	}
	if !role.IsValid() {
		return ErrInvalidRole
	}
	if _, exists := t.members[user]; exists {
		return ErrAlreadyMember
	}
	t.members[user] = role
	return nil
}

func (t *Team) RemoveMember(actor, user UserID) error {
	if err := t.requireOwner(actor); err != nil {
		return err
	}
	role, ok := t.members[user]
	if !ok {
		return ErrNotMember
	}
	if role == RoleOwner && t.countOwners() == 1 {
		return ErrLastOwner
	}
	delete(t.members, user)
	return nil
}

func (t *Team) RequireMember(user UserID) error {
	_, ok := t.members[user]
	if !ok {
		return ErrNotMember
	}
	return nil
}

func (t *Team) RequireCanEditTasks(actor UserID) error {
	role, ok := t.members[actor]
	if !ok {
		return ErrNotMember
	}
	if !role.CanEditTasks() {
		return ErrForbidden
	}
	return nil
}

// TransferOwnership makes `to` an owner and turns the actor into a member.
// Transfer to yourself changes nothing.
func (t *Team) TransferOwnership(actor, to UserID) error {
	if err := t.requireOwner(actor); err != nil {
		return err
	}
	if _, ok := t.members[to]; !ok {
		return ErrNotMember
	}
	t.members[to] = RoleOwner
	if actor != to {
		t.members[actor] = RoleMember
	}
	return nil
}
