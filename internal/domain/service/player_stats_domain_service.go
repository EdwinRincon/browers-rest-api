package service

import (
	"context"
	"fmt"

	"github.com/EdwinRincon/browersfc-api/api/constants"
	"github.com/EdwinRincon/browersfc-api/domain"
)

// PlayerStatsDomainService implements business logic for PlayerStat operations.
// It contains domain rules and validation while being infrastructure-agnostic.
type PlayerStatsDomainService struct {
	playerStatsRepository domain.PlayerStatsRepository
	playerRepository      domain.PlayerRepository
	matchRepository       domain.MatchRepository
	seasonRepository      domain.SeasonRepository
	teamRepository        domain.TeamRepository
}

func NewPlayerStatsDomainService(
	playerStatsRepository domain.PlayerStatsRepository,
	playerRepository domain.PlayerRepository,
	matchRepository domain.MatchRepository,
	seasonRepository domain.SeasonRepository,
	teamRepository domain.TeamRepository,
) *PlayerStatsDomainService {
	return &PlayerStatsDomainService{
		playerStatsRepository: playerStatsRepository,
		playerRepository:      playerRepository,
		matchRepository:       matchRepository,
		seasonRepository:      seasonRepository,
		teamRepository:        teamRepository,
	}
}

func (s *PlayerStatsDomainService) CreatePlayerStat(ctx context.Context, playerStat *domain.PlayerStat) error {
	// Validate domain rules
	if !playerStat.IsValid() {
		return constants.ErrInvalidData
	}

	// Validate Player exists
	player, err := s.playerRepository.GetPlayerByID(ctx, playerStat.PlayerID)
	if err != nil {
		return fmt.Errorf("failed to check player: %w", err)
	}
	if player == nil {
		return constants.ErrPlayerNotFound
	}

	// Validate Match exists
	match, err := s.matchRepository.GetMatchByID(ctx, playerStat.MatchID)
	if err != nil {
		return fmt.Errorf("failed to check match: %w", err)
	}
	if match == nil {
		return constants.ErrRecordNotFound
	}

	// Validate Season exists
	season, err := s.seasonRepository.GetSeasonByID(ctx, playerStat.SeasonID)
	if err != nil {
		return fmt.Errorf("failed to check season: %w", err)
	}
	if season == nil {
		return constants.ErrSeasonNotFound
	}

	// Validate Team exists if provided
	if playerStat.TeamID != nil {
		team, err := s.teamRepository.GetTeamByID(ctx, *playerStat.TeamID)
		if err != nil {
			return fmt.Errorf("failed to check team: %w", err)
		}
		if team == nil {
			return constants.ErrTeamNotFound
		}
	}

	return s.playerStatsRepository.CreatePlayerStat(ctx, playerStat)
}

func (s *PlayerStatsDomainService) GetPlayerStatByID(ctx context.Context, id uint64) (*domain.PlayerStat, error) {
	if id == 0 {
		return nil, constants.ErrInvalidData
	}

	playerStat, err := s.playerStatsRepository.GetPlayerStatByID(ctx, id)
	if err != nil {
		return nil, fmt.Errorf("failed to get player stat by ID: %w", err)
	}
	if playerStat == nil {
		return nil, constants.ErrRecordNotFound
	}

	return playerStat, nil
}

func (s *PlayerStatsDomainService) GetPlayerStatsByPlayerID(ctx context.Context, playerID uint64) ([]domain.PlayerStat, error) {
	if playerID == 0 {
		return nil, constants.ErrInvalidData
	}

	// Check if player exists
	player, err := s.playerRepository.GetPlayerByID(ctx, playerID)
	if err != nil {
		return nil, fmt.Errorf("failed to check player: %w", err)
	}
	if player == nil {
		return nil, constants.ErrPlayerNotFound
	}

	return s.playerStatsRepository.GetPlayerStatsByPlayerID(ctx, playerID)
}

func (s *PlayerStatsDomainService) GetPlayerStatsByMatchID(ctx context.Context, matchID uint64) ([]domain.PlayerStat, error) {
	if matchID == 0 {
		return nil, constants.ErrInvalidData
	}

	// Check if match exists
	match, err := s.matchRepository.GetMatchByID(ctx, matchID)
	if err != nil {
		return nil, fmt.Errorf("failed to check match: %w", err)
	}
	if match == nil {
		return nil, constants.ErrRecordNotFound
	}

	return s.playerStatsRepository.GetPlayerStatsByMatchID(ctx, matchID)
}

func (s *PlayerStatsDomainService) GetPlayerStatsBySeasonID(ctx context.Context, seasonID uint64) ([]domain.PlayerStat, error) {
	if seasonID == 0 {
		return nil, constants.ErrInvalidData
	}

	// Check if season exists
	season, err := s.seasonRepository.GetSeasonByID(ctx, seasonID)
	if err != nil {
		return nil, fmt.Errorf("failed to check season: %w", err)
	}
	if season == nil {
		return nil, constants.ErrSeasonNotFound
	}

	return s.playerStatsRepository.GetPlayerStatsBySeasonID(ctx, seasonID)
}

func (s *PlayerStatsDomainService) GetPaginatedPlayerStats(ctx context.Context, sort string, order string, page int, pageSize int) ([]domain.PlayerStat, int64, error) {
	return s.playerStatsRepository.GetPaginatedPlayerStats(ctx, sort, order, page, pageSize)
}

func (s *PlayerStatsDomainService) UpdatePlayerStat(ctx context.Context, id uint64, playerStat *domain.PlayerStat) (*domain.PlayerStat, error) {
	// Get existing player stat
	existingPlayerStat, err := s.playerStatsRepository.GetPlayerStatByID(ctx, id)
	if err != nil {
		return nil, fmt.Errorf("failed to get player stat by ID: %w", err)
	}
	if existingPlayerStat == nil {
		return nil, constants.ErrRecordNotFound
	}

	// Apply updates to existing player stat
	if playerStat.TeamID != nil {
		// Validate Team exists if provided
		team, err := s.teamRepository.GetTeamByID(ctx, *playerStat.TeamID)
		if err != nil {
			return nil, fmt.Errorf("failed to check team: %w", err)
		}
		if team == nil {
			return nil, constants.ErrTeamNotFound
		}
		existingPlayerStat.TeamID = playerStat.TeamID
	}
	if playerStat.Goals != 0 {
		existingPlayerStat.Goals = playerStat.Goals
	}
	if playerStat.Assists != 0 {
		existingPlayerStat.Assists = playerStat.Assists
	}
	if playerStat.Saves != 0 {
		existingPlayerStat.Saves = playerStat.Saves
	}
	if playerStat.YellowCards != 0 {
		existingPlayerStat.YellowCards = playerStat.YellowCards
	}
	if playerStat.RedCards != 0 {
		existingPlayerStat.RedCards = playerStat.RedCards
	}
	if playerStat.Rating != 0 {
		existingPlayerStat.Rating = playerStat.Rating
	}
	if playerStat.MinutesPlayed != 0 {
		existingPlayerStat.MinutesPlayed = playerStat.MinutesPlayed
	}
	if playerStat.Position != "" {
		existingPlayerStat.Position = playerStat.Position
	}
	// Boolean fields need special handling
	existingPlayerStat.IsStarting = playerStat.IsStarting
	existingPlayerStat.IsMVP = playerStat.IsMVP

	// Validate updated player stat
	if !existingPlayerStat.IsValid() {
		return nil, constants.ErrInvalidData
	}

	// Update the player stat
	err = s.playerStatsRepository.UpdatePlayerStat(ctx, id, existingPlayerStat)
	if err != nil {
		return nil, fmt.Errorf("failed to update player stat: %w", err)
	}

	return existingPlayerStat, nil
}

func (s *PlayerStatsDomainService) DeletePlayerStat(ctx context.Context, id uint64) error {
	if id == 0 {
		return constants.ErrInvalidData
	}

	// Check if player stat exists
	existingPlayerStat, err := s.playerStatsRepository.GetPlayerStatByID(ctx, id)
	if err != nil {
		return fmt.Errorf("failed to get player stat by ID: %w", err)
	}
	if existingPlayerStat == nil {
		return constants.ErrRecordNotFound
	}

	return s.playerStatsRepository.DeletePlayerStat(ctx, id)
}
