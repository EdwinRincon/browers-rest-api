package dto

import (
	"time"
)

type CreateLineupRequest struct {
	Position string `json:"position" binding:"required,oneof=por ceni cend lati med latd del deli deld"`
	PlayerID uint64 `json:"player_id" binding:"required"`
	MatchID  uint64 `json:"match_id" binding:"required"`
	Starting bool   `json:"starting"`
}

type UpdateLineupRequest struct {
	Position *string `json:"position,omitempty" binding:"omitempty,oneof=por ceni cend lati med latd del deli deld"`
	PlayerID *uint64 `json:"player_id,omitempty" binding:"omitempty"`
	MatchID  *uint64 `json:"match_id,omitempty" binding:"omitempty"`
	Starting *bool   `json:"starting,omitempty"`
}

type LineupResponse struct {
	ID        uint64      `json:"id"`
	Position  string      `json:"position"`
	PlayerID  uint64      `json:"player_id"`
	Player    PlayerShort `json:"player,omitempty"`
	MatchID   uint64      `json:"match_id"`
	Match     MatchShort  `json:"match,omitempty"`
	Starting  bool        `json:"starting"`
	CreatedAt time.Time   `json:"created_at"`
	UpdatedAt time.Time   `json:"updated_at"`
}

// LineupShortResponse represents a simplified version of a lineup for embedding in other responses
type LineupShortResponse struct {
	ID       uint64 `json:"id"`
	Position string `json:"position"`
	PlayerID uint64 `json:"player_id"`
	Starting bool   `json:"starting"`
}

// MatchLineupResponse represents a list of lineups organized by match
type MatchLineupResponse struct {
	MatchID        uint64                `json:"match_id"`
	HomeTeamID     uint64                `json:"home_team_id"`
	HomeTeam       string                `json:"home_team"`
	AwayTeamID     uint64                `json:"away_team_id"`
	AwayTeam       string                `json:"away_team"`
	Date           time.Time             `json:"date"`
	StartingLineup []LineupShortResponse `json:"starting_lineup"`
	Substitutes    []LineupShortResponse `json:"substitutes"`
}
