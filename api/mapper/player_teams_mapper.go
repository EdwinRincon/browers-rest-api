package mapper

import (
	"github.com/EdwinRincon/browersfc-api/api/dto"
	"github.com/EdwinRincon/browersfc-api/api/model"
)

// ToPlayerTeam converts a CreatePlayerTeamRequest DTO to a PlayerTeam model
func ToPlayerTeam(dto *dto.CreatePlayerTeamRequest) *model.PlayerTeam {
	return &model.PlayerTeam{
		PlayerID:  dto.PlayerID,
		TeamID:    dto.TeamID,
		SeasonID:  dto.SeasonID,
		StartDate: dto.StartDate,
		EndDate:   dto.EndDate,
	}
}

// UpdatePlayerTeamFromDTO updates a PlayerTeam model from an UpdatePlayerTeamRequest DTO
func UpdatePlayerTeamFromDTO(playerTeam *model.PlayerTeam, dto *dto.UpdatePlayerTeamRequest) {
	if dto.StartDate != nil {
		playerTeam.StartDate = *dto.StartDate
	}
	if dto.EndDate != nil {
		playerTeam.EndDate = dto.EndDate
	}
}

// ToPlayerTeamResponse converts a PlayerTeam model to a PlayerTeamResponse DTO
func ToPlayerTeamResponse(playerTeam *model.PlayerTeam) *dto.PlayerTeamResponse {
	if playerTeam == nil {
		return nil
	}

	response := &dto.PlayerTeamResponse{
		ID:        playerTeam.ID,
		PlayerID:  playerTeam.PlayerID,
		TeamID:    playerTeam.TeamID,
		SeasonID:  playerTeam.SeasonID,
		StartDate: playerTeam.StartDate,
		EndDate:   playerTeam.EndDate,
		CreatedAt: playerTeam.CreatedAt,
		UpdatedAt: playerTeam.UpdatedAt,
	}

	// Handle Player relationship
	if playerTeam.Player != nil {
		response.Player.ID = playerTeam.Player.ID
		response.Player.NickName = playerTeam.Player.NickName
		response.Player.Position = playerTeam.Player.Position
	}

	// Handle Team relationship
	if playerTeam.Team != nil {
		response.Team.ID = playerTeam.Team.ID
		response.Team.FullName = playerTeam.Team.FullName
		response.Team.ShortName = playerTeam.Team.ShortName
	}

	// Handle Season relationship
	if playerTeam.Season != nil {
		response.Season.ID = playerTeam.Season.ID
		response.Season.Year = playerTeam.Season.Year
	}

	return response
}

// ToPlayerTeamResponseList converts a slice of PlayerTeam models to a slice of PlayerTeamResponse DTOs
func ToPlayerTeamResponseList(playerTeams []model.PlayerTeam) []dto.PlayerTeamResponse {
	responses := make([]dto.PlayerTeamResponse, len(playerTeams))
	for i, playerTeam := range playerTeams {
		responses[i] = *ToPlayerTeamResponse(&playerTeam)
	}
	return responses
}

// ToPlayerTeamShort converts a PlayerTeam model to a PlayerTeamShort DTO
func ToPlayerTeamShort(playerTeam *model.PlayerTeam) *dto.PlayerTeamShort {
	return &dto.PlayerTeamShort{
		ID:        playerTeam.ID,
		PlayerID:  playerTeam.PlayerID,
		TeamID:    playerTeam.TeamID,
		SeasonID:  playerTeam.SeasonID,
		StartDate: playerTeam.StartDate,
		EndDate:   playerTeam.EndDate,
	}
}
