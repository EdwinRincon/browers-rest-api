package domain

import (
	"time"
)

// Match represents the core Match entity in the domain layer.
// This entity contains only business-relevant fields
type Match struct {
	ID          uint64
	Status      string
	Kickoff     time.Time
	Location    string
	HomeGoals   uint8
	AwayGoals   uint8
	HomeTeamID  uint64
	AwayTeamID  uint64
	SeasonID    uint64
	MVPPlayerID *uint64
	CreatedAt   time.Time
	UpdatedAt   time.Time

	// Related entities
	HomeTeam  *Team
	AwayTeam  *Team
	Season    *Season
	MVPPlayer *Player
}
