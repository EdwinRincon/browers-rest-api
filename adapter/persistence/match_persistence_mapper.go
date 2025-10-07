package persistence

import (
	"github.com/EdwinRincon/browersfc-api/domain"
	"github.com/EdwinRincon/browersfc-api/internal/infrastructure/persistence/model"
)

type MatchPersistenceMapper struct{}

func NewMatchPersistenceMapper() *MatchPersistenceMapper {
	return &MatchPersistenceMapper{}
}

// Domain to Model Conversions (Infrastructure layer)
func (m *MatchPersistenceMapper) DomainToModel(entity *domain.Match) *model.Match {
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

func (m *MatchPersistenceMapper) ModelToDomain(model *model.Match) *domain.Match {
	if model == nil {
		return nil
	}

	domainMatch := &domain.Match{
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

	// Map preloaded relationships if they exist
	if model.HomeTeam != nil {
		teamMapper := NewTeamPersistenceMapper()
		domainMatch.HomeTeam = teamMapper.ModelToDomain(model.HomeTeam)
	}

	if model.AwayTeam != nil {
		teamMapper := NewTeamPersistenceMapper()
		domainMatch.AwayTeam = teamMapper.ModelToDomain(model.AwayTeam)
	}

	if model.Season != nil {
		seasonMapper := NewSeasonPersistenceMapper()
		domainMatch.Season = seasonMapper.ModelToDomain(model.Season)
	}

	if model.MVPPlayer != nil {
		playerMapper := NewPlayerPersistenceMapper()
		domainMatch.MVPPlayer = playerMapper.ModelToDomain(model.MVPPlayer)
	}

	return domainMatch
}

func (m *MatchPersistenceMapper) ModelListToDomain(models []model.Match) []domain.Match {
	if models == nil {
		return nil
	}

	domains := make([]domain.Match, len(models))
	for i, model := range models {
		domain := m.ModelToDomain(&model)
		if domain != nil {
			domains[i] = *domain
		}
	}

	return domains
}
