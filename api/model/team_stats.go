package model

import (
	"time"

	"gorm.io/gorm"
)

type TeamsStats struct {
	ID        uint64 `gorm:"primaryKey"`
	Wins      uint64 `gorm:"type:int(2);not null"`
	Draws     uint64 `gorm:"type:int(2);not null"`
	Losses    uint64 `gorm:"type:int(2);not null"`
	GoalsFor  uint64 `gorm:"type:int(3);not null"`
	GoalsAg   uint64 `gorm:"type:int(3);not null"`
	Points    uint64 `gorm:"type:int(3);not null"`
	SeasonsID uint8  `form:"seasons_id"`
	Seasons   Seasons
	CreatedAt time.Time
	UpdatedAt time.Time
	DeletedAt gorm.DeletedAt `gorm:"index"`
}
