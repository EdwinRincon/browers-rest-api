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

// UserRepositoryImpl implements the UserRepository interface using GORM.
type UserRepositoryImpl struct {
	db     *gorm.DB
	mapper *mapper.UserMapper
}

// NewUserRepository creates a new UserRepository instance.
func NewUserRepository(db *gorm.DB) *UserRepositoryImpl {
	return &UserRepositoryImpl{
		db:     db,
		mapper: mapper.NewUserMapper(),
	}
}

func (ur *UserRepositoryImpl) GetUserByUsername(ctx context.Context, username string) (*domain.User, error) {
	var userModel model.User
	result := ur.db.WithContext(ctx).
		Preload("Role").
		Where("username = ?", username).
		First(&userModel)

	if errors.Is(result.Error, gorm.ErrRecordNotFound) {
		return nil, nil
	}
	if result.Error != nil {
		return nil, fmt.Errorf("error getting user by username: %w", result.Error)
	}
	
	return ur.mapper.ModelToDomain(&userModel), nil
}

// GetUserByID retrieves a user by their ID, preloading the Role.
func (ur *UserRepositoryImpl) GetUserByID(ctx context.Context, id string) (*domain.User, error) {
	var userModel model.User
	result := ur.db.WithContext(ctx).
		Preload("Role").
		Where("id = ?", id).
		First(&userModel)

	if errors.Is(result.Error, gorm.ErrRecordNotFound) {
		return nil, nil
	}
	if result.Error != nil {
		return nil, result.Error
	}
	
	return ur.mapper.ModelToDomain(&userModel), nil
}

// GetPaginatedUsers retrieves a paginated list of users with their roles and total count.
func (ur *UserRepositoryImpl) GetPaginatedUsers(ctx context.Context, sort string, order string, page int, pageSize int) ([]domain.User, int64, error) {
	var users []model.User
	var total int64

	// Count total records
	countQuery := ur.db.WithContext(ctx).Model(&model.User{})
	if err := countQuery.Count(&total).Error; err != nil {
		return nil, 0, fmt.Errorf("error counting total users: %w", err)
	}

	// Build the data query with eager loading
	query := ur.db.WithContext(ctx).Model(&model.User{}).
		Preload("Role")

	// Apply sorting if provided
	if sort != "" && (order == "asc" || order == "desc") {
		query = query.Order(fmt.Sprintf("`%s` %s", sort, order))
	}

	// Apply pagination
	offset := page * pageSize
	query = query.Offset(offset).Limit(pageSize)

	// Execute the query
	if err := query.Find(&users).Error; err != nil {
		return nil, 0, fmt.Errorf("error fetching users: %w", err)
	}

	return ur.mapper.ModelListToDomain(users), total, nil
}

func (ur *UserRepositoryImpl) CreateUser(ctx context.Context, user *domain.User) error {
	userModel := ur.mapper.DomainToModel(user)
	return ur.db.WithContext(ctx).Create(userModel).Error
}

func (ur *UserRepositoryImpl) UpdateUser(ctx context.Context, id string, user *domain.User) error {
	userModel := ur.mapper.DomainToModel(user)
	return ur.db.WithContext(ctx).
		Model(&model.User{}).
		Where("id = ?", id).
		Select("*").
		Updates(userModel).Error
}

func (ur *UserRepositoryImpl) DeleteUser(ctx context.Context, id string) error {
	return ur.db.WithContext(ctx).Delete(&model.User{}, "id = ?", id).Error
}
