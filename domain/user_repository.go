package domain

import (
	"context"
)

// UserRepository defines the interface for user persistence operations.
// This port belongs in the domain layer following hexagonal architecture principles.
// Infrastructure adapters will implement this interface.
type UserRepository interface {
	// Standard CRUD methods (updated to match business service needs)
	CreateUser(ctx context.Context, user *User) error
	GetUserByUsername(ctx context.Context, username string) (*User, error)
	GetUserByID(ctx context.Context, id string) (*User, error)
	GetPaginatedUsers(ctx context.Context, sort string, order string, page int, pageSize int) ([]User, int64, error)
	UpdateUser(ctx context.Context, id string, user *User) error
	DeleteUser(ctx context.Context, id string) error
}
