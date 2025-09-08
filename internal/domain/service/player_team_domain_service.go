package service

import (
	"context"
	"fmt"

	"github.com/EdwinRincon/browersfc-api/api/constants"
	"github.com/EdwinRincon/browersfc-api/domain"
)

type PlayerTeamDomainService struct {
	playerTeamRepository domain.PlayerTeamRepository
	playerRepository     domain.PlayerRepository
	teamRepository       domain.TeamRepository
	seasonRepository     domain.SeasonRepository
}

func NewPlayerTeamDomainService(
	playerTeamRepository domain.PlayerTeamRepository,
	playerRepository domain.PlayerRepository,
	teamRepository domain.TeamRepository,
	seasonRepository domain.SeasonRepository,
) *PlayerTeamDomainService {
	return &PlayerTeamDomainService{
		playerTeamRepository: playerTeamRepository,
		playerRepository:     playerRepository,
		teamRepository:       teamRepository,
		seasonRepository:     seasonRepository,
	}
}

func (s *PlayerTeamDomainService) CreatePlayerTeam(ctx context.Context, playerTeam *domain.PlayerTeam) (*domain.PlayerTeam, error) {
	// Validate domain entity
	if !playerTeam.IsValid() {
		return nil, constants.ErrInvalidData
	}

	// Validate that the player exists
	player, err := s.playerRepository.GetPlayerByID(ctx, playerTeam.PlayerID)
	if err != nil {
		return nil, fmt.Errorf("failed to check player existence: %w", err)
	}
	if player == nil {
		return nil, constants.ErrPlayerNotFound
	}

	// Validate that the team exists
	team, err := s.teamRepository.GetTeamByID(ctx, playerTeam.TeamID)
	if err != nil {
		return nil, fmt.Errorf("failed to check team existence: %w", err)
	}
	if team == nil {
		return nil, constants.ErrTeamNotFound
	}

	// Validate that the season exists
	season, err := s.seasonRepository.GetSeasonByID(ctx, playerTeam.SeasonID)
	if err != nil {
		return nil, fmt.Errorf("failed to check season existence: %w", err)
	}
	if season == nil {
		return nil, constants.ErrSeasonNotFound
	}

	// Check for overlapping dates using business logic
	overlapData := domain.OverlapCheckData{
		PlayerID:  playerTeam.PlayerID,
		TeamID:    playerTeam.TeamID,
		SeasonID:  playerTeam.SeasonID,
		StartDate: playerTeam.StartDate,
		EndDate:   playerTeam.EndDate,
		IsUpdate:  false,
		ID:        0,
	}

	hasOverlap, err := s.playerTeamRepository.CheckOverlappingDates(ctx, overlapData)
	if err != nil {
		return nil, fmt.Errorf("failed to check date overlaps: %w", err)
	}
	if hasOverlap {
		return nil, constants.ErrOverlappingDates
	}

	// Create the player team relationship
	if err := s.playerTeamRepository.Create(ctx, playerTeam); err != nil {
		return nil, fmt.Errorf("failed to create player team: %w", err)
	}

	return playerTeam, nil
}

func (s *PlayerTeamDomainService) GetPlayerTeamByID(ctx context.Context, id uint64) (*domain.PlayerTeam, error) {
	if id == 0 {
		return nil, constants.ErrInvalidData
	}

	playerTeam, err := s.playerTeamRepository.GetPlayerTeamByID(ctx, id)
	if err != nil {
		return nil, fmt.Errorf("failed to get player team by ID: %w", err)
	}
	if playerTeam == nil {
		return nil, constants.ErrRecordNotFound
	}

	return playerTeam, nil
}

func (s *PlayerTeamDomainService) GetPlayerTeamsByPlayerID(ctx context.Context, playerID uint64) ([]domain.PlayerTeam, error) {
	if playerID == 0 {
		return nil, constants.ErrInvalidData
	}

	return s.playerTeamRepository.GetByPlayerID(ctx, playerID)
}

func (s *PlayerTeamDomainService) GetPlayerTeamsByTeamID(ctx context.Context, teamID uint64) ([]domain.PlayerTeam, error) {
	if teamID == 0 {
		return nil, constants.ErrInvalidData
	}

	return s.playerTeamRepository.GetPlayerTeamsByTeamID(ctx, teamID)
}

func (s *PlayerTeamDomainService) GetPlayerTeamsBySeasonID(ctx context.Context, seasonID uint64) ([]domain.PlayerTeam, error) {
	if seasonID == 0 {
		return nil, constants.ErrInvalidData
	}

	return s.playerTeamRepository.GetPlayerTeamsBySeasonID(ctx, seasonID)
}

func (s *PlayerTeamDomainService) GetPaginatedPlayerTeams(ctx context.Context, sort string, order string, page int, pageSize int) ([]domain.PlayerTeam, int64, error) {
	return s.playerTeamRepository.GetPaginatedPlayerTeams(ctx, sort, order, page, pageSize)
}

func (s *PlayerTeamDomainService) UpdatePlayerTeam(ctx context.Context, id uint64, updateData *domain.PlayerTeam) (*domain.PlayerTeam, error) {
	if id == 0 {
		return nil, constants.ErrInvalidData
	}

	// Get existing player team
	existingPlayerTeam, err := s.GetPlayerTeamByID(ctx, id)
	if err != nil {
		return nil, err
	}

	// Update fields from updateData
	existingPlayerTeam.PlayerID = updateData.PlayerID
	existingPlayerTeam.TeamID = updateData.TeamID
	existingPlayerTeam.SeasonID = updateData.SeasonID
	existingPlayerTeam.StartDate = updateData.StartDate
	existingPlayerTeam.EndDate = updateData.EndDate

	// Validate updated entity
	if !existingPlayerTeam.IsValid() {
		return nil, constants.ErrInvalidData
	}

	// Check for overlapping dates (exclude current record)
	overlapData := domain.OverlapCheckData{
		PlayerID:  existingPlayerTeam.PlayerID,
		TeamID:    existingPlayerTeam.TeamID,
		SeasonID:  existingPlayerTeam.SeasonID,
		StartDate: existingPlayerTeam.StartDate,
		EndDate:   existingPlayerTeam.EndDate,
		IsUpdate:  true,
		ID:        id,
	}

	hasOverlap, err := s.playerTeamRepository.CheckOverlappingDates(ctx, overlapData)
	if err != nil {
		return nil, fmt.Errorf("failed to check date overlaps: %w", err)
	}
	if hasOverlap {
		return nil, constants.ErrOverlappingDates
	}

	// Update the player team
	if err := s.playerTeamRepository.UpdatePlayerTeam(ctx, existingPlayerTeam); err != nil {
		return nil, fmt.Errorf("failed to update player team: %w", err)
	}

	return existingPlayerTeam, nil
}

func (s *PlayerTeamDomainService) DeletePlayerTeam(ctx context.Context, id uint64) error {
	if id == 0 {
		return constants.ErrInvalidData
	}

	// Check if player team exists
	_, err := s.GetPlayerTeamByID(ctx, id)
	if err != nil {
		return err
	}

	return s.playerTeamRepository.DeletePlayerTeam(ctx, id)
}

func (s *PlayerTeamDomainService) DeletePlayerTeamsByPlayerID(ctx context.Context, playerID uint64) error {
	if playerID == 0 {
		return constants.ErrInvalidData
	}

	return s.playerTeamRepository.DeleteByPlayerID(ctx, playerID)
}
