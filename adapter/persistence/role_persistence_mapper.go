package persistence

import (
	"github.com/EdwinRincon/browersfc-api/domain"
	"github.com/EdwinRincon/browersfc-api/internal/infrastructure/persistence/model"
)

// RolePersistenceMapper handles persistence layer conversions for Role entity
type RolePersistenceMapper struct{}

func NewRolePersistenceMapper() *RolePersistenceMapper {
	return &RolePersistenceMapper{}
}

// DomainToModel converts a domain.Role to model.Role for persistence
func (m *RolePersistenceMapper) DomainToModel(domainRole *domain.Role) *model.Role {
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
func (m *RolePersistenceMapper) ModelToDomain(modelRole *model.Role) *domain.Role {
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

// ModelListToDomain converts a slice of model.Role to domain.Role for business logic
func (m *RolePersistenceMapper) ModelListToDomain(roles []model.Role) []domain.Role {
	result := make([]domain.Role, len(roles))
	for i, role := range roles {
		if domainRole := m.ModelToDomain(&role); domainRole != nil {
			result[i] = *domainRole
		}
	}
	return result
}
