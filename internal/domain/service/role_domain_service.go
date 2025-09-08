package service

import (
	"context"
	"errors"
	"fmt"

	"github.com/EdwinRincon/browersfc-api/api/constants"
	"github.com/EdwinRincon/browersfc-api/domain"
)

var (
	ErrInvalidRole            = errors.New("invalid role data")
	ErrCannotDeleteSystemRole = errors.New("cannot delete system role")
)

type RoleDomainService struct {
	roleRepository domain.RoleRepository
}

func NewRoleDomainService(roleRepository domain.RoleRepository) *RoleDomainService {
	return &RoleDomainService{
		roleRepository: roleRepository,
	}
}

// CreateRole creates a new role with domain validation.
func (s *RoleDomainService) CreateRole(ctx context.Context, role *domain.Role) (*domain.Role, error) {
	// Domain validation
	if !role.IsValid() {
		return nil, ErrInvalidRole
	}

	// Check if role already exists
	existingRole, err := s.roleRepository.GetRoleByName(ctx, role.Name)
	if err != nil {
		return nil, fmt.Errorf("failed to check existing role: %w", err)
	}
	if existingRole != nil {
		return nil, constants.ErrRecordAlreadyExists
	}

	// Create through repository (adapter handles conversion)
	err = s.roleRepository.CreateRole(ctx, role)
	if err != nil {
		return nil, fmt.Errorf("failed to create role: %w", err)
	}

	// Return the created role as domain entity
	return role, nil
}

// GetRoleByID retrieves a role by its ID with domain logic.
func (s *RoleDomainService) GetRoleByID(ctx context.Context, id uint64) (*domain.Role, error) {
	if id == 0 {
		return nil, ErrInvalidRole
	}

	domainRole, err := s.roleRepository.GetRoleByID(ctx, id)
	if err != nil {
		return nil, fmt.Errorf("failed to get role: %w", err)
	}
	if domainRole == nil {
		return nil, constants.ErrRecordNotFound
	}

	return domainRole, nil
}

// GetRoleByName retrieves a role by its name with domain logic.
func (s *RoleDomainService) GetRoleByName(ctx context.Context, name string) (*domain.Role, error) {
	if name == "" {
		return nil, ErrInvalidRole
	}

	domainRole, err := s.roleRepository.GetRoleByName(ctx, name)
	if err != nil {
		return nil, fmt.Errorf("failed to get role by name: %w", err)
	}
	if domainRole == nil {
		return nil, constants.ErrRecordNotFound
	}

	return domainRole, nil
}

// UpdateRole updates an existing role with domain validation.
func (s *RoleDomainService) UpdateRole(ctx context.Context, role *domain.Role) (*domain.Role, error) {
	// Domain validation
	if !role.IsValid() {
		return nil, ErrInvalidRole
	}

	// Check if role exists
	existingRole, err := s.roleRepository.GetRoleByID(ctx, role.ID)
	if err != nil {
		return nil, fmt.Errorf("failed to check existing role: %w", err)
	}
	if existingRole == nil {
		return nil, constants.ErrRecordNotFound
	}

	// Check if name is being changed to an existing name
	if existingRole.Name != role.Name {
		conflictingRole, err := s.roleRepository.GetRoleByName(ctx, role.Name)
		if err != nil {
			return nil, fmt.Errorf("failed to check name conflict: %w", err)
		}
		if conflictingRole != nil && conflictingRole.ID != role.ID {
			return nil, constants.ErrRecordAlreadyExists
		}
	}

	// Update through repository
	err = s.roleRepository.UpdateRole(ctx, role)
	if err != nil {
		return nil, fmt.Errorf("failed to update role: %w", err)
	}

	updatedRole, err := s.roleRepository.GetRoleByID(ctx, role.ID)
	if err != nil {
		return nil, fmt.Errorf("failed to fetch updated role: %w", err)
	}
	if updatedRole == nil {
		return nil, constants.ErrRecordNotFound
	}

	return updatedRole, nil
}

// DeleteRole deletes a role with domain business rules.
func (s *RoleDomainService) DeleteRole(ctx context.Context, id uint64) error {
	if id == 0 {
		return ErrInvalidRole
	}

	// Get role to check business rules
	domainRole, err := s.roleRepository.GetRoleByID(ctx, id)
	if err != nil {
		return fmt.Errorf("failed to get role for deletion: %w", err)
	}
	if domainRole == nil {
		return constants.ErrRecordNotFound
	}

	// Apply domain business rules
	if !domainRole.CanBeDeleted() {
		return ErrCannotDeleteSystemRole
	}

	// Proceed with deletion
	err = s.roleRepository.DeleteRole(ctx, id)
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

	domainRoles, total, err := s.roleRepository.GetPaginatedRoles(ctx, sort, order, page, pageSize)
	if err != nil {
		return nil, 0, fmt.Errorf("failed to get paginated roles: %w", err)
	}

	// Convert slice to slice of pointers
	domainRolePointers := make([]*domain.Role, len(domainRoles))
	for i := range domainRoles {
		domainRolePointers[i] = &domainRoles[i]
	}

	return domainRolePointers, total, nil
}

// ValidateRole performs comprehensive domain validation on a role.
func (s *RoleDomainService) ValidateRole(role *domain.Role) error {
	if role == nil {
		return ErrInvalidRole
	}

	// Use domain entity validation
	if !role.IsValid() {
		return ErrInvalidRole
	}

	return nil
}

// GetSystemRoles returns the default system roles.
func (s *RoleDomainService) GetSystemRoles() []*domain.Role {
	return []*domain.Role{
		{Name: domain.RoleAdmin, Description: "Administrator with full access"},
		{Name: domain.RolePlayer, Description: "Player with limited access"},
		{Name: domain.RoleCoach, Description: "Coach with team management access"},
	}
}

// IsSystemRole checks if a role is a system role.
func (s *RoleDomainService) IsSystemRole(roleName string) bool {
	systemRoles := s.GetSystemRoles()
	for _, role := range systemRoles {
		if role.Name == roleName {
			return true
		}
	}
	return false
}
