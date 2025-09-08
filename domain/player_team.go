package domain

import (
	"time"
)

// PlayerTeam represents the many-to-many relationship between players, teams, and seasons.
type PlayerTeam struct {
	ID       uint64
	PlayerID uint64
	TeamID   uint64
	SeasonID uint64

	// Related entities
	Player *Player
	Team   *Team
	Season *Season

	StartDate time.Time
	EndDate   *time.Time

	CreatedAt time.Time
	UpdatedAt time.Time
}

// IsValid validates the PlayerTeam domain entity according to business rules.
func (pt *PlayerTeam) IsValid() bool {
	return pt.PlayerID > 0 &&
		pt.TeamID > 0 &&
		pt.SeasonID > 0 &&
		!pt.StartDate.IsZero()
}

// IsActive checks if the player team relationship is currently active.
func (pt *PlayerTeam) IsActive(at time.Time) bool {
	if at.Before(pt.StartDate) {
		return false
	}
	if pt.EndDate != nil && at.After(*pt.EndDate) {
		return false
	}
	return true
}

// Duration calculates the duration of the player team relationship.
func (pt *PlayerTeam) Duration() time.Duration {
	endTime := time.Now()
	if pt.EndDate != nil {
		endTime = *pt.EndDate
	}
	return endTime.Sub(pt.StartDate)
}

// HasOverlapWith checks if this PlayerTeam has date overlap with another PlayerTeam.
func (pt *PlayerTeam) HasOverlapWith(other *PlayerTeam) bool {
	if pt.PlayerID != other.PlayerID || pt.TeamID != other.TeamID || pt.SeasonID != other.SeasonID {
		return false
	}

	ptEnd := pt.StartDate.AddDate(100, 0, 0)
	if pt.EndDate != nil {
		ptEnd = *pt.EndDate
	}

	otherEnd := other.StartDate.AddDate(100, 0, 0)
	if other.EndDate != nil {
		otherEnd = *other.EndDate
	}

	// Check for overlap: start1 <= end2 && start2 <= end1
	return pt.StartDate.Before(otherEnd.Add(time.Second)) && other.StartDate.Before(ptEnd.Add(time.Second))
}
