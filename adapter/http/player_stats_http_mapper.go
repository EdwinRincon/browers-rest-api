package http

import (
	"github.com/EdwinRincon/browersfc-api/api/dto"
	"github.com/EdwinRincon/browersfc-api/domain"
)

type PlayerStatsHTTPMapper struct{}

func NewPlayerStatsHTTPMapper() *PlayerStatsHTTPMapper {
	return &PlayerStatsHTTPMapper{}
}

// DTO to Domain Conversions (HTTP layer)
func (m *PlayerStatsHTTPMapper) DTOToDomain(dto *dto.CreatePlayerStatRequest) *domain.PlayerStat {
	if dto == nil {
		return nil
	}

	return &domain.PlayerStat{
		PlayerID:      dto.PlayerID,
		MatchID:       dto.MatchID,
		SeasonID:      dto.SeasonID,
		TeamID:        dto.TeamID,
		Goals:         dto.Goals,
		Assists:       dto.Assists,
		Saves:         dto.Saves,
		YellowCards:   dto.YellowCards,
		RedCards:      dto.RedCards,
		Rating:        dto.Rating,
		IsStarting:    dto.Starting,
		MinutesPlayed: dto.MinutesPlayed,
		IsMVP:         dto.IsMVP,
		Position:      dto.Position,
	}
}

func (m *PlayerStatsHTTPMapper) UpdateDTOToDomain(dto *dto.UpdatePlayerStatRequest) *domain.PlayerStat {
	if dto == nil {
		return nil
	}

	stat := &domain.PlayerStat{}

	if dto.Goals != nil {
		stat.Goals = *dto.Goals
	}
	if dto.Assists != nil {
		stat.Assists = *dto.Assists
	}
	if dto.Saves != nil {
		stat.Saves = *dto.Saves
	}
	if dto.YellowCards != nil {
		stat.YellowCards = *dto.YellowCards
	}
	if dto.RedCards != nil {
		stat.RedCards = *dto.RedCards
	}
	if dto.Rating != nil {
		stat.Rating = *dto.Rating
	}
	if dto.Starting != nil {
		stat.IsStarting = *dto.Starting
	}
	if dto.MinutesPlayed != nil {
		stat.MinutesPlayed = *dto.MinutesPlayed
	}
	if dto.IsMVP != nil {
		stat.IsMVP = *dto.IsMVP
	}
	if dto.Position != nil {
		stat.Position = *dto.Position
	}

	return stat
}

func (m *PlayerStatsHTTPMapper) DomainToDTO(entity *domain.PlayerStat) *dto.PlayerStatResponse {
	if entity == nil {
		return nil
	}

	response := &dto.PlayerStatResponse{
		ID:            entity.ID,
		Goals:         entity.Goals,
		Assists:       entity.Assists,
		Saves:         entity.Saves,
		YellowCards:   entity.YellowCards,
		RedCards:      entity.RedCards,
		Rating:        entity.Rating,
		Starting:      entity.IsStarting,
		MinutesPlayed: entity.MinutesPlayed,
		IsMVP:         entity.IsMVP,
		Position:      entity.Position,
		CreatedAt:     entity.CreatedAt,
		UpdatedAt:     entity.UpdatedAt,
	}

	// Player, Match, Season, and Team will be populated by handler if needed

	return response
}

func (m *PlayerStatsHTTPMapper) DomainListToDTO(entities []domain.PlayerStat) []dto.PlayerStatResponse {
	if entities == nil {
		return nil
	}

	result := make([]dto.PlayerStatResponse, len(entities))
	for i, entity := range entities {
		response := m.DomainToDTO(&entity)
		if response != nil {
			result[i] = *response
		}
	}
	return result
}
