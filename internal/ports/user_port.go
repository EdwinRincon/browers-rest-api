package ports

import (
	"context"

	"github.com/EdwinRincon/browersfc-api/api/model"
)

// UserPort defines the interface for user data access operations.
// This port separates the domain/application layer from the persistence adapter.
type UserPort interface {
	CreateUser(ctx context.Context, user *model.User) error
	GetUserByUsername(ctx context.Context, username string) (*model.User, error)
	GetUserByID(ctx context.Context, id string) (*model.User, error)
	GetPaginatedUsers(ctx context.Context, sort string, order string, page int, pageSize int) ([]model.User, int64, error)
	UpdateUser(ctx context.Context, id string, user *model.User) error
	DeleteUser(ctx context.Context, id string) error
}
