package persistence

import (
	"github.com/EdwinRincon/browersfc-api/domain"
	"github.com/EdwinRincon/browersfc-api/internal/infrastructure/persistence/model"
)

// UserPersistenceMapper handles persistence layer conversions for User entity
type UserPersistenceMapper struct{}

func NewUserPersistenceMapper() *UserPersistenceMapper {
	return &UserPersistenceMapper{}
}

// ========================================================================================
// Domain to Persistence Model Conversions (Used by infrastructure repositories)
// ========================================================================================

// DomainToModel converts a domain.User to model.User for persistence
func (m *UserPersistenceMapper) DomainToModel(domainUser *domain.User) *model.User {
	if domainUser == nil {
		return nil
	}

	modelUser := &model.User{
		ID:         domainUser.ID,
		Name:       domainUser.Name,
		LastName:   domainUser.LastName,
		Username:   domainUser.Username,
		Birthdate:  domainUser.Birthdate,
		ImgProfile: domainUser.ImgProfile,
		ImgBanner:  domainUser.ImgBanner,
		RoleID:     domainUser.RoleID,
	}

	// Only set timestamps if they have meaningful values (not zero time)
	if !domainUser.CreatedAt.IsZero() {
		modelUser.CreatedAt = domainUser.CreatedAt
	}
	if !domainUser.UpdatedAt.IsZero() {
		modelUser.UpdatedAt = domainUser.UpdatedAt
	}

	return modelUser
}

// ModelToDomain converts a model.User to domain.User for business logic
func (m *UserPersistenceMapper) ModelToDomain(modelUser *model.User) *domain.User {
	if modelUser == nil {
		return nil
	}

	var role *domain.Role
	if modelUser.Role != nil {
		roleMapper := NewRolePersistenceMapper()
		role = roleMapper.ModelToDomain(modelUser.Role)
	}

	return &domain.User{
		ID:         modelUser.ID,
		Name:       modelUser.Name,
		LastName:   modelUser.LastName,
		Username:   modelUser.Username,
		Birthdate:  modelUser.Birthdate,
		ImgProfile: modelUser.ImgProfile,
		ImgBanner:  modelUser.ImgBanner,
		RoleID:     modelUser.RoleID,
		Role:       role,
		CreatedAt:  modelUser.CreatedAt,
		UpdatedAt:  modelUser.UpdatedAt,
	}
}

// ModelListToDomain converts a slice of model.User to domain.User for business logic
func (m *UserPersistenceMapper) ModelListToDomain(users []model.User) []domain.User {
	result := make([]domain.User, len(users))
	for i, user := range users {
		if domainUser := m.ModelToDomain(&user); domainUser != nil {
			result[i] = *domainUser
		}
	}
	return result
}
