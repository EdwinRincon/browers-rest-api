package http

import (
	"github.com/EdwinRincon/browersfc-api/api/dto"
	"github.com/EdwinRincon/browersfc-api/domain"
)

type SeasonHTTPMapper struct{}

func NewSeasonHTTPMapper() *SeasonHTTPMapper {
	return &SeasonHTTPMapper{}
}

// DTO to Domain Conversions (HTTP layer)
func (m *SeasonHTTPMapper) DTOToDomain(dto *dto.CreateSeasonRequest) *domain.Season {
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

func (m *SeasonHTTPMapper) UpdateDTOToDomain(dto *dto.UpdateSeasonRequest, existingSeason *domain.Season) *domain.Season {
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

func (m *SeasonHTTPMapper) DomainToDTO(entity *domain.Season) *dto.SeasonResponse {
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

func (m *SeasonHTTPMapper) DomainListToDTO(entities []domain.Season) []dto.SeasonResponse {
	if entities == nil {
		return nil
	}

	result := make([]dto.SeasonResponse, len(entities))
	for i, entity := range entities {
		response := m.DomainToDTO(&entity)
		if response != nil {
			result[i] = *response
		}
	}
	return result
}

func (m *SeasonHTTPMapper) DomainToShortDTO(entity *domain.Season) *dto.SeasonShort {
	if entity == nil {
		return nil
	}

	return &dto.SeasonShort{
		ID:   entity.ID,
		Year: entity.Year,
	}
}
