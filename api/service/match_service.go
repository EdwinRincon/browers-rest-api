package service

import (
	"context"

	"github.com/EdwinRincon/browersfc-api/api/model"
	"github.com/EdwinRincon/browersfc-api/api/repository"
)

type MatchService interface {
	CreateMatch(ctx context.Context, match *model.Matches) error
	GetMatchByID(ctx context.Context, id uint64) (*model.Matches, error)
	ListMatches(ctx context.Context, page, pageSize uint64) ([]*model.Matches, error)
	UpdateMatch(ctx context.Context, match *model.Matches) error
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

func (s *matchService) CreateMatch(ctx context.Context, match *model.Matches) error {
	return s.MatchRepository.CreateMatch(ctx, match)
}

func (s *matchService) GetMatchByID(ctx context.Context, id uint64) (*model.Matches, error) {
	return s.MatchRepository.GetMatchByID(ctx, id)
}

func (s *matchService) ListMatches(ctx context.Context, page, pageSize uint64) ([]*model.Matches, error) {
	return s.MatchRepository.ListMatches(ctx, page, pageSize)
}

func (s *matchService) UpdateMatch(ctx context.Context, match *model.Matches) error {
	return s.MatchRepository.UpdateMatch(ctx, match)
}

func (s *matchService) DeleteMatch(ctx context.Context, id uint64) error {
	return s.MatchRepository.DeleteMatch(ctx, id)
}
