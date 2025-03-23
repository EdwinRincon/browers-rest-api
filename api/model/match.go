package model

import (
	"time"

	"gorm.io/gorm"
)

type Match struct {
	ID          uint64    `gorm:"primaryKey" json:"id" form:"id"`
	Status      string    `gorm:"type:varchar(11);not null;check:status IN ('scheduled','in_progress','completed','postponed','cancelled')" json:"status" form:"status" binding:"required,oneof=scheduled in_progress completed postponed cancelled"`
	Date        time.Time `gorm:"type:date;not null;index" json:"date" form:"date" binding:"required"`
	Time        string    `gorm:"type:varchar(5);not null" json:"time" form:"time" binding:"required,len=5"`
	Location    string    `gorm:"type:varchar(35);not null" json:"location" form:"location" binding:"required,max=35"`
	HomeGoals   uint8     `gorm:"not null;default:0" json:"home_goals" form:"home_goals"`
	AwayGoals   uint8     `gorm:"not null;default:0" json:"away_goals" form:"away_goals"`
	HomeTeamID  uint64    `gorm:"index;not null" json:"home_team_id" form:"home_team_id" binding:"required"`
	HomeTeam    *Team     `gorm:"foreignKey:HomeTeamID" json:"home_team,omitempty" form:"home_team"`
	AwayTeamID  uint64    `gorm:"index;not null" json:"away_team_id" form:"away_team_id" binding:"required"`
	AwayTeam    *Team     `gorm:"foreignKey:AwayTeamID" json:"away_team,omitempty" form:"away_team"`
	Lineups     []Lineup  `gorm:"foreignKey:MatchID" json:"lineups,omitempty" form:"lineups"`
	SeasonID    uint      `gorm:"index;not null" json:"season_id" form:"season_id" binding:"required"`
	Season      *Season   `gorm:"foreignKey:SeasonID;constraint:OnUpdate:CASCADE,OnDelete:RESTRICT;" json:"season,omitempty" form:"season"`
	MVPPlayerID *uint64   `gorm:"index" json:"mvp_player_id" form:"mvp_player_id"`
	MVPPlayer   *Player   `gorm:"foreignKey:MVPPlayerID;constraint:OnUpdate:CASCADE,OnDelete:SET NULL;" json:"mvp_player,omitempty" form:"mvp_player"`

	CreatedAt time.Time      `gorm:"autoCreateTime" json:"created_at" form:"created_at"`
	UpdatedAt time.Time      `gorm:"autoUpdateTime" json:"updated_at" form:"updated_at"`
	DeletedAt gorm.DeletedAt `gorm:"index" json:"deleted_at,omitempty" form:"-"`
}
