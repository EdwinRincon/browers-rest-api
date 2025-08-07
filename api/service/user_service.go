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
	// First check if there's an active user with this username
	activeUser, err := s.UserRepository.GetActiveUserByUsername(ctx, user.Username)
	if err != nil {
		return nil, fmt.Errorf("failed to check existing active user: %w", err)
	}
	if activeUser != nil {
		return nil, constants.ErrRecordAlreadyExists
	}

	// If no active user exists, check for soft-deleted user
	existing, err := s.UserRepository.GetUnscopedUserByUsername(ctx, user.Username)
	if err != nil {
		return nil, fmt.Errorf("failed to check existing user: %w", err)
	}

	if existing != nil && existing.DeletedAt.Valid {
		// Update all relevant fields from the new user
		existing.Name = user.Name
		existing.LastName = user.LastName
		existing.Username = user.Username
		existing.RoleID = user.RoleID
		existing.ImgProfile = user.ImgProfile
		existing.ImgBanner = user.ImgBanner
		existing.Birthdate = user.Birthdate

		// Restore and update the user in a transaction
		err := s.UserRepository.RestoreAndUpdateUser(ctx, existing)
		if err != nil {
			return nil, fmt.Errorf("failed to restore and update user: %w", err)
		}

		return mapper.ToUserShort(existing), nil
	}

	if err := s.UserRepository.CreateUser(ctx, user); err != nil {
		return nil, fmt.Errorf("failed to create user: %w", err)
	}

	return mapper.ToUserShort(user), nil
}

func (s *userService) GetUserByUsername(ctx context.Context, username string) (*model.User, error) {
	user, err := s.UserRepository.GetActiveUserByUsername(ctx, username)
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
	user, err := s.UserRepository.GetActiveUserByID(ctx, userID)
	if err != nil {
		return nil, fmt.Errorf("failed to get user by ID: %w", err)
	}
	if user == nil {
		return nil, constants.ErrRecordNotFound
	}

	if userUpdate.Username != nil && *userUpdate.Username != user.Username {
		dup, err := s.UserRepository.GetActiveUserByUsername(ctx, *userUpdate.Username)
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
