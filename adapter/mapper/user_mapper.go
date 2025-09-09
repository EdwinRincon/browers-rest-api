package mapper

import (
	"github.com/EdwinRincon/browersfc-api/api/dto"
	"github.com/EdwinRincon/browersfc-api/internal/infrastructure/persistence/model"
	"github.com/EdwinRincon/browersfc-api/domain"
)

// UserMapper handles all user conversions in one place to avoid duplication
// Follows hexagonal architecture principles with clear separation
type UserMapper struct{}

func NewUserMapper() *UserMapper {
	return &UserMapper{}
}

// ========================================================================================
// DTO to Domain Conversions (Used by HTTP handlers)
// ========================================================================================

// DTOToDomain converts a CreateUserRequest DTO to a domain.User entity
func (m *UserMapper) DTOToDomain(userDTO *dto.CreateUserRequest, roleID uint64) *domain.User {
	if userDTO == nil {
		return nil
	}
	return &domain.User{
		Name:     userDTO.Name,
		LastName: userDTO.LastName,
		Username: userDTO.Username,
		RoleID:   roleID,
	}
}

// UpdateDTOToDomain converts an UpdateUserRequest DTO to a domain.User entity
func (m *UserMapper) UpdateDTOToDomain(userDTO *dto.UpdateUserRequest) *domain.User {
	if userDTO == nil {
		return nil
	}
	domainUser := &domain.User{}

	if userDTO.Name != nil {
		domainUser.Name = *userDTO.Name
	}
	if userDTO.LastName != nil {
		domainUser.LastName = *userDTO.LastName
	}
	if userDTO.Username != nil {
		domainUser.Username = *userDTO.Username
	}
	if userDTO.Birthdate != nil {
		domainUser.Birthdate = userDTO.Birthdate
	}
	if userDTO.ImgProfile != nil {
		domainUser.ImgProfile = *userDTO.ImgProfile
	}
	if userDTO.ImgBanner != nil {
		domainUser.ImgBanner = *userDTO.ImgBanner
	}
	if userDTO.RoleID != nil {
		domainUser.RoleID = *userDTO.RoleID
	}

	return domainUser
}

// DomainToDTO converts a domain.User to UserResponse DTO
func (m *UserMapper) DomainToDTO(user *domain.User, role *dto.RoleShort) *dto.UserResponse {
	if user == nil {
		return nil
	}
	response := &dto.UserResponse{
		ID:         user.ID,
		Name:       user.Name,
		LastName:   user.LastName,
		Username:   user.Username,
		ImgProfile: user.ImgProfile,
		ImgBanner:  user.ImgBanner,
		CreatedAt:  user.CreatedAt,
		UpdatedAt:  user.UpdatedAt,
	}

	// Only set the birthdate if it's not nil
	if user.Birthdate != nil {
		response.Birthdate = *user.Birthdate
	}

	// Set role if provided
	if role != nil {
		response.Role = *role
	}

	return response
}

// DomainToShortDTO converts a domain.User to UserShort DTO
func (m *UserMapper) DomainToShortDTO(user *domain.User) *dto.UserShort {
	if user == nil {
		return nil
	}
	return &dto.UserShort{
		ID:       user.ID,
		Username: user.Username,
	}
}

// DomainListToDTO converts a slice of domain.User to UserResponse DTOs
func (m *UserMapper) DomainListToDTO(users []domain.User) []dto.UserResponse {
	userResponses := make([]dto.UserResponse, len(users))
	for i, user := range users {
		if response := m.DomainToDTO(&user, nil); response != nil {
			userResponses[i] = *response
		}
	}
	return userResponses
}

// ========================================================================================
// Domain to Persistence Model Conversions (Used by infrastructure repositories)
// ========================================================================================

// DomainToModel converts a domain.User to model.User for persistence
func (m *UserMapper) DomainToModel(domainUser *domain.User) *model.User {
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
func (m *UserMapper) ModelToDomain(modelUser *model.User) *domain.User {
	if modelUser == nil {
		return nil
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
		CreatedAt:  modelUser.CreatedAt,
		UpdatedAt:  modelUser.UpdatedAt,
	}
}

// ModelListToDomain converts a slice of model.User to domain.User for business logic
func (m *UserMapper) ModelListToDomain(users []model.User) []domain.User {
	result := make([]domain.User, len(users))
	for i, user := range users {
		if domainUser := m.ModelToDomain(&user); domainUser != nil {
			result[i] = *domainUser
		}
	}
	return result
}
