package mapper

import (
	"github.com/EdwinRincon/browersfc-api/api/dto"
	"github.com/EdwinRincon/browersfc-api/api/model"
	"github.com/EdwinRincon/browersfc-api/domain"
)

type LineupMapper struct{}

func NewLineupMapper() *LineupMapper {
	return &LineupMapper{}
}

// DTO to Domain Conversions (HTTP layer)
func (m *LineupMapper) DTOToDomain(dto *dto.CreateLineupRequest) *domain.Lineup {
	if dto == nil {
		return nil
	}

	return &domain.Lineup{
		Position: dto.Position,
		PlayerID: dto.PlayerID,
		MatchID:  dto.MatchID,
		Starting: dto.Starting,
	}
}

func (m *LineupMapper) UpdateDTOToDomain(dto *dto.UpdateLineupRequest) *domain.Lineup {
	if dto == nil {
		return nil
	}

	lineup := &domain.Lineup{}

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

	return lineup
}

func (m *LineupMapper) DomainToDTO(entity *domain.Lineup) *dto.LineupResponse {
	if entity == nil {
		return nil
	}

	response := &dto.LineupResponse{
		ID:        entity.ID,
		Position:  entity.Position,
		PlayerID:  entity.PlayerID,
		MatchID:   entity.MatchID,
		Starting:  entity.Starting,
		CreatedAt: entity.CreatedAt,
		UpdatedAt: entity.UpdatedAt,
	}

	// Map related entities if available
	if entity.Player != nil {
		playerMapper := NewPlayerMapper()
		response.Player = *playerMapper.DomainToShortDTO(entity.Player)
	}

	if entity.Match != nil {
		matchMapper := NewMatchMapper()
		response.Match = *matchMapper.DomainToShortDTO(entity.Match)
	}

	return response
}

func (m *LineupMapper) DomainListToDTO(entities []domain.Lineup) []dto.LineupResponse {
	if len(entities) == 0 {
		return nil
	}

	responses := make([]dto.LineupResponse, len(entities))
	for i, entity := range entities {
		if response := m.DomainToDTO(&entity); response != nil {
			responses[i] = *response
		}
	}
	return responses
}

func (m *LineupMapper) DomainToShortDTO(entity *domain.Lineup) *dto.LineupShortResponse {
	if entity == nil {
		return nil
	}

	return &dto.LineupShortResponse{
		ID:       entity.ID,
		Position: entity.Position,
		PlayerID: entity.PlayerID,
		Starting: entity.Starting,
	}
}

// Domain to Model Conversions (Infrastructure layer)
func (m *LineupMapper) DomainToModel(entity *domain.Lineup) *model.Lineup {
	if entity == nil {
		return nil
	}

	return &model.Lineup{
		ID:        entity.ID,
		Position:  entity.Position,
		PlayerID:  entity.PlayerID,
		MatchID:   entity.MatchID,
		Starting:  entity.Starting,
		CreatedAt: entity.CreatedAt,
		UpdatedAt: entity.UpdatedAt,
	}
}

func (m *LineupMapper) ModelToDomain(model *model.Lineup) *domain.Lineup {
	if model == nil {
		return nil
	}

	entity := &domain.Lineup{
		ID:        model.ID,
		Position:  model.Position,
		PlayerID:  model.PlayerID,
		MatchID:   model.MatchID,
		Starting:  model.Starting,
		CreatedAt: model.CreatedAt,
		UpdatedAt: model.UpdatedAt,
	}

	// Map related entities if available
	if model.Player != nil {
		playerMapper := NewPlayerMapper()
		entity.Player = playerMapper.ModelToDomain(model.Player)
	}

	if model.Match != nil {
		matchMapper := NewMatchMapper()
		entity.Match = matchMapper.ModelToDomain(model.Match)
	}

	return entity
}

func (m *LineupMapper) ModelListToDomain(models []model.Lineup) []domain.Lineup {
	if len(models) == 0 {
		return nil
	}

	entities := make([]domain.Lineup, len(models))
	for i, model := range models {
		if entity := m.ModelToDomain(&model); entity != nil {
			entities[i] = *entity
		}
	}
	return entities
}
