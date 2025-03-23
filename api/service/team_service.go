package service

import (
	"context"

	"github.com/EdwinRincon/browersfc-api/api/model"
	"github.com/EdwinRincon/browersfc-api/api/repository"
)

type TeamService interface {
	CreateTeam(ctx context.Context, team *model.Team) error
	GetTeamByID(ctx context.Context, id uint64) (*model.Team, error)
	ListTeams(ctx context.Context, page uint64) ([]*model.Team, error)
	UpdateTeam(ctx context.Context, team *model.Team) error
	DeleteTeam(ctx context.Context, id uint64) error
}

type teamService struct {
	TeamRepository repository.TeamRepository
}

func NewTeamService(teamRepo repository.TeamRepository) TeamService {
	return &teamService{
		TeamRepository: teamRepo,
	}
}

func (s *teamService) CreateTeam(ctx context.Context, team *model.Team) error {
	return s.TeamRepository.CreateTeam(ctx, team)
}

func (s *teamService) GetTeamByID(ctx context.Context, id uint64) (*model.Team, error) {
	return s.TeamRepository.GetTeamByID(ctx, id)
}

func (s *teamService) ListTeams(ctx context.Context, page uint64) ([]*model.Team, error) {
	return s.TeamRepository.ListTeams(ctx, page)
}

func (s *teamService) UpdateTeam(ctx context.Context, team *model.Team) error {
	return s.TeamRepository.UpdateTeam(ctx, team)
}

func (s *teamService) DeleteTeam(ctx context.Context, id uint64) error {
	return s.TeamRepository.DeleteTeam(ctx, id)
}
