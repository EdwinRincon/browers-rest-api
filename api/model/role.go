package model

import "gorm.io/gorm"

type Roles struct {
	gorm.Model
	ID   uint8  `gorm:"primaryKey"`
	Name string `gorm:"type:char(10);not null"`
}
