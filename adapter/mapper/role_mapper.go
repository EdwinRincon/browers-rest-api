package mapper

import (
	"github.com/EdwinRincon/browersfc-api/api/dto"
	"github.com/EdwinRincon/browersfc-api/api/model"
	"github.com/EdwinRincon/browersfc-api/domain"
)

// RoleMapper handles all role conversions in one place to avoid duplication
// Follows hexagonal architecture principles with clear separation
type RoleMapper struct{}

// NewRoleMapper creates a new RoleMapper instance
func NewRoleMapper() *RoleMapper {
	return &RoleMapper{}
}

// ========================================================================================
// DTO to Domain Conversions (Used by HTTP handlers)
// ========================================================================================

// DTOToDomain converts a CreateRoleRequest DTO to a domain.Role entity
func (m *RoleMapper) DTOToDomain(roleDTO *dto.CreateRoleRequest) *domain.Role {
	if roleDTO == nil {
		return nil
	}
	description := ""
	if roleDTO.Description != nil {
		description = *roleDTO.Description
	}

	return &domain.Role{
		Name:        roleDTO.Name,
		Description: description,
	}
}

// UpdateDTOToDomain converts an UpdateRoleRequest DTO to a domain.Role entity
func (m *RoleMapper) UpdateDTOToDomain(roleDTO *dto.UpdateRoleRequest) *domain.Role {
	if roleDTO == nil {
		return nil
	}
	role := &domain.Role{}
	if roleDTO.Name != nil {
		role.Name = *roleDTO.Name
	}
	if roleDTO.Description != nil {
		role.Description = *roleDTO.Description
	}
	return role
}

// DomainToDTO converts a domain.Role to RoleResponse DTO
func (m *RoleMapper) DomainToDTO(role *domain.Role) *dto.RoleResponse {
	if role == nil {
		return nil
	}
	return &dto.RoleResponse{
		ID:          role.ID,
		Name:        role.Name,
		Description: role.Description,
		CreatedAt:   role.CreatedAt,
		UpdatedAt:   role.UpdatedAt,
	}
}

// DomainToShortDTO converts a domain.Role to RoleShort DTO
func (m *RoleMapper) DomainToShortDTO(role *domain.Role) *dto.RoleShort {
	if role == nil {
		return nil
	}
	return &dto.RoleShort{
		ID:   role.ID,
		Name: role.Name,
	}
}

// DomainListToDTO converts a slice of domain.Role to RoleResponse DTOs
func (m *RoleMapper) DomainListToDTO(roles []*domain.Role) []dto.RoleResponse {
	roleResponses := make([]dto.RoleResponse, len(roles))
	for i, role := range roles {
		if response := m.DomainToDTO(role); response != nil {
			roleResponses[i] = *response
		}
	}
	return roleResponses
}

// ========================================================================================
// Domain to Persistence Model Conversions (Used by infrastructure repositories)
// ========================================================================================

// DomainToModel converts a domain.Role to model.Role for persistence
func (m *RoleMapper) DomainToModel(domainRole *domain.Role) *model.Role {
	if domainRole == nil {
		return nil
	}

	modelRole := &model.Role{
		ID:          domainRole.ID,
		Name:        domainRole.Name,
		Description: domainRole.Description,
	}

	// Only set timestamps if they have meaningful values (not zero time)
	if !domainRole.CreatedAt.IsZero() {
		modelRole.CreatedAt = domainRole.CreatedAt
	}
	if !domainRole.UpdatedAt.IsZero() {
		modelRole.UpdatedAt = domainRole.UpdatedAt
	}

	return modelRole
}

// ModelToDomain converts a model.Role to domain.Role for business logic
func (m *RoleMapper) ModelToDomain(modelRole *model.Role) *domain.Role {
	if modelRole == nil {
		return nil
	}

	return &domain.Role{
		ID:          modelRole.ID,
		Name:        modelRole.Name,
		Description: modelRole.Description,
		CreatedAt:   modelRole.CreatedAt,
		UpdatedAt:   modelRole.UpdatedAt,
	}
}

// ========================================================================================
// Legacy Support Functions (For backward compatibility with other entities)
// ========================================================================================

// ModelToShortDTO converts a persistence model.Role to a RoleShort DTO
// This is kept for backward compatibility with user mapper and other entities
func (m *RoleMapper) ModelToShortDTO(role *model.Role) *dto.RoleShort {
	if role == nil {
		return nil
	}
	return &dto.RoleShort{
		ID:   role.ID,
		Name: role.Name,
	}
}
