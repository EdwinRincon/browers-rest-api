package mapper

import (
	"github.com/EdwinRincon/browersfc-api/api/dto"
	"github.com/EdwinRincon/browersfc-api/api/model"
)

// ToMatch converts a CreateMatchRequest DTO to a Match model
func ToMatch(req *dto.CreateMatchRequest) *model.Match {
	match := &model.Match{
		Status:     req.Status,
		Kickoff:    req.Kickoff,
		Location:   req.Location,
		HomeGoals:  req.HomeGoals,
		AwayGoals:  req.AwayGoals,
		HomeTeamID: req.HomeTeamID,
		AwayTeamID: req.AwayTeamID,
		SeasonID:   req.SeasonID,
	}

	if req.MVPPlayerID != nil {
		match.MVPPlayerID = req.MVPPlayerID
	}

	return match
}

// UpdateMatchFromDTO updates a Match model using an UpdateMatchRequest DTO
func UpdateMatchFromDTO(match *model.Match, dto *dto.UpdateMatchRequest) {
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
}

// ToSeasonShort converts a Season model to a SeasonShort DTO
func ToSeasonShort(season *model.Season) *dto.SeasonShort {
	if season == nil {
		return nil
	}
	return &dto.SeasonShort{
		ID:   season.ID,
		Year: season.Year,
	}
}

// ToMatchResponse converts a Match model to a MatchResponse DTO
func ToMatchResponse(match *model.Match) *dto.MatchResponse {
	resp := &dto.MatchResponse{
		ID:        match.ID,
		Status:    match.Status,
		Kickoff:   match.Kickoff,
		Location:  match.Location,
		HomeGoals: match.HomeGoals,
		AwayGoals: match.AwayGoals,
		CreatedAt: match.CreatedAt,
		UpdatedAt: match.UpdatedAt,
	}

	// Map related entities if available
	if match.HomeTeam != nil {
		homeTeam := ToTeamShort(match.HomeTeam)
		resp.HomeTeam = *homeTeam
	}

	if match.AwayTeam != nil {
		awayTeam := ToTeamShort(match.AwayTeam)
		resp.AwayTeam = *awayTeam
	}

	if match.Season != nil {
		season := ToSeasonShort(match.Season)
		resp.Season = *season
	}

	if match.MVPPlayer != nil {
		mvpPlayer := ToPlayerShort(match.MVPPlayer)
		resp.MVPPlayer = mvpPlayer
	}

	return resp
}

// ToMatchResponseList converts a slice of Match models to a slice of MatchResponse DTOs
func ToMatchResponseList(matches []model.Match) []dto.MatchResponse {
	responses := make([]dto.MatchResponse, len(matches))
	for i, match := range matches {
		responses[i] = *ToMatchResponse(&match)
	}
	return responses
}

// ToMatchShort converts a Match model to a MatchShort DTO
func ToMatchShort(match *model.Match) *dto.MatchShort {
	return &dto.MatchShort{
		ID:        match.ID,
		Status:    match.Status,
		Kickoff:   match.Kickoff,
		Location:  match.Location,
		HomeGoals: match.HomeGoals,
		AwayGoals: match.AwayGoals,
	}
}

// ToMatchDetailResponse converts a Match model to a MatchDetailResponse DTO with lineups and stats
func ToMatchDetailResponse(match *model.Match) *dto.MatchDetailResponse {
	if match == nil {
		return nil
	}

	// Start with the basic match data
	detail := &dto.MatchDetailResponse{
		ID:        match.ID,
		Status:    match.Status,
		Kickoff:   match.Kickoff,
		Location:  match.Location,
		HomeGoals: match.HomeGoals,
		AwayGoals: match.AwayGoals,
		CreatedAt: match.CreatedAt,
		UpdatedAt: match.UpdatedAt,
	}

	// Map related entities if available
	if match.HomeTeam != nil {
		if homeTeam := ToTeamShort(match.HomeTeam); homeTeam != nil {
			detail.HomeTeam = *homeTeam
		}
	}

	if match.AwayTeam != nil {
		if awayTeam := ToTeamShort(match.AwayTeam); awayTeam != nil {
			detail.AwayTeam = *awayTeam
		}
	}

	if match.Season != nil {
		if season := ToSeasonShort(match.Season); season != nil {
			detail.Season = *season
		}
	}

	if match.MVPPlayer != nil {
		detail.MVPPlayer = ToPlayerShort(match.MVPPlayer)
	}

	// Map lineups if available
	detail.Lineups = mapLineups(match.Lineups)

	// Map player stats if available
	detail.PlayerStats = mapPlayerStats(match.PlayerStats)

	return detail
}

// mapLineups converts a slice of Lineup models to LineupShort DTOs
func mapLineups(lineups []model.Lineup) []dto.LineupShort {
	if len(lineups) == 0 {
		return nil
	}

	result := make([]dto.LineupShort, 0, len(lineups))
	for _, lineup := range lineups {
		if playerShort := ToPlayerShort(lineup.Player); playerShort != nil {
			result = append(result, dto.LineupShort{
				ID:     lineup.ID,
				Player: *playerShort,
			})
		}
	}
	return result
}

// mapPlayerStats converts a slice of PlayerStat models to PlayerStatShort DTOs
func mapPlayerStats(stats []model.PlayerStat) []dto.PlayerStatShort {
	if len(stats) == 0 {
		return nil
	}

	result := make([]dto.PlayerStatShort, 0, len(stats))
	for _, stat := range stats {
		if playerShort := ToPlayerShort(stat.Player); playerShort != nil {
			result = append(result, dto.PlayerStatShort{
				ID:       stat.ID,
				Player:   *playerShort,
				Goals:    stat.Goals,
				Assists:  stat.Assists,
				Minutes:  stat.MinutesPlayed,
				RedCards: stat.RC,
			})
		}
	}
	return result
}
