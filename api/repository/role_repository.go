package repository

import (
	"context"
	"errors"

	"github.com/EdwinRincon/browersfc-api/api/model"
	"gorm.io/gorm"
)

var ErrRoleNotFound = errors.New("role not found")

type RoleRepository interface {
	GetRoleByID(ctx context.Context, id uint8) (*model.Role, error)
	CreateRole(ctx context.Context, role *model.Role) error
	UpdateRole(ctx context.Context, role *model.Role) error
	DeleteRole(ctx context.Context, id uint8) error
	GetAllRoles(ctx context.Context) ([]model.Role, error)
}

type RoleRepositoryImpl struct {
	db *gorm.DB
}

func NewRoleRepository(db *gorm.DB) RoleRepository {
	return &RoleRepositoryImpl{db: db}
}

func (rr *RoleRepositoryImpl) GetRoleByID(ctx context.Context, id uint8) (*model.Role, error) {
	var role model.Role
	if err := rr.db.Where("id = ?", id).First(&role).Error; err != nil {
		return nil, err
	}
	return &role, nil
}

func (rr *RoleRepositoryImpl) CreateRole(ctx context.Context, role *model.Role) error {
	return rr.db.Create(role).Error
}

func (rr *RoleRepositoryImpl) UpdateRole(ctx context.Context, role *model.Role) error {
	return rr.db.Save(role).Error
}

func (rr *RoleRepositoryImpl) DeleteRole(ctx context.Context, id uint8) error {
	return rr.db.Delete(&model.Role{}, id).Error
}

func (rr *RoleRepositoryImpl) GetAllRoles(ctx context.Context) ([]model.Role, error) {
	var roles []model.Role
	if err := rr.db.Find(&roles).Error; err != nil {
		return nil, err
	}
	return roles, nil
}
