package model

import (
	"time"

	"gorm.io/gorm"
)

type Classifications struct {
	ID           uint64 `gorm:"primaryKey"`
	TeamsStatsID uint64 `form:"teams_stats_id"`
	TeamsStats   TeamsStats
	CreatedAt    time.Time
	UpdatedAt    time.Time
	DeletedAt    gorm.DeletedAt `gorm:"index"`
}
