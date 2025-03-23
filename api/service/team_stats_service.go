package service

import (
	"context"

	"github.com/EdwinRincon/browersfc-api/api/model"
	"github.com/EdwinRincon/browersfc-api/api/repository"
)

type TeamStatsService interface {
	CreateTeamStats(ctx context.Context, teamStats *model.TeamStat) error
	GetTeamStatsByID(ctx context.Context, id uint64) (*model.TeamStat, error)
	ListTeamStats(ctx context.Context, page uint64) ([]*model.TeamStat, error)
	UpdateTeamStats(ctx context.Context, teamStats *model.TeamStat) error
	DeleteTeamStats(ctx context.Context, id uint64) error
}

type teamStatsService struct {
	TeamStatsRepository repository.TeamStatsRepository
}

func NewTeamStatsService(teamStatsRepo repository.TeamStatsRepository) TeamStatsService {
	return &teamStatsService{
		TeamStatsRepository: teamStatsRepo,
	}
}

func (s *teamStatsService) CreateTeamStats(ctx context.Context, teamStats *model.TeamStat) error {
	return s.TeamStatsRepository.CreateTeamStats(ctx, teamStats)
}

func (s *teamStatsService) GetTeamStatsByID(ctx context.Context, id uint64) (*model.TeamStat, error) {
	return s.TeamStatsRepository.GetTeamStatsByID(ctx, id)
}

func (s *teamStatsService) ListTeamStats(ctx context.Context, page uint64) ([]*model.TeamStat, error) {
	return s.TeamStatsRepository.ListTeamStats(ctx, page)
}

func (s *teamStatsService) UpdateTeamStats(ctx context.Context, teamStats *model.TeamStat) error {
	return s.TeamStatsRepository.UpdateTeamStats(ctx, teamStats)
}

func (s *teamStatsService) DeleteTeamStats(ctx context.Context, id uint64) error {
	return s.TeamStatsRepository.DeleteTeamStats(ctx, id)
}
