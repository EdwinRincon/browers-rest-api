package server

import (
	"github.com/EdwinRincon/browersfc-api/domain"
	domainservice "github.com/EdwinRincon/browersfc-api/internal/domain/service"
)

// CreateMatchDomainService creates a match domain service with repository implementing domain interface
func CreateMatchDomainService(matchRepo domain.MatchRepository) *domainservice.MatchDomainService {
	// Repository already implements domain.MatchRepository interface
	return domainservice.NewMatchDomainService(matchRepo)
}

// CreateRoleDomainService creates a role domain service with repository implementing domain interface
func CreateRoleDomainService(roleRepo domain.RoleRepository) *domainservice.RoleDomainService {
	return domainservice.NewRoleDomainService(roleRepo)
}

// CreateSeasonDomainService creates a season domain service with repository implementing domain interface
func CreateSeasonDomainService(seasonRepo domain.SeasonRepository) *domainservice.SeasonDomainService {
	return domainservice.NewSeasonDomainService(seasonRepo)
}

// CreateUserDomainService creates a user domain service with repository implementing domain interface
func CreateUserDomainService(userRepo domain.UserRepository) *domainservice.UserDomainService {
	return domainservice.NewUserDomainService(userRepo)
}

// CreateTeamDomainService creates a team domain service with repository implementing domain interface
func CreateTeamDomainService(teamRepo domain.TeamRepository) *domainservice.TeamDomainService {
	return domainservice.NewTeamDomainService(teamRepo)
}

// CreatePlayerDomainService creates a player domain service with repository implementing domain interface
func CreatePlayerDomainService(playerRepo domain.PlayerRepository) *domainservice.PlayerDomainService {
	return domainservice.NewPlayerDomainService(playerRepo)
}

// CreatePlayerTeamDomainService creates a player team domain service with repository implementing domain interface
func CreatePlayerTeamDomainService(
	playerTeamRepo domain.PlayerTeamRepository,
	playerRepo domain.PlayerRepository,
	teamRepo domain.TeamRepository,
	seasonRepo domain.SeasonRepository,
) *domainservice.PlayerTeamDomainService {
	return domainservice.NewPlayerTeamDomainService(playerTeamRepo, playerRepo, teamRepo, seasonRepo)
}
