package service

import (
	"context"

	"github.com/EdwinRincon/browersfc-api/api/constants"
	"github.com/EdwinRincon/browersfc-api/domain"
)

type LineupDomainService struct {
	lineupRepository domain.LineupRepository
	matchRepository  domain.MatchRepository
	playerRepository domain.PlayerRepository
}

func NewLineupDomainService(
	lineupRepository domain.LineupRepository,
	matchRepository domain.MatchRepository,
	playerRepository domain.PlayerRepository,
) *LineupDomainService {
	return &LineupDomainService{
		lineupRepository: lineupRepository,
		matchRepository:  matchRepository,
		playerRepository: playerRepository,
	}
}

func (s *LineupDomainService) CreateLineup(ctx context.Context, lineup *domain.Lineup) error {
	// Business validation
	if !lineup.IsValid() {
		return constants.ErrInvalidData
	}

	// Verify match exists
	match, err := s.matchRepository.GetMatchByID(ctx, lineup.MatchID)
	if err != nil {
		return err
	}
	if match == nil {
		return constants.ErrMatchNotFound
	}

	// Verify player exists
	player, err := s.playerRepository.GetPlayerByID(ctx, lineup.PlayerID)
	if err != nil {
		return err
	}
	if player == nil {
		return constants.ErrPlayerNotFound
	}

	return s.lineupRepository.CreateLineup(ctx, lineup)
}

func (s *LineupDomainService) GetLineupByID(ctx context.Context, id uint64) (*domain.Lineup, error) {
	if id == 0 {
		return nil, constants.ErrInvalidID
	}

	lineup, err := s.lineupRepository.GetLineupByID(ctx, id)
	if err != nil {
		return nil, err
	}
	if lineup == nil {
		return nil, constants.ErrLineupNotFound
	}

	return lineup, nil
}

func (s *LineupDomainService) UpdateLineup(ctx context.Context, id uint64, lineup *domain.Lineup) error {
	if id == 0 {
		return constants.ErrInvalidID
	}

	if !lineup.IsValid() {
		return constants.ErrInvalidData
	}

	// Verify lineup exists
	existing, err := s.lineupRepository.GetLineupByID(ctx, id)
	if err != nil {
		return err
	}
	if existing == nil {
		return constants.ErrLineupNotFound
	}

	return s.lineupRepository.UpdateLineup(ctx, id, lineup)
}

func (s *LineupDomainService) DeleteLineup(ctx context.Context, id uint64) error {
	if id == 0 {
		return constants.ErrInvalidID
	}

	// Verify lineup exists
	existing, err := s.lineupRepository.GetLineupByID(ctx, id)
	if err != nil {
		return err
	}
	if existing == nil {
		return constants.ErrLineupNotFound
	}

	return s.lineupRepository.DeleteLineup(ctx, id)
}

func (s *LineupDomainService) GetPaginatedLineups(ctx context.Context, sort string, order string, page int, pageSize int) ([]domain.Lineup, int64, error) {
	if page < 0 || pageSize <= 0 {
		return nil, 0, constants.ErrInvalidPaginationParams
	}

	return s.lineupRepository.GetPaginatedLineups(ctx, sort, order, page, pageSize)
}

func (s *LineupDomainService) GetLineupsByMatchID(ctx context.Context, matchID uint64) ([]domain.Lineup, error) {
	if matchID == 0 {
		return nil, constants.ErrInvalidID
	}

	// Verify match exists
	match, err := s.matchRepository.GetMatchByID(ctx, matchID)
	if err != nil {
		return nil, err
	}
	if match == nil {
		return nil, constants.ErrMatchNotFound
	}

	return s.lineupRepository.GetLineupsByMatchID(ctx, matchID)
}

func (s *LineupDomainService) GetStartingLineupsByMatchID(ctx context.Context, matchID uint64) ([]domain.Lineup, error) {
	if matchID == 0 {
		return nil, constants.ErrInvalidID
	}

	// Verify match exists
	match, err := s.matchRepository.GetMatchByID(ctx, matchID)
	if err != nil {
		return nil, err
	}
	if match == nil {
		return nil, constants.ErrMatchNotFound
	}

	return s.lineupRepository.GetStartingLineupsByMatchID(ctx, matchID)
}

func (s *LineupDomainService) GetSubstitutesLineupsByMatchID(ctx context.Context, matchID uint64) ([]domain.Lineup, error) {
	if matchID == 0 {
		return nil, constants.ErrInvalidID
	}

	// Verify match exists
	match, err := s.matchRepository.GetMatchByID(ctx, matchID)
	if err != nil {
		return nil, err
	}
	if match == nil {
		return nil, constants.ErrMatchNotFound
	}

	return s.lineupRepository.GetSubstitutesLineupsByMatchID(ctx, matchID)
}

func (s *LineupDomainService) GetLineupsByPlayerID(ctx context.Context, playerID uint64) ([]domain.Lineup, error) {
	if playerID == 0 {
		return nil, constants.ErrInvalidID
	}

	// Verify player exists
	player, err := s.playerRepository.GetPlayerByID(ctx, playerID)
	if err != nil {
		return nil, err
	}
	if player == nil {
		return nil, constants.ErrPlayerNotFound
	}

	return s.lineupRepository.GetLineupsByPlayerID(ctx, playerID)
}
