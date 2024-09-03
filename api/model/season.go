package model

import (
	"time"

	"gorm.io/gorm"
)

type Seasons struct {
	ID        uint8     `gorm:"primaryKey"`
	Year      time.Time `gorm:"type:date;not null" form:"year"`
	CreatedAt time.Time
	UpdatedAt time.Time
	DeletedAt gorm.DeletedAt `gorm:"index"`
}
