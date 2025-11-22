package persistence

import (
	"github.com/EdwinRincon/browersfc-api/domain"
	"github.com/EdwinRincon/browersfc-api/internal/infrastructure/persistence/model"
)

type TeamStatsPersistenceMapper struct{}

func NewTeamStatsPersistenceMapper() *TeamStatsPersistenceMapper {
	return &TeamStatsPersistenceMapper{}
}

// Domain to Model Conversions (Infrastructure layer)
func (m *TeamStatsPersistenceMapper) DomainToModel(entity *domain.TeamStats) *model.TeamStat {
	if entity == nil {
		return nil
	}

	return &model.TeamStat{
		ID:           entity.ID,
		Wins:         entity.Wins,
		Draws:        entity.Draws,
		Losses:       entity.Losses,
		GoalsFor:     entity.GoalsFor,
		GoalsAgainst: entity.GoalsAgainst,
		Points:       entity.Points,
		Rank:         entity.Rank,
		SeasonID:     entity.SeasonID,
		TeamID:       entity.TeamID,
		CreatedAt:    entity.CreatedAt,
		UpdatedAt:    entity.UpdatedAt,
	}
}

func (m *TeamStatsPersistenceMapper) ModelToDomain(model *model.TeamStat) *domain.TeamStats {
	if model == nil {
		return nil
	}

	var team *domain.Team
	if model.Team != nil {
		teamMapper := NewTeamPersistenceMapper()
		team = teamMapper.ModelToDomain(model.Team)
	}

	var season *domain.Season
	if model.Season != nil {
		seasonMapper := NewSeasonPersistenceMapper()
		season = seasonMapper.ModelToDomain(model.Season)
	}

	return &domain.TeamStats{
		ID:           model.ID,
		Wins:         model.Wins,
		Draws:        model.Draws,
		Losses:       model.Losses,
		GoalsFor:     model.GoalsFor,
		GoalsAgainst: model.GoalsAgainst,
		Points:       model.Points,
		Rank:         model.Rank,
		SeasonID:     model.SeasonID,
		TeamID:       model.TeamID,
		Team:         team,
		Season:       season,
		CreatedAt:    model.CreatedAt,
		UpdatedAt:    model.UpdatedAt,
	}
}

func (m *TeamStatsPersistenceMapper) ModelListToDomain(models []model.TeamStat) []domain.TeamStats {
	if models == nil {
		return nil
	}

	domains := make([]domain.TeamStats, len(models))
	for i, model := range models {
		domain := m.ModelToDomain(&model)
		if domain != nil {
			domains[i] = *domain
		}
	}

	return domains
}
