package mapper

import (
	"github.com/EdwinRincon/browersfc-api/adapter/mapper"
	"github.com/EdwinRincon/browersfc-api/api/dto"
	"github.com/EdwinRincon/browersfc-api/api/model"
	"github.com/EdwinRincon/browersfc-api/domain"
)

func UpdateUserFromDTO(user *model.User, dto *dto.UpdateUserRequest) {
	if dto.Name != nil {
		user.Name = *dto.Name
	}
	if dto.LastName != nil {
		user.LastName = *dto.LastName
	}
	if dto.Username != nil {
		user.Username = *dto.Username
	}
	if dto.Birthdate != nil {
		*user.Birthdate = *dto.Birthdate
	}
	if dto.ImgProfile != nil {
		user.ImgProfile = *dto.ImgProfile
	}
	if dto.ImgBanner != nil {
		user.ImgBanner = *dto.ImgBanner
	}
	if dto.RoleID != nil {
		user.RoleID = *dto.RoleID
	}
}

// Domain to DTO mappings

func DomainUserToUserResponse(user *domain.User, role *dto.RoleShort) *dto.UserResponse {
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

func DomainUserToUserShort(user *domain.User) *dto.UserShort {
	return &dto.UserShort{
		ID:       user.ID,
		Username: user.Username,
	}
}

func DomainUserListToUserResponseList(users []domain.User) []dto.UserResponse {
	userResponses := make([]dto.UserResponse, len(users))
	for i, user := range users {
		userResponses[i] = *DomainUserToUserResponse(&user, nil)
	}
	return userResponses
}

func CreateUserRequestToDomain(dto *dto.CreateUserRequest, roleID uint64) *domain.User {
	return &domain.User{
		Name:     dto.Name,
		LastName: dto.LastName,
		Username: dto.Username,
		RoleID:   roleID,
	}
}

func UpdateUserRequestToDomain(dto *dto.UpdateUserRequest) *domain.User {
	domainUser := &domain.User{}
	
	if dto.Name != nil {
		domainUser.Name = *dto.Name
	}
	if dto.LastName != nil {
		domainUser.LastName = *dto.LastName
	}
	if dto.Username != nil {
		domainUser.Username = *dto.Username
	}
	if dto.Birthdate != nil {
		domainUser.Birthdate = dto.Birthdate
	}
	if dto.ImgProfile != nil {
		domainUser.ImgProfile = *dto.ImgProfile
	}
	if dto.ImgBanner != nil {
		domainUser.ImgBanner = *dto.ImgBanner
	}
	if dto.RoleID != nil {
		domainUser.RoleID = *dto.RoleID
	}
	
	return domainUser
}

// Legacy model to DTO mappings (keep for compatibility)

func ToUserResponse(user *model.User) *dto.UserResponse {
	var role dto.RoleShort
	if user.Role != nil {
		roleMapper := mapper.NewRoleMapper()
		roleShort := roleMapper.ModelToShortDTO(user.Role)
		if roleShort != nil {
			role = *roleShort
		}
	}

	response := &dto.UserResponse{
		ID:         user.ID,
		Name:       user.Name,
		LastName:   user.LastName,
		Username:   user.Username,
		ImgProfile: user.ImgProfile,
		ImgBanner:  user.ImgBanner,
		Role:       role,
		CreatedAt:  user.CreatedAt,
		UpdatedAt:  user.UpdatedAt,
	}

	// Only set the birthdate if it's not nil
	if user.Birthdate != nil {
		response.Birthdate = *user.Birthdate
	}

	return response
}

func ToUserResponseList(users []model.User) []dto.UserResponse {
	userResponses := make([]dto.UserResponse, len(users))
	for i, user := range users {
		userResponses[i] = *ToUserResponse(&user)
	}
	return userResponses
}

func ToUserShort(user *model.User) *dto.UserShort {
	return &dto.UserShort{
		ID:       user.ID,
		Username: user.Username,
	}
}

func ToUser(dto *dto.CreateUserRequest, roleID uint64) *model.User {
	return &model.User{
		Name:     dto.Name,
		LastName: dto.LastName,
		Username: dto.Username,
		RoleID:   roleID,
	}
}
