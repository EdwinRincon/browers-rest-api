package dto

import (
	"time"
)

type CreateMatchRequest struct {
	Status      string    `json:"status" binding:"required,oneof=scheduled in_progress completed postponed cancelled"`
	Kickoff     time.Time `json:"kickoff" binding:"required"`
	Location    string    `json:"location" binding:"required,max=35"`
	HomeGoals   uint8     `json:"home_goals"`
	AwayGoals   uint8     `json:"away_goals"`
	HomeTeamID  uint64    `json:"home_team_id" binding:"required"`
	AwayTeamID  uint64    `json:"away_team_id" binding:"required"`
	SeasonID    uint64    `json:"season_id" binding:"required"`
	MVPPlayerID *uint64   `json:"mvp_player_id,omitempty"`
}

type UpdateMatchRequest struct {
	Status      *string    `json:"status,omitempty" binding:"omitempty,oneof=scheduled in_progress completed postponed cancelled"`
	Kickoff     *time.Time `json:"kickoff,omitempty" binding:"omitempty"`
	Location    *string    `json:"location,omitempty" binding:"omitempty,max=35"`
	HomeGoals   *uint8     `json:"home_goals,omitempty"`
	AwayGoals   *uint8     `json:"away_goals,omitempty"`
	HomeTeamID  *uint64    `json:"home_team_id,omitempty"`
	AwayTeamID  *uint64    `json:"away_team_id,omitempty"`
	SeasonID    *uint64    `json:"season_id,omitempty"`
	MVPPlayerID *uint64    `json:"mvp_player_id,omitempty"`
}

type MatchResponse struct {
	ID        uint64       `json:"id"`
	Status    string       `json:"status"`
	Kickoff   time.Time    `json:"kickoff"`
	Location  string       `json:"location"`
	HomeGoals uint8        `json:"home_goals"`
	AwayGoals uint8        `json:"away_goals"`
	HomeTeam  TeamShort    `json:"home_team,omitempty"`
	AwayTeam  TeamShort    `json:"away_team,omitempty"`
	Season    SeasonShort  `json:"season,omitempty"`
	MVPPlayer *PlayerShort `json:"mvp_player,omitempty"`
	CreatedAt time.Time    `json:"created_at"`
	UpdatedAt time.Time    `json:"updated_at"`
}

// MatchShort is a simplified match representation for use in other responses
type MatchShort struct {
	ID        uint64    `json:"id"`
	Status    string    `json:"status"`
	Kickoff   time.Time `json:"kickoff"`
	Location  string    `json:"location"`
	HomeGoals uint8     `json:"home_goals"`
	AwayGoals uint8     `json:"away_goals"`
}

// MatchDetailResponse represents a detailed match response including lineups and stats
type MatchDetailResponse struct {
	ID          uint64            `json:"id"`
	Status      string            `json:"status"`
	Kickoff     time.Time         `json:"kickoff"`
	Location    string            `json:"location"`
	HomeGoals   uint8             `json:"home_goals"`
	AwayGoals   uint8             `json:"away_goals"`
	HomeTeam    TeamShort         `json:"home_team,omitempty"`
	AwayTeam    TeamShort         `json:"away_team,omitempty"`
	Lineups     []LineupShort     `json:"lineups,omitempty"`
	PlayerStats []PlayerStatShort `json:"player_stats,omitempty"`
	Season      SeasonShort       `json:"season,omitempty"`
	MVPPlayer   *PlayerShort      `json:"mvp_player,omitempty"`
	CreatedAt   time.Time         `json:"created_at"`
	UpdatedAt   time.Time         `json:"updated_at"`
}

// LineupShort is a simplified lineup representation used in match detail responses
// This would be defined in the actual project in lineup_dto.go
type LineupShort struct {
	ID     uint64      `json:"id"`
	Player PlayerShort `json:"player"`
}

// PlayerStatShort is a simplified player stat representation used in match detail responses
// This would be defined in the actual project in player_stats_dto.go
type PlayerStatShort struct {
	ID       uint64      `json:"id"`
	Player   PlayerShort `json:"player"`
	Goals    uint8       `json:"goals"`
	Assists  uint8       `json:"assists"`
	Minutes  uint8       `json:"minutes"`
	RedCards uint8       `json:"red_cards"`
}
