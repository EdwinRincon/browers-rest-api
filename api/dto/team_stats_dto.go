package dto

import (
	"time"
)

type CreateTeamStatsRequest struct {
	Wins         uint16 `json:"wins" binding:"gte=0" example:"10"`
	Draws        uint16 `json:"draws" binding:"gte=0" example:"5"`
	Losses       uint16 `json:"losses" binding:"gte=0" example:"3"`
	GoalsFor     uint16 `json:"goals_for" binding:"gte=0" example:"25"`
	GoalsAgainst uint16 `json:"goals_against" binding:"gte=0" example:"15"`
	Points       int16  `json:"points" binding:"gte=0" example:"35"`
	Rank         uint16 `json:"rank" binding:"gte=0" example:"3"`
	SeasonID     uint64 `json:"season_id" binding:"required" example:"1"`
	TeamID       uint64 `json:"team_id" binding:"required" example:"1"`
}

type UpdateTeamStatsRequest struct {
	Wins         *uint16 `json:"wins,omitempty" binding:"omitempty,gte=0"`
	Draws        *uint16 `json:"draws,omitempty" binding:"omitempty,gte=0"`
	Losses       *uint16 `json:"losses,omitempty" binding:"omitempty,gte=0"`
	GoalsFor     *uint16 `json:"goals_for,omitempty" binding:"omitempty,gte=0"`
	GoalsAgainst *uint16 `json:"goals_against,omitempty" binding:"omitempty,gte=0"`
	Points       *int16  `json:"points,omitempty" binding:"omitempty"`
	Rank         *uint16 `json:"rank,omitempty" binding:"omitempty,gte=0"`
	SeasonID     *uint64 `json:"season_id,omitempty" binding:"omitempty,min=1"`
	TeamID       *uint64 `json:"team_id,omitempty" binding:"omitempty,min=1"`
}

type TeamStatsResponse struct {
	ID           uint64       `json:"id"`
	Wins         uint16       `json:"wins"`
	Draws        uint16       `json:"draws"`
	Losses       uint16       `json:"losses"`
	GoalsFor     uint16       `json:"goals_for"`
	GoalsAgainst uint16       `json:"goals_against"`
	Points       int16        `json:"points"`
	Rank         uint16       `json:"rank"`
	SeasonID     uint64       `json:"season_id"`
	TeamID       uint64       `json:"team_id"`
	Team         *TeamShort   `json:"team,omitempty"`
	Season       *SeasonShort `json:"season,omitempty"`
	CreatedAt    time.Time    `json:"created_at"`
	UpdatedAt    time.Time    `json:"updated_at"`
}

type TeamStatsShort struct {
	ID     uint64 `json:"id"`
	Wins   uint16 `json:"wins"`
	Draws  uint16 `json:"draws"`
	Losses uint16 `json:"losses"`
	Points int16  `json:"points"`
	Rank   uint16 `json:"rank"`
}
