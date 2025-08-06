package mapper

import (
	"github.com/EdwinRincon/browersfc-api/api/dto"
	"github.com/EdwinRincon/browersfc-api/api/model"
)

// ToRole converts a CreateRoleRequest DTO to a Role model
func ToRole(roleDTO *dto.CreateRoleRequest) *model.Role {
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
		roleResponses[i] = *ToRoleResponse(&role)
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
