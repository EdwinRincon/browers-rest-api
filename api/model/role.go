package model

import (
	"time"

	"gorm.io/gorm"
)

type Role struct {
	ID          uint64         `gorm:"primaryKey" json:"id" form:"id"`
	Name        string         `gorm:"type:varchar(20);not null;uniqueIndex" json:"name" form:"name" binding:"required,max=20"`
	Description string         `gorm:"type:varchar(100)" json:"description,omitempty" form:"description"`
	CreatedAt   time.Time      `gorm:"autoCreateTime" json:"created_at,omitempty"`
	UpdatedAt   time.Time      `gorm:"autoUpdateTime" json:"updated_at,omitempty"`
	DeletedAt   gorm.DeletedAt `gorm:"index" json:"-" swaggerignore:"true"`
}
