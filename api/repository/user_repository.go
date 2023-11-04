package repository

import (
	"net/http"

	user "github.com/EdwinRincon/browersfc-api/api/model"
	"github.com/EdwinRincon/browersfc-api/helper"
	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
)

type UserRepository interface {
	CreateUser(c *gin.Context, user *user.Users) bool
	ListUsers(c *gin.Context, page uint64) ([]*user.UsersResponse, error)
}

type UserRepositoryImpl struct {
	db *gorm.DB
}

func NewUserRepository(db *gorm.DB) UserRepository {
	return &UserRepositoryImpl{db}
}

func (ur *UserRepositoryImpl) ListUsers(c *gin.Context, page uint64) ([]*user.UsersResponse, error) {
	var users []*user.UsersResponse
	query := ur.db.Table("users").Limit(10).Offset(int((page - 1) * 10)).Find(&users)
	if query.Error != nil {
		return nil, query.Error
	}

	if len(users) == 0 {
		return []*user.UsersResponse{}, nil
	}
	return users, nil
}

func (ur *UserRepositoryImpl) CreateUser(c *gin.Context, user *user.Users) bool {
	err := ur.db.Create(&user).Error
	if err != nil {
		helper.HandleError(c, http.StatusInternalServerError, "error al crear el usuario", err)
		return false
	}
	return true
}
