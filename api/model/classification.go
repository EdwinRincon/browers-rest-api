package model

import (
	"time"

	"gorm.io/gorm"
)

type Classifications struct {
	ID           uint64         `gorm:"primaryKey" json:"id" form:"id"`
	TeamsStatsID *uint64        `gorm:"index;constraint:OnUpdate:CASCADE,OnDelete:SET NULL;" json:"teams_stats_id" form:"teams_stats_id"` // Make this nullable
	TeamsStats   TeamsStats     `gorm:"foreignKey:TeamsStatsID" json:"teams_stats" form:"teams_stats"`
	SeasonsID    *uint8         `gorm:"index;constraint:OnUpdate:CASCADE,OnDelete:SET NULL;" json:"seasons_id" form:"seasons_id"` // Make this nullable
	Seasons      Seasons        `gorm:"foreignKey:SeasonsID" json:"seasons" form:"seasons"`
	CreatedAt    time.Time      `gorm:"autoCreateTime" json:"created_at" form:"created_at"`
	UpdatedAt    time.Time      `gorm:"autoUpdateTime" json:"updated_at" form:"updated_at"`
	DeletedAt    gorm.DeletedAt `gorm:"index" json:"deleted_at" form:"deleted_at"`
}
