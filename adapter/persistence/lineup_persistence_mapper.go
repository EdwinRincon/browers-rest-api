package persistence

import (
	"github.com/EdwinRincon/browersfc-api/domain"
	"github.com/EdwinRincon/browersfc-api/internal/infrastructure/persistence/model"
)

type LineupPersistenceMapper struct{}

func NewLineupPersistenceMapper() *LineupPersistenceMapper {
	return &LineupPersistenceMapper{}
}

// Domain to Model Conversions (Infrastructure layer)
func (m *LineupPersistenceMapper) DomainToModel(entity *domain.Lineup) *model.Lineup {
	if entity == nil {
		return nil
	}

	return &model.Lineup{
		ID:        entity.ID,
		Position:  entity.Position,
		PlayerID:  entity.PlayerID,
		MatchID:   entity.MatchID,
		Starting:  entity.Starting,
		CreatedAt: entity.CreatedAt,
		UpdatedAt: entity.UpdatedAt,
	}
}

func (m *LineupPersistenceMapper) ModelToDomain(model *model.Lineup) *domain.Lineup {
	if model == nil {
		return nil
	}

	return &domain.Lineup{
		ID:        model.ID,
		Position:  model.Position,
		PlayerID:  model.PlayerID,
		MatchID:   model.MatchID,
		Starting:  model.Starting,
		CreatedAt: model.CreatedAt,
		UpdatedAt: model.UpdatedAt,
	}
}

func (m *LineupPersistenceMapper) ModelListToDomain(models []model.Lineup) []domain.Lineup {
	if models == nil {
		return nil
	}

	domains := make([]domain.Lineup, len(models))
	for i, model := range models {
		domain := m.ModelToDomain(&model)
		if domain != nil {
			domains[i] = *domain
		}
	}

	return domains
}
