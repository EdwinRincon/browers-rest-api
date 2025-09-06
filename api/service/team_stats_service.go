package service

import (
	"context"
	"fmt"

	"github.com/EdwinRincon/browersfc-api/api/constants"
	"github.com/EdwinRincon/browersfc-api/api/dto"
	"github.com/EdwinRincon/browersfc-api/api/mapper"
	"github.com/EdwinRincon/browersfc-api/api/model"
	"github.com/EdwinRincon/browersfc-api/internal/infrastructure/persistence"
)

type TeamStatsService interface {
	CreateTeamStats(ctx context.Context, teamStats *model.TeamStat) (*dto.TeamStatsShort, error)
	GetTeamStatsByID(ctx context.Context, id uint64) (*model.TeamStat, error)
	GetTeamStatsBySeasonAndTeam(ctx context.Context, seasonID, teamID uint64) (*model.TeamStat, error)
	GetTeamStatsBySeasonID(ctx context.Context, seasonID uint64) ([]model.TeamStat, error)
	GetTeamStatsByTeamID(ctx context.Context, teamID uint64) ([]model.TeamStat, error)
	GetPaginatedTeamStats(ctx context.Context, sort string, order string, page int, pageSize int) ([]model.TeamStat, int64, error)
	UpdateTeamStats(ctx context.Context, teamStatsUpdate *dto.UpdateTeamStatsRequest, teamStatsID uint64) (*model.TeamStat, error)
	DeleteTeamStats(ctx context.Context, id uint64) error
}

type teamStatsService struct {
	TeamStatsRepository persistence.TeamStatsRepository
	TeamRepository      persistence.TeamRepository
	SeasonRepository    persistence.SeasonRepository
}

func NewTeamStatsService(teamStatsRepo persistence.TeamStatsRepository, teamRepo persistence.TeamRepository, seasonRepo persistence.SeasonRepository) TeamStatsService {
	return &teamStatsService{
		TeamStatsRepository: teamStatsRepo,
		TeamRepository:      teamRepo,
		SeasonRepository:    seasonRepo,
	}
}

func (s *teamStatsService) CreateTeamStats(ctx context.Context, teamStats *model.TeamStat) (*dto.TeamStatsShort, error) {
	// Validate that the team exists
	team, err := s.TeamRepository.GetTeamByID(ctx, teamStats.TeamID)
	if err != nil {
		return nil, fmt.Errorf("failed to check team existence: %w", err)
	}
	if team == nil {
		return nil, constants.ErrTeamNotFound
	}

	// Validate that the season exists
	season, err := s.SeasonRepository.GetSeasonByID(ctx, teamStats.SeasonID)
	if err != nil {
		return nil, fmt.Errorf("failed to check season existence: %w", err)
	}
	if season == nil {
		return nil, constants.ErrSeasonNotFound
	}

	// Check if team stats already exist for this season and team combination
	existing, err := s.TeamStatsRepository.GetTeamStatsBySeasonAndTeam(ctx, teamStats.SeasonID, teamStats.TeamID)
	if err != nil {
		return nil, fmt.Errorf("failed to check existing team stats: %w", err)
	}
	if existing != nil {
		return nil, constants.ErrRecordAlreadyExists
	}

	if err := s.TeamStatsRepository.CreateTeamStats(ctx, teamStats); err != nil {
		return nil, fmt.Errorf("failed to create team stats: %w", err)
	}

	return mapper.ToTeamStatsShort(teamStats), nil
}

func (s *teamStatsService) GetTeamStatsByID(ctx context.Context, id uint64) (*model.TeamStat, error) {
	teamStats, err := s.TeamStatsRepository.GetTeamStatsByID(ctx, id)
	if err != nil {
		return nil, err
	}
	if teamStats == nil {
		return nil, constants.ErrRecordNotFound
	}
	return teamStats, nil
}

func (s *teamStatsService) GetTeamStatsBySeasonAndTeam(ctx context.Context, seasonID, teamID uint64) (*model.TeamStat, error) {
	teamStats, err := s.TeamStatsRepository.GetTeamStatsBySeasonAndTeam(ctx, seasonID, teamID)
	if err != nil {
		return nil, err
	}
	if teamStats == nil {
		return nil, constants.ErrRecordNotFound
	}
	return teamStats, nil
}

func (s *teamStatsService) GetTeamStatsBySeasonID(ctx context.Context, seasonID uint64) ([]model.TeamStat, error) {

	season, err := s.SeasonRepository.GetSeasonByID(ctx, seasonID)
	if err != nil {
		return nil, fmt.Errorf("failed to check season existence: %w", err)
	}
	if season == nil {
		return nil, constants.ErrSeasonNotFound
	}

	return s.TeamStatsRepository.GetTeamStatsBySeasonID(ctx, seasonID)
}

func (s *teamStatsService) GetTeamStatsByTeamID(ctx context.Context, teamID uint64) ([]model.TeamStat, error) {

	team, err := s.TeamRepository.GetTeamByID(ctx, teamID)
	if err != nil {
		return nil, fmt.Errorf("failed to check team existence: %w", err)
	}
	if team == nil {
		return nil, constants.ErrTeamNotFound
	}

	return s.TeamStatsRepository.GetTeamStatsByTeamID(ctx, teamID)
}

func (s *teamStatsService) GetPaginatedTeamStats(ctx context.Context, sort string, order string, page int, pageSize int) ([]model.TeamStat, int64, error) {
	return s.TeamStatsRepository.GetPaginatedTeamStats(ctx, sort, order, page, pageSize)
}

func (s *teamStatsService) UpdateTeamStats(ctx context.Context, teamStatsUpdate *dto.UpdateTeamStatsRequest, teamStatsID uint64) (*model.TeamStat, error) {
	teamStats, err := s.TeamStatsRepository.GetTeamStatsByID(ctx, teamStatsID)
	if err != nil {
		return nil, fmt.Errorf("failed to get team stats by ID: %w", err)
	}
	if teamStats == nil {
		return nil, constants.ErrRecordNotFound
	}

	// Validate team if being updated
	if teamStatsUpdate.TeamID != nil && *teamStatsUpdate.TeamID != teamStats.TeamID {
		team, err := s.TeamRepository.GetTeamByID(ctx, *teamStatsUpdate.TeamID)
		if err != nil {
			return nil, fmt.Errorf("failed to check team existence: %w", err)
		}
		if team == nil {
			return nil, constants.ErrTeamNotFound
		}
	}

	// Validate season if being updated
	if teamStatsUpdate.SeasonID != nil && *teamStatsUpdate.SeasonID != teamStats.SeasonID {
		season, err := s.SeasonRepository.GetSeasonByID(ctx, *teamStatsUpdate.SeasonID)
		if err != nil {
			return nil, fmt.Errorf("failed to check season existence: %w", err)
		}
		if season == nil {
			return nil, constants.ErrSeasonNotFound
		}
	}

	// Check if updating season/team combination would create a duplicate
	if (teamStatsUpdate.SeasonID != nil && *teamStatsUpdate.SeasonID != teamStats.SeasonID) ||
		(teamStatsUpdate.TeamID != nil && *teamStatsUpdate.TeamID != teamStats.TeamID) {

		newSeasonID := teamStats.SeasonID
		newTeamID := teamStats.TeamID

		if teamStatsUpdate.SeasonID != nil {
			newSeasonID = *teamStatsUpdate.SeasonID
		}
		if teamStatsUpdate.TeamID != nil {
			newTeamID = *teamStatsUpdate.TeamID
		}

		existing, err := s.TeamStatsRepository.GetTeamStatsBySeasonAndTeam(ctx, newSeasonID, newTeamID)
		if err != nil {
			return nil, fmt.Errorf("failed to check duplicate team stats: %w", err)
		}
		if existing != nil && existing.ID != teamStats.ID {
			return nil, constants.ErrRecordAlreadyExists
		}
	}

	mapper.UpdateTeamStatsFromDTO(teamStats, teamStatsUpdate)

	if err := s.TeamStatsRepository.UpdateTeamStats(ctx, teamStatsID, teamStats); err != nil {
		return nil, fmt.Errorf("failed to update team stats: %w", err)
	}

	return teamStats, nil
}

func (s *teamStatsService) DeleteTeamStats(ctx context.Context, id uint64) error {
	return s.TeamStatsRepository.DeleteTeamStats(ctx, id)
}
