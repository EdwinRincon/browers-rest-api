package persistence

import (
	"context"
	"errors"
	"fmt"

	"github.com/EdwinRincon/browersfc-api/adapter/mapper"
	"github.com/EdwinRincon/browersfc-api/api/model"
	"github.com/EdwinRincon/browersfc-api/domain"
	"gorm.io/gorm"
)

type RoleRepositoryImpl struct {
	db     *gorm.DB
	mapper *mapper.RoleMapper
}

func NewRoleRepository(db *gorm.DB) domain.RoleRepository {
	return &RoleRepositoryImpl{
		db:     db,
		mapper: mapper.NewRoleMapper(),
	}
}

func (rr *RoleRepositoryImpl) GetRoleByName(ctx context.Context, name string) (*domain.Role, error) {
	var roleModel model.Role
	err := rr.db.WithContext(ctx).
		Where("name = ?", name).
		First(&roleModel).Error

	if errors.Is(err, gorm.ErrRecordNotFound) {
		return nil, nil
	}
	if err != nil {
		return nil, err
	}

	return rr.mapper.ModelToDomain(&roleModel), nil
}

func (rr *RoleRepositoryImpl) GetRoleByID(ctx context.Context, id uint64) (*domain.Role, error) {
	var roleModel model.Role
	result := rr.db.WithContext(ctx).Where("id = ?", id).First(&roleModel)
	if errors.Is(result.Error, gorm.ErrRecordNotFound) {
		return nil, nil
	}
	if result.Error != nil {
		return nil, result.Error
	}
	return rr.mapper.ModelToDomain(&roleModel), nil
}

func (rr *RoleRepositoryImpl) CreateRole(ctx context.Context, role *domain.Role) error {
	roleModel := rr.mapper.DomainToModel(role)
	err := rr.db.WithContext(ctx).Create(roleModel).Error
	if err != nil {
		return err
	}

	// Update the domain entity with generated ID and timestamps
	domainRole := rr.mapper.ModelToDomain(roleModel)
	*role = *domainRole
	return nil
}

func (rr *RoleRepositoryImpl) UpdateRole(ctx context.Context, role *domain.Role) error {
	result := rr.db.WithContext(ctx).Model(&model.Role{}).
		Where("id = ?", role.ID).
		Updates(map[string]interface{}{
			"name":        role.Name,
			"description": role.Description,
		})
	return result.Error
}

func (rr *RoleRepositoryImpl) DeleteRole(ctx context.Context, id uint64) error {
	return rr.db.WithContext(ctx).Delete(&model.Role{}, id).Error
}

func (rr *RoleRepositoryImpl) GetPaginatedRoles(ctx context.Context, sort string, order string, page int, pageSize int) ([]domain.Role, int64, error) {
	var roleModels []model.Role
	var total int64

	countQuery := rr.db.WithContext(ctx).Model(&model.Role{})
	if err := countQuery.Count(&total).Error; err != nil {
		return nil, 0, err
	}

	dataQuery := rr.db.WithContext(ctx).Model(&model.Role{})

	if sort != "" && (order == "asc" || order == "desc") {
		// Escape the sort field with backticks to handle reserved words
		dataQuery = dataQuery.Order(fmt.Sprintf("`%s` %s", sort, order))
	}

	offset := page * pageSize
	if err := dataQuery.Offset(offset).Limit(pageSize).Find(&roleModels).Error; err != nil {
		return nil, 0, err
	}

	// Convert to domain entities
	domainRoles := make([]domain.Role, len(roleModels))
	for i, roleModel := range roleModels {
		domainRole := rr.mapper.ModelToDomain(&roleModel)
		domainRoles[i] = *domainRole
	}

	return domainRoles, total, nil
}
