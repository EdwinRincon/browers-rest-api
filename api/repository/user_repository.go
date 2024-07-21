package repository

import (
	"context"

	"github.com/EdwinRincon/browersfc-api/api/model"
	"gorm.io/gorm"
)

type UserRepository interface {
	CreateUser(ctx context.Context, user *model.Users) (*model.UsersResponse, error)
	ListUsers(ctx context.Context, page uint64) ([]*model.UsersResponse, error)
}

type UserRepositoryImpl struct {
	db *gorm.DB
}

func NewUserRepository(db *gorm.DB) UserRepository {
	return &UserRepositoryImpl{db: db}
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

func (ur *UserRepositoryImpl) CreateUser(ctx context.Context, user *model.Users) (*model.UsersResponse, error) {
	err := ur.db.WithContext(ctx).Create(user).Error
	if err != nil {
		return nil, err
	}

	userResponse := &model.UsersResponse{
		ID:         user.ID,
		Name:       user.Name,
		LastName:   user.LastName,
		Username:   user.Username,
		IsActive:   user.IsActive,
		Birthdate:  user.Birthdate,
		ImgProfile: user.ImgProfile,
		ImgBanner:  user.ImgBanner,
		RoleName:   user.Roles.Name,
	}

	return userResponse, nil
}
