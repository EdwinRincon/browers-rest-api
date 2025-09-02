package mapper

import (
	"github.com/EdwinRincon/browersfc-api/api/dto"
	"github.com/EdwinRincon/browersfc-api/api/model"
)

// ToLineup maps CreateLineupRequest DTO to Lineup model
func ToLineup(dto *dto.CreateLineupRequest) *model.Lineup {
	return &model.Lineup{
		Position: dto.Position,
		PlayerID: dto.PlayerID,
		MatchID:  dto.MatchID,
		Starting: dto.Starting,
	}
}

// UpdateLineupFromDTO updates a Lineup model from UpdateLineupRequest DTO
func UpdateLineupFromDTO(lineup *model.Lineup, dto *dto.UpdateLineupRequest) {
	if dto.Position != nil {
		lineup.Position = *dto.Position
	}
	if dto.PlayerID != nil {
		lineup.PlayerID = *dto.PlayerID
	}
	if dto.MatchID != nil {
		lineup.MatchID = *dto.MatchID
	}
	if dto.Starting != nil {
		lineup.Starting = *dto.Starting
	}
}

// ToLineupResponse maps Lineup model to LineupResponse DTO
func ToLineupResponse(lineup *model.Lineup) *dto.LineupResponse {
	response := &dto.LineupResponse{
		ID:        lineup.ID,
		Position:  lineup.Position,
		PlayerID:  lineup.PlayerID,
		MatchID:   lineup.MatchID,
		Starting:  lineup.Starting,
		CreatedAt: lineup.CreatedAt,
		UpdatedAt: lineup.UpdatedAt,
	}

	// Add player information if available
	if lineup.Player != nil {
		playerShort := ToPlayerShort(lineup.Player)
		if playerShort != nil {
			response.Player = *playerShort
		}
	}

	// Add match information if available
	if lineup.Match != nil {
		matchShort := ToMatchShort(lineup.Match)
		if matchShort != nil {
			response.Match = *matchShort
		}
	}

	return response
}

// ToLineupResponseList maps a slice of Lineup models to a slice of LineupResponse DTOs
func ToLineupResponseList(lineups []model.Lineup) []dto.LineupResponse {
	lineupResponses := make([]dto.LineupResponse, len(lineups))
	for i, lineup := range lineups {
		lineupResponses[i] = *ToLineupResponse(&lineup)
	}
	return lineupResponses
}

// ToLineupShortResponse maps Lineup model to LineupShortResponse DTO
func ToLineupShortResponse(lineup *model.Lineup) *dto.LineupShortResponse {
	return &dto.LineupShortResponse{
		ID:       lineup.ID,
		Position: lineup.Position,
		PlayerID: lineup.PlayerID,
		Starting: lineup.Starting,
	}
}

// ToLineupShortResponseList maps a slice of Lineup models to a slice of LineupShortResponse DTOs
func ToLineupShortResponseList(lineups []model.Lineup) []dto.LineupShortResponse {
	lineupResponses := make([]dto.LineupShortResponse, len(lineups))
	for i, lineup := range lineups {
		lineupResponses[i] = *ToLineupShortResponse(&lineup)
	}
	return lineupResponses
}

// OrganizeLineupsByMatchID groups lineups by whether they are starting or substitutes
func OrganizeLineupsByMatchID(match *model.Match, lineups []model.Lineup) *dto.MatchLineupResponse {
	if match == nil {
		return nil
	}

	response := &dto.MatchLineupResponse{
		MatchID:        match.ID,
		HomeTeamID:     match.HomeTeamID,
		AwayTeamID:     match.AwayTeamID,
		Date:           match.Kickoff,
		StartingLineup: make([]dto.LineupShortResponse, 0),
		Substitutes:    make([]dto.LineupShortResponse, 0),
	}

	if match.HomeTeam != nil {
		response.HomeTeam = match.HomeTeam.FullName
	}

	if match.AwayTeam != nil {
		response.AwayTeam = match.AwayTeam.FullName
	}

	for _, lineup := range lineups {
		lineupShort := ToLineupShortResponse(&lineup)
		if lineup.Starting {
			response.StartingLineup = append(response.StartingLineup, *lineupShort)
		} else {
			response.Substitutes = append(response.Substitutes, *lineupShort)
		}
	}

	return response
}
