package mapper

import (
	"github.com/EdwinRincon/browersfc-api/api/dto"
	"github.com/EdwinRincon/browersfc-api/domain"
)

// CreatePlayerRequestToDomain converts a CreatePlayerRequest DTO to a domain Player.
func CreatePlayerRequestToDomain(dto *dto.CreatePlayerRequest) *domain.Player {
	if dto == nil {
		return nil
	}

	return &domain.Player{
		NickName:      dto.NickName,
		Height:        dto.Height,
		Country:       dto.Country,
		Country2:      dto.Country2,
		Foot:          dto.Foot,
		Age:           dto.Age,
		SquadNumber:   dto.SquadNumber,
		Position:      dto.Position,
		CareerSummary: dto.CareerSummary,
		UserID:        dto.UserID,
		// ID, statistics, and timestamps will be set by the persistence layer
	}
}

// PlayerDomainToResponse converts a domain Player to a PlayerResponse DTO.
func PlayerDomainToResponse(player *domain.Player) *dto.PlayerResponse {
	if player == nil {
		return nil
	}

	return &dto.PlayerResponse{
		ID:            player.ID,
		NickName:      player.NickName,
		Height:        player.Height,
		Country:       player.Country,
		Country2:      player.Country2,
		Foot:          player.Foot,
		Age:           player.Age,
		SquadNumber:   player.SquadNumber,
		Rating:        player.Rating,
		Matches:       player.Matches,
		YCards:        player.YCards,
		RCards:        player.RCards,
		Goals:         player.Goals,
		Assists:       player.Assists,
		Saves:         player.Saves,
		Position:      player.Position,
		Injured:       player.Injured,
		CareerSummary: player.CareerSummary,
		MVPCount:      player.MVPCount,
		// User and Teams will need to be handled separately based on requirements
		CreatedAt: player.CreatedAt,
		UpdatedAt: player.UpdatedAt,
	}
}

// PlayerDomainListToResponse converts a slice of domain Players to PlayerResponse DTOs.
func PlayerDomainListToResponse(players []domain.Player) []dto.PlayerResponse {
	responses := make([]dto.PlayerResponse, len(players))
	for i, player := range players {
		responses[i] = *PlayerDomainToResponse(&player)
	}
	return responses
}

// PlayerDomainToShort converts a domain Player to a PlayerShort DTO.
func PlayerDomainToShort(player *domain.Player) *dto.PlayerShort {
	if player == nil {
		return nil
	}

	return &dto.PlayerShort{
		ID:       player.ID,
		NickName: player.NickName,
		Position: player.Position,
	}
}

// UpdatePlayerRequestToDomain updates a domain Player with values from UpdatePlayerRequest DTO.
func UpdatePlayerRequestToDomain(player *domain.Player, dto *dto.UpdatePlayerRequest) {
	if dto == nil || player == nil {
		return
	}

	if dto.NickName != nil {
		player.NickName = *dto.NickName
	}
	if dto.Height != nil {
		player.Height = *dto.Height
	}
	if dto.Country != nil {
		player.Country = *dto.Country
	}
	if dto.Country2 != nil {
		player.Country2 = *dto.Country2
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
}
