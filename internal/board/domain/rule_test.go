package domain

import "testing"

func TestRole_Permissions(t *testing.T) {
	cases := []struct {
		role                Role
		manageMembers, edit bool
	}{
		{RoleOwner, true, true},
		{RoleMember, false, true},
		{RoleViewer, false, false},
	}
	for _, tc := range cases {
		if got := tc.role.CanManageMembers(); got != tc.manageMembers {
			t.Errorf("%s CanManageMembers = %v, want %v", tc.role, got, tc.manageMembers)
		}
		if got := tc.role.CanEditTasks(); got != tc.edit {
			t.Errorf("%s CanEditTasks = %v, want %v", tc.role, got, tc.edit)
		}
	}
}
