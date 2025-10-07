package service

import (
	"context"
	"fmt"

	"github.com/EdwinRincon/browersfc-api/api/constants"
	"github.com/EdwinRincon/browersfc-api/domain"
)

// Error constants
const (
	ErrTeamIDCannotBeZero = "team ID cannot be zero"
)

// TeamDomainService encapsulates business logic for team operations.
type TeamDomainService struct {
	teamRepository domain.TeamRepository
}

// NewTeamDomainService creates a new TeamDomainService instance.
func NewTeamDomainService(teamRepository domain.TeamRepository) *TeamDomainService {
	return &TeamDomainService{
		teamRepository: teamRepository,
	}
}

func (s *TeamDomainService) CreateTeam(ctx context.Context, team *domain.Team) (*domain.Team, error) {
	// Validate domain entity
	if !team.IsValid() {
		return nil, fmt.Errorf("invalid team data")
	}

	// Check if a team with this name already exists
	existing, err := s.teamRepository.GetTeamByName(ctx, team.ShortName)
	if err != nil {
		return nil, fmt.Errorf("failed to check existing team: %w", err)
	}
	if existing != nil {
		return nil, constants.ErrRecordAlreadyExists
	}

	// Create the team
	if err := s.teamRepository.CreateTeam(ctx, team); err != nil {
		return nil, fmt.Errorf("failed to create team: %w", err)
	}

	return team, nil
}

func (s *TeamDomainService) GetTeamByID(ctx context.Context, id uint64) (*domain.Team, error) {
	if id == 0 {
		return nil, fmt.Errorf("%s", ErrTeamIDCannotBeZero)
	}

	team, err := s.teamRepository.GetTeamByID(ctx, id)
	if err != nil {
		return nil, fmt.Errorf("failed to get team by ID: %w", err)
	}
	if team == nil {
		return nil, constants.ErrRecordNotFound
	}

	return team, nil
}

func (s *TeamDomainService) GetTeamByName(ctx context.Context, name string) (*domain.Team, error) {
	if name == "" {
		return nil, fmt.Errorf("team name cannot be empty")
	}

	team, err := s.teamRepository.GetTeamByName(ctx, name)
	if err != nil {
		return nil, fmt.Errorf("failed to get team by name: %w", err)
	}
	if team == nil {
		return nil, constants.ErrRecordNotFound
	}

	return team, nil
}

func (s *TeamDomainService) GetPaginatedTeams(ctx context.Context, sort string, order string, page int, pageSize int) ([]domain.Team, int64, error) {
	teams, total, err := s.teamRepository.GetPaginatedTeams(ctx, sort, order, page, pageSize)
	if err != nil {
		return nil, 0, fmt.Errorf("failed to get paginated teams: %w", err)
	}

	return teams, total, nil
}

func (s *TeamDomainService) UpdateTeam(ctx context.Context, id uint64, updates *domain.Team) (*domain.Team, error) {
	if id == 0 {
		return nil, fmt.Errorf("%s", ErrTeamIDCannotBeZero)
	}

	// Get existing team
	existingTeam, err := s.teamRepository.GetTeamByID(ctx, id)
	if err != nil {
		return nil, fmt.Errorf("failed to get existing team: %w", err)
	}
	if existingTeam == nil {
		return nil, constants.ErrRecordNotFound
	}

	// Update fields
	if updates.FullName != "" {
		existingTeam.FullName = updates.FullName
	}
	if updates.ShortName != "" {
		existingTeam.ShortName = updates.ShortName
	}
	if updates.Shield != "" {
		existingTeam.Shield = updates.Shield
	}
	if updates.PrimaryColor != "" {
		existingTeam.PrimaryColor = updates.PrimaryColor
	}
	if updates.SecondaryColor != "" {
		existingTeam.SecondaryColor = updates.SecondaryColor
	}
	if updates.NextMatchID != nil {
		existingTeam.NextMatchID = updates.NextMatchID
	}

	// Validate updated entity
	if !existingTeam.IsValid() {
		return nil, fmt.Errorf("invalid team data after update")
	}

	// Save updated team
	if err := s.teamRepository.UpdateTeam(ctx, existingTeam); err != nil {
		return nil, fmt.Errorf("failed to update team: %w", err)
	}

	return existingTeam, nil
}

func (s *TeamDomainService) DeleteTeam(ctx context.Context, id uint64) error {
	if id == 0 {
		return fmt.Errorf("%s", ErrTeamIDCannotBeZero)
	}

	// Check if team exists
	existing, err := s.teamRepository.GetTeamByID(ctx, id)
	if err != nil {
		return fmt.Errorf("failed to check existing team: %w", err)
	}
	if existing == nil {
		return constants.ErrRecordNotFound
	}

	// Delete the team
	if err := s.teamRepository.DeleteTeam(ctx, id); err != nil {
		return fmt.Errorf("failed to delete team: %w", err)
	}

	return nil
}
