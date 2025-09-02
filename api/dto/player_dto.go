package dto

import (
	"time"
)

type CreatePlayerRequest struct {
	NickName      string   `json:"nick_name" binding:"required,max=20" example:"messi"`
	Height        uint16   `json:"height" binding:"required" example:"170"`
	Country       string   `json:"country" binding:"required,len=3" example:"ESP"`
	Country2      string   `json:"country2,omitempty" binding:"omitempty,len=3" example:"ARG"`
	Foot          string   `json:"foot" binding:"required,oneof=L R" example:"R"`
	Age           uint8    `json:"age" binding:"required,gte=16,lte=50" example:"30"`
	SquadNumber   uint8    `json:"squad_number" binding:"required,gte=1,lte=99" example:"10"`
	Position      string   `json:"position" binding:"required,oneof=por ceni cend lati med latd del deli deld" example:"del"`
	CareerSummary string   `json:"career_summary" example:"Veteran striker with good aerial ability."`
	UserID        *string  `json:"user_id,omitempty" example:"123e4567-e89b-12d3-a456-426614174000"`
	TeamIDs       []uint64 `json:"team_ids,omitempty" example:"1"`
}

type UpdatePlayerRequest struct {
	NickName      *string  `json:"nick_name,omitempty" binding:"omitempty,max=20"`
	Height        *uint16  `json:"height,omitempty" binding:"omitempty,gte=100,lte=250"`
	Country       *string  `json:"country,omitempty" binding:"omitempty,len=3"`
	Country2      *string  `json:"country2,omitempty" binding:"omitempty,len=3"`
	Foot          *string  `json:"foot,omitempty" binding:"omitempty,oneof=L R"`
	Age           *uint8   `json:"age,omitempty" binding:"omitempty,gte=16,lte=50"`
	SquadNumber   *uint8   `json:"squad_number,omitempty" binding:"omitempty,gte=1,lte=99"`
	Rating        *uint8   `json:"rating,omitempty" binding:"omitempty,max=100"`
	Position      *string  `json:"position,omitempty" binding:"omitempty,oneof=por ceni cend lati med latd del deli deld"`
	Injured       *bool    `json:"injured,omitempty"`
	CareerSummary *string  `json:"career_summary"`
	UserID        *string  `json:"user_id,omitempty"`
	TeamIDs       []uint64 `json:"team_ids,omitempty"`
}

type PlayerResponse struct {
	ID            uint64                  `json:"id"`
	NickName      string                  `json:"nick_name"`
	Height        uint16                  `json:"height"`
	Country       string                  `json:"country"`
	Country2      string                  `json:"country2,omitempty"`
	Foot          string                  `json:"foot"`
	Age           uint8                   `json:"age"`
	SquadNumber   uint8                   `json:"squad_number"`
	Rating        uint8                   `json:"rating"`
	Matches       uint16                  `json:"matches"`
	YCards        uint8                   `json:"y_cards"`
	RCards        uint8                   `json:"r_cards"`
	Goals         uint16                  `json:"goals"`
	Assists       uint16                  `json:"assists"`
	Saves         uint16                  `json:"saves"`
	Position      string                  `json:"position"`
	Injured       bool                    `json:"injured"`
	CareerSummary string                  `json:"career_summary,omitempty"`
	MVPCount      uint8                   `json:"mvp_count"`
	User          *UserShort             `json:"user,omitempty"`
	Teams         []TeamShort            `json:"teams,omitempty"`
	PlayerStats   []PlayerStatResponse   `json:"player_stats,omitempty"`
	CreatedAt     time.Time               `json:"created_at"`
	UpdatedAt     time.Time               `json:"updated_at"`
}

type PlayerShort struct {
	ID       uint64 `json:"id"`
	NickName string `json:"nick_name"`
	Position string `json:"position"`
}

type PlayerStats struct {
	ID       uint64 `json:"id"`
	NickName string `json:"nick_name"`
	Matches  uint16 `json:"matches"`
	Goals    uint16 `json:"goals"`
	Assists  uint16 `json:"assists"`
	YCards   uint8  `json:"y_cards"`
	RCards   uint8  `json:"r_cards"`
	Saves    uint16 `json:"saves"`
	Position string `json:"position"`
	MVPCount uint8  `json:"mvp_count"`
}
