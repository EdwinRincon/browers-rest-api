package service

import (
	"context"
	"fmt"

	"github.com/EdwinRincon/browersfc-api/api/constants"
	"github.com/EdwinRincon/browersfc-api/domain"
)

type TeamStatsDomainService struct {
	teamStatsRepository domain.TeamStatsRepository
	teamRepository      domain.TeamRepository
	seasonRepository    domain.SeasonRepository
}

func NewTeamStatsDomainService(
	teamStatsRepository domain.TeamStatsRepository,
	teamRepository domain.TeamRepository,
	seasonRepository domain.SeasonRepository,
) *TeamStatsDomainService {
	return &TeamStatsDomainService{
		teamStatsRepository: teamStatsRepository,
		teamRepository:      teamRepository,
		seasonRepository:    seasonRepository,
	}
}

func (s *TeamStatsDomainService) CreateTeamStats(ctx context.Context, teamStats *domain.TeamStats) error {
	// Validate that the team exists
	team, err := s.teamRepository.GetTeamByID(ctx, teamStats.TeamID)
	if err != nil {
		return fmt.Errorf("failed to check team existence: %w", err)
	}
	if team == nil {
		return constants.ErrTeamNotFound
	}

	// Validate that the season exists
	season, err := s.seasonRepository.GetSeasonByID(ctx, teamStats.SeasonID)
	if err != nil {
		return fmt.Errorf("failed to check season existence: %w", err)
	}
	if season == nil {
		return constants.ErrSeasonNotFound
	}

	// Check if team stats already exist for this season and team combination
	existing, err := s.teamStatsRepository.GetTeamStatsBySeasonAndTeam(ctx, teamStats.SeasonID, teamStats.TeamID)
	if err != nil {
		return fmt.Errorf("failed to check existing team stats: %w", err)
	}
	if existing != nil {
		return constants.ErrRecordAlreadyExists
	}

	if err := s.teamStatsRepository.CreateTeamStats(ctx, teamStats); err != nil {
		return fmt.Errorf("failed to create team stats: %w", err)
	}

	return nil
}

// GetTeamStatsByID retrieves team statistics by ID
func (s *TeamStatsDomainService) GetTeamStatsByID(ctx context.Context, id uint64) (*domain.TeamStats, error) {
	teamStats, err := s.teamStatsRepository.GetTeamStatsByID(ctx, id)
	if err != nil {
		return nil, err
	}
	if teamStats == nil {
		return nil, constants.ErrRecordNotFound
	}
	return teamStats, nil
}

// GetTeamStatsBySeasonAndTeam retrieves team statistics by season and team
func (s *TeamStatsDomainService) GetTeamStatsBySeasonAndTeam(ctx context.Context, seasonID, teamID uint64) (*domain.TeamStats, error) {
	teamStats, err := s.teamStatsRepository.GetTeamStatsBySeasonAndTeam(ctx, seasonID, teamID)
	if err != nil {
		return nil, err
	}
	if teamStats == nil {
		return nil, constants.ErrRecordNotFound
	}
	return teamStats, nil
}

// GetTeamStatsBySeasonID retrieves all team statistics for a season
func (s *TeamStatsDomainService) GetTeamStatsBySeasonID(ctx context.Context, seasonID uint64) ([]domain.TeamStats, error) {
	season, err := s.seasonRepository.GetSeasonByID(ctx, seasonID)
	if err != nil {
		return nil, fmt.Errorf("failed to check season existence: %w", err)
	}
	if season == nil {
		return nil, constants.ErrSeasonNotFound
	}

	return s.teamStatsRepository.GetTeamStatsBySeasonID(ctx, seasonID)
}

// GetTeamStatsByTeamID retrieves all team statistics for a team
func (s *TeamStatsDomainService) GetTeamStatsByTeamID(ctx context.Context, teamID uint64) ([]domain.TeamStats, error) {
	team, err := s.teamRepository.GetTeamByID(ctx, teamID)
	if err != nil {
		return nil, fmt.Errorf("failed to check team existence: %w", err)
	}
	if team == nil {
		return nil, constants.ErrTeamNotFound
	}

	return s.teamStatsRepository.GetTeamStatsByTeamID(ctx, teamID)
}

// GetPaginatedTeamStats retrieves paginated team statistics.
func (s *TeamStatsDomainService) GetPaginatedTeamStats(ctx context.Context, sort string, order string, page int, pageSize int) ([]domain.TeamStats, int64, error) {
	return s.teamStatsRepository.GetPaginatedTeamStats(ctx, sort, order, page, pageSize)
}

// UpdateTeamStats updates team statistics
func (s *TeamStatsDomainService) UpdateTeamStats(ctx context.Context, teamStatsID uint64, updatedTeamStats *domain.TeamStats) (*domain.TeamStats, error) {
	currentTeamStats, err := s.teamStatsRepository.GetTeamStatsByID(ctx, teamStatsID)
	if err != nil {
		return nil, fmt.Errorf("failed to get team stats by ID: %w", err)
	}
	if currentTeamStats == nil {
		return nil, constants.ErrRecordNotFound
	}

	// Validate team if being updated
	if updatedTeamStats.TeamID != currentTeamStats.TeamID {
		team, err := s.teamRepository.GetTeamByID(ctx, updatedTeamStats.TeamID)
		if err != nil {
			return nil, fmt.Errorf("failed to check team existence: %w", err)
		}
		if team == nil {
			return nil, constants.ErrTeamNotFound
		}
	}

	// Validate season if being updated
	if updatedTeamStats.SeasonID != currentTeamStats.SeasonID {
		season, err := s.seasonRepository.GetSeasonByID(ctx, updatedTeamStats.SeasonID)
		if err != nil {
			return nil, fmt.Errorf("failed to check season existence: %w", err)
		}
		if season == nil {
			return nil, constants.ErrSeasonNotFound
		}
	}

	// Check if updating season/team combination would create a duplicate
	if updatedTeamStats.SeasonID != currentTeamStats.SeasonID || updatedTeamStats.TeamID != currentTeamStats.TeamID {
		existing, err := s.teamStatsRepository.GetTeamStatsBySeasonAndTeam(ctx, updatedTeamStats.SeasonID, updatedTeamStats.TeamID)
		if err != nil {
			return nil, fmt.Errorf("failed to check duplicate team stats: %w", err)
		}
		if existing != nil && existing.ID != currentTeamStats.ID {
			return nil, constants.ErrRecordAlreadyExists
		}
	}

	// Set the ID to ensure we're updating the correct record
	updatedTeamStats.ID = teamStatsID

	if err := s.teamStatsRepository.UpdateTeamStats(ctx, teamStatsID, updatedTeamStats); err != nil {
		return nil, fmt.Errorf("failed to update team stats: %w", err)
	}

	// Return the updated team stats
	return s.teamStatsRepository.GetTeamStatsByID(ctx, teamStatsID)
}

// DeleteTeamStats deletes team statistics by ID.
func (s *TeamStatsDomainService) DeleteTeamStats(ctx context.Context, id uint64) error {
	return s.teamStatsRepository.DeleteTeamStats(ctx, id)
}
