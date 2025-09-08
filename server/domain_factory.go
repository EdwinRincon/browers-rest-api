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

// CreateSeasonDomainService creates a season domain service with repository implementing domain interface
func CreateSeasonDomainService(seasonRepo *persistence.SeasonRepositoryImpl) *domainservice.SeasonDomainService {
	// Repository implements domain.SeasonRepository interface directly
	var seasonRepository domain.SeasonRepository = seasonRepo
	return domainservice.NewSeasonDomainService(seasonRepository)
}

// CreateUserDomainService creates a user domain service with repository implementing domain interface
func CreateUserDomainService(userRepo *persistence.UserRepositoryImpl) *domainservice.UserDomainService {
	// Repository implements domain.UserRepository interface directly
	var userRepository domain.UserRepository = userRepo
	return domainservice.NewUserDomainService(userRepository)
}

// CreateTeamDomainService creates a team domain service with repository implementing domain interface
func CreateTeamDomainService(teamRepo *persistence.TeamRepositoryImpl) *domainservice.TeamDomainService {
	// Repository implements domain.TeamRepository interface directly
	var teamRepository domain.TeamRepository = teamRepo
	return domainservice.NewTeamDomainService(teamRepository)
}

// CreatePlayerDomainService creates a player domain service with repository implementing domain interface
func CreatePlayerDomainService(playerRepo *persistence.PlayerRepositoryImpl) *domainservice.PlayerDomainService {
	// Repository implements domain.PlayerRepository interface directly
	var playerRepository domain.PlayerRepository = playerRepo
	return domainservice.NewPlayerDomainService(playerRepository)
}

// CreatePlayerTeamDomainService creates a player team domain service with repository implementing domain interface
func CreatePlayerTeamDomainService(
	playerTeamRepo *persistence.PlayerTeamRepositoryImpl,
	playerRepo *persistence.PlayerRepositoryImpl,
	teamRepo *persistence.TeamRepositoryImpl,
	seasonRepo *persistence.SeasonRepositoryImpl,
) *domainservice.PlayerTeamDomainService {
	// Repositories implement domain interfaces directly
	var playerTeamRepository domain.PlayerTeamRepository = playerTeamRepo
	var playerRepository domain.PlayerRepository = playerRepo
	var teamRepository domain.TeamRepository = teamRepo
	var seasonRepository domain.SeasonRepository = seasonRepo
	return domainservice.NewPlayerTeamDomainService(playerTeamRepository, playerRepository, teamRepository, seasonRepository)
}
