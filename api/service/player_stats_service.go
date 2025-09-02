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

type PlayerStatsService interface {
	CreatePlayerStat(ctx context.Context, playerStatDTO *dto.CreatePlayerStatRequest) (*model.PlayerStat, error)
	GetPlayerStatByID(ctx context.Context, id uint64) (*model.PlayerStat, error)
	GetPlayerStatsByPlayerID(ctx context.Context, playerID uint64) ([]model.PlayerStat, error)
	GetPlayerStatsByMatchID(ctx context.Context, matchID uint64) ([]model.PlayerStat, error)
	GetPlayerStatsBySeasonID(ctx context.Context, seasonID uint64) ([]model.PlayerStat, error)
	GetPaginatedPlayerStats(ctx context.Context, sort string, order string, page int, pageSize int) ([]model.PlayerStat, int64, error)
	UpdatePlayerStat(ctx context.Context, id uint64, updateDTO *dto.UpdatePlayerStatRequest) (*model.PlayerStat, error)
	DeletePlayerStat(ctx context.Context, id uint64) error
}

type playerStatsService struct {
	PlayerStatsRepository repository.PlayerStatsRepository
	PlayerRepository      repository.PlayerRepository
	MatchRepository       repository.MatchRepository
	SeasonRepository      repository.SeasonRepository
	TeamRepository        repository.TeamRepository
}

func NewPlayerStatsService(
	playerStatsRepo repository.PlayerStatsRepository,
	playerRepo repository.PlayerRepository,
	matchRepo repository.MatchRepository,
	seasonRepo repository.SeasonRepository,
	teamRepo repository.TeamRepository,
) PlayerStatsService {
	return &playerStatsService{
		PlayerStatsRepository: playerStatsRepo,
		PlayerRepository:      playerRepo,
		MatchRepository:       matchRepo,
		SeasonRepository:      seasonRepo,
		TeamRepository:        teamRepo,
	}
}

func (s *playerStatsService) CreatePlayerStat(ctx context.Context, playerStatDTO *dto.CreatePlayerStatRequest) (*model.PlayerStat, error) {
	// Validate Player exists
	player, err := s.PlayerRepository.GetPlayerByID(ctx, playerStatDTO.PlayerID)
	if err != nil {
		return nil, fmt.Errorf("failed to check player: %w", err)
	}
	if player == nil {
		return nil, constants.ErrPlayerNotFound
	}

	// Validate Match exists
	match, err := s.MatchRepository.GetMatchByID(ctx, playerStatDTO.MatchID)
	if err != nil {
		return nil, fmt.Errorf("failed to check match: %w", err)
	}
	if match == nil {
		return nil, constants.ErrRecordNotFound
	}

	// Validate Season exists
	season, err := s.SeasonRepository.GetSeasonByID(ctx, playerStatDTO.SeasonID)
	if err != nil {
		return nil, fmt.Errorf("failed to check season: %w", err)
	}
	if season == nil {
		return nil, constants.ErrSeasonNotFound
	}

	// Validate Team exists if provided
	if playerStatDTO.TeamID != nil {
		team, err := s.TeamRepository.GetTeamByID(ctx, *playerStatDTO.TeamID)
		if err != nil {
			return nil, fmt.Errorf("failed to check team: %w", err)
		}
		if team == nil {
			return nil, constants.ErrTeamNotFound
		}
	}

	// Create PlayerStat
	playerStat := mapper.ToPlayerStat(playerStatDTO)

	if err := s.PlayerStatsRepository.CreatePlayerStat(ctx, playerStat); err != nil {
		return nil, fmt.Errorf("failed to create player stat: %w", err)
	}

	// Return created PlayerStat
	return playerStat, nil
}

func (s *playerStatsService) GetPlayerStatByID(ctx context.Context, id uint64) (*model.PlayerStat, error) {
	playerStat, err := s.PlayerStatsRepository.GetPlayerStatByID(ctx, id)
	if err != nil {
		return nil, fmt.Errorf("failed to get player stat by ID: %w", err)
	}
	if playerStat == nil {
		return nil, constants.ErrRecordNotFound
	}
	return playerStat, nil
}

// GetPlayerStatsByPlayerID retrieves all player stats for a specific player.
func (s *playerStatsService) GetPlayerStatsByPlayerID(ctx context.Context, playerID uint64) ([]model.PlayerStat, error) {
	// Check if player exists
	player, err := s.PlayerRepository.GetPlayerByID(ctx, playerID)
	if err != nil {
		return nil, fmt.Errorf("failed to check player: %w", err)
	}
	if player == nil {
		return nil, constants.ErrPlayerNotFound
	}

	playerStats, err := s.PlayerStatsRepository.GetPlayerStatsByPlayerID(ctx, playerID)
	if err != nil {
		return nil, fmt.Errorf("failed to get player stats by player ID: %w", err)
	}
	return playerStats, nil
}

// GetPlayerStatsByMatchID retrieves all player stats for a specific match.
func (s *playerStatsService) GetPlayerStatsByMatchID(ctx context.Context, matchID uint64) ([]model.PlayerStat, error) {
	// Check if match exists
	match, err := s.MatchRepository.GetMatchByID(ctx, matchID)
	if err != nil {
		return nil, fmt.Errorf("failed to check match: %w", err)
	}
	if match == nil {
		return nil, constants.ErrRecordNotFound
	}

	playerStats, err := s.PlayerStatsRepository.GetPlayerStatsByMatchID(ctx, matchID)
	if err != nil {
		return nil, fmt.Errorf("failed to get player stats by match ID: %w", err)
	}
	return playerStats, nil
}

// GetPlayerStatsBySeasonID retrieves all player stats for a specific season.
func (s *playerStatsService) GetPlayerStatsBySeasonID(ctx context.Context, seasonID uint64) ([]model.PlayerStat, error) {
	// Check if season exists
	season, err := s.SeasonRepository.GetSeasonByID(ctx, seasonID)
	if err != nil {
		return nil, fmt.Errorf("failed to check season: %w", err)
	}
	if season == nil {
		return nil, constants.ErrSeasonNotFound
	}

	playerStats, err := s.PlayerStatsRepository.GetPlayerStatsBySeasonID(ctx, seasonID)
	if err != nil {
		return nil, fmt.Errorf("failed to get player stats by season ID: %w", err)
	}
	return playerStats, nil
}

func (s *playerStatsService) GetPaginatedPlayerStats(ctx context.Context, sort string, order string, page int, pageSize int) ([]model.PlayerStat, int64, error) {
	return s.PlayerStatsRepository.GetPaginatedPlayerStats(ctx, sort, order, page, pageSize)
}

func (s *playerStatsService) UpdatePlayerStat(ctx context.Context, id uint64, updateDTO *dto.UpdatePlayerStatRequest) (*model.PlayerStat, error) {
	playerStat, err := s.PlayerStatsRepository.GetPlayerStatByID(ctx, id)
	if err != nil {
		return nil, fmt.Errorf("failed to get player stat by ID: %w", err)
	}
	if playerStat == nil {
		return nil, constants.ErrRecordNotFound
	}

	// Validate Team exists if provided
	if updateDTO.TeamID != nil {
		team, err := s.TeamRepository.GetTeamByID(ctx, *updateDTO.TeamID)
		if err != nil {
			return nil, fmt.Errorf("failed to check team: %w", err)
		}
		if team == nil {
			return nil, constants.ErrTeamNotFound
		}
	}

	// Update player stat
	mapper.UpdatePlayerStatFromDTO(playerStat, updateDTO)

	if err := s.PlayerStatsRepository.UpdatePlayerStat(ctx, id, playerStat); err != nil {
		return nil, fmt.Errorf("failed to update player stat: %w", err)
	}

	// Return updated player stat
	return playerStat, nil
}

func (s *playerStatsService) DeletePlayerStat(ctx context.Context, id uint64) error {
	playerStat, err := s.PlayerStatsRepository.GetPlayerStatByID(ctx, id)
	if err != nil {
		return fmt.Errorf("failed to get player stat by ID: %w", err)
	}
	if playerStat == nil {
		return constants.ErrRecordNotFound
	}

	return s.PlayerStatsRepository.DeletePlayerStat(ctx, id)
}
