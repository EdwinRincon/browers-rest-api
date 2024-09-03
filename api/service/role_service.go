package service

import (
	"context"

	"github.com/EdwinRincon/browersfc-api/api/model"
	"github.com/EdwinRincon/browersfc-api/api/repository"
)

// RoleService es la interfaz que define los métodos relacionados con roles
type RoleService interface {
	GetRoleByID(ctx context.Context, id uint8) (*model.Roles, error)
}

// roleService es la implementación concreta de la interfaz RoleService
type roleService struct {
	RoleRepository repository.RoleRepository
}

// NewRoleService crea una nueva instancia de RoleService
func NewRoleService(roleRepo repository.RoleRepository) RoleService {
	return &roleService{
		RoleRepository: roleRepo,
	}
}

// GetRoleByID obtiene un rol por su ID
func (s *roleService) GetRoleByID(ctx context.Context, id uint8) (*model.Roles, error) {
	return s.RoleRepository.GetRoleByID(ctx, id)
}
