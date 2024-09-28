package model

import (
	"time"

	"gorm.io/gorm"
)

type Teams struct {
	ID           uint64         `gorm:"primaryKey" json:"id" form:"id"`
	FullName     string         `gorm:"type:varchar(35);not null" json:"full_name" form:"full_name"`
	ShortName    string         `gorm:"type:varchar(5);unique;not null" json:"short_name" form:"short_name"`
	Color        string         `gorm:"type:varchar(10);not null" json:"color" form:"color"`
	Color2       string         `gorm:"type:varchar(10);not null" json:"color2" form:"color2"`
	Shield       string         `gorm:"type:varchar(200);not null" json:"shield" form:"shield"`
	NextMatch    time.Time      `gorm:"type:date" json:"next_match" form:"next_match"`
	Players      []Players      `gorm:"many2many:player_teams;" json:"players" form:"players"`
	TeamsStatsID uint64         `gorm:"not null" json:"teams_stats_id" form:"teams_stats_id"`
	TeamsStats   TeamsStats     `gorm:"foreignKey:TeamsStatsID" json:"teams_stats" form:"teams_stats"`
	CreatedAt    time.Time      `gorm:"autoCreateTime" json:"created_at" form:"created_at"`
	UpdatedAt    time.Time      `gorm:"autoUpdateTime" json:"updated_at" form:"updated_at"`
	DeletedAt    gorm.DeletedAt `gorm:"index" json:"deleted_at" form:"deleted_at"`
}
