package mapper

import (
	"github.com/EdwinRincon/browersfc-api/api/dto"
	"github.com/EdwinRincon/browersfc-api/api/model"
	"github.com/EdwinRincon/browersfc-api/domain"
)

type TeamStatsMapper struct{}

func NewTeamStatsMapper() *TeamStatsMapper {
	return &TeamStatsMapper{}
}

// DTO to Domain Conversions (HTTP layer)
func (m *TeamStatsMapper) CreateRequestToDomain(dto *dto.CreateTeamStatsRequest) *domain.TeamStats {
	return &domain.TeamStats{
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

func (m *TeamStatsMapper) UpdateRequestToDomain(current *domain.TeamStats, dto *dto.UpdateTeamStatsRequest) *domain.TeamStats {
	updated := *current // Create a copy

	if dto.Wins != nil {
		updated.Wins = *dto.Wins
	}
	if dto.Draws != nil {
		updated.Draws = *dto.Draws
	}
	if dto.Losses != nil {
		updated.Losses = *dto.Losses
	}
	if dto.GoalsFor != nil {
		updated.GoalsFor = *dto.GoalsFor
	}
	if dto.GoalsAgainst != nil {
		updated.GoalsAgainst = *dto.GoalsAgainst
	}
	if dto.Points != nil {
		updated.Points = *dto.Points
	}
	if dto.Rank != nil {
		updated.Rank = *dto.Rank
	}
	if dto.SeasonID != nil {
		updated.SeasonID = *dto.SeasonID
	}
	if dto.TeamID != nil {
		updated.TeamID = *dto.TeamID
	}

	return &updated
}

func (m *TeamStatsMapper) DomainToResponse(entity *domain.TeamStats) *dto.TeamStatsResponse {
	var team *dto.TeamShort
	if entity.Team != nil {
		teamMapper := NewTeamMapper()
		team = teamMapper.DomainToShortDTO(entity.Team)
	}

	var season *dto.SeasonShort
	if entity.Season != nil {
		seasonMapper := NewSeasonMapper()
		season = seasonMapper.DomainToShortDTO(entity.Season)
	}

	return &dto.TeamStatsResponse{
		ID:           entity.ID,
		Wins:         entity.Wins,
		Draws:        entity.Draws,
		Losses:       entity.Losses,
		GoalsFor:     entity.GoalsFor,
		GoalsAgainst: entity.GoalsAgainst,
		Points:       entity.Points,
		Rank:         entity.Rank,
		SeasonID:     entity.SeasonID,
		TeamID:       entity.TeamID,
		Team:         team,
		Season:       season,
		CreatedAt:    entity.CreatedAt,
		UpdatedAt:    entity.UpdatedAt,
	}
}

func (m *TeamStatsMapper) DomainListToResponse(entities []domain.TeamStats) []dto.TeamStatsResponse {
	responses := make([]dto.TeamStatsResponse, len(entities))
	for i, entity := range entities {
		responses[i] = *m.DomainToResponse(&entity)
	}
	return responses
}

func (m *TeamStatsMapper) DomainToShortDTO(entity *domain.TeamStats) *dto.TeamStatsShort {
	return &dto.TeamStatsShort{
		ID:     entity.ID,
		Wins:   entity.Wins,
		Draws:  entity.Draws,
		Losses: entity.Losses,
		Points: entity.Points,
		Rank:   entity.Rank,
	}
}

// Domain to Model Conversions (Infrastructure layer)
func (m *TeamStatsMapper) DomainToModel(entity *domain.TeamStats) *model.TeamStat {
	modelEntity := &model.TeamStat{
		ID:           entity.ID,
		Wins:         entity.Wins,
		Draws:        entity.Draws,
		Losses:       entity.Losses,
		GoalsFor:     entity.GoalsFor,
		GoalsAgainst: entity.GoalsAgainst,
		Points:       entity.Points,
		Rank:         entity.Rank,
		SeasonID:     entity.SeasonID,
		TeamID:       entity.TeamID,
		CreatedAt:    entity.CreatedAt,
		UpdatedAt:    entity.UpdatedAt,
	}

	// Convert related entities
	if entity.Team != nil {
		teamMapper := NewTeamMapper()
		modelEntity.Team = teamMapper.DomainToModel(entity.Team)
	}

	if entity.Season != nil {
		seasonMapper := NewSeasonMapper()
		modelEntity.Season = seasonMapper.DomainToModel(entity.Season)
	}

	return modelEntity
}

func (m *TeamStatsMapper) ModelToDomain(model *model.TeamStat) *domain.TeamStats {
	domainEntity := &domain.TeamStats{
		ID:           model.ID,
		Wins:         model.Wins,
		Draws:        model.Draws,
		Losses:       model.Losses,
		GoalsFor:     model.GoalsFor,
		GoalsAgainst: model.GoalsAgainst,
		Points:       model.Points,
		Rank:         model.Rank,
		SeasonID:     model.SeasonID,
		TeamID:       model.TeamID,
		CreatedAt:    model.CreatedAt,
		UpdatedAt:    model.UpdatedAt,
	}

	// Convert related entities
	if model.Team != nil {
		teamMapper := NewTeamMapper()
		domainEntity.Team = teamMapper.ModelToDomain(model.Team)
	}

	if model.Season != nil {
		seasonMapper := NewSeasonMapper()
		domainEntity.Season = seasonMapper.ModelToDomain(model.Season)
	}

	return domainEntity
}

func (m *TeamStatsMapper) ModelListToDomain(models []model.TeamStat) []domain.TeamStats {
	entities := make([]domain.TeamStats, len(models))
	for i, model := range models {
		entities[i] = *m.ModelToDomain(&model)
	}
	return entities
}
