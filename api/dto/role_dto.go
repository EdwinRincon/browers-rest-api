package dto

import "time"

type CreateRoleRequest struct {
	Name        string  `json:"name" binding:"required,max=20"`
	Description *string `json:"description,omitempty"`
}

type UpdateRoleRequest struct {
	Name        *string `json:"name,omitempty" binding:"omitempty,min=1,max=20"`
	Description *string `json:"description,omitempty"`
}

type RoleResponse struct {
	ID          uint64    `json:"id"`
	Name        string    `json:"name"`
	Description string    `json:"description"`
	CreatedAt   time.Time `json:"created_at"`
	UpdatedAt   time.Time `json:"updated_at"`
}

type RoleShort struct {
	ID   uint64 `json:"id"`
	Name string `json:"name"`
}
