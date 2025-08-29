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

type MatchService interface {
	CreateMatch(ctx context.Context, matchDTO *dto.CreateMatchRequest) (*model.Match, error)
	GetMatchByID(ctx context.Context, id uint64) (*model.Match, error)
	GetDetailedMatchByID(ctx context.Context, id uint64) (*model.Match, error)
	GetPaginatedMatches(ctx context.Context, sort string, order string, page int, pageSize int) ([]model.Match, int64, error)
	GetMatchesBySeasonID(ctx context.Context, seasonID uint64, sort string, order string, page int, pageSize int) ([]model.Match, int64, error)
	GetMatchesByTeamID(ctx context.Context, teamID uint64, sort string, order string, page int, pageSize int) ([]model.Match, int64, error)
	GetNextMatchByTeamID(ctx context.Context, teamID uint64) (*model.Match, error)
	UpdateMatch(ctx context.Context, id uint64, matchDTO *dto.UpdateMatchRequest) (*model.Match, error)
	DeleteMatch(ctx context.Context, id uint64) error
}

type matchService struct {
	MatchRepository repository.MatchRepository
}

func NewMatchService(matchRepo repository.MatchRepository) MatchService {
	return &matchService{
		MatchRepository: matchRepo,
	}
}

// CreateMatch creates a new match
func (s *matchService) CreateMatch(ctx context.Context, matchDTO *dto.CreateMatchRequest) (*model.Match, error) {
	match := mapper.ToMatch(matchDTO)

	if err := s.MatchRepository.CreateMatch(ctx, match); err != nil {
		return nil, fmt.Errorf("failed to create match: %w", err)
	}

	return match, nil
}

// GetMatchByID retrieves a match by its ID
func (s *matchService) GetMatchByID(ctx context.Context, id uint64) (*model.Match, error) {
	match, err := s.MatchRepository.GetActiveMatchByID(ctx, id)
	if err != nil {
		return nil, err
	}
	if match == nil {
		return nil, constants.ErrRecordNotFound
	}
	return match, nil
}

// GetDetailedMatchByID retrieves a match with all its related data
func (s *matchService) GetDetailedMatchByID(ctx context.Context, id uint64) (*model.Match, error) {
	match, err := s.MatchRepository.GetDetailedMatchByID(ctx, id)
	if err != nil {
		return nil, err
	}
	if match == nil {
		return nil, constants.ErrRecordNotFound
	}
	return match, nil
}

// GetPaginatedMatches retrieves paginated matches with sorting
func (s *matchService) GetPaginatedMatches(ctx context.Context, sort string, order string, page int, pageSize int) ([]model.Match, int64, error) {
	return s.MatchRepository.GetPaginatedMatches(ctx, sort, order, page, pageSize)
}

// GetMatchesBySeasonID retrieves matches for a specific season
func (s *matchService) GetMatchesBySeasonID(ctx context.Context, seasonID uint64, sort string, order string, page int, pageSize int) ([]model.Match, int64, error) {
	return s.MatchRepository.GetMatchesBySeasonID(ctx, seasonID, sort, order, page, pageSize)
}

// GetMatchesByTeamID retrieves matches for a specific team
func (s *matchService) GetMatchesByTeamID(ctx context.Context, teamID uint64, sort string, order string, page int, pageSize int) ([]model.Match, int64, error) {
	return s.MatchRepository.GetMatchesByTeamID(ctx, teamID, sort, order, page, pageSize)
}

// GetNextMatchByTeamID retrieves the next scheduled match for a team
func (s *matchService) GetNextMatchByTeamID(ctx context.Context, teamID uint64) (*model.Match, error) {
	return s.MatchRepository.GetNextMatchByTeamID(ctx, teamID)
}

// UpdateMatch updates an existing match
func (s *matchService) UpdateMatch(ctx context.Context, id uint64, matchDTO *dto.UpdateMatchRequest) (*model.Match, error) {
	// First get the existing match
	existingMatch, err := s.MatchRepository.GetActiveMatchByID(ctx, id)
	if err != nil {
		return nil, fmt.Errorf("failed to get match by ID: %w", err)
	}
	if existingMatch == nil {
		return nil, constants.ErrRecordNotFound
	}

	// Update the match with the new data
	mapper.UpdateMatchFromDTO(existingMatch, matchDTO)

	// Save the updated match
	if err := s.MatchRepository.UpdateMatch(ctx, id, existingMatch); err != nil {
		return nil, fmt.Errorf("failed to update match: %w", err)
	}

	// Reload the match to get the updated version with associations
	updatedMatch, err := s.MatchRepository.GetActiveMatchByID(ctx, id)
	if err != nil {
		return nil, fmt.Errorf("failed to reload match after update: %w", err)
	}

	return updatedMatch, nil
}

// DeleteMatch deletes a match by its ID
func (s *matchService) DeleteMatch(ctx context.Context, id uint64) error {
	return s.MatchRepository.DeleteMatch(ctx, id)
}
