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

type LineupService interface {
	CreateLineup(ctx context.Context, lineup *model.Lineup) (*model.Lineup, error)
	GetLineupByID(ctx context.Context, id uint64) (*model.Lineup, error)
	GetLineupsByMatchID(ctx context.Context, matchID uint64) ([]model.Lineup, error)
	GetMatchLineups(ctx context.Context, matchID uint64) (*dto.MatchLineupResponse, error)
	GetLineupsByPlayerID(ctx context.Context, playerID uint64) ([]model.Lineup, error)
	GetPaginatedLineups(ctx context.Context, sort string, order string, page int, pageSize int) ([]model.Lineup, int64, error)
	UpdateLineup(ctx context.Context, id uint64, lineupUpdate *dto.UpdateLineupRequest) (*model.Lineup, error)
	DeleteLineup(ctx context.Context, id uint64) error
}

type lineupService struct {
	LineupRepository repository.LineupRepository
	MatchService     MatchService
}

func NewLineupService(lineupRepo repository.LineupRepository, matchService MatchService) LineupService {
	return &lineupService{
		LineupRepository: lineupRepo,
		MatchService:     matchService,
	}
}

func (s *lineupService) CreateLineup(ctx context.Context, lineup *model.Lineup) (*model.Lineup, error) {
	if err := s.LineupRepository.CreateLineup(ctx, lineup); err != nil {
		return nil, fmt.Errorf("failed to create lineup: %w", err)
	}
	return lineup, nil
}

func (s *lineupService) GetLineupByID(ctx context.Context, id uint64) (*model.Lineup, error) {
	lineup, err := s.LineupRepository.GetLineupByID(ctx, id)
	if err != nil {
		return nil, fmt.Errorf("failed to get lineup by ID: %w", err)
	}
	if lineup == nil {
		return nil, constants.ErrRecordNotFound
	}
	return lineup, nil
}

func (s *lineupService) GetLineupsByMatchID(ctx context.Context, matchID uint64) ([]model.Lineup, error) {
	lineups, err := s.LineupRepository.GetLineupsByMatchID(ctx, matchID)
	if err != nil {
		return nil, fmt.Errorf("failed to get lineups by match ID: %w", err)
	}
	return lineups, nil
}

func (s *lineupService) GetMatchLineups(ctx context.Context, matchID uint64) (*dto.MatchLineupResponse, error) {
	// Get the match details
	match, err := s.MatchService.GetMatchByID(ctx, matchID)
	if err != nil {
		return nil, err
	}
	if match == nil {
		return nil, constants.ErrRecordNotFound
	}

	// Get all lineups for the match
	lineups, err := s.LineupRepository.GetLineupsByMatchID(ctx, matchID)
	if err != nil {
		return nil, fmt.Errorf("failed to get lineups by match ID: %w", err)
	}

	// Organize lineups into starting XI and substitutes
	return mapper.OrganizeLineupsByMatchID(match, lineups), nil
}

func (s *lineupService) GetLineupsByPlayerID(ctx context.Context, playerID uint64) ([]model.Lineup, error) {
	lineups, err := s.LineupRepository.GetLineupsByPlayerID(ctx, playerID)
	if err != nil {
		return nil, fmt.Errorf("failed to get lineups by player ID: %w", err)
	}
	return lineups, nil
}

func (s *lineupService) GetPaginatedLineups(ctx context.Context, sort string, order string, page int, pageSize int) ([]model.Lineup, int64, error) {
	return s.LineupRepository.GetPaginatedLineups(ctx, sort, order, page, pageSize)
}

func (s *lineupService) UpdateLineup(ctx context.Context, id uint64, lineupUpdate *dto.UpdateLineupRequest) (*model.Lineup, error) {
	// Get the existing lineup
	lineup, err := s.LineupRepository.GetLineupByID(ctx, id)
	if err != nil {
		return nil, fmt.Errorf("failed to get lineup by ID: %w", err)
	}
	if lineup == nil {
		return nil, constants.ErrRecordNotFound
	}

	// Update the lineup with the provided DTO
	mapper.UpdateLineupFromDTO(lineup, lineupUpdate)

	// Save the updated lineup
	if err := s.LineupRepository.UpdateLineup(ctx, id, lineup); err != nil {
		return nil, fmt.Errorf("failed to update lineup: %w", err)
	}

	// Retrieve the updated lineup with all relationships loaded
	updatedLineup, err := s.LineupRepository.GetLineupByID(ctx, id)
	if err != nil {
		return nil, fmt.Errorf("failed to get updated lineup: %w", err)
	}

	return updatedLineup, nil
}

func (s *lineupService) DeleteLineup(ctx context.Context, id uint64) error {
	// Check if the lineup exists first
	lineup, err := s.LineupRepository.GetLineupByID(ctx, id)
	if err != nil {
		return fmt.Errorf("failed to check if lineup exists: %w", err)
	}
	if lineup == nil {
		return constants.ErrRecordNotFound
	}

	return s.LineupRepository.DeleteLineup(ctx, id)
}
