package adapters

import (
	"context"

	"github.com/EdwinRincon/browersfc-api/api/mapper"
	"github.com/EdwinRincon/browersfc-api/domain"
	"github.com/EdwinRincon/browersfc-api/internal/infrastructure/persistence"
	"github.com/EdwinRincon/browersfc-api/internal/ports"
)

// TeamDomainAdapter implements the domain operations for teams.
// This adapter bridges the domain layer with the persistence layer.
type TeamDomainAdapter struct {
	persistenceRepo persistence.TeamRepository
}

// NewTeamDomainAdapter creates a new team domain adapter.
func NewTeamDomainAdapter(persistenceRepo persistence.TeamRepository) ports.TeamPort {
	return &TeamDomainAdapter{
		persistenceRepo: persistenceRepo,
	}
}

// CreateTeam creates a new team using domain entities.
func (a *TeamDomainAdapter) CreateTeam(ctx context.Context, team *domain.Team) error {
	persistenceTeam := mapper.TeamDomainToPersistence(team)
	return a.persistenceRepo.CreateTeam(ctx, persistenceTeam)
}

// GetTeamByID retrieves a team by ID as domain entity.
func (a *TeamDomainAdapter) GetTeamByID(ctx context.Context, id uint64) (*domain.Team, error) {
	persistenceTeam, err := a.persistenceRepo.GetTeamByID(ctx, id)
	if err != nil {
		return nil, err
	}
	if persistenceTeam == nil {
		return nil, nil
	}
	return mapper.TeamPersistenceToDomain(persistenceTeam), nil
}

// GetTeamByName retrieves a team by name as domain entity.
func (a *TeamDomainAdapter) GetTeamByName(ctx context.Context, name string) (*domain.Team, error) {
	persistenceTeam, err := a.persistenceRepo.GetTeamByName(ctx, name)
	if err != nil {
		return nil, err
	}
	if persistenceTeam == nil {
		return nil, nil
	}
	return mapper.TeamPersistenceToDomain(persistenceTeam), nil
}

// GetPaginatedTeams retrieves paginated teams as domain entities.
func (a *TeamDomainAdapter) GetPaginatedTeams(ctx context.Context, sort string, order string, page int, pageSize int) ([]domain.Team, int64, error) {
	persistenceTeams, total, err := a.persistenceRepo.GetPaginatedTeams(ctx, sort, order, page, pageSize)
	if err != nil {
		return nil, 0, err
	}
	domainTeams := mapper.TeamPersistenceListToDomain(persistenceTeams)
	return domainTeams, total, nil
}

// UpdateTeam updates a team using domain entities.
func (a *TeamDomainAdapter) UpdateTeam(ctx context.Context, team *domain.Team) error {
	persistenceTeam := mapper.TeamDomainToPersistence(team)
	return a.persistenceRepo.UpdateTeam(ctx, persistenceTeam)
}

// DeleteTeam deletes a team by ID.
func (a *TeamDomainAdapter) DeleteTeam(ctx context.Context, id uint64) error {
	return a.persistenceRepo.DeleteTeam(ctx, id)
}
