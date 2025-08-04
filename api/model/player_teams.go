package model

import (
	"time"

	"gorm.io/gorm"
)

// PlayerTeam represents the many-to-many relationship between players, teams, and seasons
type PlayerTeam struct {
	PlayerID uint64 `gorm:"primaryKey;autoIncrement:false" json:"player_id"`
	TeamID   uint64 `gorm:"primaryKey;autoIncrement:false" json:"team_id"`
	SeasonID uint64 `gorm:"primaryKey;autoIncrement:false" json:"season_id"`

	Player *Player `gorm:"foreignKey:PlayerID" json:"player,omitempty"`
	Team   *Team   `gorm:"foreignKey:TeamID" json:"team,omitempty"`
	Season *Season `gorm:"foreignKey:SeasonID" json:"season,omitempty"`

	StartDate time.Time `gorm:"not null" json:"start_date"`
	EndDate   time.Time `json:"end_date,omitempty"`

	CreatedAt time.Time      `gorm:"autoCreateTime" json:"created_at,omitempty"`
	UpdatedAt time.Time      `gorm:"autoUpdateTime" json:"updated_at,omitempty"`
	DeletedAt gorm.DeletedAt `gorm:"index" json:"-" swaggerignore:"true"`
}
