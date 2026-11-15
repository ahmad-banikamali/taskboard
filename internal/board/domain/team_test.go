package domain

import (
	"errors"
	"testing"
	"time"
)

// A time with a +03:00 offset, to prove NewTeam converts to UTC.
var testNow = time.Date(2026, 11, 11, 12, 0, 0, 0, time.FixedZone("plus3", 3*3600))

// teamWith builds a team: alice is OWNER, plus the other members given.
func teamWith(t *testing.T, others map[UserID]Role) *Team {
	t.Helper()
	team, err := NewTeam("t1", "Backend", "alice", testNow)
	if err != nil {
		t.Fatal(err)
	}
	for u, r := range others {
		if err := team.AddMember("alice", u, r); err != nil {
			t.Fatal(err)
		}
	}
	return team
}

func wantErr(t *testing.T, got, want error) {
	t.Helper()
	if !errors.Is(got, want) {
		t.Fatalf("err = %v, want %v", got, want)
	}
}

func wantRole(t *testing.T, team *Team, user UserID, want Role) {
	t.Helper()
	got, ok := team.RoleOf(user)
	if !ok || got != want {
		t.Fatalf("role of %s = %q (member=%v), want %q", user, got, ok, want)
	}
}

func wantNotMember(t *testing.T, team *Team, user UserID) {
	t.Helper()
	if _, ok := team.RoleOf(user); ok {
		t.Fatalf("%s should not be a member", user)
	}
}

func TestNewTeam_CreatorBecomesOwner_TimeIsUTC(t *testing.T) {
	team := teamWith(t, nil)
	wantRole(t, team, "alice", RoleOwner)
	if team.CreatedAt().Location() != time.UTC {
		t.Errorf("CreatedAt location = %v, want UTC", team.CreatedAt().Location())
	}
	if !team.CreatedAt().Equal(testNow) {
		t.Errorf("CreatedAt changed the moment in time")
	}
}

func TestNewTeam_BadName_ReturnsError(t *testing.T) {
	for _, name := range []string{"", "   ", string(make([]rune, 101))} {
		_, err := NewTeam("t1", name, "alice", testNow)
		wantErr(t, err, ErrInvalidTeamName)
	}
}

func TestAddMember_OwnerAddsMember_Works(t *testing.T) {
	team := teamWith(t, nil)
	if err := team.AddMember("alice", "bob", RoleMember); err != nil {
		t.Fatal(err)
	}
	wantRole(t, team, "bob", RoleMember)
}

func TestAddMember_ActorIsNotOwner_Forbidden(t *testing.T) {
	for _, role := range []Role{RoleMember, RoleViewer} {
		team := teamWith(t, map[UserID]Role{"bob": role})
		err := team.AddMember("bob", "carol", RoleMember)
		wantErr(t, err, ErrForbidden)
		wantNotMember(t, team, "carol")
	}
}

func TestAddMember_ActorIsNotInTeam_NotMember(t *testing.T) {
	team := teamWith(t, nil)
	wantErr(t, team.AddMember("stranger", "carol", RoleMember), ErrNotMember)
	wantNotMember(t, team, "carol")
}

func TestAddMember_AlreadyMember_ReturnsError_RoleUnchanged(t *testing.T) {
	team := teamWith(t, map[UserID]Role{"bob": RoleViewer})
	wantErr(t, team.AddMember("alice", "bob", RoleOwner), ErrAlreadyMember)
	wantRole(t, team, "bob", RoleViewer)
}

func TestAddMember_UnknownRole_ReturnsError(t *testing.T) {
	team := teamWith(t, nil)
	wantErr(t, team.AddMember("alice", "bob", Role("BOSS")), ErrInvalidRole)
	wantNotMember(t, team, "bob")
}

func TestRemoveMember_OwnerRemovesMember_Works(t *testing.T) {
	team := teamWith(t, map[UserID]Role{"bob": RoleMember})
	if err := team.RemoveMember("alice", "bob"); err != nil {
		t.Fatal(err)
	}
	wantNotMember(t, team, "bob")
}

func TestRemoveMember_ActorIsNotOwner_Forbidden(t *testing.T) {
	team := teamWith(t, map[UserID]Role{"bob": RoleMember, "carol": RoleViewer})
	wantErr(t, team.RemoveMember("bob", "carol"), ErrForbidden)
	wantRole(t, team, "carol", RoleViewer)
}

func TestRemoveMember_TargetIsNotMember_ReturnsError(t *testing.T) {
	team := teamWith(t, nil)
	wantErr(t, team.RemoveMember("alice", "ghost"), ErrNotMember)
}

func TestRemoveMember_OnlyOwnerRemovesSelf_LastOwner(t *testing.T) {
	team := teamWith(t, map[UserID]Role{"bob": RoleMember})
	wantErr(t, team.RemoveMember("alice", "alice"), ErrLastOwner)
	wantRole(t, team, "alice", RoleOwner)
}

func TestRemoveMember_TwoOwners_OneCanLeave(t *testing.T) {
	team := teamWith(t, map[UserID]Role{"bob": RoleOwner})
	if err := team.RemoveMember("alice", "alice"); err != nil {
		t.Fatal(err)
	}
	wantNotMember(t, team, "alice")
	wantRole(t, team, "bob", RoleOwner)
}

func TestTransferOwnership_ToMember_SwapsRoles(t *testing.T) {
	team := teamWith(t, map[UserID]Role{"bob": RoleMember})
	if err := team.TransferOwnership("alice", "bob"); err != nil {
		t.Fatal(err)
	}
	wantRole(t, team, "bob", RoleOwner)
	wantRole(t, team, "alice", RoleMember)
}

func TestTransferOwnership_ActorIsNotOwner_Forbidden(t *testing.T) {
	team := teamWith(t, map[UserID]Role{"bob": RoleMember, "carol": RoleMember})
	wantErr(t, team.TransferOwnership("bob", "carol"), ErrForbidden)
	wantRole(t, team, "alice", RoleOwner)
	wantRole(t, team, "carol", RoleMember)
}

func TestTransferOwnership_TargetIsNotMember_ReturnsError_AliceStaysOwner(t *testing.T) {
	team := teamWith(t, nil)
	wantErr(t, team.TransferOwnership("alice", "ghost"), ErrNotMember)
	wantRole(t, team, "alice", RoleOwner)
}

func TestTransferOwnership_ToSelf_ChangesNothing(t *testing.T) {
	team := teamWith(t, nil)
	if err := team.TransferOwnership("alice", "alice"); err != nil {
		t.Fatal(err)
	}
	wantRole(t, team, "alice", RoleOwner)
}

func TestMembers_ReturnsCopy_ChangingItDoesNotChangeTeam(t *testing.T) {
	team := teamWith(t, nil)
	m := team.Members()
	delete(m, "alice")
	wantRole(t, team, "alice", RoleOwner)
}
