package adapters

import (
	"context"

	"github.com/EdwinRincon/browersfc-api/api/mapper"
	"github.com/EdwinRincon/browersfc-api/domain"
	"github.com/EdwinRincon/browersfc-api/internal/ports"
)

// TeamDomainAdapter implements the domain operations for teams.
// This adapter bridges the domain layer with the persistence layer.
type TeamDomainAdapter struct {
	teamPort ports.TeamPort
}

// NewTeamDomainAdapter creates a new team domain adapter.
func NewTeamDomainAdapter(teamPort ports.TeamPort) *TeamDomainAdapter {
	return &TeamDomainAdapter{
		teamPort: teamPort,
	}
}

// CreateTeam creates a new team using domain entities.
func (a *TeamDomainAdapter) CreateTeam(ctx context.Context, team *domain.Team) (*domain.Team, error) {
	// Convert domain entity to model
	modelTeam := mapper.TeamDomainToPersistence(team)

	// Create through port
	err := a.teamPort.CreateTeam(ctx, modelTeam)
	if err != nil {
		return nil, err
	}

	// Convert back to domain entity
	return mapper.TeamPersistenceToDomain(modelTeam), nil
}

// GetTeamByID retrieves a team by ID as domain entity.
func (a *TeamDomainAdapter) GetTeamByID(ctx context.Context, id uint64) (*domain.Team, error) {
	// Get model from port
	modelTeam, err := a.teamPort.GetTeamByID(ctx, id)
	if err != nil {
		return nil, err
	}

	// Convert to domain entity
	return mapper.TeamPersistenceToDomain(modelTeam), nil
}

// GetTeamByName retrieves a team by name as domain entity.
func (a *TeamDomainAdapter) GetTeamByName(ctx context.Context, name string) (*domain.Team, error) {
	// Get model from port
	modelTeam, err := a.teamPort.GetTeamByName(ctx, name)
	if err != nil {
		return nil, err
	}

	// Convert to domain entity
	return mapper.TeamPersistenceToDomain(modelTeam), nil
}

// UpdateTeam updates a team using domain entities.
func (a *TeamDomainAdapter) UpdateTeam(ctx context.Context, team *domain.Team) (*domain.Team, error) {
	// Convert domain entity to model
	modelTeam := mapper.TeamDomainToPersistence(team)

	// Update through port
	err := a.teamPort.UpdateTeam(ctx, modelTeam)
	if err != nil {
		return nil, err
	}

	// Convert back to domain entity
	return mapper.TeamPersistenceToDomain(modelTeam), nil
}

// DeleteTeam deletes a team by ID.
func (a *TeamDomainAdapter) DeleteTeam(ctx context.Context, id uint64) error {
	return a.teamPort.DeleteTeam(ctx, id)
}

// GetAllTeams retrieves all teams as domain entities.
func (a *TeamDomainAdapter) GetAllTeams(ctx context.Context) ([]domain.Team, error) {
	// Get paginated teams (using a large page size to get all)
	modelTeams, _, err := a.teamPort.GetPaginatedTeams(ctx, "id", "asc", 1, 1000)
	if err != nil {
		return nil, err
	}

	// Convert to domain entities
	return mapper.TeamPersistenceListToDomain(modelTeams), nil
}
