package http

import (
	"github.com/EdwinRincon/browersfc-api/api/dto"
	"github.com/EdwinRincon/browersfc-api/domain"
)

type PlayerTeamHTTPMapper struct{}

func NewPlayerTeamHTTPMapper() *PlayerTeamHTTPMapper {
	return &PlayerTeamHTTPMapper{}
}

// DTO to Domain Conversions (HTTP layer)
func (m *PlayerTeamHTTPMapper) DTOToDomain(dto *dto.CreatePlayerTeamRequest) *domain.PlayerTeam {
	if dto == nil {
		return nil
	}

	return &domain.PlayerTeam{
		PlayerID: dto.PlayerID,
		TeamID:   dto.TeamID,
		SeasonID: dto.SeasonID,
	}
}

func (m *PlayerTeamHTTPMapper) DomainToDTO(entity *domain.PlayerTeam) *dto.PlayerTeamResponse {
	if entity == nil {
		return nil
	}

	response := &dto.PlayerTeamResponse{
		ID:        entity.ID,
		PlayerID:  entity.PlayerID,
		TeamID:    entity.TeamID,
		SeasonID:  entity.SeasonID,
		StartDate: entity.StartDate,
		EndDate:   entity.EndDate,
		CreatedAt: entity.CreatedAt,
		UpdatedAt: entity.UpdatedAt,
	}

	// Map related entities if present
	if entity.Player != nil {
		playerMapper := NewPlayerHTTPMapper()
		playerShort := playerMapper.DomainToShortDTO(entity.Player)
		if playerShort != nil {
			response.Player = *playerShort
		}
	}

	if entity.Team != nil {
		teamMapper := NewTeamHTTPMapper()
		teamShort := teamMapper.DomainToShortDTO(entity.Team)
		if teamShort != nil {
			response.Team = *teamShort
		}
	}

	if entity.Season != nil {
		seasonMapper := NewSeasonHTTPMapper()
		seasonShort := seasonMapper.DomainToShortDTO(entity.Season)
		if seasonShort != nil {
			response.Season = *seasonShort
		}
	}

	return response
}

func (m *PlayerTeamHTTPMapper) DomainListToDTO(entities []domain.PlayerTeam) []dto.PlayerTeamResponse {
	if entities == nil {
		return nil
	}

	result := make([]dto.PlayerTeamResponse, len(entities))
	for i, entity := range entities {
		response := m.DomainToDTO(&entity)
		if response != nil {
			result[i] = *response
		}
	}
	return result
}
