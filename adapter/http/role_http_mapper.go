package http

import (
	"github.com/EdwinRincon/browersfc-api/api/dto"
	"github.com/EdwinRincon/browersfc-api/domain"
)

// RoleHTTPMapper handles HTTP layer conversions for Role entity
type RoleHTTPMapper struct{}

func NewRoleHTTPMapper() *RoleHTTPMapper {
	return &RoleHTTPMapper{}
}

// DTOToDomain converts a CreateRoleRequest DTO to a domain.Role entity
func (m *RoleHTTPMapper) DTOToDomain(roleDTO *dto.CreateRoleRequest) *domain.Role {
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
func (m *RoleHTTPMapper) UpdateDTOToDomain(roleDTO *dto.UpdateRoleRequest) *domain.Role {
	if roleDTO == nil {
		return nil
	}

	domainRole := &domain.Role{}

	if roleDTO.Name != nil {
		domainRole.Name = *roleDTO.Name
	}
	if roleDTO.Description != nil {
		domainRole.Description = *roleDTO.Description
	}

	return domainRole
}

// DomainToDTO converts a domain.Role to RoleResponse DTO
func (m *RoleHTTPMapper) DomainToDTO(role *domain.Role) *dto.RoleResponse {
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
func (m *RoleHTTPMapper) DomainToShortDTO(role *domain.Role) *dto.RoleShort {
	if role == nil {
		return nil
	}

	return &dto.RoleShort{
		ID:   role.ID,
		Name: role.Name,
	}
}

// DomainListToDTO converts a slice of domain.Role to RoleResponse DTOs
func (m *RoleHTTPMapper) DomainListToDTO(roles []domain.Role) []dto.RoleResponse {
	roleResponses := make([]dto.RoleResponse, len(roles))
	for i, role := range roles {
		if response := m.DomainToDTO(&role); response != nil {
			roleResponses[i] = *response
		}
	}
	return roleResponses
}
