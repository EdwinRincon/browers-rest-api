package mapper

import (
	"github.com/EdwinRincon/browersfc-api/api/model"
	"github.com/EdwinRincon/browersfc-api/domain"
)

// ToDomain converts a persistence model.User to a domain.User
func UserToDomain(user *model.User) *domain.User {
	if user == nil {
		return nil
	}

	return &domain.User{
		ID:         user.ID,
		Name:       user.Name,
		LastName:   user.LastName,
		Username:   user.Username,
		Birthdate:  user.Birthdate,
		ImgProfile: user.ImgProfile,
		ImgBanner:  user.ImgBanner,
		RoleID:     user.RoleID,
		CreatedAt:  user.CreatedAt,
		UpdatedAt:  user.UpdatedAt,
	}
}

// ToModel converts a domain.User to a persistence model.User
func UserToModel(user *domain.User) *model.User {
	if user == nil {
		return nil
	}

	return &model.User{
		ID:         user.ID,
		Name:       user.Name,
		LastName:   user.LastName,
		Username:   user.Username,
		Birthdate:  user.Birthdate,
		ImgProfile: user.ImgProfile,
		ImgBanner:  user.ImgBanner,
		RoleID:     user.RoleID,
		CreatedAt:  user.CreatedAt,
		UpdatedAt:  user.UpdatedAt,
	}
}

// UserListToDomain converts a slice of persistence model.User to a slice of domain.User
func UserListToDomain(users []model.User) []domain.User {
	result := make([]domain.User, len(users))
	for i, user := range users {
		domainUser := UserToDomain(&user)
		if domainUser != nil {
			result[i] = *domainUser
		}
	}
	return result
}
