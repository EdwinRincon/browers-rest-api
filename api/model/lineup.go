package model

import (
	"time"

	"gorm.io/gorm"
)

type Lineups struct {
	ID        uint64         `gorm:"primaryKey" json:"id" form:"id"`
	Position  string         `gorm:"type:varchar(5);not null" json:"position" form:"position" binding:"required,oneof=por ceni cend lati med latd del deli deld"`
	PlayerID  *uint64        `gorm:"index;constraint:OnUpdate:CASCADE,OnDelete:SET NULL;uniqueIndex:idx_player_match" json:"player_id" form:"player_id"`
	Player    Players        `gorm:"foreignKey:PlayerID" json:"player" form:"player"`
	MatchID   *uint64        `gorm:"index;constraint:OnUpdate:CASCADE,OnDelete:SET NULL" json:"match_id" form:"match_id"`
	Match     Matches        `gorm:"foreignKey:MatchID" json:"match" form:"match"`
	CreatedAt time.Time      `gorm:"autoCreateTime" json:"created_at" form:"created_at"`
	UpdatedAt time.Time      `gorm:"autoUpdateTime" json:"updated_at" form:"updated_at"`
	DeletedAt gorm.DeletedAt `gorm:"index" json:"deleted_at" form:"deleted_at"`
}
