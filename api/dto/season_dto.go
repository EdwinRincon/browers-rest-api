package dto

import (
	"time"
)

type CreateSeasonRequest struct {
	Year      uint16    `json:"year" binding:"required,gte=1999,lte=2100" example:"2025"`
	StartDate time.Time `json:"start_date" binding:"required" example:"2025-08-01T00:00:00Z"`
	EndDate   time.Time `json:"end_date" binding:"required" example:"2026-06-30T00:00:00Z"`
	IsCurrent bool      `json:"is_current" example:"true"`
}

type UpdateSeasonRequest struct {
	Year      *uint16    `json:"year,omitempty" binding:"omitempty,gte=1999,lte=2100"`
	StartDate *time.Time `json:"start_date,omitempty"`
	EndDate   *time.Time `json:"end_date,omitempty"`
	IsCurrent *bool      `json:"is_current,omitempty"`
}

type SeasonResponse struct {
	ID        uint64    `json:"id"`
	Year      uint16    `json:"year"`
	StartDate time.Time `json:"start_date"`
	EndDate   time.Time `json:"end_date"`
	IsCurrent bool      `json:"is_current"`
	CreatedAt time.Time `json:"created_at"`
	UpdatedAt time.Time `json:"updated_at"`
}

type SeasonStatsResponse struct {
	ID            uint64    `json:"id"`
	Year          uint16    `json:"year"`
	StartDate     time.Time `json:"start_date"`
	EndDate       time.Time `json:"end_date"`
	IsCurrent     bool      `json:"is_current"`
	MatchCount    int       `json:"match_count"`
	TeamCount     int       `json:"team_count"`
	PlayerCount   int       `json:"player_count"`
	ArticleCount  int       `json:"article_count"`
	TotalGoals    int       `json:"total_goals"`
	TotalRedCards int       `json:"total_red_cards"`
	CreatedAt     time.Time `json:"created_at"`
	UpdatedAt     time.Time `json:"updated_at"`
}

type SeasonShort struct {
	ID   uint64 `json:"id"`
	Year uint16 `json:"year"`
}
