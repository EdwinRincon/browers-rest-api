package domain

import "time"

// TeamNextMatch represents a simplified match for the team's next match
type TeamNextMatch struct {
	ID        uint64
	Status    string
	Kickoff   time.Time
	Location  string
	HomeGoals uint8
	AwayGoals uint8
}

// Team represents a football team in the domain layer.
// This is a pure business entity without infrastructure concerns.
type Team struct {
	ID             uint64
	FullName       string
	ShortName      string
	PrimaryColor   string
	SecondaryColor string
	Shield         string
	NextMatchID    *uint64
	NextMatch      *TeamNextMatch
	CreatedAt      time.Time
	UpdatedAt      time.Time
}

// IsValid performs basic domain validation for the team.
func (t *Team) IsValid() bool {
	return t.FullName != "" &&
		t.ShortName != "" &&
		len(t.ShortName) <= 5 &&
		len(t.FullName) <= 35 &&
		t.PrimaryColor != "" &&
		t.SecondaryColor != "" &&
		t.PrimaryColor != t.SecondaryColor &&
		t.Shield != ""
}
