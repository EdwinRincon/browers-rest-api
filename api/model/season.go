package model

import (
	"time"

	"gorm.io/gorm"
)

type Seasons struct {
	gorm.Model
	ID   uint8     `gorm:"primaryKey"`
	Year time.Time `gorm:"type:date;not null" form:"year"`
}
