package service

import (
	"context"
	"fmt"

	"github.com/EdwinRincon/browersfc-api/api/constants"
	"github.com/EdwinRincon/browersfc-api/domain"
)

// UserDomainService encapsulates business logic for user operations.
type UserDomainService struct {
	userRepository domain.UserRepository
}

// NewUserDomainService creates a new UserDomainService instance.
func NewUserDomainService(userRepository domain.UserRepository) *UserDomainService {
	return &UserDomainService{
		userRepository: userRepository,
	}
}

// CreateUser creates a new user after validating business rules.
func (s *UserDomainService) CreateUser(ctx context.Context, user *domain.User) (*domain.User, error) {
	// Validate domain entity
	if !user.IsValid() {
		return nil, constants.ErrInvalidData
	}

	// Check if a user with this username already exists
	existing, err := s.userRepository.GetUserByUsername(ctx, user.Username)
	if err != nil {
		return nil, fmt.Errorf("failed to check existing user: %w", err)
	}
	if existing != nil {
		return nil, constants.ErrRecordAlreadyExists
	}

	// Create the user
	err = s.userRepository.CreateUser(ctx, user)
	if err != nil {
		return nil, fmt.Errorf("failed to create user: %w", err)
	}

	return user, nil
}

// GetUserByUsername retrieves a user by username.
func (s *UserDomainService) GetUserByUsername(ctx context.Context, username string) (*domain.User, error) {
	user, err := s.userRepository.GetUserByUsername(ctx, username)
	if err != nil {
		return nil, err
	}
	if user == nil {
		return nil, constants.ErrRecordNotFound
	}
	return user, nil
}

// GetUserByID retrieves a user by ID.
func (s *UserDomainService) GetUserByID(ctx context.Context, id string) (*domain.User, error) {
	user, err := s.userRepository.GetUserByID(ctx, id)
	if err != nil {
		return nil, err
	}
	if user == nil {
		return nil, constants.ErrRecordNotFound
	}
	return user, nil
}

// GetPaginatedUsers retrieves a paginated list of users.
func (s *UserDomainService) GetPaginatedUsers(ctx context.Context, sort string, order string, page int, pageSize int) ([]domain.User, int64, error) {
	return s.userRepository.GetPaginatedUsers(ctx, sort, order, page, pageSize)
}

// UpdateUser updates an existing user after validating business rules.
func (s *UserDomainService) UpdateUser(ctx context.Context, id string, updates *domain.User) (*domain.User, error) {
	// Get existing user
	existingUser, err := s.userRepository.GetUserByID(ctx, id)
	if err != nil {
		return nil, fmt.Errorf("failed to get user by ID: %w", err)
	}
	if existingUser == nil {
		return nil, constants.ErrRecordNotFound
	}

	// Apply updates to existing user
	if updates.Name != "" {
		existingUser.Name = updates.Name
	}
	if updates.LastName != "" {
		existingUser.LastName = updates.LastName
	}
	if updates.Username != "" && updates.Username != existingUser.Username {
		// Check for duplicate username
		dup, err := s.userRepository.GetUserByUsername(ctx, updates.Username)
		if err != nil {
			return nil, fmt.Errorf("failed to check duplicate username: %w", err)
		}
		if dup != nil && dup.ID != existingUser.ID {
			return nil, constants.ErrRecordAlreadyExists
		}
		existingUser.Username = updates.Username
	}
	if updates.Birthdate != nil {
		existingUser.Birthdate = updates.Birthdate
	}
	if updates.ImgProfile != "" {
		existingUser.ImgProfile = updates.ImgProfile
	}
	if updates.ImgBanner != "" {
		existingUser.ImgBanner = updates.ImgBanner
	}
	if updates.RoleID != 0 {
		existingUser.RoleID = updates.RoleID
	}

	// Validate updated user
	if !existingUser.IsValid() {
		return nil, constants.ErrInvalidData
	}

	// Update the user
	err = s.userRepository.UpdateUser(ctx, id, existingUser)
	if err != nil {
		return nil, fmt.Errorf("failed to update user: %w", err)
	}

	return existingUser, nil
}

// DeleteUser deletes a user by ID.
func (s *UserDomainService) DeleteUser(ctx context.Context, id string) error {
	// Check if user exists
	user, err := s.userRepository.GetUserByID(ctx, id)
	if err != nil {
		return fmt.Errorf("failed to get user by ID: %w", err)
	}
	if user == nil {
		return constants.ErrRecordNotFound
	}

	return s.userRepository.DeleteUser(ctx, id)
}
