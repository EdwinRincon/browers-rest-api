package http

import (
	"github.com/EdwinRincon/browersfc-api/api/dto"
	"github.com/EdwinRincon/browersfc-api/domain"
)

// UserHTTPMapper handles HTTP layer conversions for User entity
type UserHTTPMapper struct{}

func NewUserHTTPMapper() *UserHTTPMapper {
	return &UserHTTPMapper{}
}

// ========================================================================================
// DTO to Domain Conversions (Used by HTTP handlers)
// ========================================================================================

// DTOToDomain converts a CreateUserRequest DTO to a domain.User entity
func (m *UserHTTPMapper) DTOToDomain(userDTO *dto.CreateUserRequest, roleID uint64) *domain.User {
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
func (m *UserHTTPMapper) UpdateDTOToDomain(userDTO *dto.UpdateUserRequest) *domain.User {
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
func (m *UserHTTPMapper) DomainToDTO(user *domain.User, role *dto.RoleShort) *dto.UserResponse {
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
func (m *UserHTTPMapper) DomainToShortDTO(user *domain.User) *dto.UserShort {
	if user == nil {
		return nil
	}
	return &dto.UserShort{
		ID:       user.ID,
		Username: user.Username,
	}
}

// DomainListToDTO converts a slice of domain.User to UserResponse DTOs
func (m *UserHTTPMapper) DomainListToDTO(users []domain.User) []dto.UserResponse {
	userResponses := make([]dto.UserResponse, len(users))
	for i, user := range users {
		if response := m.DomainToDTO(&user, nil); response != nil {
			userResponses[i] = *response
		}
	}
	return userResponses
}
