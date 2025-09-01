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
	// Check if a team with this name already exists
	existing, err := s.TeamRepository.GetTeamByName(ctx, team.FullName)
	if err != nil {
		return nil, fmt.Errorf("failed to check existing team: %w", err)
	}
	if existing != nil {
		return nil, constants.ErrRecordAlreadyExists
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
	if team == nil {
		return nil, constants.ErrRecordNotFound
	}
	return team, nil
}

func (s *teamService) GetTeamByName(ctx context.Context, fullName string) (*model.Team, error) {
	team, err := s.TeamRepository.GetTeamByName(ctx, fullName)
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
		dup, err := s.TeamRepository.GetTeamByName(ctx, *teamUpdate.FullName)
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
