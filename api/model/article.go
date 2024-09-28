package model

import (
	"time"

	"gorm.io/gorm"
)

type Articles struct {
	ID        uint64         `gorm:"primaryKey" json:"id" form:"id"`
	Title     string         `gorm:"type:varchar(100);not null" json:"title" form:"title"`
	Content   string         `gorm:"type:text;not null" json:"content" form:"content"`
	ImgBanner string         `gorm:"type:varchar(255);not null" json:"img_banner" form:"img_banner"`
	Date      time.Time      `gorm:"type:date;not null" json:"date" form:"date"`
	SeasonsID uint8          `gorm:"not null" json:"seasons_id" form:"seasons_id"`
	Seasons   Seasons        `gorm:"foreignKey:SeasonsID;constraint:OnUpdate:CASCADE,OnDelete:SET NULL;" json:"seasons" form:"seasons"`
	CreatedAt time.Time      `gorm:"autoCreateTime" json:"created_at" form:"created_at"`
	UpdatedAt time.Time      `gorm:"autoUpdateTime" json:"updated_at" form:"updated_at"`
	DeletedAt gorm.DeletedAt `gorm:"index" json:"deleted_at" form:"deleted_at"`
}
