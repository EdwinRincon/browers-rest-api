package http

import (
	"github.com/EdwinRincon/browersfc-api/api/dto"
	"github.com/EdwinRincon/browersfc-api/domain"
)

type PlayerHTTPMapper struct{}

func NewPlayerHTTPMapper() *PlayerHTTPMapper {
	return &PlayerHTTPMapper{}
}

// DTO to Domain Conversions (HTTP layer)
func (m *PlayerHTTPMapper) DTOToDomain(dto *dto.CreatePlayerRequest) *domain.Player {
	if dto == nil {
		return nil
	}

	return &domain.Player{
		NickName:         dto.NickName,
		Height:           dto.Height,
		Country:          dto.Country,
		SecondaryCountry: dto.SecondaryCountry,
		Foot:             dto.Foot,
		Age:              dto.Age,
		SquadNumber:      dto.SquadNumber,
		Position:         dto.Position,
		CareerSummary:    dto.CareerSummary,
		UserID:           dto.UserID,
	}
}

func (m *PlayerHTTPMapper) UpdateDTOToDomain(dto *dto.UpdatePlayerRequest) *domain.Player {
	if dto == nil {
		return nil
	}

	player := &domain.Player{}

	if dto.NickName != nil {
		player.NickName = *dto.NickName
	}
	if dto.Height != nil {
		player.Height = *dto.Height
	}
	if dto.Country != nil {
		player.Country = *dto.Country
	}
	if dto.SecondaryCountry != nil {
		player.SecondaryCountry = *dto.SecondaryCountry
	}
	if dto.Foot != nil {
		player.Foot = *dto.Foot
	}
	if dto.Age != nil {
		player.Age = *dto.Age
	}
	if dto.SquadNumber != nil {
		player.SquadNumber = *dto.SquadNumber
	}
	if dto.Rating != nil {
		player.Rating = *dto.Rating
	}
	if dto.Position != nil {
		player.Position = *dto.Position
	}
	if dto.Injured != nil {
		player.Injured = *dto.Injured
	}
	if dto.CareerSummary != nil {
		player.CareerSummary = *dto.CareerSummary
	}
	if dto.UserID != nil {
		player.UserID = dto.UserID
	}

	return player
}

func (m *PlayerHTTPMapper) DomainToDTO(entity *domain.Player) *dto.PlayerResponse {
	if entity == nil {
		return nil
	}

	return &dto.PlayerResponse{
		ID:               entity.ID,
		NickName:         entity.NickName,
		Height:           entity.Height,
		Country:          entity.Country,
		SecondaryCountry: entity.SecondaryCountry,
		Foot:             entity.Foot,
		Age:              entity.Age,
		SquadNumber:      entity.SquadNumber,
		Rating:           entity.Rating,
		Matches:          entity.Matches,
		YCards:           entity.YCards,
		RCards:           entity.RCards,
		Goals:            entity.Goals,
		Assists:          entity.Assists,
		Saves:            entity.Saves,
		Position:         entity.Position,
		Injured:          entity.Injured,
		CareerSummary:    entity.CareerSummary,
		MVPCount:         entity.MVPCount,
		CreatedAt:        entity.CreatedAt,
		UpdatedAt:        entity.UpdatedAt,
	}
}

func (m *PlayerHTTPMapper) DomainListToDTO(entities []domain.Player) []dto.PlayerResponse {
	if entities == nil {
		return nil
	}

	result := make([]dto.PlayerResponse, len(entities))
	for i, entity := range entities {
		response := m.DomainToDTO(&entity)
		if response != nil {
			result[i] = *response
		}
	}
	return result
}

// DomainToShortDTO converts a domain.Player to PlayerShort DTO
// Used for operations that need basic player information
func (m *PlayerHTTPMapper) DomainToShortDTO(entity *domain.Player) *dto.PlayerShort {
	if entity == nil {
		return nil
	}

	return &dto.PlayerShort{
		ID:       entity.ID,
		NickName: entity.NickName,
		Position: entity.Position,
	}
}
