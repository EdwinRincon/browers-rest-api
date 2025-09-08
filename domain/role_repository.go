package domain

import (
	"context"
)

// RoleRepository defines the interface for role persistence operations.
// This port belongs in the domain layer following hexagonal architecture principles.
// Infrastructure adapters will implement this interface.
type RoleRepository interface {
	GetRoleByID(ctx context.Context, id uint64) (*Role, error)
	GetRoleByName(ctx context.Context, name string) (*Role, error)
	CreateRole(ctx context.Context, role *Role) error
	UpdateRole(ctx context.Context, role *Role) error
	DeleteRole(ctx context.Context, id uint64) error
	GetPaginatedRoles(ctx context.Context, sort string, order string, page int, pageSize int) ([]Role, int64, error)
}
