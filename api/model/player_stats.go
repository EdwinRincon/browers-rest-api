package model

import (
	"time"

	"gorm.io/gorm"
)

// PlayerStat represents the statistics of a player in a specific match and season.
type PlayerStat struct {
	ID       uint64  `gorm:"primaryKey" json:"id"`
	PlayerID uint64  `gorm:"index;not null;uniqueIndex:idx_player_match" json:"player_id"`
	MatchID  uint64  `gorm:"index;not null;uniqueIndex:idx_player_match" json:"match_id"`
	SeasonID uint64  `gorm:"index;not null;index:idx_player_season" json:"season_id"`
	TeamID   *uint64 `gorm:"index" json:"team_id,omitempty"`

	Goals         uint8  `gorm:"type:tinyint;not null;default:0" json:"goals"`
	Assists       uint8  `gorm:"type:tinyint;not null;default:0" json:"assists"`
	Saves         uint8  `gorm:"type:tinyint;not null;default:0" json:"saves"`
	YC            uint8  `gorm:"type:tinyint;not null;default:0" json:"yellow_cards"`
	RC            uint8  `gorm:"type:tinyint;not null;default:0" json:"red_cards"`
	Rating        uint8  `gorm:"type:tinyint;not null;default:0;check:rating <= 100" json:"rating"`
	Starting      bool   `gorm:"default:false" json:"starting"`
	MinutesPlayed uint8  `gorm:"type:tinyint;not null;default:0" json:"minutes_played"`
	IsMVP         bool   `gorm:"default:false" json:"is_mvp"`
	Position      string `gorm:"type:varchar(5)" json:"position"`

	Player *Player `gorm:"foreignKey:PlayerID;constraint:OnUpdate:CASCADE,OnDelete:RESTRICT" json:"player,omitempty"`
	Match  *Match  `gorm:"foreignKey:MatchID;constraint:OnUpdate:CASCADE,OnDelete:RESTRICT" json:"match,omitempty"`
	Season *Season `gorm:"foreignKey:SeasonID;constraint:OnUpdate:CASCADE,OnDelete:RESTRICT" json:"season,omitempty"`
	Team   *Team   `gorm:"foreignKey:TeamID;constraint:OnUpdate:CASCADE,OnDelete:SET NULL" json:"team,omitempty"`

	CreatedAt time.Time      `gorm:"type:timestamp;autoCreateTime" json:"created_at,omitempty"`
	UpdatedAt time.Time      `gorm:"type:timestamp;autoUpdateTime" json:"updated_at,omitempty"`
	DeletedAt gorm.DeletedAt `gorm:"type:timestamp;index" json:",omitempty" swaggerignore:"true"`
}
