package mapper

import (
	"github.com/EdwinRincon/browersfc-api/api/dto"
	"github.com/EdwinRincon/browersfc-api/api/model"
)

// ToPlayerStat converts a CreatePlayerStatRequest DTO to a PlayerStat model
func ToPlayerStat(dto *dto.CreatePlayerStatRequest) *model.PlayerStat {
	return &model.PlayerStat{
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

// UpdatePlayerStatFromDTO updates a PlayerStat model with values from an UpdatePlayerStatRequest DTO
func UpdatePlayerStatFromDTO(playerStat *model.PlayerStat, dto *dto.UpdatePlayerStatRequest) {
	if dto.TeamID != nil {
		playerStat.TeamID = dto.TeamID
	}
	if dto.Goals != nil {
		playerStat.Goals = *dto.Goals
	}
	if dto.Assists != nil {
		playerStat.Assists = *dto.Assists
	}
	if dto.Saves != nil {
		playerStat.Saves = *dto.Saves
	}
	if dto.YellowCards != nil {
		playerStat.YellowCards = *dto.YellowCards
	}
	if dto.RedCards != nil {
		playerStat.RedCards = *dto.RedCards
	}
	if dto.Rating != nil {
		playerStat.Rating = *dto.Rating
	}
	if dto.Starting != nil {
		playerStat.IsStarting = *dto.Starting
	}
	if dto.MinutesPlayed != nil {
		playerStat.MinutesPlayed = *dto.MinutesPlayed
	}
	if dto.IsMVP != nil {
		playerStat.IsMVP = *dto.IsMVP
	}
	if dto.Position != nil {
		playerStat.Position = *dto.Position
	}
}

// ToPlayerStatResponse converts a PlayerStat model to a PlayerStatResponse DTO
func ToPlayerStatResponse(playerStat *model.PlayerStat) *dto.PlayerStatResponse {
	response := &dto.PlayerStatResponse{
		ID:            playerStat.ID,
		PlayerID:      playerStat.PlayerID,
		MatchID:       playerStat.MatchID,
		SeasonID:      playerStat.SeasonID,
		TeamID:        playerStat.TeamID,
		Goals:         playerStat.Goals,
		Assists:       playerStat.Assists,
		Saves:         playerStat.Saves,
		YellowCards:   playerStat.YellowCards,
		RedCards:      playerStat.RedCards,
		Rating:        playerStat.Rating,
		Starting:      playerStat.IsStarting,
		MinutesPlayed: playerStat.MinutesPlayed,
		IsMVP:         playerStat.IsMVP,
		Position:      playerStat.Position,
		CreatedAt:     playerStat.CreatedAt,
		UpdatedAt:     playerStat.UpdatedAt,
	}

	// Add nested objects when available
	if playerStat.Player != nil {
		response.Player = *ToPlayerShort(playerStat.Player)
	}

	if playerStat.Match != nil {
		response.Match = dto.MatchShort{
			ID:        playerStat.Match.ID,
			Status:    playerStat.Match.Status,
			Kickoff:   playerStat.Match.Kickoff,
			Location:  playerStat.Match.Location,
			HomeGoals: playerStat.Match.HomeGoals,
			AwayGoals: playerStat.Match.AwayGoals,
		}
	}

	if playerStat.Season != nil {
		response.Season = dto.SeasonShort{
			ID:   playerStat.Season.ID,
			Year: playerStat.Season.Year,
		}
	}

	if playerStat.Team != nil {
		response.Team = &dto.TeamShort{
			ID:        playerStat.Team.ID,
			FullName:  playerStat.Team.FullName,
			ShortName: playerStat.Team.ShortName,
		}
	}

	return response
}

// ToPlayerStatResponseList converts a slice of PlayerStat models to a slice of PlayerStatResponse DTOs
func ToPlayerStatResponseList(playerStats []model.PlayerStat) []dto.PlayerStatResponse {
	playerStatResponses := make([]dto.PlayerStatResponse, len(playerStats))
	for i, playerStat := range playerStats {
		playerStatResponses[i] = *ToPlayerStatResponse(&playerStat)
	}
	return playerStatResponses
}
