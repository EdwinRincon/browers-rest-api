package service

import (
	"context"

	"github.com/EdwinRincon/browersfc-api/api/model"
	"github.com/EdwinRincon/browersfc-api/api/repository"
)

type TeamStatsService interface {
	CreateTeamStat(ctx context.Context, teamStat *model.TeamStat) error
	GetTeamStatByID(ctx context.Context, id uint64) (*model.TeamStat, error)
	ListTeamStats(ctx context.Context, page uint64) ([]*model.TeamStat, error)
	UpdateTeamStat(ctx context.Context, teamStat *model.TeamStat) error
	DeleteTeamStat(ctx context.Context, id uint64) error
}

type teamStatsService struct {
	TeamStatsRepository repository.TeamStatsRepository
}

func NewTeamStatsService(teamStatsRepo repository.TeamStatsRepository) TeamStatsService {
	return &teamStatsService{
		TeamStatsRepository: teamStatsRepo,
	}
}

func (s *teamStatsService) CreateTeamStat(ctx context.Context, teamStat *model.TeamStat) error {
	return s.TeamStatsRepository.CreateTeamStat(ctx, teamStat)
}

func (s *teamStatsService) GetTeamStatByID(ctx context.Context, id uint64) (*model.TeamStat, error) {
	return s.TeamStatsRepository.GetTeamStatByID(ctx, id)
}

func (s *teamStatsService) ListTeamStats(ctx context.Context, page uint64) ([]*model.TeamStat, error) {
	return s.TeamStatsRepository.ListTeamStats(ctx, page)
}

func (s *teamStatsService) UpdateTeamStat(ctx context.Context, teamStat *model.TeamStat) error {
	return s.TeamStatsRepository.UpdateTeamStat(ctx, teamStat)
}

func (s *teamStatsService) DeleteTeamStat(ctx context.Context, id uint64) error {
	return s.TeamStatsRepository.DeleteTeamStat(ctx, id)
}
