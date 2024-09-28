package model

import (
	"time"

	"gorm.io/gorm"
)

type TeamsStats struct {
	ID        uint64         `gorm:"primaryKey" json:"id" form:"id"`
	Wins      uint8          `gorm:"type:int(2);not null" json:"wins" form:"wins"`
	Draws     uint8          `gorm:"type:int(2);not null" json:"draws" form:"draws"`
	Losses    uint8          `gorm:"type:int(2);not null" json:"losses" form:"losses"`
	GoalsFor  uint16         `gorm:"type:int(3);not null" json:"goals_for" form:"goals_for"`
	GoalsAg   uint16         `gorm:"type:int(3);not null" json:"goals_against" form:"goals_against"`
	Points    uint16         `gorm:"type:int(3);not null" json:"points" form:"points"`
	SeasonsID uint8          `gorm:"not null" json:"seasons_id" form:"seasons_id"`
	Seasons   Seasons        `gorm:"foreignKey:SeasonsID" json:"seasons" form:"seasons"`
	CreatedAt time.Time      `gorm:"autoCreateTime" json:"created_at" form:"created_at"`
	UpdatedAt time.Time      `gorm:"autoUpdateTime" json:"updated_at" form:"updated_at"`
	DeletedAt gorm.DeletedAt `gorm:"index" json:"deleted_at" form:"deleted_at"`
}
