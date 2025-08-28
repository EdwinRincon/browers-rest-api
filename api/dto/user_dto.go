package dto

import (
	"time"
)

type CreateUserRequest struct {
	Name     string `json:"name" binding:"required,min=2,max=35"`
	LastName string `json:"last_name" binding:"required,min=2,max=35"`
	Username string `json:"username" binding:"required,min=3,max=50,alphanum"`
	RoleID   uint64 `json:"role_id,omitempty" binding:"omitempty,gte=0,lte=255"`
}

type UpdateUserRequest struct {
	Name       *string    `json:"name,omitempty"`
	LastName   *string    `json:"last_name,omitempty"`
	Username   *string    `json:"username,omitempty"`
	Birthdate  *time.Time `json:"birthdate,omitempty"`
	ImgProfile *string    `json:"img_profile,omitempty"`
	ImgBanner  *string    `json:"img_banner,omitempty"`
	RoleID     *uint64    `json:"role_id,omitempty" binding:"omitempty,gte=0,lte=255"`
}

type UserResponse struct {
	ID         string    `json:"id"`
	Name       string    `json:"name"`
	LastName   string    `json:"last_name"`
	Username   string    `json:"username"`
	Birthdate  time.Time `json:"birthdate,omitempty"`
	ImgProfile string    `json:"img_profile,omitempty"`
	ImgBanner  string    `json:"img_banner,omitempty"`
	Role       RoleShort `json:"role,omitempty"`
	CreatedAt  time.Time `json:"created_at"`
	UpdatedAt  time.Time `json:"updated_at"`
}

type UserShort struct {
	ID       string `json:"id"`
	Username string `json:"username"`
}

type AuthUserResponse struct {
	ID         string `json:"id"`
	Name       string `json:"name"`
	LastName   string `json:"last_name"`
	Username   string `json:"username"`
	ImgProfile string `json:"img_profile,omitempty"`
	RoleName   string `json:"role_name"`
}
