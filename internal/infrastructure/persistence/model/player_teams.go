package model

import (
	"time"
)

// PlayerTeam represents the many-to-many relationship between players, teams, and seasons
type PlayerTeam struct {
	ID       uint64 `gorm:"primaryKey;autoIncrement" json:"id"`
	PlayerID uint64 `gorm:"uniqueIndex:idx_player_team_unique;not null;index:fk_players_player_teams" json:"player_id"`
	TeamID   uint64 `gorm:"uniqueIndex:idx_player_team_unique;not null" json:"team_id"`
	SeasonID uint64 `gorm:"uniqueIndex:idx_player_team_unique;not null" json:"season_id"`

	Player *Player `gorm:"foreignKey:PlayerID;constraint:OnUpdate:CASCADE,OnDelete:RESTRICT" json:"player,omitempty" swaggerignore:"true"`
	Team   *Team   `gorm:"foreignKey:TeamID;constraint:OnUpdate:CASCADE,OnDelete:RESTRICT" json:"team,omitempty" swaggerignore:"true"`
	Season *Season `gorm:"foreignKey:SeasonID;constraint:OnUpdate:CASCADE,OnDelete:RESTRICT" json:"season,omitempty" swaggerignore:"true"`

	StartDate time.Time  `gorm:"type:timestamp;uniqueIndex:idx_player_team_unique;not null" json:"start_date"`
	EndDate   *time.Time `gorm:"type:timestamp" json:"end_date,omitempty"`

	CreatedAt time.Time `gorm:"type:timestamp;autoCreateTime" json:"created_at,omitempty"`
	UpdatedAt time.Time `gorm:"type:timestamp;autoUpdateTime" json:"updated_at,omitempty"`
}
