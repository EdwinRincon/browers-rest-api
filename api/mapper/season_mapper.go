package mapper

import (
	"github.com/EdwinRincon/browersfc-api/api/dto"
	"github.com/EdwinRincon/browersfc-api/api/model"
	"github.com/EdwinRincon/browersfc-api/domain"
)

// ToSeason converts a CreateSeasonRequest DTO to a Season model
func ToSeason(dto *dto.CreateSeasonRequest) *model.Season {
	return &model.Season{
		Year:      dto.Year,
		StartDate: dto.StartDate,
		EndDate:   dto.EndDate,
		IsCurrent: dto.IsCurrent,
	}
}

// UpdateSeasonFromDTO updates an existing Season model with data from an UpdateSeasonRequest DTO
func UpdateSeasonFromDTO(season *model.Season, dto *dto.UpdateSeasonRequest) {
	if dto.Year != nil {
		season.Year = *dto.Year
	}
	if dto.StartDate != nil {
		season.StartDate = *dto.StartDate
	}
	if dto.EndDate != nil {
		season.EndDate = *dto.EndDate
	}
	if dto.IsCurrent != nil {
		season.IsCurrent = *dto.IsCurrent
	}
}

// ToSeasonResponse converts a Season model to a SeasonResponse DTO
func ToSeasonResponse(season *model.Season) *dto.SeasonResponse {
	return &dto.SeasonResponse{
		ID:        season.ID,
		Year:      season.Year,
		StartDate: season.StartDate,
		EndDate:   season.EndDate,
		IsCurrent: season.IsCurrent,
		CreatedAt: season.CreatedAt,
		UpdatedAt: season.UpdatedAt,
	}
}

// DomainSeasonToResponse converts a domain Season directly to SeasonResponse DTO
// This avoids coupling the HTTP layer to persistence structure
func DomainSeasonToResponse(domainSeason *domain.Season) *dto.SeasonResponse {
	if domainSeason == nil {
		return nil
	}

	return &dto.SeasonResponse{
		ID:        domainSeason.ID,
		Year:      domainSeason.Year,
		StartDate: domainSeason.StartDate,
		EndDate:   domainSeason.EndDate,
		IsCurrent: domainSeason.IsCurrent,
		CreatedAt: domainSeason.CreatedAt,
		UpdatedAt: domainSeason.UpdatedAt,
	}
}

// CreateRequestToDomain converts CreateSeasonRequest DTO to domain Season
func CreateRequestToDomain(createRequest *dto.CreateSeasonRequest) *domain.Season {
	if createRequest == nil {
		return nil
	}

	return &domain.Season{
		Year:      createRequest.Year,
		StartDate: createRequest.StartDate,
		EndDate:   createRequest.EndDate,
		IsCurrent: createRequest.IsCurrent,
	}
}

// UpdateRequestToDomain converts UpdateSeasonRequest DTO to domain Season
// Note: This creates a new domain season with only the provided fields set
func UpdateRequestToDomain(updateRequest *dto.UpdateSeasonRequest, existingSeason *domain.Season) *domain.Season {
	if updateRequest == nil || existingSeason == nil {
		return nil
	}

	// Start with a copy of the existing season
	updatedSeason := *existingSeason

	// Apply updates from the request
	if updateRequest.Year != nil {
		updatedSeason.Year = *updateRequest.Year
	}
	if updateRequest.StartDate != nil {
		updatedSeason.StartDate = *updateRequest.StartDate
	}
	if updateRequest.EndDate != nil {
		updatedSeason.EndDate = *updateRequest.EndDate
	}
	if updateRequest.IsCurrent != nil {
		updatedSeason.IsCurrent = *updateRequest.IsCurrent
	}

	return &updatedSeason
}

// DomainSeasonListToResponse converts a slice of domain Seasons to SeasonResponse DTOs
func DomainSeasonListToResponse(domainSeasons []domain.Season) []dto.SeasonResponse {
	responses := make([]dto.SeasonResponse, len(domainSeasons))
	for i, domainSeason := range domainSeasons {
		responses[i] = *DomainSeasonToResponse(&domainSeason)
	}
	return responses
}
