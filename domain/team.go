package domain

import "time"

// Team represents a football team in the domain layer.
// This is a pure business entity without infrastructure concerns.
type Team struct {
	ID          uint64
	FullName    string
	ShortName   string
	Color       string
	Color2      string
	Shield      string
	NextMatchID *uint64
	CreatedAt   time.Time
	UpdatedAt   time.Time
}

// IsValid performs basic domain validation for the team.
func (t *Team) IsValid() bool {
	return t.FullName != "" &&
		t.ShortName != "" &&
		len(t.ShortName) <= 5 &&
		len(t.FullName) <= 35 &&
		t.Color != "" &&
		t.Color2 != "" &&
		t.Color != t.Color2 &&
		t.Shield != ""
}
