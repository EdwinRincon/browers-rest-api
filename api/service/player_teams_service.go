package service

import (
	"context"
	"fmt"

	"github.com/EdwinRincon/browersfc-api/api/constants"
	"github.com/EdwinRincon/browersfc-api/api/dto"
	"github.com/EdwinRincon/browersfc-api/api/mapper"
	"github.com/EdwinRincon/browersfc-api/api/repository"
)

type PlayerTeamService interface {
	CreatePlayerTeam(ctx context.Context, createRequest *dto.CreatePlayerTeamRequest) (*dto.PlayerTeamResponse, error)
	GetPlayerTeamByID(ctx context.Context, id uint64) (*dto.PlayerTeamResponse, error)
	GetPlayerTeamsByPlayerID(ctx context.Context, playerID uint64) ([]dto.PlayerTeamResponse, error)
	GetPlayerTeamsByTeamID(ctx context.Context, teamID uint64) ([]dto.PlayerTeamResponse, error)
	GetPlayerTeamsBySeasonID(ctx context.Context, seasonID uint64) ([]dto.PlayerTeamResponse, error)
	GetPaginatedPlayerTeams(ctx context.Context, sort string, order string, page int, pageSize int) ([]dto.PlayerTeamResponse, int64, error)
	UpdatePlayerTeam(ctx context.Context, id uint64, updateRequest *dto.UpdatePlayerTeamRequest) (*dto.PlayerTeamResponse, error)
	DeletePlayerTeam(ctx context.Context, id uint64) error
}

type PlayerTeamServiceImpl struct {
	PlayerTeamRepository repository.PlayerTeamRepository
	PlayerRepository     repository.PlayerRepository
	TeamRepository       repository.TeamRepository
	SeasonRepository     repository.SeasonRepository
}

func NewPlayerTeamService(
	playerTeamRepo repository.PlayerTeamRepository,
	playerRepo repository.PlayerRepository,
	teamRepo repository.TeamRepository,
	seasonRepo repository.SeasonRepository,
) PlayerTeamService {
	return &PlayerTeamServiceImpl{
		PlayerTeamRepository: playerTeamRepo,
		PlayerRepository:     playerRepo,
		TeamRepository:       teamRepo,
		SeasonRepository:     seasonRepo,
	}
}

// CreatePlayerTeam creates a new player-team relationship
func (s *PlayerTeamServiceImpl) CreatePlayerTeam(ctx context.Context, createRequest *dto.CreatePlayerTeamRequest) (*dto.PlayerTeamResponse, error) {
	// Validate that the player, team, and season exist
	player, err := s.PlayerRepository.GetActivePlayerByID(ctx, createRequest.PlayerID)
	if err != nil || player == nil {
		return nil, constants.ErrPlayerNotFound
	}

	team, err := s.TeamRepository.GetTeamByID(ctx, createRequest.TeamID)
	if err != nil || team == nil {
		return nil, constants.ErrTeamNotFound
	}

	season, err := s.SeasonRepository.GetActiveSeasonByID(ctx, createRequest.SeasonID)
	if err != nil || season == nil {
		return nil, constants.ErrSeasonNotFound
	}

	// Check for date overlaps
	hasOverlap, err := s.PlayerTeamRepository.CheckOverlappingDates(
		ctx,
		repository.OverlapCheckData{
			PlayerID:  createRequest.PlayerID,
			TeamID:    createRequest.TeamID,
			SeasonID:  createRequest.SeasonID,
			StartDate: createRequest.StartDate,
			EndDate:   createRequest.EndDate,
			IsUpdate:  false,
			ID:        0, // Zero ID for new records
		},
	)
	if err != nil {
		return nil, fmt.Errorf("failed to check date overlaps: %w", err)
	}
	if hasOverlap {
		return nil, fmt.Errorf("player already has an overlapping relationship with this team during the specified period: %w",
			constants.ErrOverlappingDates)
	}

	// Create the player-team relationship
	playerTeam := mapper.ToPlayerTeam(createRequest)

	if err := s.PlayerTeamRepository.Create(ctx, playerTeam); err != nil {
		return nil, fmt.Errorf("failed to create player team relationship: %w", err)
	}

	// Get the full record with associations
	createdPlayerTeam, err := s.PlayerTeamRepository.GetPlayerTeamByID(ctx, playerTeam.ID)
	if err != nil {
		return nil, fmt.Errorf("failed to retrieve created player team: %w", err)
	}

	return mapper.ToPlayerTeamResponse(createdPlayerTeam), nil
}

// GetPlayerTeamByID retrieves a player-team relationship by ID
func (s *PlayerTeamServiceImpl) GetPlayerTeamByID(ctx context.Context, id uint64) (*dto.PlayerTeamResponse, error) {
	playerTeam, err := s.PlayerTeamRepository.GetPlayerTeamByID(ctx, id)
	if err != nil {
		return nil, fmt.Errorf("failed to get player team relationship: %w", err)
	}
	if playerTeam == nil {
		return nil, constants.ErrRecordNotFound
	}

	return mapper.ToPlayerTeamResponse(playerTeam), nil
}

// GetPlayerTeamsByPlayerID retrieves all team relationships for a player
func (s *PlayerTeamServiceImpl) GetPlayerTeamsByPlayerID(ctx context.Context, playerID uint64) ([]dto.PlayerTeamResponse, error) {
	// Validate player exists
	player, err := s.PlayerRepository.GetActivePlayerByID(ctx, playerID)
	if err != nil || player == nil {
		return nil, constants.ErrPlayerNotFound
	}

	playerTeams, err := s.PlayerTeamRepository.GetByPlayerID(ctx, playerID)
	if err != nil {
		return nil, fmt.Errorf("failed to get player team relationships: %w", err)
	}

	return mapper.ToPlayerTeamResponseList(playerTeams), nil
}

// GetPlayerTeamsByTeamID retrieves all player relationships for a team
func (s *PlayerTeamServiceImpl) GetPlayerTeamsByTeamID(ctx context.Context, teamID uint64) ([]dto.PlayerTeamResponse, error) {
	// Validate team exists
	team, err := s.TeamRepository.GetTeamByID(ctx, teamID)
	if err != nil || team == nil {
		return nil, constants.ErrTeamNotFound
	}

	playerTeams, err := s.PlayerTeamRepository.GetPlayerTeamsByTeamID(ctx, teamID)
	if err != nil {
		return nil, fmt.Errorf("failed to get player team relationships: %w", err)
	}

	return mapper.ToPlayerTeamResponseList(playerTeams), nil
}

// GetPlayerTeamsBySeasonID retrieves all player-team relationships for a season
func (s *PlayerTeamServiceImpl) GetPlayerTeamsBySeasonID(ctx context.Context, seasonID uint64) ([]dto.PlayerTeamResponse, error) {
	// Validate season exists
	season, err := s.SeasonRepository.GetActiveSeasonByID(ctx, seasonID)
	if err != nil || season == nil {
		return nil, constants.ErrSeasonNotFound
	}

	playerTeams, err := s.PlayerTeamRepository.GetPlayerTeamsBySeasonID(ctx, seasonID)
	if err != nil {
		return nil, fmt.Errorf("failed to get player team relationships: %w", err)
	}

	return mapper.ToPlayerTeamResponseList(playerTeams), nil
}

// GetPaginatedPlayerTeams retrieves a paginated list of player-team relationships
func (s *PlayerTeamServiceImpl) GetPaginatedPlayerTeams(ctx context.Context, sort string, order string, page int, pageSize int) ([]dto.PlayerTeamResponse, int64, error) {
	playerTeams, total, err := s.PlayerTeamRepository.GetPaginatedPlayerTeams(ctx, sort, order, page, pageSize)
	if err != nil {
		return nil, 0, fmt.Errorf("failed to get paginated player teams: %w", err)
	}

	return mapper.ToPlayerTeamResponseList(playerTeams), total, nil
}

// UpdatePlayerTeam updates an existing player-team relationship
func (s *PlayerTeamServiceImpl) UpdatePlayerTeam(ctx context.Context, id uint64, updateRequest *dto.UpdatePlayerTeamRequest) (*dto.PlayerTeamResponse, error) {
	// Check if the player-team relationship exists
	existingPlayerTeam, err := s.PlayerTeamRepository.GetPlayerTeamByID(ctx, id)
	if err != nil {
		return nil, fmt.Errorf("failed to get player team relationship: %w", err)
	}
	if existingPlayerTeam == nil {
		return nil, constants.ErrRecordNotFound
	}

	// Prepare the update
	if updateRequest.StartDate != nil {
		existingPlayerTeam.StartDate = *updateRequest.StartDate
	}
	if updateRequest.EndDate != nil {
		existingPlayerTeam.EndDate = updateRequest.EndDate
	}

	// Check for date overlaps (exclude the current record being updated)
	hasOverlap, err := s.PlayerTeamRepository.CheckOverlappingDates(
		ctx,
		repository.OverlapCheckData{
			PlayerID:  existingPlayerTeam.PlayerID,
			TeamID:    existingPlayerTeam.TeamID,
			SeasonID:  existingPlayerTeam.SeasonID,
			StartDate: existingPlayerTeam.StartDate,
			EndDate:   existingPlayerTeam.EndDate,
			IsUpdate:  true,
			ID:        existingPlayerTeam.ID,
		},
	)
	if err != nil {
		return nil, fmt.Errorf("failed to check date overlaps: %w", err)
	}
	if hasOverlap {
		return nil, constants.ErrOverlappingDates
	}

	if err := s.PlayerTeamRepository.UpdatePlayerTeam(ctx, existingPlayerTeam); err != nil {
		return nil, fmt.Errorf("failed to update player team relationship: %w", err)
	}

	// Get the updated record with associations
	updatedPlayerTeam, err := s.PlayerTeamRepository.GetPlayerTeamByID(ctx, id)
	if err != nil {
		return nil, fmt.Errorf("failed to retrieve updated player team: %w", err)
	}

	return mapper.ToPlayerTeamResponse(updatedPlayerTeam), nil
}

// DeletePlayerTeam deletes a player-team relationship
func (s *PlayerTeamServiceImpl) DeletePlayerTeam(ctx context.Context, id uint64) error {
	// Check if the player-team relationship exists
	existingPlayerTeam, err := s.PlayerTeamRepository.GetPlayerTeamByID(ctx, id)
	if err != nil {
		return fmt.Errorf("failed to get player team relationship: %w", err)
	}
	if existingPlayerTeam == nil {
		return constants.ErrRecordNotFound
	}

	if err := s.PlayerTeamRepository.DeletePlayerTeam(ctx, id); err != nil {
		return fmt.Errorf("failed to delete player team relationship: %w", err)
	}

	return nil
}
