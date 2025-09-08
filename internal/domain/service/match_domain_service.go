package service

import (
	"context"

	"github.com/EdwinRincon/browersfc-api/api/constants"
	"github.com/EdwinRincon/browersfc-api/domain"
)

type MatchDomainService struct {
	matchRepository domain.MatchRepository
}

func NewMatchDomainService(matchRepository domain.MatchRepository) *MatchDomainService {
	return &MatchDomainService{
		matchRepository: matchRepository,
	}
}

func (s *MatchDomainService) CreateMatch(ctx context.Context, match *domain.Match) (*domain.Match, error) {
	if err := s.matchRepository.CreateMatch(ctx, match); err != nil {
		return nil, err
	}
	return match, nil
}

func (s *MatchDomainService) GetMatchByID(ctx context.Context, id uint64) (*domain.Match, error) {
	match, err := s.matchRepository.GetMatchByID(ctx, id)
	if err != nil {
		return nil, err
	}
	if match == nil {
		return nil, constants.ErrRecordNotFound
	}
	return match, nil
}

// GetDetailedMatchByID retrieves a match with all its related data
func (s *MatchDomainService) GetDetailedMatchByID(ctx context.Context, id uint64) (*domain.Match, error) {
	match, err := s.matchRepository.GetDetailedMatchByID(ctx, id)
	if err != nil {
		return nil, err
	}
	if match == nil {
		return nil, constants.ErrRecordNotFound
	}
	return match, nil
}

// GetPaginatedMatches retrieves paginated matches with sorting
func (s *MatchDomainService) GetPaginatedMatches(ctx context.Context, sort string, order string, page int, pageSize int) ([]domain.Match, int64, error) {
	return s.matchRepository.GetPaginatedMatches(ctx, sort, order, page, pageSize)
}

// GetMatchesBySeasonID retrieves matches for a specific season
func (s *MatchDomainService) GetMatchesBySeasonID(ctx context.Context, seasonID uint64, sort string, order string, page int, pageSize int) ([]domain.Match, int64, error) {
	return s.matchRepository.GetMatchesBySeasonID(ctx, seasonID, sort, order, page, pageSize)
}

// GetMatchesByTeamID retrieves matches for a specific team
func (s *MatchDomainService) GetMatchesByTeamID(ctx context.Context, teamID uint64, sort string, order string, page int, pageSize int) ([]domain.Match, int64, error) {
	return s.matchRepository.GetMatchesByTeamID(ctx, teamID, sort, order, page, pageSize)
}

// GetNextMatchByTeamID retrieves the next scheduled match for a team
func (s *MatchDomainService) GetNextMatchByTeamID(ctx context.Context, teamID uint64) (*domain.Match, error) {
	return s.matchRepository.GetNextMatchByTeamID(ctx, teamID)
}

func (s *MatchDomainService) UpdateMatch(ctx context.Context, id uint64, match *domain.Match) (*domain.Match, error) {
	// First check if match exists
	existingMatch, err := s.matchRepository.GetMatchByID(ctx, id)
	if err != nil {
		return nil, err
	}
	if existingMatch == nil {
		return nil, constants.ErrRecordNotFound
	}

	if err := s.matchRepository.UpdateMatch(ctx, id, match); err != nil {
		return nil, err
	}

	// Return the updated match
	return s.matchRepository.GetMatchByID(ctx, id)
}

// DeleteMatch deletes a match by its ID
func (s *MatchDomainService) DeleteMatch(ctx context.Context, id uint64) error {
	// First check if match exists
	existingMatch, err := s.matchRepository.GetMatchByID(ctx, id)
	if err != nil {
		return err
	}
	if existingMatch == nil {
		return constants.ErrRecordNotFound
	}

	return s.matchRepository.DeleteMatch(ctx, id)
}
