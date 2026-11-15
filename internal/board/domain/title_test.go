package domain

import (
	"errors"
	"strings"
	"testing"
)

func TestNewTitle(t *testing.T) {
	cases := []struct {
		name    string
		input   string
		want    string // expected title when valid
		wantErr bool
	}{
		{"empty", "", "", true},
		{"only spaces", "   ", "", true},
		{"one character", "a", "a", false},
		{"spaces are trimmed", "  Fix bug  ", "Fix bug", false},
		{"exactly 100 characters", strings.Repeat("a", 100), strings.Repeat("a", 100), false},
		{"101 characters", strings.Repeat("a", 101), "", true},
		{"100 two-byte characters (counts characters, not bytes)", strings.Repeat("é", 100), strings.Repeat("é", 100), false},
		{"101 two-byte characters", strings.Repeat("é", 101), "", true},
	}
	for _, tc := range cases {
		t.Run(tc.name, func(t *testing.T) {
			got, err := NewTitle(tc.input)
			if tc.wantErr {
				if !errors.Is(err, ErrInvalidTitle) {
					t.Fatalf("err = %v, want ErrInvalidTitle", err)
				}
				return
			}
			if err != nil {
				t.Fatalf("unexpected error: %v", err)
			}
			if got.String() != tc.want {
				t.Errorf("title = %q, want %q", got.String(), tc.want)
			}
		})
	}
}
