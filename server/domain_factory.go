package server

import (
	"github.com/EdwinRincon/browersfc-api/domain"
	domainservice "github.com/EdwinRincon/browersfc-api/internal/domain/service"
	"github.com/EdwinRincon/browersfc-api/internal/infrastructure/persistence"
)

// CreateRoleDomainService creates a role domain service with repository implementing domain interface
func CreateRoleDomainService(roleRepo *persistence.RoleRepositoryImpl) *domainservice.RoleDomainService {
	// Repository implements domain.RoleRepository interface directly
	var roleRepository domain.RoleRepository = roleRepo
	return domainservice.NewRoleDomainService(roleRepository)
}
