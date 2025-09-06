package mapper

import (
	"github.com/EdwinRincon/browersfc-api/api/dto"
	"github.com/EdwinRincon/browersfc-api/api/model"
	"github.com/EdwinRincon/browersfc-api/domain"
)

// ToRole converts a CreateRoleRequest DTO to a Role model
func ToRole(roleDTO *dto.CreateRoleRequest) *model.Role {
	if roleDTO == nil {
		return nil
	}
	description := ""
	if roleDTO.Description != nil {
		description = *roleDTO.Description
	}

	return &model.Role{
		Name:        roleDTO.Name,
		Description: description,
	}
}

// ToRoleResponse converts a Role model to a RoleResponse DTO
func ToRoleResponse(role *model.Role) *dto.RoleResponse {
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

// ToRoleResponseList converts a slice of Role models to a slice of RoleResponse DTOs
func ToRoleResponseList(roles []model.Role) []dto.RoleResponse {
	roleResponses := make([]dto.RoleResponse, len(roles))
	for i, role := range roles {
		if response := ToRoleResponse(&role); response != nil {
			roleResponses[i] = *response
		}
	}
	return roleResponses
}

// ToRoleFromUpdate converts an UpdateRoleRequest DTO to a Role model
func ToRoleFromUpdate(roleDTO *dto.UpdateRoleRequest) *model.Role {
	if roleDTO == nil {
		return nil
	}
	role := &model.Role{}
	if roleDTO.Name != nil {
		role.Name = *roleDTO.Name
	}
	if roleDTO.Description != nil {
		role.Description = *roleDTO.Description
	}
	return role
}

// ToRoleShort converts a Role model to a RoleShort DTO
func ToRoleShort(role *model.Role) *dto.RoleShort {
	if role == nil {
		return nil
	}
	return &dto.RoleShort{
		ID:   role.ID,
		Name: role.Name,
	}
}

// ==============================================================================
// Direct Domain-to-DTO Mappers (Hexagonal Architecture - Clean Boundaries)
// ==============================================================================

// DomainToRoleResponse converts a domain.Role directly to RoleResponse DTO
// This avoids coupling the HTTP layer to persistence models
func DomainToRoleResponse(role *domain.Role) *dto.RoleResponse {
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

// DomainToRoleResponseList converts a slice of domain.Role to RoleResponse DTOs
func DomainToRoleResponseList(roles []*domain.Role) []dto.RoleResponse {
	roleResponses := make([]dto.RoleResponse, len(roles))
	for i, role := range roles {
		if response := DomainToRoleResponse(role); response != nil {
			roleResponses[i] = *response
		}
	}
	return roleResponses
}

// DomainToRoleShort converts a domain.Role directly to RoleShort DTO
func DomainToRoleShort(role *domain.Role) *dto.RoleShort {
	if role == nil {
		return nil
	}
	return &dto.RoleShort{
		ID:   role.ID,
		Name: role.Name,
	}
}
