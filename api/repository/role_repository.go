package repository

import (
	"context"

	"github.com/EdwinRincon/browersfc-api/api/model"
	"gorm.io/gorm"
)

type RoleRepository interface {
	GetRoleByID(ctx context.Context, id uint8) (*model.Roles, error)
}

type RoleRepositoryImpl struct {
	db *gorm.DB
}

func NewRoleRepository(db *gorm.DB) RoleRepository {
	return &RoleRepositoryImpl{db: db}
}

func (rr *RoleRepositoryImpl) GetRoleByID(ctx context.Context, id uint8) (*model.Roles, error) {
	var role model.Roles
	if err := rr.db.Where("id = ?", id).First(&role).Error; err != nil {
		return nil, err
	}
	return &role, nil
}
