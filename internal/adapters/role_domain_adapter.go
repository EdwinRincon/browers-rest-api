package adapters

import (
	"context"

	"github.com/EdwinRincon/browersfc-api/api/model"
	"github.com/EdwinRincon/browersfc-api/internal/infrastructure/persistence"
	"github.com/EdwinRincon/browersfc-api/internal/ports"
)

// RoleDomainAdapter implements the RolePort interface and handles data persistence.
// It acts as an adapter between the domain layer and the infrastructure layer.
type RoleDomainAdapter struct {
	repository persistence.RoleRepository
}

// NewRoleDomainAdapter creates a new RoleDomainAdapter instance.
func NewRoleDomainAdapter(repository persistence.RoleRepository) ports.RolePort {
	return &RoleDomainAdapter{
		repository: repository,
	}
}

// CreateRole creates a new role in the data store.
func (a *RoleDomainAdapter) CreateRole(ctx context.Context, role *model.Role) error {
	return a.repository.CreateRole(ctx, role)
}

// GetRoleByID retrieves a role by its ID.
func (a *RoleDomainAdapter) GetRoleByID(ctx context.Context, id uint64) (*model.Role, error) {
	return a.repository.GetRoleByID(ctx, id)
}

// GetRoleByName retrieves a role by its name.
func (a *RoleDomainAdapter) GetRoleByName(ctx context.Context, name string) (*model.Role, error) {
	return a.repository.GetRoleByName(ctx, name)
}

// UpdateRole updates an existing role in the data store.
func (a *RoleDomainAdapter) UpdateRole(ctx context.Context, role *model.Role) error {
	return a.repository.UpdateRole(ctx, role)
}

// DeleteRole deletes a role from the data store.
func (a *RoleDomainAdapter) DeleteRole(ctx context.Context, id uint64) error {
	return a.repository.DeleteRole(ctx, id)
}

// GetPaginatedRoles retrieves paginated roles from the data store.
func (a *RoleDomainAdapter) GetPaginatedRoles(ctx context.Context, sort string, order string, page int, pageSize int) ([]model.Role, int64, error) {
	return a.repository.GetPaginatedRoles(ctx, sort, order, page, pageSize)
}
