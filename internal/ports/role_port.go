package ports

import (
	"context"

	"github.com/EdwinRincon/browersfc-api/api/model"
)

// RolePort defines the interface for role data access operations.
// This port separates the domain/application layer from the persistence adapter.
type RolePort interface {
	GetRoleByID(ctx context.Context, id uint64) (*model.Role, error)
	GetRoleByName(ctx context.Context, name string) (*model.Role, error)
	CreateRole(ctx context.Context, role *model.Role) error
	UpdateRole(ctx context.Context, role *model.Role) error
	DeleteRole(ctx context.Context, id uint64) error
	GetPaginatedRoles(ctx context.Context, sort string, order string, page int, pageSize int) ([]model.Role, int64, error)
}
