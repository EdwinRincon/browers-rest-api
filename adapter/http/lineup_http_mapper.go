package http

import (
	"github.com/EdwinRincon/browersfc-api/api/dto"
	"github.com/EdwinRincon/browersfc-api/domain"
)

type LineupHTTPMapper struct{}

func NewLineupHTTPMapper() *LineupHTTPMapper {
	return &LineupHTTPMapper{}
}

// DTO to Domain Conversions (HTTP layer)
func (m *LineupHTTPMapper) CreateRequestToDomain(request dto.CreateLineupRequest) *domain.Lineup {
	return &domain.Lineup{
		Position: request.Position,
		PlayerID: request.PlayerID,
		MatchID:  request.MatchID,
		Starting: request.Starting,
	}
}

func (m *LineupHTTPMapper) UpdateRequestToDomain(request dto.UpdateLineupRequest, existing *domain.Lineup) *domain.Lineup {
	updated := *existing // Create a copy

	if request.Position != nil {
		updated.Position = *request.Position
	}
	if request.PlayerID != nil {
		updated.PlayerID = *request.PlayerID
	}
	if request.MatchID != nil {
		updated.MatchID = *request.MatchID
	}
	if request.Starting != nil {
		updated.Starting = *request.Starting
	}

	return &updated
}

// Domain to DTO Conversions (HTTP layer)
func (m *LineupHTTPMapper) DomainToResponse(entity *domain.Lineup) dto.LineupResponse {
	response := dto.LineupResponse{
		ID:        entity.ID,
		Position:  entity.Position,
		PlayerID:  entity.PlayerID,
		MatchID:   entity.MatchID,
		Starting:  entity.Starting,
		CreatedAt: entity.CreatedAt,
		UpdatedAt: entity.UpdatedAt,
	}

	// Include related entities if they exist
	if entity.Player != nil {
		playerMapper := NewPlayerHTTPMapper()
		playerShort := playerMapper.DomainToShortDTO(entity.Player)
		if playerShort != nil {
			response.Player = *playerShort
		}
	}

	if entity.Match != nil {
		matchMapper := NewMatchHTTPMapper()
		matchShort := matchMapper.DomainToShortDTO(entity.Match)
		if matchShort != nil {
			response.Match = *matchShort
		}
	}

	return response
}

func (m *LineupHTTPMapper) DomainToShortResponse(entity *domain.Lineup) dto.LineupShortResponse {
	return dto.LineupShortResponse{
		ID:       entity.ID,
		Position: entity.Position,
		PlayerID: entity.PlayerID,
		Starting: entity.Starting,
	}
}

func (m *LineupHTTPMapper) DomainListToResponse(entities []domain.Lineup) []dto.LineupResponse {
	if entities == nil {
		return nil
	}

	responses := make([]dto.LineupResponse, len(entities))
	for i, entity := range entities {
		responses[i] = m.DomainToResponse(&entity)
	}

	return responses
}

func (m *LineupHTTPMapper) DomainListToShortResponse(entities []domain.Lineup) []dto.LineupShortResponse {
	if entities == nil {
		return nil
	}

	responses := make([]dto.LineupShortResponse, len(entities))
	for i, entity := range entities {
		responses[i] = m.DomainToShortResponse(&entity)
	}

	return responses
}
