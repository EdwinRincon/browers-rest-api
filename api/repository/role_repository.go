package repository

import (
	"context"
	"errors"
	"fmt"

	"github.com/EdwinRincon/browersfc-api/api/constants"
	"github.com/EdwinRincon/browersfc-api/api/model"
	"gorm.io/gorm"
)

type RoleRepository interface {
	GetRoleByID(ctx context.Context, id uint8) (*model.Role, error)
	GetActiveRoleByName(ctx context.Context, name string) (*model.Role, error)
	GetUnscopedRoleByName(ctx context.Context, name string) (*model.Role, error)
	CreateRole(ctx context.Context, role *model.Role) error
	UpdateRole(ctx context.Context, role *model.Role) error
	DeleteRole(ctx context.Context, id uint8) error
	GetPaginatedRoles(ctx context.Context, sort string, order string, page int, pageSize int) ([]model.Role, int64, error)
}

type RoleRepositoryImpl struct {
	db *gorm.DB
}

func NewRoleRepository(db *gorm.DB) RoleRepository {
	return &RoleRepositoryImpl{db: db}
}

func (rr *RoleRepositoryImpl) GetActiveRoleByName(ctx context.Context, name string) (*model.Role, error) {
	var role model.Role
	err := rr.db.WithContext(ctx).
		Where("name = ?", name).
		First(&role).Error

	if errors.Is(err, gorm.ErrRecordNotFound) {
		return nil, nil
	}
	return &role, err
}

func (rr *RoleRepositoryImpl) GetUnscopedRoleByName(ctx context.Context, name string) (*model.Role, error) {
	var role model.Role
	err := rr.db.WithContext(ctx).
		Unscoped(). // include soft-deleted records
		Where("name = ?", name).
		First(&role).Error

	if errors.Is(err, gorm.ErrRecordNotFound) {
		return nil, nil
	}
	return &role, err
}

func (rr *RoleRepositoryImpl) GetRoleByID(ctx context.Context, id uint8) (*model.Role, error) {
	var role model.Role
	result := rr.db.WithContext(ctx).Where("id = ?", id).First(&role)
	if errors.Is(result.Error, gorm.ErrRecordNotFound) {
		return nil, constants.ErrRecordNotFound
	}
	if result.Error != nil {
		return nil, result.Error
	}
	return &role, nil
}

func (rr *RoleRepositoryImpl) CreateRole(ctx context.Context, role *model.Role) error {
	return rr.db.WithContext(ctx).Create(role).Error
}

func (rr *RoleRepositoryImpl) UpdateRole(ctx context.Context, role *model.Role) error {
	return rr.db.WithContext(ctx).Save(role).Error
}

func (rr *RoleRepositoryImpl) DeleteRole(ctx context.Context, id uint8) error {
	return rr.db.WithContext(ctx).Delete(&model.Role{}, id).Error
}

func (rr *RoleRepositoryImpl) GetPaginatedRoles(ctx context.Context, sort string, order string, page int, pageSize int) ([]model.Role, int64, error) {
	var roles []model.Role
	var total int64

	countQuery := rr.db.WithContext(ctx).Model(&model.Role{})
	if err := countQuery.Count(&total).Error; err != nil {
		return nil, 0, err
	}

	dataQuery := rr.db.WithContext(ctx).Model(&model.Role{})

	if sort != "" && (order == "asc" || order == "desc") {
		dataQuery = dataQuery.Order(fmt.Sprintf("%s %s", sort, order))
	}

	offset := page * pageSize
	if err := dataQuery.Offset(offset).Limit(pageSize).Find(&roles).Error; err != nil {
		return nil, 0, err
	}

	return roles, total, nil
}
