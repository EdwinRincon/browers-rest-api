package model

import (
	"time"

	"gorm.io/gorm"
)

type Roles struct {
	ID        uint8          `gorm:"primaryKey" json:"id" form:"id"`
	Name      string         `gorm:"type:varchar(20);not null" json:"name" form:"name"`
	CreatedAt time.Time      `gorm:"autoCreateTime" json:"created_at" form:"created_at"`
	UpdatedAt time.Time      `gorm:"autoUpdateTime" json:"updated_at" form:"updated_at"`
	DeletedAt gorm.DeletedAt `gorm:"index" json:"deleted_at" form:"deleted_at" swaggerignore:"true"`
}
