package mapper

import (
	"github.com/EdwinRincon/browersfc-api/api/dto"
	"github.com/EdwinRincon/browersfc-api/api/model"
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

// ToSeasonResponseList converts a slice of Season models to a slice of SeasonResponse DTOs
func ToSeasonResponseList(seasons []model.Season) []dto.SeasonResponse {
	seasonResponses := make([]dto.SeasonResponse, len(seasons))
	for i, season := range seasons {
		seasonResponses[i] = *ToSeasonResponse(&season)
	}
	return seasonResponses
}
