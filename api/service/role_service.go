package service

import (
	"context"
	"fmt"

	"github.com/EdwinRincon/browersfc-api/api/constants"
	"github.com/EdwinRincon/browersfc-api/api/model"
	"github.com/EdwinRincon/browersfc-api/api/repository"
)

type RoleService interface {
	GetRoleByID(ctx context.Context, id uint64) (*model.Role, error)
	GetRoleByName(ctx context.Context, name string) (*model.Role, error)
	CreateRole(ctx context.Context, role *model.Role) (*model.Role, error)
	UpdateRole(ctx context.Context, id uint64, updated *model.Role) error
	DeleteRole(ctx context.Context, id uint64) error
	GetPaginatedRoles(ctx context.Context, sort string, order string, page int, pageSize int) ([]model.Role, int64, error)
}

type roleService struct {
	RoleRepository repository.RoleRepository
}

func NewRoleService(roleRepo repository.RoleRepository) RoleService {
	return &roleService{
		RoleRepository: roleRepo,
	}
}

func (s *roleService) GetRoleByID(ctx context.Context, id uint64) (*model.Role, error) {
	role, err := s.RoleRepository.GetRoleByID(ctx, id)
	if err != nil {
		return nil, err
	}
	if role == nil {
		return nil, constants.ErrRecordNotFound
	}
	return role, nil
}

func (s *roleService) GetRoleByName(ctx context.Context, name string) (*model.Role, error) {
	role, err := s.RoleRepository.GetRoleByName(ctx, name)
	if err != nil {
		return nil, err
	}
	if role == nil {
		return nil, constants.ErrRecordNotFound
	}
	return role, nil
}

func (s *roleService) CreateRole(ctx context.Context, role *model.Role) (*model.Role, error) {
	existing, err := s.RoleRepository.GetRoleByName(ctx, role.Name)
	if err != nil {
		return nil, fmt.Errorf("failed to check existing role: %w", err)
	}

	if existing != nil {
		return nil, constants.ErrRecordAlreadyExists
	}

	err = s.RoleRepository.CreateRole(ctx, role)
	if err != nil {
		return nil, fmt.Errorf("failed to create role: %w", err)
	}
	return role, nil
}

func (s *roleService) UpdateRole(ctx context.Context, id uint64, updated *model.Role) error {
	existing, err := s.RoleRepository.GetRoleByID(ctx, id)
	if err != nil {
		return err
	}
	if existing == nil {
		return constants.ErrRecordNotFound
	}

	if existing.Name != updated.Name {
		duplicate, err := s.RoleRepository.GetRoleByName(ctx, updated.Name)
		if err != nil {
			return fmt.Errorf("failed to check duplicate role: %w", err)
		}
		if duplicate != nil && duplicate.ID != existing.ID {
			return constants.ErrRecordAlreadyExists
		}
	}

	existing.Name = updated.Name
	existing.Description = updated.Description

	return s.RoleRepository.UpdateRole(ctx, existing)
}

func (s *roleService) DeleteRole(ctx context.Context, id uint64) error {
	return s.RoleRepository.DeleteRole(ctx, id)
}

func (s *roleService) GetPaginatedRoles(ctx context.Context, sort, order string, page, pageSize int) ([]model.Role, int64, error) {
	return s.RoleRepository.GetPaginatedRoles(ctx, sort, order, page, pageSize)
}
