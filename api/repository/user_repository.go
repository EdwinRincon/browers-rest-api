package repository

import (
	"context"
	"fmt"
	"time"

	"github.com/EdwinRincon/browersfc-api/api/model"
	"gorm.io/gorm"
)

type UserRepository interface {
	CreateUser(ctx context.Context, user *model.Users) (*model.UserMin, error)
	GetUserByUsername(ctx context.Context, username string) (*model.Users, error)
	ListUsers(ctx context.Context, page uint64) ([]*model.UsersResponse, error)
	UpdateUser(ctx context.Context, user *model.Users) (*model.UserMin, error)
	DeleteUser(ctx context.Context, username string) error
}

type UserRepositoryImpl struct {
	db *gorm.DB
}

func NewUserRepository(db *gorm.DB) UserRepository {
	return &UserRepositoryImpl{db: db}
}

func (ur *UserRepositoryImpl) GetUserByUsername(ctx context.Context, username string) (*model.Users, error) {
	var user model.Users
	if err := ur.db.WithContext(ctx).Preload("Roles").Where("username = ?", username).First(&user).Error; err != nil {
		if err == gorm.ErrRecordNotFound {
			return nil, err
		}
		return nil, err
	}
	return &user, nil
}

func (ur *UserRepositoryImpl) ListUsers(ctx context.Context, page uint64) ([]*model.UsersResponse, error) {
	var users []*model.UsersResponse
	// Especifica las columnas que deseas recuperar
	query := ur.db.Table("users").
		Select("users.id, users.name, users.username, users.last_name, users.birthdate, users.is_active, users.img_profile, users.img_banner, roles.name AS role_name").
		Joins("left join roles on users.roles_id = roles.id").
		Limit(10).
		Offset(int((page - 1) * 10)).
		Find(&users)
	if query.Error != nil {
		return nil, query.Error
	}

	if len(users) == 0 {
		return []*model.UsersResponse{}, nil
	}
	return users, nil
}

func (ur *UserRepositoryImpl) CreateUser(ctx context.Context, user *model.Users) (*model.UserMin, error) {
	err := ur.db.WithContext(ctx).Create(user).Error
	if err != nil {
		return nil, err
	}

	userResponse := &model.UserMin{
		ID:       user.ID,
		Username: user.Username,
	}

	return userResponse, nil
}

func (ur *UserRepositoryImpl) UpdateUser(ctx context.Context, user *model.Users) (*model.UserMin, error) {
	err := ur.db.WithContext(ctx).Updates(user).Error
	if err != nil {
		return nil, err
	}

	userResponse := &model.UserMin{
		ID:       user.ID,
		Username: user.Username,
	}

	return userResponse, nil
}

// DeleteUser updates the is_active and deleted_at fields for a user identified by username
func (ur *UserRepositoryImpl) DeleteUser(ctx context.Context, username string) error {
	// Load the location for Europe/Madrid
	location, err := time.LoadLocation("Europe/Madrid")
	if err != nil {
		return fmt.Errorf("could not load time zone: %v", err)
	}

	// Get the current time in the specified time zone
	nowLocal := time.Now().In(location)
	fmt.Println("Current time in Europe/Madrid:", nowLocal)

	// Prepare the update map with the local time
	updates := map[string]interface{}{
		"is_active":  "N",
		"deleted_at": gorm.DeletedAt{Time: nowLocal, Valid: true},
	}

	// Perform the update
	err = ur.db.WithContext(ctx).Table("Users").Where("username = ?", username).UpdateColumns(updates).Error
	if err != nil {
		return err
	}
	return nil
}
