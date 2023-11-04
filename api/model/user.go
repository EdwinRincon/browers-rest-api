package model

import (
	"time"

	"gorm.io/gorm"
)

type Users struct {
	gorm.Model
	Name       string    `gorm:"type:varchar(35);not null" json:"name" binding:"required"`
	Username   string    `gorm:"type:varchar(15);not null;unique" json:"username" binding:"required"`
	LastName   string    `gorm:"type:varchar(35);not null" json:"lastname" binding:"required"`
	Password   string    `gorm:"type:varchar(22);not null" json:"-" binding:"required"`
	Birthdate  time.Time `gorm:"type:date;not null" json:"birthdate" binding:"required"`
	Active     bool      `gorm:"type:boolean;default:true" json:"is_active" binding:"required"`
	ImgProfile string    `gorm:"type:varchar(255)" json:"img_profile"`
	ImgBanner  string    `gorm:"type:varchar(255)" json:"img_banner"`
	RolesID    uint64    `json:"roles_id" binding:"required"`
	Roles      Roles
}

type UsersResponse struct {
	ID         uint64    `json:"id"`
	Name       string    `json:"name"`
	Username   string    `json:"username"`
	LastName   string    `json:"lastname"`
	Birthdate  time.Time `json:"birthdate"`
	Active     bool      `json:"is_active"`
	ImgProfile string    `json:"img_profile"`
	ImgBanner  string    `json:"img_banner"`
	RolesID    uint64    `json:"roles_id"`
	Roles      Roles
}
