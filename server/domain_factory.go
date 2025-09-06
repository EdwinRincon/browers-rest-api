package server

import (
	"github.com/EdwinRincon/browersfc-api/api/mapper"
	"github.com/EdwinRincon/browersfc-api/internal/adapters"
	domainservice "github.com/EdwinRincon/browersfc-api/internal/domain/service"
	"github.com/EdwinRincon/browersfc-api/internal/infrastructure/persistence"
)

// CreateRoleDomainService creates a role domain service with all dependencies
func CreateRoleDomainService(roleRepo persistence.RoleRepository) *domainservice.RoleDomainService {
	roleAdapter := adapters.NewRoleDomainAdapter(roleRepo)
	roleMapper := mapper.NewRoleDomainMapper()
	return domainservice.NewRoleDomainService(roleAdapter, roleMapper)
}
