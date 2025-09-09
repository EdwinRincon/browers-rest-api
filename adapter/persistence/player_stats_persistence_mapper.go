package persistence

import (
	"github.com/EdwinRincon/browersfc-api/domain"
	"github.com/EdwinRincon/browersfc-api/internal/infrastructure/persistence/model"
)

type PlayerStatsPersistenceMapper struct{}

func NewPlayerStatsPersistenceMapper() *PlayerStatsPersistenceMapper {
	return &PlayerStatsPersistenceMapper{}
}

// Domain to Model Conversions (Infrastructure layer)
func (m *PlayerStatsPersistenceMapper) DomainToModel(entity *domain.PlayerStat) *model.PlayerStat {
	if entity == nil {
		return nil
	}

	return &model.PlayerStat{
		ID:            entity.ID,
		PlayerID:      entity.PlayerID,
		MatchID:       entity.MatchID,
		SeasonID:      entity.SeasonID,
		TeamID:        entity.TeamID,
		Goals:         entity.Goals,
		Assists:       entity.Assists,
		Saves:         entity.Saves,
		YellowCards:   entity.YellowCards,
		RedCards:      entity.RedCards,
		Rating:        entity.Rating,
		IsStarting:    entity.IsStarting,
		MinutesPlayed: entity.MinutesPlayed,
		IsMVP:         entity.IsMVP,
		Position:      entity.Position,
		CreatedAt:     entity.CreatedAt,
		UpdatedAt:     entity.UpdatedAt,
	}
}

func (m *PlayerStatsPersistenceMapper) ModelToDomain(model *model.PlayerStat) *domain.PlayerStat {
	if model == nil {
		return nil
	}

	return &domain.PlayerStat{
		ID:            model.ID,
		PlayerID:      model.PlayerID,
		MatchID:       model.MatchID,
		SeasonID:      model.SeasonID,
		TeamID:        model.TeamID,
		Goals:         model.Goals,
		Assists:       model.Assists,
		Saves:         model.Saves,
		YellowCards:   model.YellowCards,
		RedCards:      model.RedCards,
		Rating:        model.Rating,
		IsStarting:    model.IsStarting,
		MinutesPlayed: model.MinutesPlayed,
		IsMVP:         model.IsMVP,
		Position:      model.Position,
		CreatedAt:     model.CreatedAt,
		UpdatedAt:     model.UpdatedAt,
	}
}

func (m *PlayerStatsPersistenceMapper) ModelListToDomain(models []model.PlayerStat) []domain.PlayerStat {
	if models == nil {
		return nil
	}

	domains := make([]domain.PlayerStat, len(models))
	for i, model := range models {
		domain := m.ModelToDomain(&model)
		if domain != nil {
			domains[i] = *domain
		}
	}

	return domains
}
