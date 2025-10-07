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
		ID:             entity.ID,
		FullName:       entity.FullName,
		ShortName:      entity.ShortName,
		PrimaryColor:   entity.PrimaryColor,
		SecondaryColor: entity.SecondaryColor,
		Shield:         entity.Shield,
		NextMatchID:    entity.NextMatchID,
		CreatedAt:      entity.CreatedAt,
		UpdatedAt:      entity.UpdatedAt,
	}
}

func (m *TeamPersistenceMapper) ModelToDomain(model *model.Team) *domain.Team {
	if model == nil {
		return nil
	}

	team := &domain.Team{
		ID:             model.ID,
		FullName:       model.FullName,
		ShortName:      model.ShortName,
		PrimaryColor:   model.PrimaryColor,
		SecondaryColor: model.SecondaryColor,
		Shield:         model.Shield,
		NextMatchID:    model.NextMatchID,
		CreatedAt:      model.CreatedAt,
		UpdatedAt:      model.UpdatedAt,
	}

	// Convert NextMatch if it exists
	if model.NextMatch != nil {
		team.NextMatch = &domain.TeamNextMatch{
			ID:        model.NextMatch.ID,
			Status:    model.NextMatch.Status,
			Kickoff:   model.NextMatch.Kickoff,
			Location:  model.NextMatch.Location,
			HomeGoals: model.NextMatch.HomeGoals,
			AwayGoals: model.NextMatch.AwayGoals,
		}
	}

	return team
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
