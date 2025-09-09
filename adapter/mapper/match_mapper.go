package mapper

import (
	"github.com/EdwinRincon/browersfc-api/api/dto"
	"github.com/EdwinRincon/browersfc-api/internal/infrastructure/persistence/model"
	"github.com/EdwinRincon/browersfc-api/domain"
)

type MatchMapper struct{}

func NewMatchMapper() *MatchMapper {
	return &MatchMapper{}
}

// DTO to Domain Conversions (HTTP layer)
func (m *MatchMapper) DTOToDomain(dto *dto.CreateMatchRequest) *domain.Match {
	if dto == nil {
		return nil
	}

	return &domain.Match{
		Status:      dto.Status,
		Kickoff:     dto.Kickoff,
		Location:    dto.Location,
		HomeGoals:   dto.HomeGoals,
		AwayGoals:   dto.AwayGoals,
		HomeTeamID:  dto.HomeTeamID,
		AwayTeamID:  dto.AwayTeamID,
		SeasonID:    dto.SeasonID,
		MVPPlayerID: dto.MVPPlayerID,
	}
}

func (m *MatchMapper) UpdateDTOToDomain(dto *dto.UpdateMatchRequest) *domain.Match {
	if dto == nil {
		return nil
	}

	match := &domain.Match{}

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

	return match
}

func (m *MatchMapper) DomainToDTO(entity *domain.Match) *dto.MatchResponse {
	if entity == nil {
		return nil
	}

	resp := &dto.MatchResponse{
		ID:        entity.ID,
		Status:    entity.Status,
		Kickoff:   entity.Kickoff,
		Location:  entity.Location,
		HomeGoals: entity.HomeGoals,
		AwayGoals: entity.AwayGoals,
		CreatedAt: entity.CreatedAt,
		UpdatedAt: entity.UpdatedAt,
	}

	// Map related entities if available
	if entity.HomeTeam != nil {
		teamMapper := NewTeamMapper()
		homeTeam := teamMapper.DomainToShortDTO(entity.HomeTeam)
		if homeTeam != nil {
			resp.HomeTeam = *homeTeam
		}
	}

	if entity.AwayTeam != nil {
		teamMapper := NewTeamMapper()
		awayTeam := teamMapper.DomainToShortDTO(entity.AwayTeam)
		if awayTeam != nil {
			resp.AwayTeam = *awayTeam
		}
	}

	if entity.Season != nil {
		seasonMapper := NewSeasonMapper()
		season := seasonMapper.DomainToShortDTO(entity.Season)
		if season != nil {
			resp.Season = *season
		}
	}

	if entity.MVPPlayer != nil {
		playerMapper := NewPlayerMapper()
		mvpPlayer := playerMapper.DomainToShortDTO(entity.MVPPlayer)
		resp.MVPPlayer = mvpPlayer
	}

	return resp
}

func (m *MatchMapper) DomainListToDTO(entities []domain.Match) []dto.MatchResponse {
	if len(entities) == 0 {
		return nil
	}

	responses := make([]dto.MatchResponse, len(entities))
	for i, match := range entities {
		if response := m.DomainToDTO(&match); response != nil {
			responses[i] = *response
		}
	}
	return responses
}

func (m *MatchMapper) DomainToShortDTO(entity *domain.Match) *dto.MatchShort {
	if entity == nil {
		return nil
	}

	return &dto.MatchShort{
		ID:        entity.ID,
		Status:    entity.Status,
		Kickoff:   entity.Kickoff,
		Location:  entity.Location,
		HomeGoals: entity.HomeGoals,
		AwayGoals: entity.AwayGoals,
	}
}

func (m *MatchMapper) DomainToDetailDTO(entity *domain.Match) *dto.MatchDetailResponse {
	if entity == nil {
		return nil
	}

	// Start with the basic match data
	detail := &dto.MatchDetailResponse{
		ID:        entity.ID,
		Status:    entity.Status,
		Kickoff:   entity.Kickoff,
		Location:  entity.Location,
		HomeGoals: entity.HomeGoals,
		AwayGoals: entity.AwayGoals,
		CreatedAt: entity.CreatedAt,
		UpdatedAt: entity.UpdatedAt,
	}

	// Map related entities if available
	if entity.HomeTeam != nil {
		teamMapper := NewTeamMapper()
		if homeTeam := teamMapper.DomainToShortDTO(entity.HomeTeam); homeTeam != nil {
			detail.HomeTeam = *homeTeam
		}
	}

	if entity.AwayTeam != nil {
		teamMapper := NewTeamMapper()
		if awayTeam := teamMapper.DomainToShortDTO(entity.AwayTeam); awayTeam != nil {
			detail.AwayTeam = *awayTeam
		}
	}

	if entity.Season != nil {
		seasonMapper := NewSeasonMapper()
		if season := seasonMapper.DomainToShortDTO(entity.Season); season != nil {
			detail.Season = *season
		}
	}

	if entity.MVPPlayer != nil {
		playerMapper := NewPlayerMapper()
		detail.MVPPlayer = playerMapper.DomainToShortDTO(entity.MVPPlayer)
	}

	return detail
}

// Domain to Model Conversions (Infrastructure layer)
func (m *MatchMapper) DomainToModel(entity *domain.Match) *model.Match {
	if entity == nil {
		return nil
	}

	return &model.Match{
		ID:          entity.ID,
		Status:      entity.Status,
		Kickoff:     entity.Kickoff,
		Location:    entity.Location,
		HomeGoals:   entity.HomeGoals,
		AwayGoals:   entity.AwayGoals,
		HomeTeamID:  entity.HomeTeamID,
		AwayTeamID:  entity.AwayTeamID,
		SeasonID:    entity.SeasonID,
		MVPPlayerID: entity.MVPPlayerID,
		CreatedAt:   entity.CreatedAt,
		UpdatedAt:   entity.UpdatedAt,
	}
}

func (m *MatchMapper) ModelToDomain(model *model.Match) *domain.Match {
	if model == nil {
		return nil
	}

	entity := &domain.Match{
		ID:          model.ID,
		Status:      model.Status,
		Kickoff:     model.Kickoff,
		Location:    model.Location,
		HomeGoals:   model.HomeGoals,
		AwayGoals:   model.AwayGoals,
		HomeTeamID:  model.HomeTeamID,
		AwayTeamID:  model.AwayTeamID,
		SeasonID:    model.SeasonID,
		MVPPlayerID: model.MVPPlayerID,
		CreatedAt:   model.CreatedAt,
		UpdatedAt:   model.UpdatedAt,
	}

	// Map related entities if available
	if model.HomeTeam != nil {
		teamMapper := NewTeamMapper()
		entity.HomeTeam = teamMapper.ModelToDomain(model.HomeTeam)
	}

	if model.AwayTeam != nil {
		teamMapper := NewTeamMapper()
		entity.AwayTeam = teamMapper.ModelToDomain(model.AwayTeam)
	}

	if model.Season != nil {
		seasonMapper := NewSeasonMapper()
		entity.Season = seasonMapper.ModelToDomain(model.Season)
	}

	if model.MVPPlayer != nil {
		playerMapper := NewPlayerMapper()
		entity.MVPPlayer = playerMapper.ModelToDomain(model.MVPPlayer)
	}

	return entity
}

func (m *MatchMapper) ModelListToDomain(models []model.Match) []domain.Match {
	if len(models) == 0 {
		return nil
	}

	entities := make([]domain.Match, len(models))
	for i, model := range models {
		if entity := m.ModelToDomain(&model); entity != nil {
			entities[i] = *entity
		}
	}
	return entities
}
