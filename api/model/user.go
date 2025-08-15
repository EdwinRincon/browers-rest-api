package model

import (
	"time"

	"github.com/google/uuid"
	"gorm.io/gorm"
)

type User struct {
	ID         string         `gorm:"type:char(36);primaryKey" json:"id"`
	Name       string         `gorm:"type:varchar(35);not null" json:"name" binding:"required,min=2,max=35"`
	LastName   string         `gorm:"type:varchar(35);not null" json:"last_name" binding:"required,min=2,max=35"`
	Username   string         `gorm:"type:varchar(50);not null;uniqueIndex" json:"username" binding:"required,safe_email,allowed_domain"`
	Birthdate  *time.Time     `gorm:"type:date" json:"birthdate" example:"1990-01-01"`
	ImgProfile string         `gorm:"type:varchar(255)" json:"img_profile,omitempty" binding:"omitempty,url"`
	ImgBanner  string         `gorm:"type:varchar(255)" json:"img_banner,omitempty" binding:"omitempty,url"`
	RoleID     uint8          `json:"role_id" binding:"required,min=1"`
	Role       *Role          `gorm:"foreignKey:RoleID" json:"role,omitempty" binding:"-"`
	CreatedAt  time.Time      `json:"created_at,omitempty"`
	UpdatedAt  time.Time      `json:"updated_at,omitempty"`
	DeletedAt  gorm.DeletedAt `gorm:"index" json:"-" swaggerignore:"true"`
}

// BeforeCreate will set a UUID rather than numeric ID.
func (u *User) BeforeCreate(tx *gorm.DB) error {

	if u.ID == "" {
		u.ID = uuid.New().String()
	}
	return nil
}
