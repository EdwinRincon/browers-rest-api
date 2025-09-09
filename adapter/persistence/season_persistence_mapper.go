package persistence

import (
	"github.com/EdwinRincon/browersfc-api/domain"
	"github.com/EdwinRincon/browersfc-api/internal/infrastructure/persistence/model"
)

type SeasonPersistenceMapper struct{}

func NewSeasonPersistenceMapper() *SeasonPersistenceMapper {
	return &SeasonPersistenceMapper{}
}

// Domain to Model Conversions (Infrastructure layer)
func (m *SeasonPersistenceMapper) DomainToModel(entity *domain.Season) *model.Season {
	if entity == nil {
		return nil
	}

	return &model.Season{
		ID:        entity.ID,
		Year:      entity.Year,
		StartDate: entity.StartDate,
		EndDate:   entity.EndDate,
		IsCurrent: entity.IsCurrent,
		CreatedAt: entity.CreatedAt,
		UpdatedAt: entity.UpdatedAt,
	}
}

func (m *SeasonPersistenceMapper) ModelToDomain(model *model.Season) *domain.Season {
	if model == nil {
		return nil
	}

	return &domain.Season{
		ID:        model.ID,
		Year:      model.Year,
		StartDate: model.StartDate,
		EndDate:   model.EndDate,
		IsCurrent: model.IsCurrent,
		CreatedAt: model.CreatedAt,
		UpdatedAt: model.UpdatedAt,
	}
}

func (m *SeasonPersistenceMapper) ModelListToDomain(models []model.Season) []domain.Season {
	if models == nil {
		return nil
	}

	domains := make([]domain.Season, len(models))
	for i, model := range models {
		domain := m.ModelToDomain(&model)
		if domain != nil {
			domains[i] = *domain
		}
	}

	return domains
}
