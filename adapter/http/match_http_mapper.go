package http

import (
	"github.com/EdwinRincon/browersfc-api/api/dto"
	"github.com/EdwinRincon/browersfc-api/domain"
)

type MatchHTTPMapper struct{}

func NewMatchHTTPMapper() *MatchHTTPMapper {
	return &MatchHTTPMapper{}
}

// DTO to Domain Conversions (HTTP layer)
func (m *MatchHTTPMapper) DTOToDomain(dto *dto.CreateMatchRequest) *domain.Match {
	if dto == nil {
		return nil
	}

	return &domain.Match{
		Status:      dto.Status,
		Kickoff:     dto.Kickoff,
		Location:    dto.Location,
		HomeGoals:   dto.HomeGoals,
		AwayGoals:   dto.AwayGoals,
		HomeTeamID:  dto.HomeTeamID,
		AwayTeamID:  dto.AwayTeamID,
		SeasonID:    dto.SeasonID,
		MVPPlayerID: dto.MVPPlayerID,
	}
}

func (m *MatchHTTPMapper) UpdateDTOToDomain(dto *dto.UpdateMatchRequest) *domain.Match {
	if dto == nil {
		return nil
	}

	match := &domain.Match{}

	if dto.Status != nil {
		match.Status = *dto.Status
	}
	if dto.Kickoff != nil {
		match.Kickoff = *dto.Kickoff
	}
	if dto.Location != nil {
		match.Location = *dto.Location
	}
	if dto.HomeGoals != nil {
		match.HomeGoals = *dto.HomeGoals
	}
	if dto.AwayGoals != nil {
		match.AwayGoals = *dto.AwayGoals
	}
	if dto.HomeTeamID != nil {
		match.HomeTeamID = *dto.HomeTeamID
	}
	if dto.AwayTeamID != nil {
		match.AwayTeamID = *dto.AwayTeamID
	}
	if dto.SeasonID != nil {
		match.SeasonID = *dto.SeasonID
	}
	if dto.MVPPlayerID != nil {
		match.MVPPlayerID = dto.MVPPlayerID
	}

	return match
}

func (m *MatchHTTPMapper) DomainToDTO(entity *domain.Match) *dto.MatchResponse {
	if entity == nil {
		return nil
	}

	response := &dto.MatchResponse{
		ID:        entity.ID,
		Status:    entity.Status,
		Kickoff:   entity.Kickoff,
		Location:  entity.Location,
		HomeGoals: entity.HomeGoals,
		AwayGoals: entity.AwayGoals,
		CreatedAt: entity.CreatedAt,
		UpdatedAt: entity.UpdatedAt,
	}

	// Note: HomeTeam, AwayTeam, Season, and MVPPlayer will be populated
	// by the handler using separate mappers when needed

	return response
}

func (m *MatchHTTPMapper) DomainListToDTO(entities []domain.Match) []dto.MatchResponse {
	if entities == nil {
		return nil
	}

	result := make([]dto.MatchResponse, len(entities))
	for i, entity := range entities {
		response := m.DomainToDTO(&entity)
		if response != nil {
			result[i] = *response
		}
	}
	return result
}

func (m *MatchHTTPMapper) DomainToShortDTO(entity *domain.Match) *dto.MatchShort {
	if entity == nil {
		return nil
	}

	return &dto.MatchShort{
		ID:        entity.ID,
		Status:    entity.Status,
		Kickoff:   entity.Kickoff,
		Location:  entity.Location,
		HomeGoals: entity.HomeGoals,
		AwayGoals: entity.AwayGoals,
	}
}
