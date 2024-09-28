package model

import (
	"time"

	"gorm.io/gorm"
)

type Lineups struct {
	ID        uint64         `gorm:"primaryKey" json:"id" form:"id"`
	Position  string         `gorm:"type:varchar(5);not null" json:"position" form:"position"`
	PlayersID *uint64        `gorm:"index;constraint:OnUpdate:CASCADE,OnDelete:SET NULL;" json:"players_id" form:"players_id"`
	Players   Players        `gorm:"foreignKey:PlayersID" json:"players" form:"players"`
	MatchesID *uint64        `gorm:"index;constraint:OnUpdate:CASCADE,OnDelete:SET NULL;" json:"matches_id" form:"matches_id"`
	Matches   Matches        `gorm:"foreignKey:MatchesID" json:"matches" form:"matches"`
	CreatedAt time.Time      `gorm:"autoCreateTime" json:"created_at" form:"created_at"`
	UpdatedAt time.Time      `gorm:"autoUpdateTime" json:"updated_at" form:"updated_at"`
	DeletedAt gorm.DeletedAt `gorm:"index" json:"deleted_at" form:"deleted_at"`
}
