package repository

import (
	"context"
	"errors"
	"fmt"
	"log/slog"
	"os"
	"time"

	"github.com/EdwinRincon/browersfc-api/api/constants"
	"github.com/EdwinRincon/browersfc-api/api/model"
	"github.com/google/uuid"
	"gorm.io/gorm"
)

// init initializes the Europe/Madrid timezone.
var defaultTimezone = os.Getenv("DEFAULT_TIMEZONE")
var europeMadrid *time.Location

func init() {
	var err error
	if defaultTimezone == "" {
		defaultTimezone = "Europe/Madrid"
	}
	europeMadrid, err = time.LoadLocation(defaultTimezone)
	if err != nil {
		// Fallback to UTC if the timezone cannot be loaded
		slog.Error("failed to load timezone, falling back to UTC", "timezone", defaultTimezone, "error", err)
		europeMadrid = time.UTC
	}
}

// UserRepository defines the interface for user data access.
type UserRepository interface {
	CreateUser(ctx context.Context, user *model.User) (*model.UserMin, error)
	GetUserByUsername(ctx context.Context, username string) (*model.User, error)
	GetUserByID(ctx context.Context, id string) (*model.User, error)
	ListUsers(ctx context.Context, page uint64) ([]*model.UserResponse, error)
	UpdateUser(ctx context.Context, user *model.User) (*model.UserMin, error)
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

// GetUserByUsername retrieves a user by their username, preloading the Role.
func (ur *UserRepositoryImpl) GetUserByUsername(ctx context.Context, username string) (*model.User, error) {
	var user model.User
	result := ur.db.WithContext(ctx).Preload("Role").Unscoped().Where("username = ?", username).First(&user)
	if errors.Is(result.Error, gorm.ErrRecordNotFound) {
		return nil, constants.ErrUserNotFound
	} else if result.Error != nil {
		return nil, result.Error
	}
	return &user, nil
}

// GetUserByID retrieves a user by their ID, preloading the Role.
func (ur *UserRepositoryImpl) GetUserByID(ctx context.Context, id string) (*model.User, error) {
	var user model.User
	result := ur.db.WithContext(ctx).Preload("Role").Unscoped().Where("id = ?", id).First(&user)
	if errors.Is(result.Error, gorm.ErrRecordNotFound) {
		return nil, constants.ErrUserNotFound
	} else if result.Error != nil {
		return nil, result.Error
	}
	return &user, nil
}

// ListUsers retrieves a paginated list of users with their roles.
func (ur *UserRepositoryImpl) ListUsers(ctx context.Context, page uint64) ([]*model.UserResponse, error) {
	pageSize := 10
	if page == 0 {
		page = 1
	}
	offset := int((page - 1) * uint64(pageSize))
	var users []*model.UserResponse
	query := ur.db.WithContext(ctx).Table("users").
		Select("users.id, users.name, users.last_name, users.username, users.birthdate, users.is_active, users.img_profile, users.img_banner, roles.name AS role_name").
		Joins("left join roles on users.role_id = roles.id").
		Limit(pageSize).
		Offset(offset).
		Scan(&users)

	if query.Error != nil {
		return nil, query.Error
	}
	// No need to check for empty slice; return the (possibly empty) slice.

	return users, nil
}

// CreateUser creates a new user in the database.  Returns specific errors.
func (ur *UserRepositoryImpl) CreateUser(ctx context.Context, user *model.User) (*model.UserMin, error) {
	err := ur.db.WithContext(ctx).Create(user).Error
	if err != nil {
		if errors.Is(err, gorm.ErrDuplicatedKey) { // Check for *specific* GORM error
			return nil, constants.ErrDuplicatedUsername
		}
		return nil, fmt.Errorf("%w: %v", constants.ErrCreateUser, err) // Wrap the original error
	}

	userResponse := &model.UserMin{
		ID:       user.ID,
		Username: user.Username,
	}

	return userResponse, nil
}

// UpdateUser updates an existing user's information.
func (ur *UserRepositoryImpl) UpdateUser(ctx context.Context, user *model.User) (*model.UserMin, error) {
	// Only update specific fields to avoid updating relationships
	updateFields := map[string]interface{}{
		"name":        user.Name,
		"last_name":   user.LastName,
		"username":    user.Username,
		"is_active":   user.IsActive,
		"birthdate":   user.Birthdate,
		"img_profile": user.ImgProfile,
		"img_banner":  user.ImgBanner,
	}

	result := ur.db.WithContext(ctx).Model(user).Updates(updateFields)
	if result.Error != nil {
		return nil, fmt.Errorf("%w: %v", constants.ErrUserUpdate, result.Error)
	}
	if result.RowsAffected == 0 {
		return nil, constants.ErrUserNotFound
	}

	userResponse := &model.UserMin{
		ID:       user.ID,
		Username: user.Username,
	}

	return userResponse, nil
}

// DeleteUser soft-deletes a user by setting their is_active to false and deleted_at to the current time.
func (ur *UserRepositoryImpl) DeleteUser(ctx context.Context, id string) error {
	_, err := uuid.Parse(id)
	if err != nil {
		return fmt.Errorf("%w: %v", constants.ErrInvalidUUID, err) // Wrap original error for context
	}

	nowLocal := time.Now().In(europeMadrid) // Use preloaded timezone

	updateFields := map[string]interface{}{
		"is_active":  false,
		"deleted_at": gorm.DeletedAt{Time: nowLocal, Valid: true},
	}

	result := ur.db.WithContext(ctx).
		Model(&model.User{}).
		Where("id = ?", id).
		Updates(updateFields)

	if result.Error != nil {
		return fmt.Errorf("%w: %v", constants.ErrUserDelete, result.Error) // Wrap the original error
	}

	if result.RowsAffected == 0 {
		return constants.ErrUserNotFound
	}

	return nil
}
