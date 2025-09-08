package service

import (
	"context"
	"errors"
	"fmt"

	"github.com/EdwinRincon/browersfc-api/api/constants"
	"github.com/EdwinRincon/browersfc-api/api/mapper"
	"github.com/EdwinRincon/browersfc-api/domain"
	"github.com/EdwinRincon/browersfc-api/internal/ports"
)

var (
	ErrInvalidRole            = errors.New("invalid role data")
	ErrCannotDeleteSystemRole = errors.New("cannot delete system role")
)

type RoleDomainService struct {
	rolePort   ports.RolePort
	roleMapper *mapper.RoleDomainMapper
}

func NewRoleDomainService(
	rolePort ports.RolePort,
	roleMapper *mapper.RoleDomainMapper,
) *RoleDomainService {
	return &RoleDomainService{
		rolePort:   rolePort,
		roleMapper: roleMapper,
	}
}

// CreateRole creates a new role with domain validation.
func (s *RoleDomainService) CreateRole(ctx context.Context, role *domain.Role) (*domain.Role, error) {
	// Domain validation
	if !role.IsValid() {
		return nil, ErrInvalidRole
	}

	// Check if role already exists
	existingRole, err := s.rolePort.GetRoleByName(ctx, role.Name)
	if err != nil {
		return nil, fmt.Errorf("failed to check existing role: %w", err)
	}
	if existingRole != nil {
		return nil, constants.ErrRecordAlreadyExists
	}

	// Convert to model and create
	modelRole := s.roleMapper.ToModel(role)
	err = s.rolePort.CreateRole(ctx, modelRole)
	if err != nil {
		return nil, fmt.Errorf("failed to create role: %w", err)
	}

	// Return the created role as domain entity
	return s.roleMapper.ToDomain(modelRole), nil
}

// GetRoleByID retrieves a role by its ID with domain logic.
func (s *RoleDomainService) GetRoleByID(ctx context.Context, id uint64) (*domain.Role, error) {
	if id == 0 {
		return nil, ErrInvalidRole
	}

	modelRole, err := s.rolePort.GetRoleByID(ctx, id)
	if err != nil {
		return nil, fmt.Errorf("failed to get role: %w", err)
	}
	if modelRole == nil {
		return nil, constants.ErrRecordNotFound
	}

	return s.roleMapper.ToDomain(modelRole), nil
}

// GetRoleByName retrieves a role by its name with domain logic.
func (s *RoleDomainService) GetRoleByName(ctx context.Context, name string) (*domain.Role, error) {
	if name == "" {
		return nil, ErrInvalidRole
	}

	modelRole, err := s.rolePort.GetRoleByName(ctx, name)
	if err != nil {
		return nil, fmt.Errorf("failed to get role by name: %w", err)
	}
	if modelRole == nil {
		return nil, constants.ErrRecordNotFound
	}

	return s.roleMapper.ToDomain(modelRole), nil
}

// UpdateRole updates an existing role with domain validation.
func (s *RoleDomainService) UpdateRole(ctx context.Context, role *domain.Role) (*domain.Role, error) {
	// Domain validation
	if !role.IsValid() {
		return nil, ErrInvalidRole
	}

	// Check if role exists
	existingRole, err := s.rolePort.GetRoleByID(ctx, role.ID)
	if err != nil {
		return nil, fmt.Errorf("failed to check existing role: %w", err)
	}
	if existingRole == nil {
		return nil, constants.ErrRecordNotFound
	}

	// Check if name is being changed to an existing name
	if existingRole.Name != role.Name {
		conflictingRole, err := s.rolePort.GetRoleByName(ctx, role.Name)
		if err != nil {
			return nil, fmt.Errorf("failed to check name conflict: %w", err)
		}
		if conflictingRole != nil && conflictingRole.ID != role.ID {
			return nil, constants.ErrRecordAlreadyExists
		}
	}

	// Convert to model and update
	modelRole := s.roleMapper.ToModel(role)
	err = s.rolePort.UpdateRole(ctx, modelRole)
	if err != nil {
		return nil, fmt.Errorf("failed to update role: %w", err)
	}

	// Return the updated role as domain entity
	return s.roleMapper.ToDomain(modelRole), nil
}

// DeleteRole deletes a role with domain business rules.
func (s *RoleDomainService) DeleteRole(ctx context.Context, id uint64) error {
	if id == 0 {
		return ErrInvalidRole
	}

	// Get the role to check if it can be deleted
	modelRole, err := s.rolePort.GetRoleByID(ctx, id)
	if err != nil {
		return fmt.Errorf("failed to get role for deletion: %w", err)
	}
	if modelRole == nil {
		return constants.ErrRecordNotFound
	}

	// Apply domain business rules
	domainRole := s.roleMapper.ToDomain(modelRole)
	if !domainRole.CanBeDeleted() {
		return ErrCannotDeleteSystemRole
	}

	// Proceed with deletion
	err = s.rolePort.DeleteRole(ctx, id)
	if err != nil {
		return fmt.Errorf("failed to delete role: %w", err)
	}

	return nil
}

// GetPaginatedRoles retrieves paginated roles with domain conversion.
func (s *RoleDomainService) GetPaginatedRoles(ctx context.Context, sort string, order string, page int, pageSize int) ([]*domain.Role, int64, error) {

	if page < 0 {
		page = 0
	}
	if pageSize < 1 || pageSize > 100 {
		pageSize = 10
	}

	// Default sorting
	if sort == "" {
		sort = "name"
	}
	if order != "asc" && order != "desc" {
		order = "asc"
	}

	modelRoles, total, err := s.rolePort.GetPaginatedRoles(ctx, sort, order, page, pageSize)
	if err != nil {
		return nil, 0, fmt.Errorf("failed to get paginated roles: %w", err)
	}

	// Convert to domain entities
	domainRoles := s.roleMapper.ToDomainSlice(modelRoles)

	return domainRoles, total, nil
}

// ValidateRole performs comprehensive domain validation on a role.
func (s *RoleDomainService) ValidateRole(role *domain.Role) error {
	if role == nil {
		return ErrInvalidRole
	}

	if !role.IsValid() {
		return ErrInvalidRole
	}

	return nil
}

// GetSystemRoles returns all system-defined roles.
func (s *RoleDomainService) GetSystemRoles(ctx context.Context) ([]*domain.Role, error) {
	allRoles, _, err := s.GetPaginatedRoles(ctx, "name", "asc", 1, 100)
	if err != nil {
		return nil, err
	}

	var systemRoles []*domain.Role
	for _, role := range allRoles {
		if role.IsSystemRole() {
			systemRoles = append(systemRoles, role)
		}
	}

	return systemRoles, nil
}

// GetUserDefinedRoles returns all user-defined (non-system) roles.
func (s *RoleDomainService) GetUserDefinedRoles(ctx context.Context) ([]*domain.Role, error) {
	allRoles, _, err := s.GetPaginatedRoles(ctx, "name", "asc", 1, 100)
	if err != nil {
		return nil, err
	}

	var userRoles []*domain.Role
	for _, role := range allRoles {
		if !role.IsSystemRole() {
			userRoles = append(userRoles, role)
		}
	}

	return userRoles, nil
}
