package model

import "gorm.io/gorm"

type Lineups struct {
	gorm.Model
	ID        uint64 `gorm:"primaryKey"`
	Position  string `gorm:"type:varchar(5);not null" form:"position"`
	PlayersID uint64 `form:"players_id"`
	Players   Players
	MatchesID uint64 `form:"matches_id"`
	Matches   Matches
}
