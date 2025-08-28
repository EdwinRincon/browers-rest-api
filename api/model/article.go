package model

import (
	"time"

	"gorm.io/gorm"
)

type Article struct {
	ID        uint64    `gorm:"primaryKey" json:"id" form:"id"`
	Title     string    `gorm:"type:varchar(100);not null;uniqueIndex" json:"title" form:"title" binding:"required,max=100"`
	Content   string    `gorm:"type:text;not null" json:"content" form:"content" binding:"required"`
	ImgBanner string    `gorm:"type:varchar(255);default:null" json:"img_banner,omitempty" form:"img_banner" binding:"omitempty,url"`
	Date      time.Time `gorm:"type:date;not null;index" json:"date" form:"date" binding:"required"`
	SeasonID  uint64    `gorm:"not null;index" json:"season_id" form:"season_id" binding:"required"`
	Season    *Season   `gorm:"foreignKey:SeasonID;constraint:OnUpdate:CASCADE,OnDelete:RESTRICT;" json:"season,omitempty"`

	CreatedAt time.Time      `gorm:"autoCreateTime" json:"created_at,omitempty"`
	UpdatedAt time.Time      `gorm:"autoUpdateTime" json:"updated_at,omitempty"`
	DeletedAt gorm.DeletedAt `gorm:"index" json:"-" swaggerignore:"true"`
}
