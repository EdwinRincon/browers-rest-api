package model

import (
	"time"

	"gorm.io/gorm"
)

type Lineups struct {
	ID        uint64 `gorm:"primaryKey"`
	Position  string `gorm:"type:varchar(5);not null" form:"position"`
	PlayersID uint64 `form:"players_id"`
	Players   Players
	MatchesID uint64 `form:"matches_id"`
	Matches   Matches
	CreatedAt time.Time
	UpdatedAt time.Time
	DeletedAt gorm.DeletedAt `gorm:"index"`
}
