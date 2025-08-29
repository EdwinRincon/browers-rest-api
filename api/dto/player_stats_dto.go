package dto

type PlayerStatCreateDTO struct {
	PlayerID      uint64  `json:"player_id" binding:"required"`
	MatchID       uint64  `json:"match_id" binding:"required"`
	SeasonID      uint64  `json:"season_id" binding:"required"`
	TeamID        *uint64 `json:"team_id,omitempty"`
	Goals         uint8   `json:"goals" binding:"gte=0"`
	Assists       uint8   `json:"assists" binding:"gte=0"`
	Saves         uint8   `json:"saves" binding:"gte=0"`
	YellowCards   uint8   `json:"yellow_cards" binding:"gte=0"`
	RedCards      uint8   `json:"red_cards" binding:"gte=0"`
	Rating        uint8   `json:"rating" binding:"gte=0,lte=100"`
	Starting      bool    `json:"starting"`
	MinutesPlayed uint8   `json:"minutes_played" binding:"gte=0,lte=120"`
	Position      string  `json:"position" binding:"omitempty,oneof=por ceni cend lati med latd del deli deld"`
}

type PlayerStatUpdateDTO struct {
	ID            uint64  `json:"id" binding:"required"`
	TeamID        *uint64 `json:"team_id,omitempty"`
	Goals         *uint8  `json:"goals,omitempty" binding:"omitempty,gte=0"`
	Assists       *uint8  `json:"assists,omitempty" binding:"omitempty,gte=0"`
	Saves         *uint8  `json:"saves,omitempty" binding:"omitempty,gte=0"`
	YellowCards   *uint8  `json:"yellow_cards,omitempty" binding:"omitempty,gte=0"`
	RedCards      *uint8  `json:"red_cards,omitempty" binding:"omitempty,gte=0"`
	Rating        *uint8  `json:"rating,omitempty" binding:"omitempty,gte=0,lte=100"`
	Starting      *bool   `json:"starting,omitempty"`
	MinutesPlayed *uint8  `json:"minutes_played,omitempty" binding:"omitempty,gte=0,lte=120"`
	IsMVP         *bool   `json:"is_mvp,omitempty"`
	Position      *string `json:"position,omitempty" binding:"omitempty,oneof=por ceni cend lati med latd del deli deld"`
}

type PlayerStatResponseDTO struct {
	ID            uint64  `json:"id"`
	PlayerID      uint64  `json:"player_id"`
	PlayerName    string  `json:"player_name,omitempty"`
	MatchID       uint64  `json:"match_id"`
	MatchDate     string  `json:"match_date,omitempty"`
	SeasonID      uint64  `json:"season_id"`
	SeasonYear    uint16  `json:"season_year,omitempty"`
	TeamID        *uint64 `json:"team_id,omitempty"`
	TeamName      string  `json:"team_name,omitempty"`
	Goals         uint8   `json:"goals"`
	Assists       uint8   `json:"assists"`
	Saves         uint8   `json:"saves"`
	YellowCards   uint8   `json:"yellow_cards"`
	RedCards      uint8   `json:"red_cards"`
	Rating        uint8   `json:"rating"`
	Starting      bool    `json:"starting"`
	MinutesPlayed uint8   `json:"minutes_played"`
	IsMVP         bool    `json:"is_mvp"`
	Position      string  `json:"position,omitempty"`
	CreatedAt     string  `json:"created_at,omitempty"`
	UpdatedAt     string  `json:"updated_at,omitempty"`
}

type PlayerStatFilterDTO struct {
	PlayerID  *uint64 `form:"player_id"`
	MatchID   *uint64 `form:"match_id"`
	SeasonID  *uint64 `form:"season_id"`
	TeamID    *uint64 `form:"team_id"`
	Rating    *uint8  `form:"rating"`
	MinGoals  *uint8  `form:"min_goals"`
	IsMVP     *bool   `form:"is_mvp"`
	Position  *string `form:"position"`
	StartDate *string `form:"start_date"`
	EndDate   *string `form:"end_date"`
}

type PlayerStatSeasonAggregateDTO struct {
	PlayerID      uint64  `json:"player_id"`
	PlayerName    string  `json:"player_name"`
	SeasonID      uint64  `json:"season_id"`
	SeasonYear    uint16  `json:"season_year"`
	MatchesPlayed uint16  `json:"matches_played"`
	TotalMinutes  uint16  `json:"total_minutes"`
	TotalGoals    uint16  `json:"total_goals"`
	TotalAssists  uint16  `json:"total_assists"`
	TotalSaves    uint16  `json:"total_saves"`
	YellowCards   uint8   `json:"yellow_cards"`
	RedCards      uint8   `json:"red_cards"`
	MVPCount      uint8   `json:"mvp_count"`
	AvgRating     float32 `json:"avg_rating"`
}
