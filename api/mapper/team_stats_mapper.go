package mapper

import (
	"github.com/EdwinRincon/browersfc-api/adapter/mapper"
	"github.com/EdwinRincon/browersfc-api/api/dto"
	"github.com/EdwinRincon/browersfc-api/api/model"
)

func UpdateTeamStatsFromDTO(teamStats *model.TeamStat, dto *dto.UpdateTeamStatsRequest) {
	if dto.Wins != nil {
		teamStats.Wins = *dto.Wins
	}
	if dto.Draws != nil {
		teamStats.Draws = *dto.Draws
	}
	if dto.Losses != nil {
		teamStats.Losses = *dto.Losses
	}
	if dto.GoalsFor != nil {
		teamStats.GoalsFor = *dto.GoalsFor
	}
	if dto.GoalsAgainst != nil {
		teamStats.GoalsAgainst = *dto.GoalsAgainst
	}
	if dto.Points != nil {
		teamStats.Points = *dto.Points
	}
	if dto.Rank != nil {
		teamStats.Rank = *dto.Rank
	}
	if dto.SeasonID != nil {
		teamStats.SeasonID = *dto.SeasonID
	}
	if dto.TeamID != nil {
		teamStats.TeamID = *dto.TeamID
	}
}

func ToTeamStatsResponse(teamStats *model.TeamStat) *dto.TeamStatsResponse {
	var team *dto.TeamShort
	if teamStats.Team != nil {
		team = ToTeamShort(teamStats.Team)
	}

	var season *dto.SeasonShort
	if teamStats.Season != nil {
		seasonMapper := mapper.NewSeasonMapper()
		season = seasonMapper.ModelToShortDTO(teamStats.Season)
	}

	return &dto.TeamStatsResponse{
		ID:           teamStats.ID,
		Wins:         teamStats.Wins,
		Draws:        teamStats.Draws,
		Losses:       teamStats.Losses,
		GoalsFor:     teamStats.GoalsFor,
		GoalsAgainst: teamStats.GoalsAgainst,
		Points:       teamStats.Points,
		Rank:         teamStats.Rank,
		SeasonID:     teamStats.SeasonID,
		TeamID:       teamStats.TeamID,
		Team:         team,
		Season:       season,
		CreatedAt:    teamStats.CreatedAt,
		UpdatedAt:    teamStats.UpdatedAt,
	}
}

func ToTeamStatsResponseList(teamStats []model.TeamStat) []dto.TeamStatsResponse {
	teamStatsResponses := make([]dto.TeamStatsResponse, len(teamStats))
	for i, stats := range teamStats {
		teamStatsResponses[i] = *ToTeamStatsResponse(&stats)
	}
	return teamStatsResponses
}

func ToTeamStatsShort(teamStats *model.TeamStat) *dto.TeamStatsShort {
	return &dto.TeamStatsShort{
		ID:     teamStats.ID,
		Wins:   teamStats.Wins,
		Draws:  teamStats.Draws,
		Losses: teamStats.Losses,
		Points: teamStats.Points,
		Rank:   teamStats.Rank,
	}
}

func ToTeamStats(dto *dto.CreateTeamStatsRequest) *model.TeamStat {
	return &model.TeamStat{
		Wins:         dto.Wins,
		Draws:        dto.Draws,
		Losses:       dto.Losses,
		GoalsFor:     dto.GoalsFor,
		GoalsAgainst: dto.GoalsAgainst,
		Points:       dto.Points,
		Rank:         dto.Rank,
		SeasonID:     dto.SeasonID,
		TeamID:       dto.TeamID,
	}
}
