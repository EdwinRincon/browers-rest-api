package mapper

import (
	"github.com/EdwinRincon/browersfc-api/api/dto"
	"github.com/EdwinRincon/browersfc-api/api/model"
)

func UpdateTeamFromDTO(team *model.Team, dto *dto.UpdateTeamRequest) {
	if dto.FullName != nil {
		team.FullName = *dto.FullName
	}
	if dto.ShortName != nil {
		team.ShortName = *dto.ShortName
	}
	if dto.Color != nil {
		team.Color = *dto.Color
	}
	if dto.Color2 != nil {
		team.Color2 = *dto.Color2
	}
	if dto.Shield != nil {
		team.Shield = *dto.Shield
	}
	if dto.NextMatch != nil {
		team.NextMatch = *dto.NextMatch
	}
}

func ToTeamResponse(team *model.Team) *dto.TeamResponse {
	return &dto.TeamResponse{
		ID:        team.ID,
		FullName:  team.FullName,
		ShortName: team.ShortName,
		Color:     team.Color,
		Color2:    team.Color2,
		Shield:    team.Shield,
		NextMatch: team.NextMatch,
		CreatedAt: team.CreatedAt,
		UpdatedAt: team.UpdatedAt,
	}
}

func ToTeamResponseList(teams []model.Team) []dto.TeamResponse {
	teamResponses := make([]dto.TeamResponse, len(teams))
	for i, team := range teams {
		teamResponses[i] = *ToTeamResponse(&team)
	}
	return teamResponses
}

func ToTeamShort(team *model.Team) *dto.TeamShort {
	return &dto.TeamShort{
		ID:        team.ID,
		FullName:  team.FullName,
		ShortName: team.ShortName,
	}
}

func ToTeam(dto *dto.CreateTeamRequest) *model.Team {
	return &model.Team{
		FullName:  dto.FullName,
		ShortName: dto.ShortName,
		Color:     dto.Color,
		Color2:    dto.Color2,
		Shield:    dto.Shield,
		NextMatch: dto.NextMatch,
	}
}
