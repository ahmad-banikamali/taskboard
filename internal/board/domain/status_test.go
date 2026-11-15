package domain

import (
	"errors"
	"testing"
)

// Every pair of statuses (3 x 3 = 9) is listed on purpose.
// We never rely on a default.
func TestStatus_CanMoveTo(t *testing.T) {
	cases := []struct {
		from, to Status
		want     bool
	}{
		{StatusTodo, StatusTodo, false},
		{StatusTodo, StatusInProgress, true},
		{StatusTodo, StatusDone, false},

		{StatusInProgress, StatusTodo, true},
		{StatusInProgress, StatusInProgress, false},
		{StatusInProgress, StatusDone, true},

		{StatusDone, StatusTodo, false},
		{StatusDone, StatusInProgress, false},
		{StatusDone, StatusDone, false},
	}
	for _, tc := range cases {
		t.Run(string(tc.from)+"->"+string(tc.to), func(t *testing.T) {
			if got := tc.from.CanMoveTo(tc.to); got != tc.want {
				t.Errorf("CanMoveTo = %v, want %v", got, tc.want)
			}
		})
	}
}

func TestStatus_MoveTo_NotAllowed_ReturnsDomainError(t *testing.T) {
	_, err := StatusDone.MoveTo(StatusTodo)
	if !errors.Is(err, ErrInvalidTransition) {
		t.Fatalf("err = %v, want ErrInvalidTransition", err)
	}
}

func TestParseStatus(t *testing.T) {
	if got, err := ParseStatus("IN_PROGRESS"); err != nil || got != StatusInProgress {
		t.Errorf("ParseStatus(IN_PROGRESS) = %v, %v", got, err)
	}
	for _, bad := range []string{"", "todo", "BLOCKED"} {
		if _, err := ParseStatus(bad); !errors.Is(err, ErrInvalidStatus) {
			t.Errorf("ParseStatus(%q) err = %v, want ErrInvalidStatus", bad, err)
		}
	}
}
