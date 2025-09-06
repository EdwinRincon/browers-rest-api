package domain

import "time"

// Season represents a football season in the domain layer.
// This is a pure business entity without infrastructure concerns.
type Season struct {
	ID        uint64
	Year      uint16
	StartDate time.Time
	EndDate   time.Time
	IsCurrent bool
	CreatedAt time.Time
	UpdatedAt time.Time
}

// IsValid performs basic domain validation for the season.
func (s *Season) IsValid() bool {
	return s.Year >= 1999 &&
		s.Year <= 2100 &&
		s.StartDate.Before(s.EndDate) &&
		s.EndDate.Sub(s.StartDate) >= 24*time.Hour // At least one day
}
