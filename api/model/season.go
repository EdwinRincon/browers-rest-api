package model

import (
	"time"

	"gorm.io/gorm"
)

type Seasons struct {
	ID         uint8          `gorm:"primaryKey" json:"id" form:"id"`
	Year       uint16         `gorm:"type:int(4);not null" json:"year" form:"year"`
	Matches    []Matches      `gorm:"foreignKey:SeasonsID" json:"matches" form:"matches"`         // Has Many (Matches)
	Articles   []Articles     `gorm:"foreignKey:SeasonsID" json:"articles" form:"articles"`       // Has Many (Articles)
	TeamsStats []TeamsStats   `gorm:"foreignKey:SeasonsID" json:"teams_stats" form:"teams_stats"` // Has Many (TeamsStats)
	CreatedAt  time.Time      `gorm:"autoCreateTime" json:"created_at" form:"created_at"`
	UpdatedAt  time.Time      `gorm:"autoUpdateTime" json:"updated_at" form:"updated_at"`
	DeletedAt  gorm.DeletedAt `gorm:"index" json:"deleted_at" form:"deleted_at"`
}
