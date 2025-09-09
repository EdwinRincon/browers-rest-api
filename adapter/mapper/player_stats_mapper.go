package mapper

import (
	"github.com/EdwinRincon/browersfc-api/api/dto"
	"github.com/EdwinRincon/browersfc-api/internal/infrastructure/persistence/model"
	"github.com/EdwinRincon/browersfc-api/domain"
)

type PlayerStatsMapper struct{}

func NewPlayerStatsMapper() *PlayerStatsMapper {
	return &PlayerStatsMapper{}
}

// DTO to Domain Conversions (HTTP layer)
func (m *PlayerStatsMapper) DTOToDomain(dto *dto.CreatePlayerStatRequest) *domain.PlayerStat {
	if dto == nil {
		return nil
	}

	return &domain.PlayerStat{
		PlayerID:      dto.PlayerID,
		MatchID:       dto.MatchID,
		SeasonID:      dto.SeasonID,
		TeamID:        dto.TeamID,
		Goals:         dto.Goals,
		Assists:       dto.Assists,
		Saves:         dto.Saves,
		YellowCards:   dto.YellowCards,
		RedCards:      dto.RedCards,
		Rating:        dto.Rating,
		IsStarting:    dto.Starting,
		MinutesPlayed: dto.MinutesPlayed,
		IsMVP:         dto.IsMVP,
		Position:      dto.Position,
	}
}

func (m *PlayerStatsMapper) UpdateDTOToDomain(dto *dto.UpdatePlayerStatRequest) *domain.PlayerStat {
	if dto == nil {
		return nil
	}

	playerStat := &domain.PlayerStat{}

	if dto.TeamID != nil {
		playerStat.TeamID = dto.TeamID
	}
	if dto.Goals != nil {
		playerStat.Goals = *dto.Goals
	}
	if dto.Assists != nil {
		playerStat.Assists = *dto.Assists
	}
	if dto.Saves != nil {
		playerStat.Saves = *dto.Saves
	}
	if dto.YellowCards != nil {
		playerStat.YellowCards = *dto.YellowCards
	}
	if dto.RedCards != nil {
		playerStat.RedCards = *dto.RedCards
	}
	if dto.Rating != nil {
		playerStat.Rating = *dto.Rating
	}
	if dto.Starting != nil {
		playerStat.IsStarting = *dto.Starting
	}
	if dto.MinutesPlayed != nil {
		playerStat.MinutesPlayed = *dto.MinutesPlayed
	}
	if dto.IsMVP != nil {
		playerStat.IsMVP = *dto.IsMVP
	}
	if dto.Position != nil {
		playerStat.Position = *dto.Position
	}

	return playerStat
}

func (m *PlayerStatsMapper) DomainToDTO(entity *domain.PlayerStat) *dto.PlayerStatResponse {
	if entity == nil {
		return nil
	}

	return &dto.PlayerStatResponse{
		ID:            entity.ID,
		PlayerID:      entity.PlayerID,
		MatchID:       entity.MatchID,
		SeasonID:      entity.SeasonID,
		TeamID:        entity.TeamID,
		Goals:         entity.Goals,
		Assists:       entity.Assists,
		Saves:         entity.Saves,
		YellowCards:   entity.YellowCards,
		RedCards:      entity.RedCards,
		Rating:        entity.Rating,
		Starting:      entity.IsStarting,
		MinutesPlayed: entity.MinutesPlayed,
		IsMVP:         entity.IsMVP,
		Position:      entity.Position,
		CreatedAt:     entity.CreatedAt,
		UpdatedAt:     entity.UpdatedAt,
	}
}

func (m *PlayerStatsMapper) DomainListToDTO(entities []domain.PlayerStat) []dto.PlayerStatResponse {
	if entities == nil {
		return nil
	}

	result := make([]dto.PlayerStatResponse, len(entities))
	for i, entity := range entities {
		response := m.DomainToDTO(&entity)
		if response != nil {
			result[i] = *response
		}
	}
	return result
}

// Domain to Model Conversions (Infrastructure layer)
func (m *PlayerStatsMapper) DomainToModel(entity *domain.PlayerStat) *model.PlayerStat {
	if entity == nil {
		return nil
	}

	return &model.PlayerStat{
		ID:            entity.ID,
		PlayerID:      entity.PlayerID,
		MatchID:       entity.MatchID,
		SeasonID:      entity.SeasonID,
		TeamID:        entity.TeamID,
		Goals:         entity.Goals,
		Assists:       entity.Assists,
		Saves:         entity.Saves,
		YellowCards:   entity.YellowCards,
		RedCards:      entity.RedCards,
		Rating:        entity.Rating,
		IsStarting:    entity.IsStarting,
		MinutesPlayed: entity.MinutesPlayed,
		IsMVP:         entity.IsMVP,
		Position:      entity.Position,
		CreatedAt:     entity.CreatedAt,
		UpdatedAt:     entity.UpdatedAt,
	}
}

func (m *PlayerStatsMapper) ModelToDomain(model *model.PlayerStat) *domain.PlayerStat {
	if model == nil {
		return nil
	}

	return &domain.PlayerStat{
		ID:            model.ID,
		PlayerID:      model.PlayerID,
		MatchID:       model.MatchID,
		SeasonID:      model.SeasonID,
		TeamID:        model.TeamID,
		Goals:         model.Goals,
		Assists:       model.Assists,
		Saves:         model.Saves,
		YellowCards:   model.YellowCards,
		RedCards:      model.RedCards,
		Rating:        model.Rating,
		IsStarting:    model.IsStarting,
		MinutesPlayed: model.MinutesPlayed,
		IsMVP:         model.IsMVP,
		Position:      model.Position,
		CreatedAt:     model.CreatedAt,
		UpdatedAt:     model.UpdatedAt,
	}
}

func (m *PlayerStatsMapper) ModelListToDomain(models []model.PlayerStat) []domain.PlayerStat {
	if models == nil {
		return nil
	}

	result := make([]domain.PlayerStat, len(models))
	for i, model := range models {
		domain := m.ModelToDomain(&model)
		if domain != nil {
			result[i] = *domain
		}
	}
	return result
}
