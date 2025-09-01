package service

import (
	"context"
	"fmt"

	"github.com/EdwinRincon/browersfc-api/api/constants"
	"github.com/EdwinRincon/browersfc-api/api/dto"
	"github.com/EdwinRincon/browersfc-api/api/mapper"
	"github.com/EdwinRincon/browersfc-api/api/model"
	"github.com/EdwinRincon/browersfc-api/api/repository"
)

type UserService interface {
	CreateUser(ctx context.Context, user *model.User) (*dto.UserShort, error)
	GetUserByUsername(ctx context.Context, username string) (*model.User, error)
	GetPaginatedUsers(ctx context.Context, sort string, order string, page int, pageSize int) ([]model.User, int64, error)
	UpdateUser(ctx context.Context, userUpdate *dto.UpdateUserRequest, userID string) (*model.User, error)
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

func (s *userService) CreateUser(ctx context.Context, user *model.User) (*dto.UserShort, error) {
	// Check if a user with this username already exists
	existing, err := s.UserRepository.GetUserByUsername(ctx, user.Username)
	if err != nil {
		return nil, fmt.Errorf("failed to check existing user: %w", err)
	}
	if existing != nil {
		return nil, constants.ErrRecordAlreadyExists
	}

	if err := s.UserRepository.CreateUser(ctx, user); err != nil {
		return nil, fmt.Errorf("failed to create user: %w", err)
	}

	return mapper.ToUserShort(user), nil
}

func (s *userService) GetUserByUsername(ctx context.Context, username string) (*model.User, error) {
	user, err := s.UserRepository.GetUserByUsername(ctx, username)
	if err != nil {
		return nil, err
	}
	if user == nil {
		return nil, constants.ErrRecordNotFound
	}
	return user, nil
}

func (s *userService) GetPaginatedUsers(ctx context.Context, sort string, order string, page int, pageSize int) ([]model.User, int64, error) {
	return s.UserRepository.GetPaginatedUsers(ctx, sort, order, page, pageSize)
}

func (s *userService) UpdateUser(ctx context.Context, userUpdate *dto.UpdateUserRequest, userID string) (*model.User, error) {
	user, err := s.UserRepository.GetUserByID(ctx, userID)
	if err != nil {
		return nil, fmt.Errorf("failed to get user by ID: %w", err)
	}
	if user == nil {
		return nil, constants.ErrRecordNotFound
	}

	if userUpdate.Username != nil && *userUpdate.Username != user.Username {
		dup, err := s.UserRepository.GetUserByUsername(ctx, *userUpdate.Username)
		if err != nil {
			return nil, fmt.Errorf("failed to check duplicate username: %w", err)
		}
		if dup != nil && dup.ID != user.ID {
			return nil, constants.ErrRecordAlreadyExists
		}
	}

	mapper.UpdateUserFromDTO(user, userUpdate)

	if err := s.UserRepository.UpdateUser(ctx, userID, user); err != nil {
		return nil, fmt.Errorf("failed to update user: %w", err)
	}

	return user, nil
}

func (s *userService) DeleteUser(ctx context.Context, id string) error {
	return s.UserRepository.DeleteUser(ctx, id)
}
