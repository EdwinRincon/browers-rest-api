package mapper

import (
	"github.com/EdwinRincon/browersfc-api/api/dto"
	"github.com/EdwinRincon/browersfc-api/api/model"
	"github.com/EdwinRincon/browersfc-api/domain"
)

type SeasonMapper struct{}

func NewSeasonMapper() *SeasonMapper {
	return &SeasonMapper{}
}

// DTO to Domain Conversions (HTTP layer)
func (m *SeasonMapper) DTOToDomain(dto *dto.CreateSeasonRequest) *domain.Season {
	if dto == nil {
		return nil
	}

	return &domain.Season{
		Year:      dto.Year,
		StartDate: dto.StartDate,
		EndDate:   dto.EndDate,
		IsCurrent: dto.IsCurrent,
	}
}

func (m *SeasonMapper) UpdateDTOToDomain(dto *dto.UpdateSeasonRequest, existingSeason *domain.Season) *domain.Season {
	if dto == nil || existingSeason == nil {
		return nil
	}

	// Start with a copy of the existing season
	updatedSeason := *existingSeason

	// Apply updates from the request
	if dto.Year != nil {
		updatedSeason.Year = *dto.Year
	}
	if dto.StartDate != nil {
		updatedSeason.StartDate = *dto.StartDate
	}
	if dto.EndDate != nil {
		updatedSeason.EndDate = *dto.EndDate
	}
	if dto.IsCurrent != nil {
		updatedSeason.IsCurrent = *dto.IsCurrent
	}

	return &updatedSeason
}

func (m *SeasonMapper) DomainToDTO(entity *domain.Season) *dto.SeasonResponse {
	if entity == nil {
		return nil
	}

	return &dto.SeasonResponse{
		ID:        entity.ID,
		Year:      entity.Year,
		StartDate: entity.StartDate,
		EndDate:   entity.EndDate,
		IsCurrent: entity.IsCurrent,
		CreatedAt: entity.CreatedAt,
		UpdatedAt: entity.UpdatedAt,
	}
}

func (m *SeasonMapper) DomainListToDTO(entities []domain.Season) []dto.SeasonResponse {
	responses := make([]dto.SeasonResponse, len(entities))
	for i, entity := range entities {
		responses[i] = *m.DomainToDTO(&entity)
	}
	return responses
}

// Domain to Model Conversions (Infrastructure layer)
func (m *SeasonMapper) DomainToModel(entity *domain.Season) *model.Season {
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
		// Note: Relations (Matches, Articles, TeamStats, PlayerTeams, PlayerStats)
		// are not mapped as they should be handled separately when needed
	}
}

func (m *SeasonMapper) ModelToDomain(model *model.Season) *domain.Season {
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

func (m *SeasonMapper) DomainListToModel(entities []domain.Season) []model.Season {
	models := make([]model.Season, len(entities))
	for i, entity := range entities {
		models[i] = *m.DomainToModel(&entity)
	}
	return models
}

func (m *SeasonMapper) ModelListToDomain(models []model.Season) []domain.Season {
	entities := make([]domain.Season, len(models))
	for i, model := range models {
		entities[i] = *m.ModelToDomain(&model)
	}
	return entities
}

// ModelToShortDTO converts Season model to SeasonShort DTO for cross-entity relationships
func (m *SeasonMapper) ModelToShortDTO(model *model.Season) *dto.SeasonShort {
	if model == nil {
		return nil
	}

	return &dto.SeasonShort{
		ID:   model.ID,
		Year: model.Year,
	}
}
