package model

import (
	"time"

	"gorm.io/gorm"
)

type Roles struct {
	ID        uint8  `gorm:"primarykey"`
	Name      string `gorm:"type:char(20);not null" json:"name"`
	CreatedAt time.Time
	UpdatedAt time.Time
	DeletedAt gorm.DeletedAt `gorm:"index"`
}
