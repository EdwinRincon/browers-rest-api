package service

import (
	"context"
	"errors"
	"fmt"

	"github.com/EdwinRincon/browersfc-api/api/constants"
	"github.com/EdwinRincon/browersfc-api/domain"
	"github.com/EdwinRincon/browersfc-api/internal/ports"
)

var (
	ErrInvalidUserData = errors.New("invalid user data")
)

// UserDomainService encapsulates business logic for user operations.
type UserDomainService struct {
	userPort ports.UserPort
}

// Compile-time check to ensure UserDomainService implements UserPort
var _ ports.UserPort = (*UserDomainService)(nil)

// NewUserDomainService creates a new UserDomainService instance.
func NewUserDomainService(userPort ports.UserPort) *UserDomainService {
	return &UserDomainService{
		userPort: userPort,
	}
}

// CreateUser creates a new user after validating business rules.
func (s *UserDomainService) CreateUser(ctx context.Context, user *domain.User) (*domain.User, error) {
	// Validate domain entity
	if !user.IsValid() {
		return nil, ErrInvalidUserData
	}

	// Check if a user with this username already exists
	existing, err := s.userPort.GetUserByUsername(ctx, user.Username)
	if err != nil {
		return nil, fmt.Errorf("failed to check existing user: %w", err)
	}
	if existing != nil {
		return nil, constants.ErrRecordAlreadyExists
	}

	// Create the user
	createdUser, err := s.userPort.CreateUser(ctx, user)
	if err != nil {
		return nil, fmt.Errorf("failed to create user: %w", err)
	}

	return createdUser, nil
}

// GetUserByUsername retrieves a user by username.
func (s *UserDomainService) GetUserByUsername(ctx context.Context, username string) (*domain.User, error) {
	user, err := s.userPort.GetUserByUsername(ctx, username)
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
	user, err := s.userPort.GetUserByID(ctx, id)
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
	return s.userPort.GetPaginatedUsers(ctx, sort, order, page, pageSize)
}

// UpdateUser updates an existing user after validating business rules.
func (s *UserDomainService) UpdateUser(ctx context.Context, id string, updates *domain.User) (*domain.User, error) {
	// Get existing user
	existingUser, err := s.userPort.GetUserByID(ctx, id)
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
		dup, err := s.userPort.GetUserByUsername(ctx, updates.Username)
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
		return nil, ErrInvalidUserData
	}

	// Update the user
	updatedUser, err := s.userPort.UpdateUser(ctx, id, existingUser)
	if err != nil {
		return nil, fmt.Errorf("failed to update user: %w", err)
	}

	return updatedUser, nil
}

// DeleteUser deletes a user by ID.
func (s *UserDomainService) DeleteUser(ctx context.Context, id string) error {
	// Check if user exists
	user, err := s.userPort.GetUserByID(ctx, id)
	if err != nil {
		return fmt.Errorf("failed to get user by ID: %w", err)
	}
	if user == nil {
		return constants.ErrRecordNotFound
	}

	return s.userPort.DeleteUser(ctx, id)
}
