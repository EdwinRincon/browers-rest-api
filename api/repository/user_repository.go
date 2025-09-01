package repository

import (
	"context"
	"errors"
	"fmt"

	"github.com/EdwinRincon/browersfc-api/api/model"
	"gorm.io/gorm"
)

// UserRepository defines the interface for user data access.
type UserRepository interface {
	CreateUser(ctx context.Context, user *model.User) error
	GetUserByUsername(ctx context.Context, username string) (*model.User, error)
	GetUserByID(ctx context.Context, id string) (*model.User, error)
	GetPaginatedUsers(ctx context.Context, sort string, order string, page int, pageSize int) ([]model.User, int64, error)
	UpdateUser(ctx context.Context, id string, user *model.User) error
	DeleteUser(ctx context.Context, id string) error
}

// UserRepositoryImpl implements the UserRepository interface using GORM.
type UserRepositoryImpl struct {
	db *gorm.DB
}

// NewUserRepository creates a new UserRepository instance.
func NewUserRepository(db *gorm.DB) UserRepository {
	return &UserRepositoryImpl{db: db}
}

func (ur *UserRepositoryImpl) GetUserByUsername(ctx context.Context, username string) (*model.User, error) {
	var user model.User
	result := ur.db.WithContext(ctx).
		Preload("Role").
		Where("username = ?", username).
		First(&user)

	if errors.Is(result.Error, gorm.ErrRecordNotFound) {
		return nil, nil
	}
	if result.Error != nil {
		return nil, fmt.Errorf("error getting user by username: %w", result.Error)
	}
	return &user, nil
}

// GetUserByID retrieves a user by their ID, preloading the Role.
func (ur *UserRepositoryImpl) GetUserByID(ctx context.Context, id string) (*model.User, error) {
	var user model.User
	result := ur.db.WithContext(ctx).
		Preload("Role").
		Where("id = ?", id).
		First(&user)

	if errors.Is(result.Error, gorm.ErrRecordNotFound) {
		return nil, nil
	}
	return &user, result.Error
}

// GetPaginatedUsers retrieves a paginated list of users with their roles and total count.
func (ur *UserRepositoryImpl) GetPaginatedUsers(ctx context.Context, sort string, order string, page int, pageSize int) ([]model.User, int64, error) {
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
		query = query.Order(fmt.Sprintf("%s %s", sort, order))
	}

	// Apply pagination
	offset := page * pageSize
	query = query.Offset(offset).Limit(pageSize)

	// Execute the query
	if err := query.Find(&users).Error; err != nil {
		return nil, 0, fmt.Errorf("error fetching users: %w", err)
	}

	return users, total, nil
}

func (ur *UserRepositoryImpl) CreateUser(ctx context.Context, user *model.User) error {
	return ur.db.WithContext(ctx).Create(user).Error
}

func (ur *UserRepositoryImpl) UpdateUser(ctx context.Context, id string, user *model.User) error {
	return ur.db.WithContext(ctx).
		Model(&model.User{}).
		Where("id = ?", id).
		Select("*").
		Updates(user).Error
}

func (ur *UserRepositoryImpl) DeleteUser(ctx context.Context, id string) error {
	return ur.db.WithContext(ctx).Delete(&model.User{}, "id = ?", id).Error
}
