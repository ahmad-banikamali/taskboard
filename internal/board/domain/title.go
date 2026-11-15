package domain

import (
	"strings"
	"unicode/utf8"
)

const (
	MinTitleLength = 1
	MaxTitleLength = 100
)

// Title is a value object. Create it only with NewTitle.
type Title struct {
	value string
}

func NewTitle(raw string) (Title, error) {
	s := strings.TrimSpace(raw)
	n := utf8.RuneCountInString(s)
	if n < MinTitleLength || n > MaxTitleLength {
		return Title{}, ErrInvalidTitle
	}
	return Title{value: s}, nil
}

func (t Title) String() string { return t.value }
