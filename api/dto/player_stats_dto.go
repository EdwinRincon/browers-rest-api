package dto

import "time"

type CreatePlayerStatRequest struct {
	PlayerID      uint64  `json:"player_id" binding:"required" example:"1"`
	MatchID       uint64  `json:"match_id" binding:"required" example:"1"`
	SeasonID      uint64  `json:"season_id" binding:"required" example:"1"`
	TeamID        *uint64 `json:"team_id,omitempty" example:"1"`
	Goals         uint8   `json:"goals" binding:"gte=0" example:"2"`
	Assists       uint8   `json:"assists" binding:"gte=0" example:"1"`
	Saves         uint8   `json:"saves" binding:"gte=0" example:"0"`
	YellowCards   uint8   `json:"yellow_cards" binding:"gte=0" example:"1"`
	RedCards      uint8   `json:"red_cards" binding:"gte=0" example:"0"`
	Rating        uint8   `json:"rating" binding:"gte=0,lte=100" example:"85"`
	Starting      bool    `json:"starting" example:"true"`
	MinutesPlayed uint8   `json:"minutes_played" binding:"gte=0,lte=120" example:"90"`
	IsMVP         bool    `json:"is_mvp" example:"false"`
	Position      string  `json:"position" binding:"omitempty,oneof=por ceni cend lati med latd del deli deld" example:"del"`
}

type UpdatePlayerStatRequest struct {
	TeamID        *uint64 `json:"team_id,omitempty" example:"2"`
	Goals         *uint8  `json:"goals,omitempty" binding:"omitempty,gte=0" example:"3"`
	Assists       *uint8  `json:"assists,omitempty" binding:"omitempty,gte=0" example:"2"`
	Saves         *uint8  `json:"saves,omitempty" binding:"omitempty,gte=0" example:"0"`
	YellowCards   *uint8  `json:"yellow_cards,omitempty" binding:"omitempty,gte=0" example:"1"`
	RedCards      *uint8  `json:"red_cards,omitempty" binding:"omitempty,gte=0" example:"0"`
	Rating        *uint8  `json:"rating,omitempty" binding:"omitempty,gte=0,lte=100" example:"88"`
	Starting      *bool   `json:"starting,omitempty" example:"true"`
	MinutesPlayed *uint8  `json:"minutes_played,omitempty" binding:"omitempty,gte=0,lte=120" example:"85"`
	IsMVP         *bool   `json:"is_mvp,omitempty" example:"true"`
	Position      *string `json:"position,omitempty" binding:"omitempty,oneof=por ceni cend lati med latd del deli deld" example:"del"`
}

type PlayerStatResponse struct {
	ID            uint64      `json:"id"`
	PlayerID      uint64      `json:"player_id"`
	MatchID       uint64      `json:"match_id"`
	SeasonID      uint64      `json:"season_id"`
	TeamID        *uint64     `json:"team_id,omitempty"`
	Goals         uint8       `json:"goals"`
	Assists       uint8       `json:"assists"`
	Saves         uint8       `json:"saves"`
	YellowCards   uint8       `json:"yellow_cards"`
	RedCards      uint8       `json:"red_cards"`
	Rating        uint8       `json:"rating"`
	Starting      bool        `json:"starting"`
	MinutesPlayed uint8       `json:"minutes_played"`
	IsMVP         bool        `json:"is_mvp"`
	Position      string      `json:"position,omitempty"`
	Player        PlayerShort `json:"player,omitempty"`
	Match         MatchShort  `json:"match,omitempty"`
	Season        SeasonShort `json:"season,omitempty"`
	Team          *TeamShort  `json:"team,omitempty"`
	CreatedAt     time.Time   `json:"created_at"`
	UpdatedAt     time.Time   `json:"updated_at"`
}
