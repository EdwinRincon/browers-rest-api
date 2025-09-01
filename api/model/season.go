package model

import (
	"time"
)

type Season struct {
	ID          uint64         `gorm:"primaryKey" json:"id"`
	Year        uint16         `gorm:"not null;uniqueIndex;check:year >= 1999 AND year <= 2100" json:"year"`
	StartDate   time.Time      `gorm:"type:timestamp;not null" json:"start_date"`
	EndDate     time.Time      `gorm:"type:timestamp;not null" json:"end_date"`
	IsCurrent   bool           `gorm:"default:false" json:"is_current"`
	Matches     []Match        `gorm:"foreignKey:SeasonID" json:"matches" swaggerignore:"true"`
	Articles    []Article      `gorm:"foreignKey:SeasonID" json:"articles" swaggerignore:"true"`
	TeamStats   []TeamStat     `gorm:"foreignKey:SeasonID" json:"team_stats" swaggerignore:"true"`
	PlayerTeams []PlayerTeam   `gorm:"foreignKey:SeasonID" json:"player_teams" swaggerignore:"true"`
	PlayerStats []PlayerStat `gorm:"foreignKey:SeasonID" json:"player_stats" swaggerignore:"true"`
	CreatedAt   time.Time    `gorm:"type:timestamp;autoCreateTime" json:"created_at"`
	UpdatedAt   time.Time    `gorm:"type:timestamp;autoUpdateTime" json:"updated_at"`
}
