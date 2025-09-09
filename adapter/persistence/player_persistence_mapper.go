package persistence

import (
	"github.com/EdwinRincon/browersfc-api/domain"
	"github.com/EdwinRincon/browersfc-api/internal/infrastructure/persistence/model"
)

type PlayerPersistenceMapper struct{}

func NewPlayerPersistenceMapper() *PlayerPersistenceMapper {
	return &PlayerPersistenceMapper{}
}

// Domain to Model Conversions (Infrastructure layer)
func (m *PlayerPersistenceMapper) DomainToModel(entity *domain.Player) *model.Player {
	if entity == nil {
		return nil
	}

	return &model.Player{
		ID:            entity.ID,
		NickName:      entity.NickName,
		Height:        entity.Height,
		Country:       entity.Country,
		Country2:      entity.Country2,
		Foot:          entity.Foot,
		Age:           entity.Age,
		SquadNumber:   entity.SquadNumber,
		Rating:        entity.Rating,
		Matches:       entity.Matches,
		YCards:        entity.YCards,
		RCards:        entity.RCards,
		Goals:         entity.Goals,
		Assists:       entity.Assists,
		Saves:         entity.Saves,
		Position:      entity.Position,
		Injured:       entity.Injured,
		CareerSummary: entity.CareerSummary,
		MVPCount:      entity.MVPCount,
		UserID:        entity.UserID,
		CreatedAt:     entity.CreatedAt,
		UpdatedAt:     entity.UpdatedAt,
	}
}

func (m *PlayerPersistenceMapper) ModelToDomain(model *model.Player) *domain.Player {
	if model == nil {
		return nil
	}

	return &domain.Player{
		ID:            model.ID,
		NickName:      model.NickName,
		Height:        model.Height,
		Country:       model.Country,
		Country2:      model.Country2,
		Foot:          model.Foot,
		Age:           model.Age,
		SquadNumber:   model.SquadNumber,
		Rating:        model.Rating,
		Matches:       model.Matches,
		YCards:        model.YCards,
		RCards:        model.RCards,
		Goals:         model.Goals,
		Assists:       model.Assists,
		Saves:         model.Saves,
		Position:      model.Position,
		Injured:       model.Injured,
		CareerSummary: model.CareerSummary,
		MVPCount:      model.MVPCount,
		UserID:        model.UserID,
		CreatedAt:     model.CreatedAt,
		UpdatedAt:     model.UpdatedAt,
	}
}

func (m *PlayerPersistenceMapper) ModelListToDomain(models []model.Player) []domain.Player {
	if models == nil {
		return nil
	}

	result := make([]domain.Player, len(models))
	for i, model := range models {
		domain := m.ModelToDomain(&model)
		if domain != nil {
			result[i] = *domain
		}
	}
	return result
}
