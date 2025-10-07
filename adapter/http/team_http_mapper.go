package http

import (
	"time"

	"github.com/EdwinRincon/browersfc-api/api/dto"
	"github.com/EdwinRincon/browersfc-api/domain"
)

type TeamHTTPMapper struct{}

func NewTeamHTTPMapper() *TeamHTTPMapper {
	return &TeamHTTPMapper{}
}

// DTO to Domain Conversions (HTTP layer)
func (m *TeamHTTPMapper) DTOToDomain(dto *dto.CreateTeamRequest) *domain.Team {
	if dto == nil {
		return nil
	}

	return &domain.Team{
		FullName:       dto.FullName,
		ShortName:      dto.ShortName,
		PrimaryColor:   dto.PrimaryColor,
		SecondaryColor: dto.SecondaryColor,
		Shield:         dto.Shield,
		NextMatchID:    dto.NextMatchID,
		CreatedAt:      time.Now(),
		UpdatedAt:      time.Now(),
	}
}

func (m *TeamHTTPMapper) UpdateDTOToDomain(dto *dto.UpdateTeamRequest) *domain.Team {
	if dto == nil {
		return nil
	}

	team := &domain.Team{}

	if dto.FullName != nil {
		team.FullName = *dto.FullName
	}
	if dto.ShortName != nil {
		team.ShortName = *dto.ShortName
	}
	if dto.PrimaryColor != nil {
		team.PrimaryColor = *dto.PrimaryColor
	}
	if dto.SecondaryColor != nil {
		team.SecondaryColor = *dto.SecondaryColor
	}
	if dto.Shield != nil {
		team.Shield = *dto.Shield
	}
	if dto.NextMatchID != nil {
		team.NextMatchID = dto.NextMatchID
	}

	return team
}

func (m *TeamHTTPMapper) DomainToDTO(entity *domain.Team) *dto.TeamResponse {
	if entity == nil {
		return nil
	}

	response := &dto.TeamResponse{
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

	// Convert NextMatch if it exists
	if entity.NextMatch != nil {
		response.NextMatch = &dto.MatchShort{
			ID:        entity.NextMatch.ID,
			Status:    entity.NextMatch.Status,
			Kickoff:   entity.NextMatch.Kickoff,
			Location:  entity.NextMatch.Location,
			HomeGoals: entity.NextMatch.HomeGoals,
			AwayGoals: entity.NextMatch.AwayGoals,
		}
	}

	return response
}

func (m *TeamHTTPMapper) DomainListToDTO(entities []domain.Team) []dto.TeamResponse {
	if entities == nil {
		return nil
	}

	result := make([]dto.TeamResponse, len(entities))
	for i, entity := range entities {
		response := m.DomainToDTO(&entity)
		if response != nil {
			result[i] = *response
		}
	}
	return result
}

func (m *TeamHTTPMapper) DomainToShortDTO(entity *domain.Team) *dto.TeamShort {
	if entity == nil {
		return nil
	}

	return &dto.TeamShort{
		ID:        entity.ID,
		FullName:  entity.FullName,
		ShortName: entity.ShortName,
	}
}
