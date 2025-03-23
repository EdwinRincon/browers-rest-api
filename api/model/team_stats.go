package model

import (
	"time"

	"gorm.io/gorm"
)

type TeamStat struct {
	ID           uint64         `gorm:"primaryKey" json:"id"`
	Wins         uint8          `gorm:"not null;default:0" json:"wins" form:"wins" binding:"gte=0"`
	Draws        uint8          `gorm:"not null;default:0" json:"draws" form:"draws" binding:"gte=0"`
	Losses       uint8          `gorm:"not null;default:0" json:"losses" form:"losses" binding:"gte=0"`
	GoalsFor     uint16         `gorm:"not null;default:0" json:"goals_for" form:"goals_for" binding:"gte=0"`
	GoalsAgainst uint16         `gorm:"not null;default:0" json:"goals_against" form:"goals_against" binding:"gte=0"`
	Points       uint16         `gorm:"not null;default:0" json:"points" form:"points" binding:"gte=0"`
	Rank         uint8          `gorm:"not null;default:0" json:"rank" form:"rank" binding:"gte=0"`
	TeamID       uint64         `gorm:"index;not null" json:"team_id" form:"team_id" binding:"required"`
	Team         *Team          `gorm:"foreignKey:TeamID" json:"team,omitempty"`
	SeasonID     uint           `gorm:"index;not null" json:"season_id" form:"season_id" binding:"required"`
	Season       *Season        `gorm:"foreignKey:SeasonID" json:"season,omitempty"`
	CreatedAt    time.Time      `gorm:"autoCreateTime" json:"created_at,omitempty"`
	UpdatedAt    time.Time      `gorm:"autoUpdateTime" json:"updated_at,omitempty"`
	DeletedAt    gorm.DeletedAt `gorm:"index" json:"-"`
}
