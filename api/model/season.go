package model

import (
	"time"

	"gorm.io/gorm"
)

type Season struct {
	ID          uint           `gorm:"primaryKey" json:"id"`
	Year        uint16         `gorm:"not null;uniqueIndex;check:year >= 1999 AND year <= 2100" json:"year"`
	StartDate   time.Time      `gorm:"not null" json:"start_date"`
	EndDate     time.Time      `gorm:"not null" json:"end_date"`
	IsCurrent   bool           `gorm:"default:false" json:"is_current"`
	Matches     []Match        `gorm:"foreignKey:SeasonID" json:"matches"`
	Articles    []Article      `gorm:"foreignKey:SeasonID" json:"articles"`
	TeamStats   []TeamStat     `gorm:"foreignKey:SeasonID" json:"team_stats"`
	PlayerTeams []PlayerTeam   `gorm:"foreignKey:SeasonID" json:"player_teams"`
	CreatedAt   time.Time      `gorm:"autoCreateTime" json:"created_at"`
	UpdatedAt   time.Time      `gorm:"autoUpdateTime" json:"updated_at"`
	DeletedAt   gorm.DeletedAt `gorm:"index" json:"-"`
}
