package domain

type Role string

const (
	RoleOwner  Role = "OWNER"
	RoleMember Role = "MEMBER"
	RoleViewer Role = "VIEWER"
)

func ParseRole(raw string) (Role, error) {
	switch r := Role(raw); r {
	case RoleOwner, RoleMember, RoleViewer:
		return r, nil
	}
	return "", ErrInvalidRole
}

func (r Role) IsValid() bool {
	_, err := ParseRole(string(r))
	return err == nil
}

// CanManageMembers Only owners manage members and ownership.
func (r Role) CanManageMembers() bool { return r == RoleOwner }

// CanEditTasks Viewers can read but not create or edit tasks. (Used in a later step.)
func (r Role) CanEditTasks() bool { return r == RoleOwner || r == RoleMember }
