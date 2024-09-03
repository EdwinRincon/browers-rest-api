package service

import (
	"context"

	"github.com/EdwinRincon/browersfc-api/api/model"
	"github.com/EdwinRincon/browersfc-api/api/repository"
)

type UserService interface {
	CreateUser(ctx context.Context, user *model.Users) (*model.UserMin, error)
	GetUserByUsername(ctx context.Context, username string) (*model.Users, error)
	ListUsers(ctx context.Context, page uint64) ([]*model.UsersResponse, error)
	UpdateUser(ctx context.Context, user *model.Users) (*model.UserMin, error)
	DeleteUser(ctx context.Context, username string) error
}

type userService struct {
	UserRepository repository.UserRepository
}

func NewUserService(userRepo repository.UserRepository) UserService {
	return &userService{
		UserRepository: userRepo,
	}
}

func (s *userService) CreateUser(ctx context.Context, user *model.Users) (*model.UserMin, error) {
	createdUser, err := s.UserRepository.CreateUser(ctx, user)
	if err != nil {
		return nil, err
	}
	return createdUser, nil
}

func (s *userService) GetUserByUsername(ctx context.Context, username string) (*model.Users, error) {
	user, err := s.UserRepository.GetUserByUsername(ctx, username)
	if err != nil {
		return nil, err
	}
	return user, nil
}

func (s *userService) ListUsers(ctx context.Context, page uint64) ([]*model.UsersResponse, error) {
	users, err := s.UserRepository.ListUsers(ctx, page)
	if err != nil {
		return nil, err
	}
	return users, nil
}

func (s *userService) UpdateUser(ctx context.Context, user *model.Users) (*model.UserMin, error) {
	updatedUser, err := s.UserRepository.UpdateUser(ctx, user)
	if err != nil {
		return nil, err
	}
	return updatedUser, nil
}

func (s *userService) DeleteUser(ctx context.Context, username string) error {
	err := s.UserRepository.DeleteUser(ctx, username)
	if err != nil {
		return err
	}
	return nil
}
