package model

import (
	"time"

	"gorm.io/gorm"
)

type Team struct {
	ID        uint64    `gorm:"primaryKey" json:"id"`
	FullName  string    `gorm:"type:varchar(35);not null;uniqueIndex" json:"full_name" form:"full_name" binding:"required,max=35"`
	ShortName string    `gorm:"type:varchar(5);unique;not null" json:"short_name" form:"short_name" binding:"required,max=5"`
	Color     string    `gorm:"type:varchar(10);not null" json:"color" form:"color" binding:"required,max=10"`
	Color2    string    `gorm:"type:varchar(10);not null" json:"color2" form:"color2" binding:"required,max=10"`
	Shield    string    `gorm:"type:varchar(200);not null" json:"shield" form:"shield" binding:"required,url"`
	NextMatch time.Time `gorm:"type:date" json:"next_match,omitempty" form:"next_match"`

	Players     []Player   `gorm:"many2many:player_teams;" json:"players,omitempty"`
	HomeMatches []Match    `gorm:"foreignKey:HomeTeamID" json:"home_matches,omitempty"`
	AwayMatches []Match    `gorm:"foreignKey:AwayTeamID" json:"away_matches,omitempty"`
	TeamStats   []TeamStat `gorm:"foreignKey:TeamID" json:"team_stats,omitempty"`

	CreatedAt time.Time      `gorm:"autoCreateTime" json:"created_at,omitempty"`
	UpdatedAt time.Time      `gorm:"autoUpdateTime" json:"updated_at,omitempty"`
	DeletedAt gorm.DeletedAt `gorm:"index" json:"-"`
}
