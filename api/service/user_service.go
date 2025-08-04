package service

import (
	"context"

	"github.com/EdwinRincon/browersfc-api/api/constants"
	"github.com/EdwinRincon/browersfc-api/api/model"
	"github.com/EdwinRincon/browersfc-api/api/repository"
)

type UserService interface {
	CreateUser(ctx context.Context, user *model.User) (*model.UserMin, error)
	GetUserByUsername(ctx context.Context, username string) (*model.User, error)
	ListUsers(ctx context.Context, page uint64) ([]*model.UserResponse, error)
	UpdateUser(ctx context.Context, userUpdate *model.UserUpdate, userID string) (*model.UserMin, error)
	DeleteUser(ctx context.Context, id string) error
}

type userService struct {
	UserRepository repository.UserRepository
}

func NewUserService(userRepo repository.UserRepository) UserService {
	return &userService{
		UserRepository: userRepo,
	}
}

func (s *userService) CreateUser(ctx context.Context, user *model.User) (*model.UserMin, error) {
	return s.UserRepository.CreateUser(ctx, user)
}

func (s *userService) GetUserByUsername(ctx context.Context, username string) (*model.User, error) {
	return s.UserRepository.GetUserByUsername(ctx, username)
}

func (s *userService) ListUsers(ctx context.Context, page uint64) ([]*model.UserResponse, error) {
	return s.UserRepository.ListUsers(ctx, page)
}

func (s *userService) UpdateUser(ctx context.Context, userUpdate *model.UserUpdate, userID string) (*model.UserMin, error) {
	// Obtener el usuario existente
	user, err := s.UserRepository.GetUserByID(ctx, userID)
	if err != nil {
		return nil, err
	}
	if user == nil {
		return nil, constants.ErrUserNotFound
	}

	// Actualizar solo los campos que no son nil
	if userUpdate.Name != nil {
		user.Name = *userUpdate.Name
	}
	if userUpdate.LastName != nil {
		user.LastName = *userUpdate.LastName
	}
	if userUpdate.Username != nil {
		user.Username = *userUpdate.Username
	}
	if userUpdate.Birthdate != nil {
		user.Birthdate = *userUpdate.Birthdate
	}
	if userUpdate.IsActive != nil {
		user.IsActive = *userUpdate.IsActive
	}
	if userUpdate.ImgProfile != nil {
		user.ImgProfile = *userUpdate.ImgProfile
	}
	if userUpdate.ImgBanner != nil {
		user.ImgBanner = *userUpdate.ImgBanner
	}

	updatedUser, err := s.UserRepository.UpdateUser(ctx, user)
	if err != nil {
		return nil, err
	}
	return updatedUser, nil
}

func (s *userService) DeleteUser(ctx context.Context, id string) error {
	return s.UserRepository.DeleteUser(ctx, id)
}
