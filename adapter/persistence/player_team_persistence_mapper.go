package persistence

import (
	"github.com/EdwinRincon/browersfc-api/domain"
	"github.com/EdwinRincon/browersfc-api/internal/infrastructure/persistence/model"
)

type PlayerTeamPersistenceMapper struct{}

func NewPlayerTeamPersistenceMapper() *PlayerTeamPersistenceMapper {
	return &PlayerTeamPersistenceMapper{}
}

// Domain to Model Conversions (Infrastructure layer)
func (m *PlayerTeamPersistenceMapper) DomainToModel(entity *domain.PlayerTeam) *model.PlayerTeam {
	if entity == nil {
		return nil
	}

	return &model.PlayerTeam{
		ID:        entity.ID,
		PlayerID:  entity.PlayerID,
		TeamID:    entity.TeamID,
		SeasonID:  entity.SeasonID,
		CreatedAt: entity.CreatedAt,
		UpdatedAt: entity.UpdatedAt,
	}
}

func (m *PlayerTeamPersistenceMapper) ModelToDomain(model *model.PlayerTeam) *domain.PlayerTeam {
	if model == nil {
		return nil
	}

	var player *domain.Player
	if model.Player != nil {
		playerMapper := NewPlayerPersistenceMapper()
		player = playerMapper.ModelToDomain(model.Player)
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

	return &domain.PlayerTeam{
		ID:        model.ID,
		PlayerID:  model.PlayerID,
		TeamID:    model.TeamID,
		SeasonID:  model.SeasonID,
		Player:    player,
		Team:      team,
		Season:    season,
		StartDate: model.StartDate,
		EndDate:   model.EndDate,
		CreatedAt: model.CreatedAt,
		UpdatedAt: model.UpdatedAt,
	}
}

func (m *PlayerTeamPersistenceMapper) ModelListToDomain(models []model.PlayerTeam) []domain.PlayerTeam {
	if models == nil {
		return nil
	}

	domains := make([]domain.PlayerTeam, len(models))
	for i, model := range models {
		domain := m.ModelToDomain(&model)
		if domain != nil {
			domains[i] = *domain
		}
	}

	return domains
}
