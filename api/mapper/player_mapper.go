package mapper

import (
	"github.com/EdwinRincon/browersfc-api/api/dto"
	"github.com/EdwinRincon/browersfc-api/api/model"
)

func UpdatePlayerFromDTO(player *model.Player, dto *dto.UpdatePlayerRequest) {
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

func ToPlayerResponse(player *model.Player) *dto.PlayerResponse {
	var userShort *dto.UserShort
	if player.User != nil {
		userShort = ToUserShort(player.User)
	}

	// Map PlayerTeams to TeamShort objects
	var teamShorts []dto.TeamShort
	if len(player.PlayerTeams) > 0 {
		// Create a slice with capacity for only those player teams that have an associated team
		validTeams := []dto.TeamShort{}

		for _, playerTeam := range player.PlayerTeams {
			if playerTeam.Team != nil {
				validTeams = append(validTeams, *ToTeamShort(playerTeam.Team))
			}
		}

		teamShorts = validTeams
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
		User:          userShort,
		Teams:         teamShorts,
		CreatedAt:     player.CreatedAt,
		UpdatedAt:     player.UpdatedAt,
	}
}

func ToPlayerResponseList(players []model.Player) []dto.PlayerResponse {
	playerResponses := make([]dto.PlayerResponse, len(players))
	for i, player := range players {
		playerResponses[i] = *ToPlayerResponse(&player)
	}
	return playerResponses
}

func ToPlayerShort(player *model.Player) *dto.PlayerShort {
	return &dto.PlayerShort{
		ID:       player.ID,
		NickName: player.NickName,
		Position: player.Position,
	}
}

func ToPlayer(dto *dto.CreatePlayerRequest) *model.Player {
	return &model.Player{
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
	}
}
