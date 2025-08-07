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
	GetActiveUserByUsername(ctx context.Context, username string) (*model.User, error)
	GetUnscopedUserByUsername(ctx context.Context, username string) (*model.User, error)
	GetActiveUserByID(ctx context.Context, id string) (*model.User, error)
	GetUnscopedUserByID(ctx context.Context, id string) (*model.User, error)
	GetPaginatedUsers(ctx context.Context, sort string, order string, page int, pageSize int) ([]model.User, int64, error)
	UpdateUser(ctx context.Context, id string, user *model.User) error
	DeleteUser(ctx context.Context, id string) error
	RestoreAndUpdateUser(ctx context.Context, user *model.User) error
}

// UserRepositoryImpl implements the UserRepository interface using GORM.
type UserRepositoryImpl struct {
	db *gorm.DB
}

// NewUserRepository creates a new UserRepository instance.
func NewUserRepository(db *gorm.DB) UserRepository {
	return &UserRepositoryImpl{db: db}
}

// GetActiveUserByUsername retrieves an active (not deleted) user by their username.
func (ur *UserRepositoryImpl) GetActiveUserByUsername(ctx context.Context, username string) (*model.User, error) {
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

// GetUnscopedUserByUsername retrieves a user by their username, including soft-deleted records.
func (ur *UserRepositoryImpl) GetUnscopedUserByUsername(ctx context.Context, username string) (*model.User, error) {
	var user model.User
	result := ur.db.WithContext(ctx).
		Preload("Role").
		Unscoped().
		Where("username = ?", username).
		First(&user)

	if errors.Is(result.Error, gorm.ErrRecordNotFound) {
		return nil, nil
	}
	return &user, result.Error
}

// GetUserByID retrieves an active (not deleted) user by their ID, preloading the Role.
func (ur *UserRepositoryImpl) GetActiveUserByID(ctx context.Context, id string) (*model.User, error) {
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

// GetUnscopedUserByID retrieves a user by their ID including soft-deleted records, preloading the Role.
func (ur *UserRepositoryImpl) GetUnscopedUserByID(ctx context.Context, id string) (*model.User, error) {
	var user model.User
	result := ur.db.WithContext(ctx).
		Preload("Role").
		Unscoped().
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

// RestoreAndUpdateUser restores a soft-deleted user and updates their information in a single transaction
func (ur *UserRepositoryImpl) RestoreAndUpdateUser(ctx context.Context, user *model.User) error {
	return ur.db.WithContext(ctx).Transaction(func(tx *gorm.DB) error {
		// First verify the user exists and is soft-deleted
		var existingUser model.User
		if err := tx.Unscoped().Where("id = ?", user.ID).First(&existingUser).Error; err != nil {
			if errors.Is(err, gorm.ErrRecordNotFound) {
				return fmt.Errorf("user not found: %w", err)
			}
			return fmt.Errorf("failed to find user: %w", err)
		}

		if !existingUser.DeletedAt.Valid {
			return fmt.Errorf("user is not soft-deleted")
		}

		// Preserve important metadata
		user.CreatedAt = existingUser.CreatedAt
		user.DeletedAt = gorm.DeletedAt{} // Explicitly set to zero value to restore

		// Save the entire user object (this will update all fields including deleted_at)
		if err := tx.Unscoped().Save(user).Error; err != nil {
			return fmt.Errorf("failed to restore and update user: %w", err)
		}

		// Verify the restoration was successful
		var restoredUser model.User
		if err := tx.Where("id = ?", user.ID).First(&restoredUser).Error; err != nil {
			return fmt.Errorf("failed to verify user restoration: %w", err)
		}

		if restoredUser.DeletedAt.Valid {
			return fmt.Errorf("user restoration failed: deleted_at field is still set")
		}

		return nil
	})
}
