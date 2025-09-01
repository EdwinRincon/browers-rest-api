package model

import (
	"time"
)

type Player struct {
	ID            uint64         `gorm:"primaryKey" json:"id" form:"id"`
	NickName      string         `gorm:"type:varchar(20)" json:"nick_name" form:"nick_name" binding:"required,max=20"`
	Height        uint16         `gorm:"type:smallint;not null;check:height >= 100 AND height <= 250" json:"height" form:"height" binding:"required,gte=100,lte=250"`
	Country       string         `gorm:"type:varchar(3);not null" json:"country" form:"country" binding:"required,len=3"`
	Country2      string         `gorm:"type:varchar(3)" json:"country2,omitempty" form:"country2"`
	Foot          string         `gorm:"type:varchar(1);not null" json:"foot" form:"foot" binding:"required,oneof=L R"`
	Age           uint8          `gorm:"type:tinyint;not null;check:age >= 16 AND age <= 50" json:"age" form:"age" binding:"required,gte=16,lte=50"`
	SquadNumber   uint8          `gorm:"type:tinyint;not null;check:squad_number >= 1 AND squad_number <= 99" json:"squad_number" form:"squad_number" binding:"required,gte=1,lte=99"`
	Rating        uint8          `gorm:"type:tinyint;not null;default:0;check:rating <= 100" json:"rating" form:"rating"`
	Matches       uint16         `gorm:"type:smallint;not null;default:0;" json:"matches" form:"matches"`
	YCards        uint8          `gorm:"type:tinyint;not null;default:0;" json:"y_cards" form:"y_cards"`
	RCards        uint8          `gorm:"type:tinyint;not null;default:0;" json:"r_cards" form:"r_cards"`
	Goals         uint16         `gorm:"type:smallint;not null;default:0;" json:"goals" form:"goals"`
	Assists       uint16         `gorm:"type:smallint;not null;default:0;" json:"assists" form:"assists"`
	Saves         uint16         `gorm:"type:smallint;not null;default:0;" json:"saves" form:"saves"`
	Position      string         `gorm:"type:varchar(5);not null;" json:"position" form:"position" binding:"required,oneof=por ceni cend lati med latd del deli deld"`
	Injured       bool           `gorm:"default:false;" json:"injured" form:"injured"`
	CareerSummary string         `gorm:"type:varchar(1000);not null;" json:"career_summary,omitempty" form:"career_summary"`
	PlayerTeams   []PlayerTeam   `json:"player_teams,omitempty" swaggerignore:"true"`
	Lineups       []Lineup       `gorm:"foreignKey:PlayerID;constraint:OnDelete:RESTRICT;" json:"lineups,omitempty" form:"lineups" swaggerignore:"true"`
	PlayerStats   []PlayerStat   `gorm:"foreignKey:PlayerID;constraint:OnDelete:RESTRICT;" json:"player_stats,omitempty" swaggerignore:"true"`
	MVPCount      uint8          `gorm:"type:tinyint;not null;default:0;" json:"mvp_count" form:"mvp_count"`
	UserID        *string        `gorm:"index;" json:"user_id,omitempty" form:"user_id"`
	User          *User       `gorm:"foreignKey:UserID;constraint:OnUpdate:CASCADE,OnDelete:SET NULL" json:"user,omitempty" form:"user" swaggerignore:"true"`
	CreatedAt     time.Time   `gorm:"type:timestamp;autoCreateTime;" json:"created_at,omitempty"`
	UpdatedAt     time.Time   `gorm:"type:timestamp;autoUpdateTime;" json:"updated_at,omitempty"`
}
