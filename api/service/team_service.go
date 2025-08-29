package service

import (
	"context"
	"fmt"

	"github.com/EdwinRincon/browersfc-api/api/constants"
	"github.com/EdwinRincon/browersfc-api/api/dto"
	"github.com/EdwinRincon/browersfc-api/api/mapper"
	"github.com/EdwinRincon/browersfc-api/api/model"
	"github.com/EdwinRincon/browersfc-api/api/repository"
)

type TeamService interface {
	CreateTeam(ctx context.Context, team *model.Team) (*dto.TeamShort, error)
	GetTeamByID(ctx context.Context, id uint64) (*model.Team, error)
	GetTeamByName(ctx context.Context, fullName string) (*model.Team, error)
	GetPaginatedTeams(ctx context.Context, sort string, order string, page int, pageSize int) ([]model.Team, int64, error)
	UpdateTeam(ctx context.Context, teamUpdate *dto.UpdateTeamRequest, teamID uint64) (*model.Team, error)
	DeleteTeam(ctx context.Context, id uint64) error
}

type teamService struct {
	TeamRepository repository.TeamRepository
}

func NewTeamService(teamRepo repository.TeamRepository) TeamService {
	return &teamService{
		TeamRepository: teamRepo,
	}
}

func (s *teamService) CreateTeam(ctx context.Context, team *model.Team) (*dto.TeamShort, error) {
	// First check if there's an active team with this name
	activeTeam, err := s.TeamRepository.GetActiveTeamByName(ctx, team.FullName)
	if err != nil {
		return nil, fmt.Errorf("failed to check existing active team: %w", err)
	}
	if activeTeam != nil {
		return nil, constants.ErrRecordAlreadyExists
	}

	// If no active team exists, check for soft-deleted team
	existing, err := s.TeamRepository.GetUnscopedTeamByName(ctx, team.FullName)
	if err != nil {
		return nil, fmt.Errorf("failed to check existing team: %w", err)
	}

	if existing != nil && existing.DeletedAt.Valid {
		// Update all relevant fields from the new team
		existing.FullName = team.FullName
		existing.ShortName = team.ShortName
		existing.Color = team.Color
		existing.Color2 = team.Color2
		existing.Shield = team.Shield
		existing.NextMatchID = team.NextMatchID

		// Restore and update the team in a transaction
		err := s.TeamRepository.RestoreAndUpdateTeam(ctx, existing)
		if err != nil {
			return nil, fmt.Errorf("failed to restore and update team: %w", err)
		}

		return mapper.ToTeamShort(existing), nil
	}

	if err := s.TeamRepository.CreateTeam(ctx, team); err != nil {
		return nil, fmt.Errorf("failed to create team: %w", err)
	}

	return mapper.ToTeamShort(team), nil
}

func (s *teamService) GetTeamByID(ctx context.Context, id uint64) (*model.Team, error) {
	team, err := s.TeamRepository.GetTeamByID(ctx, id)
	if err != nil {
		return nil, err
	}
	return team, nil
}

func (s *teamService) GetTeamByName(ctx context.Context, fullName string) (*model.Team, error) {
	team, err := s.TeamRepository.GetActiveTeamByName(ctx, fullName)
	if err != nil {
		return nil, err
	}
	if team == nil {
		return nil, constants.ErrRecordNotFound
	}
	return team, nil
}

func (s *teamService) GetPaginatedTeams(ctx context.Context, sort string, order string, page int, pageSize int) ([]model.Team, int64, error) {
	return s.TeamRepository.GetPaginatedTeams(ctx, sort, order, page, pageSize)
}

func (s *teamService) UpdateTeam(ctx context.Context, teamUpdate *dto.UpdateTeamRequest, teamID uint64) (*model.Team, error) {
	team, err := s.TeamRepository.GetTeamByID(ctx, teamID)
	if err != nil {
		return nil, fmt.Errorf("failed to get team by ID: %w", err)
	}
	if team == nil {
		return nil, constants.ErrRecordNotFound
	}

	if teamUpdate.FullName != nil && *teamUpdate.FullName != team.FullName {
		dup, err := s.TeamRepository.GetActiveTeamByName(ctx, *teamUpdate.FullName)
		if err != nil {
			return nil, fmt.Errorf("failed to check duplicate team name: %w", err)
		}
		if dup != nil && dup.ID != team.ID {
			return nil, constants.ErrRecordAlreadyExists
		}
	}

	mapper.UpdateTeamFromDTO(team, teamUpdate)

	if err := s.TeamRepository.UpdateTeam(ctx, team); err != nil {
		return nil, fmt.Errorf("failed to update team: %w", err)
	}

	return team, nil
}

func (s *teamService) DeleteTeam(ctx context.Context, id uint64) error {
	return s.TeamRepository.DeleteTeam(ctx, id)
}
