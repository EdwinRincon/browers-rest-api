package dto

import (
	"time"
)

// CreatePlayerTeamRequest represents the DTO for creating a new player-team relationship
type CreatePlayerTeamRequest struct {
	PlayerID  uint64     `json:"player_id" binding:"required,min=1"`
	TeamID    uint64     `json:"team_id" binding:"required,min=1"`
	SeasonID  uint64     `json:"season_id" binding:"required,min=1"`
	StartDate time.Time  `json:"start_date" binding:"required"`
	EndDate   *time.Time `json:"end_date,omitempty"`
}

// UpdatePlayerTeamRequest represents the DTO for updating a player-team relationship
type UpdatePlayerTeamRequest struct {
	StartDate *time.Time `json:"start_date,omitempty"`
	EndDate   *time.Time `json:"end_date,omitempty"`
}

type PlayerTeamResponse struct {
	ID        uint64      `json:"id"`
	PlayerID  uint64      `json:"player_id"`
	TeamID    uint64      `json:"team_id"`
	SeasonID  uint64      `json:"season_id"`
	Player    PlayerShort `json:"player,omitempty"`
	Team      TeamShort   `json:"team,omitempty"`
	Season    SeasonShort `json:"season,omitempty"`
	StartDate time.Time   `json:"start_date"`
	EndDate   *time.Time  `json:"end_date,omitempty"`
	CreatedAt time.Time   `json:"created_at"`
	UpdatedAt time.Time   `json:"updated_at"`
}

type PlayerTeamShort struct {
	ID        uint64     `json:"id"`
	PlayerID  uint64     `json:"player_id"`
	TeamID    uint64     `json:"team_id"`
	SeasonID  uint64     `json:"season_id"`
	StartDate time.Time  `json:"start_date"`
	EndDate   *time.Time `json:"end_date,omitempty"`
}
