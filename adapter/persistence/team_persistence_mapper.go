package persistence

import (
	"github.com/EdwinRincon/browersfc-api/domain"
	"github.com/EdwinRincon/browersfc-api/internal/infrastructure/persistence/model"
)

type TeamPersistenceMapper struct{}

func NewTeamPersistenceMapper() *TeamPersistenceMapper {
	return &TeamPersistenceMapper{}
}

// Domain to Model Conversions (Infrastructure layer)
func (m *TeamPersistenceMapper) DomainToModel(entity *domain.Team) *model.Team {
	if entity == nil {
		return nil
	}

	return &model.Team{
		ID:          entity.ID,
		FullName:    entity.FullName,
		ShortName:   entity.ShortName,
		Color:       entity.Color,
		Color2:      entity.Color2,
		Shield:      entity.Shield,
		NextMatchID: entity.NextMatchID,
		CreatedAt:   entity.CreatedAt,
		UpdatedAt:   entity.UpdatedAt,
	}
}

func (m *TeamPersistenceMapper) ModelToDomain(model *model.Team) *domain.Team {
	if model == nil {
		return nil
	}

	return &domain.Team{
		ID:          model.ID,
		FullName:    model.FullName,
		ShortName:   model.ShortName,
		Color:       model.Color,
		Color2:      model.Color2,
		Shield:      model.Shield,
		NextMatchID: model.NextMatchID,
		CreatedAt:   model.CreatedAt,
		UpdatedAt:   model.UpdatedAt,
	}
}

func (m *TeamPersistenceMapper) ModelListToDomain(models []model.Team) []domain.Team {
	if models == nil {
		return nil
	}

	domains := make([]domain.Team, len(models))
	for i, model := range models {
		domain := m.ModelToDomain(&model)
		if domain != nil {
			domains[i] = *domain
		}
	}

	return domains
}
