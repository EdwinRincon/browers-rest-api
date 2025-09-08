package adapters

import (
	"context"

	"github.com/EdwinRincon/browersfc-api/api/mapper"
	"github.com/EdwinRincon/browersfc-api/domain"
	"github.com/EdwinRincon/browersfc-api/internal/infrastructure/persistence"
	"github.com/EdwinRincon/browersfc-api/internal/ports"
)

// UserDomainAdapter implements the UserPort interface and bridges domain and persistence layers.
type UserDomainAdapter struct {
	userRepo persistence.UserRepository
}

// NewUserDomainAdapter creates a new UserDomainAdapter instance.
func NewUserDomainAdapter(userRepo persistence.UserRepository) ports.UserPort {
	return &UserDomainAdapter{
		userRepo: userRepo,
	}
}

func (a *UserDomainAdapter) CreateUser(ctx context.Context, user *domain.User) (*domain.User, error) {
	modelUser := mapper.UserToModel(user)
	err := a.userRepo.CreateUser(ctx, modelUser)
	if err != nil {
		return nil, err
	}
	// Return the domain user with any updates from persistence (like generated ID, timestamps)
	return mapper.UserToDomain(modelUser), nil
}

func (a *UserDomainAdapter) GetUserByUsername(ctx context.Context, username string) (*domain.User, error) {
	modelUser, err := a.userRepo.GetUserByUsername(ctx, username)
	if err != nil {
		return nil, err
	}
	return mapper.UserToDomain(modelUser), nil
}

func (a *UserDomainAdapter) GetUserByID(ctx context.Context, id string) (*domain.User, error) {
	modelUser, err := a.userRepo.GetUserByID(ctx, id)
	if err != nil {
		return nil, err
	}
	return mapper.UserToDomain(modelUser), nil
}

func (a *UserDomainAdapter) GetPaginatedUsers(ctx context.Context, sort string, order string, page int, pageSize int) ([]domain.User, int64, error) {
	modelUsers, total, err := a.userRepo.GetPaginatedUsers(ctx, sort, order, page, pageSize)
	if err != nil {
		return nil, 0, err
	}
	return mapper.UserListToDomain(modelUsers), total, nil
}

func (a *UserDomainAdapter) UpdateUser(ctx context.Context, id string, user *domain.User) (*domain.User, error) {
	modelUser := mapper.UserToModel(user)
	err := a.userRepo.UpdateUser(ctx, id, modelUser)
	if err != nil {
		return nil, err
	}
	// Return the updated domain user
	return mapper.UserToDomain(modelUser), nil
}

func (a *UserDomainAdapter) DeleteUser(ctx context.Context, id string) error {
	return a.userRepo.DeleteUser(ctx, id)
}
