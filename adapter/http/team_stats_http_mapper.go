package http

import (
	"github.com/EdwinRincon/browersfc-api/api/dto"
	"github.com/EdwinRincon/browersfc-api/domain"
)

type TeamStatsHTTPMapper struct{}

func NewTeamStatsHTTPMapper() *TeamStatsHTTPMapper {
	return &TeamStatsHTTPMapper{}
}

// DTO to Domain Conversions (HTTP layer)
func (m *TeamStatsHTTPMapper) CreateRequestToDomain(request dto.CreateTeamStatsRequest) *domain.TeamStats {
	return &domain.TeamStats{
		Wins:         request.Wins,
		Draws:        request.Draws,
		Losses:       request.Losses,
		GoalsFor:     request.GoalsFor,
		GoalsAgainst: request.GoalsAgainst,
		Points:       request.Points,
		Rank:         request.Rank,
		SeasonID:     request.SeasonID,
		TeamID:       request.TeamID,
	}
}

func (m *TeamStatsHTTPMapper) UpdateRequestToDomain(request dto.UpdateTeamStatsRequest, existing *domain.TeamStats) *domain.TeamStats {
	updated := *existing // Create a copy

	if request.Wins != nil {
		updated.Wins = *request.Wins
	}
	if request.Draws != nil {
		updated.Draws = *request.Draws
	}
	if request.Losses != nil {
		updated.Losses = *request.Losses
	}
	if request.GoalsFor != nil {
		updated.GoalsFor = *request.GoalsFor
	}
	if request.GoalsAgainst != nil {
		updated.GoalsAgainst = *request.GoalsAgainst
	}
	if request.Points != nil {
		updated.Points = *request.Points
	}
	if request.Rank != nil {
		updated.Rank = *request.Rank
	}
	if request.SeasonID != nil {
		updated.SeasonID = *request.SeasonID
	}
	if request.TeamID != nil {
		updated.TeamID = *request.TeamID
	}

	return &updated
}

// Domain to DTO Conversions (HTTP layer)
func (m *TeamStatsHTTPMapper) DomainToDTO(entity *domain.TeamStats) dto.TeamStatsResponse {
	response := dto.TeamStatsResponse{
		ID:           entity.ID,
		Wins:         entity.Wins,
		Draws:        entity.Draws,
		Losses:       entity.Losses,
		GoalsFor:     entity.GoalsFor,
		GoalsAgainst: entity.GoalsAgainst,
		Points:       entity.Points,
		Rank:         entity.Rank,
		SeasonID:     entity.SeasonID,
		TeamID:       entity.TeamID,
		CreatedAt:    entity.CreatedAt,
		UpdatedAt:    entity.UpdatedAt,
	}

	// Include related entities if they exist
	if entity.Team != nil {
		teamMapper := NewTeamHTTPMapper()
		teamShort := teamMapper.DomainToShortDTO(entity.Team)
		if teamShort != nil {
			response.Team = teamShort
		}
	}

	if entity.Season != nil {
		seasonMapper := NewSeasonHTTPMapper()
		seasonShort := seasonMapper.DomainToShortDTO(entity.Season)
		if seasonShort != nil {
			response.Season = seasonShort
		}
	}

	return response
}

func (m *TeamStatsHTTPMapper) DomainToShortDTO(entity *domain.TeamStats) *dto.TeamStatsShort {
	return &dto.TeamStatsShort{
		ID:     entity.ID,
		Wins:   entity.Wins,
		Draws:  entity.Draws,
		Losses: entity.Losses,
		Points: entity.Points,
		Rank:   entity.Rank,
	}
}

func (m *TeamStatsHTTPMapper) DomainListToDTO(entities []domain.TeamStats) []dto.TeamStatsResponse {
	if entities == nil {
		return nil
	}

	responses := make([]dto.TeamStatsResponse, len(entities))
	for i, entity := range entities {
		responses[i] = m.DomainToDTO(&entity)
	}

	return responses
}
