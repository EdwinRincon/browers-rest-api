package model

import (
	"time"
)

type Match struct {
	ID          uint64       `gorm:"primaryKey" json:"id" form:"id"`
	Status      string       `gorm:"type:varchar(11);not null;check:status IN ('scheduled','in_progress','completed','postponed','cancelled')" json:"status" form:"status" binding:"required,oneof=scheduled in_progress completed postponed cancelled"`
	Kickoff     time.Time    `gorm:"type:timestamp;not null" json:"kickoff" form:"kickoff" binding:"required"`
	Location    string       `gorm:"type:varchar(35);not null" json:"location" form:"location" binding:"required,max=35"`
	HomeGoals   uint8        `gorm:"not null;default:0" json:"home_goals" form:"home_goals"`
	AwayGoals   uint8        `gorm:"not null;default:0" json:"away_goals" form:"away_goals"`
	HomeTeamID  uint64       `gorm:"index;not null" json:"home_team_id" form:"home_team_id" binding:"required"`
	HomeTeam    *Team        `gorm:"foreignKey:HomeTeamID;constraint:OnUpdate:CASCADE,OnDelete:RESTRICT;" json:"home_team,omitempty" form:"home_team" swaggerignore:"true"`
	AwayTeamID  uint64       `gorm:"index;not null" json:"away_team_id" form:"away_team_id" binding:"required"`
	AwayTeam    *Team        `gorm:"foreignKey:AwayTeamID;constraint:OnUpdate:CASCADE,OnDelete:RESTRICT;" json:"away_team,omitempty" form:"away_team" swaggerignore:"true"`
	Lineups     []Lineup     `gorm:"foreignKey:MatchID;constraint:OnDelete:CASCADE;" json:"lineups,omitempty" form:"lineups" swaggerignore:"true"`
	PlayerStats []PlayerStat `gorm:"foreignKey:MatchID;constraint:OnDelete:CASCADE;" json:"player_stats,omitempty" swaggerignore:"true"`
	SeasonID    uint64       `gorm:"index;not null" json:"season_id" form:"season_id" binding:"required"`
	Season      *Season      `gorm:"foreignKey:SeasonID;constraint:OnUpdate:CASCADE,OnDelete:RESTRICT;" json:"season,omitempty" form:"season" swaggerignore:"true"`
	MVPPlayerID *uint64      `gorm:"index" json:"mvp_player_id" form:"mvp_player_id"`
	MVPPlayer   *Player      `gorm:"foreignKey:MVPPlayerID;constraint:OnUpdate:CASCADE,OnDelete:SET NULL;" json:"mvp_player,omitempty" form:"mvp_player" swaggerignore:"true"`

	CreatedAt time.Time `gorm:"type:timestamp;autoCreateTime" json:"created_at" form:"created_at"`
	UpdatedAt time.Time `gorm:"type:timestamp;autoUpdateTime" json:"updated_at" form:"updated_at"`
}
