package model

import (
	"time"

	"gorm.io/gorm"
)

type Matches struct {
	ID          uint64         `gorm:"primaryKey" json:"id" form:"id"`
	Status      string         `gorm:"type:varchar(11);not null" json:"status" form:"status"`
	Date        time.Time      `gorm:"type:date;not null" json:"date" form:"date"`
	Time        string         `gorm:"type:varchar(5);not null" json:"time" form:"time"`
	Home        string         `gorm:"type:varchar(35);not null" json:"home" form:"home"`
	Away        string         `gorm:"type:varchar(35);not null" json:"away" form:"away"`
	Location    string         `gorm:"type:varchar(35);not null" json:"location" form:"location"`
	HomeGoals   uint8          `gorm:"not null" json:"home_goals" form:"home_goals"`
	AwayGoals   uint8          `gorm:"not null" json:"away_goals" form:"away_goals"`
	HomeTeamID  uint64         `gorm:"index" json:"home_team_id" form:"home_team_id"`
	HomeTeam    Teams          `gorm:"foreignKey:HomeTeamID" json:"home_team" form:"home_team"`
	AwayTeamID  uint64         `gorm:"index" json:"away_team_id" form:"away_team_id"`
	AwayTeam    Teams          `gorm:"foreignKey:AwayTeamID" json:"away_team" form:"away_team"`
	Lineups     []Lineups      `gorm:"foreignKey:MatchID" json:"lineups" form:"lineups"`
	SeasonsID   uint8          `gorm:"not null" json:"seasons_id" form:"seasons_id"`
	Seasons     Seasons        `gorm:"foreignKey:SeasonsID;constraint:OnUpdate:CASCADE,OnDelete:SET NULL;" json:"seasons" form:"seasons"`
	MVPPlayerID *uint64        `gorm:"index;constraint:OnDelete:SET NULL,OnUpdate:CASCADE;" json:"mvp_player_id" form:"mvp_player_id"`
	MVPPlayer   Players        `gorm:"foreignKey:MVPPlayerID;constraint:OnUpdate:CASCADE,OnDelete:SET NULL;" json:"mvp_player" form:"mvp_player"`
	CreatedAt   time.Time      `gorm:"autoCreateTime" json:"created_at" form:"created_at"`
	UpdatedAt   time.Time      `gorm:"autoUpdateTime" json:"updated_at" form:"updated_at"`
	DeletedAt   gorm.DeletedAt `gorm:"index" json:"deleted_at" form:"deleted_at"`
}
