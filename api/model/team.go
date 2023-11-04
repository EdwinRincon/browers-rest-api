package model

import (
	"time"

	"gorm.io/gorm"
)

type Teams struct {
	gorm.Model
	ID           uint64    `gorm:"primaryKey"`
	FullName     string    `gorm:"type:varchar(35);not null"`
	ShortName    string    `gorm:"type:varchar(5);unique"`
	Color        string    `gorm:"type:varchar(10);not null"`
	Color2       string    `gorm:"type:varchar(10);not null"`
	Shield       string    `gorm:"type:varchar(200);not null"`
	NextMatch    time.Time `gorm:"type:date" form:"next_match"`
	TeamsStatsID uint64    `form:"teams_stats_id"`
	TeamsStats   TeamsStats
}
