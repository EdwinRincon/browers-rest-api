package model

import (
	"time"

	"gorm.io/gorm"
)

type Articles struct {
	gorm.Model
	ID        uint64    `gorm:"primaryKey"`
	Title     string    `gorm:"type:varchar(50);not null" form:"title"`
	Content   string    `gorm:"type:varchar(500);not null" form:"content"`
	ImgBanner string    `gorm:"type:varchar(200);not null" form:"img_banner"`
	Date      time.Time `gorm:"type:date" form:"date"`
	SeasonsID uint8     `form:"seasons_id"`
	Seasons   Seasons
}
