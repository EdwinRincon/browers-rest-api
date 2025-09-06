package mapper

import (
	"github.com/EdwinRincon/browersfc-api/api/model"
	"github.com/EdwinRincon/browersfc-api/domain"
)

// RoleDomainMapper handles conversions between domain.Role and model.Role
type RoleDomainMapper struct{}

// NewRoleDomainMapper creates a new RoleDomainMapper instance
func NewRoleDomainMapper() *RoleDomainMapper {
	return &RoleDomainMapper{}
}

// ToModel converts a domain.Role to model.Role for persistence
func (m *RoleDomainMapper) ToModel(domainRole *domain.Role) *model.Role {
	if domainRole == nil {
		return nil
	}

	return &model.Role{
		ID:          domainRole.ID,
		Name:        domainRole.Name,
		Description: domainRole.Description,
		CreatedAt:   domainRole.CreatedAt,
		UpdatedAt:   domainRole.UpdatedAt,
	}
}

// ToDomain converts a model.Role to domain.Role for business logic
func (m *RoleDomainMapper) ToDomain(modelRole *model.Role) *domain.Role {
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

// ToModelSlice converts a slice of domain.Role to model.Role
func (m *RoleDomainMapper) ToModelSlice(domainRoles []*domain.Role) []*model.Role {
	if domainRoles == nil {
		return nil
	}

	modelRoles := make([]*model.Role, len(domainRoles))
	for i, role := range domainRoles {
		modelRoles[i] = m.ToModel(role)
	}
	return modelRoles
}

// ToDomainSlice converts a slice of model.Role to domain.Role
func (m *RoleDomainMapper) ToDomainSlice(modelRoles []model.Role) []*domain.Role {
	if modelRoles == nil {
		return nil
	}

	domainRoles := make([]*domain.Role, len(modelRoles))
	for i, role := range modelRoles {
		domainRoles[i] = m.ToDomain(&role)
	}
	return domainRoles
}
