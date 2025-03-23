package service

import (
	"context"

	"github.com/EdwinRincon/browersfc-api/api/model"
	"github.com/EdwinRincon/browersfc-api/api/repository"
)

// RoleService es la interfaz que define los métodos relacionados con roles
type RoleService interface {
	GetRoleByID(ctx context.Context, id uint8) (*model.Role, error)
	CreateRole(ctx context.Context, role *model.Role) error
	UpdateRole(ctx context.Context, role *model.Role) error
	DeleteRole(ctx context.Context, id uint8) error
	GetAllRoles(ctx context.Context) ([]model.Role, error)
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
func (s *roleService) GetRoleByID(ctx context.Context, id uint8) (*model.Role, error) {
	return s.RoleRepository.GetRoleByID(ctx, id)
}

func (s *roleService) CreateRole(ctx context.Context, role *model.Role) error {
	return s.RoleRepository.CreateRole(ctx, role)
}

func (s *roleService) UpdateRole(ctx context.Context, role *model.Role) error {
	return s.RoleRepository.UpdateRole(ctx, role)
}

func (s *roleService) DeleteRole(ctx context.Context, id uint8) error {
	return s.RoleRepository.DeleteRole(ctx, id)
}

func (s *roleService) GetAllRoles(ctx context.Context) ([]model.Role, error) {
	return s.RoleRepository.GetAllRoles(ctx)
}
