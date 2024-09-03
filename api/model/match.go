package model

import (
	"time"

	"gorm.io/gorm"
)

type Matches struct {
	ID          uint64    `gorm:"primaryKey"`
	Date        time.Time `gorm:"not null;unique" form:"date"`
	Time        string    `gorm:"type:varchar(5);not null" form:"time"`
	Home        string    `gorm:"type:varchar(35);not null" form:"home"`
	Away        string    `gorm:"type:varchar(35);not null" form:"away"`
	Location    string    `gorm:"type:varchar(35);not null" form:"location"`
	HomeGoals   uint8     `gorm:"not null" form:"home_goals"`
	AwayGoals   uint8     `gorm:"not null" form:"away_goals"`
	SeasonsID   uint8     `form:"seasons_id"`
	Seasons     Seasons
	MVPPlayerID uint64  `form:"players_id"`
	Players     Players `gorm:"foreignKey:MVPPlayerID"`
	CreatedAt   time.Time
	UpdatedAt   time.Time
	DeletedAt   gorm.DeletedAt `gorm:"index"`
}
