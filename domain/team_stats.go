package domain

import (
	"time"
)

// TeamStats represents a team's statistics for a specific season.
// This is the domain entity
type TeamStats struct {
	ID           uint64
	Wins         uint16
	Draws        uint16
	Losses       uint16
	GoalsFor     uint16
	GoalsAgainst uint16
	Points       int16
	Rank         uint16
	SeasonID     uint64
	TeamID       uint64
	Team         *Team
	Season       *Season
	CreatedAt    time.Time
	UpdatedAt    time.Time
}
