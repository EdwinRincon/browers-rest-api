package model

import "gorm.io/gorm"

type Classifications struct {
	gorm.Model
	ID           uint64 `gorm:"primaryKey"`
	TeamsStatsID uint64 `form:"teams_stats_id"`
	TeamsStats   TeamsStats
}
