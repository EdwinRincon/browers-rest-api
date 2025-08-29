package dto

import (
	"time"
)

type CreateTeamRequest struct {
	FullName    string  `json:"full_name" binding:"required,max=35"`
	ShortName   string  `json:"short_name" binding:"required,max=5"`
	Color       string  `json:"color" binding:"required,max=10"`
	Color2      string  `json:"color2" binding:"required,max=10"`
	Shield      string  `json:"shield" binding:"required,url"`
	NextMatchID *uint64 `json:"next_match_id,omitempty"`
}

type UpdateTeamRequest struct {
	FullName    *string  `json:"full_name,omitempty" binding:"omitempty,max=35"`
	ShortName   *string  `json:"short_name,omitempty" binding:"omitempty,max=5"`
	Color       *string  `json:"color,omitempty" binding:"omitempty,max=10"`
	Color2      *string  `json:"color2,omitempty" binding:"omitempty,max=10"`
	Shield      *string  `json:"shield,omitempty" binding:"omitempty,url"`
	NextMatchID *uint64  `json:"next_match_id,omitempty"`
}

type TeamResponse struct {
	ID          uint64      `json:"id"`
	FullName    string      `json:"full_name"`
	ShortName   string      `json:"short_name"`
	Color       string      `json:"color"`
	Color2      string      `json:"color2"`
	Shield      string      `json:"shield"`
	NextMatchID *uint64     `json:"next_match_id,omitempty"`
	NextMatch   *MatchShort `json:"next_match,omitempty"`
	CreatedAt   time.Time   `json:"created_at"`
	UpdatedAt   time.Time   `json:"updated_at"`
}

type TeamShort struct {
	ID        uint64 `json:"id"`
	FullName  string `json:"full_name"`
	ShortName string `json:"short_name"`
}
