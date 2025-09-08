package mapper

import (
	"github.com/EdwinRincon/browersfc-api/api/dto"
	"github.com/EdwinRincon/browersfc-api/api/model"
	"github.com/EdwinRincon/browersfc-api/domain"
)

type PlayerMapper struct{}

func NewPlayerMapper() *PlayerMapper {
	return &PlayerMapper{}
}

// DTO to Domain Conversions (HTTP layer)
func (m *PlayerMapper) DTOToDomain(dto *dto.CreatePlayerRequest) *domain.Player {
	if dto == nil {
		return nil
	}

	return &domain.Player{
		NickName:      dto.NickName,
		Height:        dto.Height,
		Country:       dto.Country,
		Country2:      dto.Country2,
		Foot:          dto.Foot,
		Age:           dto.Age,
		SquadNumber:   dto.SquadNumber,
		Position:      dto.Position,
		CareerSummary: dto.CareerSummary,
		UserID:        dto.UserID,
	}
}

func (m *PlayerMapper) UpdateDTOToDomain(dto *dto.UpdatePlayerRequest) *domain.Player {
	if dto == nil {
		return nil
	}

	player := &domain.Player{}

	if dto.NickName != nil {
		player.NickName = *dto.NickName
	}
	if dto.Height != nil {
		player.Height = *dto.Height
	}
	if dto.Country != nil {
		player.Country = *dto.Country
	}
	if dto.Country2 != nil {
		player.Country2 = *dto.Country2
	}
	if dto.Foot != nil {
		player.Foot = *dto.Foot
	}
	if dto.Age != nil {
		player.Age = *dto.Age
	}
	if dto.SquadNumber != nil {
		player.SquadNumber = *dto.SquadNumber
	}
	if dto.Rating != nil {
		player.Rating = *dto.Rating
	}
	if dto.Position != nil {
		player.Position = *dto.Position
	}
	if dto.Injured != nil {
		player.Injured = *dto.Injured
	}
	if dto.CareerSummary != nil {
		player.CareerSummary = *dto.CareerSummary
	}
	if dto.UserID != nil {
		player.UserID = dto.UserID
	}

	return player
}

func (m *PlayerMapper) DomainToDTO(entity *domain.Player) *dto.PlayerResponse {
	if entity == nil {
		return nil
	}

	return &dto.PlayerResponse{
		ID:            entity.ID,
		NickName:      entity.NickName,
		Height:        entity.Height,
		Country:       entity.Country,
		Country2:      entity.Country2,
		Foot:          entity.Foot,
		Age:           entity.Age,
		SquadNumber:   entity.SquadNumber,
		Rating:        entity.Rating,
		Matches:       entity.Matches,
		YCards:        entity.YCards,
		RCards:        entity.RCards,
		Goals:         entity.Goals,
		Assists:       entity.Assists,
		Saves:         entity.Saves,
		Position:      entity.Position,
		Injured:       entity.Injured,
		CareerSummary: entity.CareerSummary,
		MVPCount:      entity.MVPCount,
		CreatedAt:     entity.CreatedAt,
		UpdatedAt:     entity.UpdatedAt,
	}
}

func (m *PlayerMapper) DomainListToDTO(entities []domain.Player) []dto.PlayerResponse {
	if entities == nil {
		return nil
	}

	result := make([]dto.PlayerResponse, len(entities))
	for i, entity := range entities {
		response := m.DomainToDTO(&entity)
		if response != nil {
			result[i] = *response
		}
	}
	return result
}

// Domain to Model Conversions (Infrastructure layer)
func (m *PlayerMapper) DomainToModel(entity *domain.Player) *model.Player {
	if entity == nil {
		return nil
	}

	return &model.Player{
		ID:            entity.ID,
		NickName:      entity.NickName,
		Height:        entity.Height,
		Country:       entity.Country,
		Country2:      entity.Country2,
		Foot:          entity.Foot,
		Age:           entity.Age,
		SquadNumber:   entity.SquadNumber,
		Rating:        entity.Rating,
		Matches:       entity.Matches,
		YCards:        entity.YCards,
		RCards:        entity.RCards,
		Goals:         entity.Goals,
		Assists:       entity.Assists,
		Saves:         entity.Saves,
		Position:      entity.Position,
		Injured:       entity.Injured,
		CareerSummary: entity.CareerSummary,
		MVPCount:      entity.MVPCount,
		UserID:        entity.UserID,
		CreatedAt:     entity.CreatedAt,
		UpdatedAt:     entity.UpdatedAt,
	}
}

func (m *PlayerMapper) ModelToDomain(model *model.Player) *domain.Player {
	if model == nil {
		return nil
	}

	return &domain.Player{
		ID:            model.ID,
		NickName:      model.NickName,
		Height:        model.Height,
		Country:       model.Country,
		Country2:      model.Country2,
		Foot:          model.Foot,
		Age:           model.Age,
		SquadNumber:   model.SquadNumber,
		Rating:        model.Rating,
		Matches:       model.Matches,
		YCards:        model.YCards,
		RCards:        model.RCards,
		Goals:         model.Goals,
		Assists:       model.Assists,
		Saves:         model.Saves,
		Position:      model.Position,
		Injured:       model.Injured,
		CareerSummary: model.CareerSummary,
		MVPCount:      model.MVPCount,
		UserID:        model.UserID,
		CreatedAt:     model.CreatedAt,
		UpdatedAt:     model.UpdatedAt,
	}
}

func (m *PlayerMapper) ModelListToDomain(models []model.Player) []domain.Player {
	if models == nil {
		return nil
	}

	result := make([]domain.Player, len(models))
	for i, model := range models {
		domain := m.ModelToDomain(&model)
		if domain != nil {
			result[i] = *domain
		}
	}
	return result
}

// Legacy Support (For backward compatibility)
func (m *PlayerMapper) ModelToShortDTO(model *model.Player) *dto.PlayerShort {
	if model == nil {
		return nil
	}

	return &dto.PlayerShort{
		ID:       model.ID,
		NickName: model.NickName,
		Position: model.Position,
	}
}

func (m *PlayerMapper) DomainToShortDTO(entity *domain.Player) *dto.PlayerShort {
	if entity == nil {
		return nil
	}

	return &dto.PlayerShort{
		ID:       entity.ID,
		NickName: entity.NickName,
		Position: entity.Position,
	}
}

func (m *PlayerMapper) ModelToStatsDTO(model *model.Player) *dto.PlayerStats {
	if model == nil {
		return nil
	}

	return &dto.PlayerStats{
		ID:       model.ID,
		NickName: model.NickName,
		Matches:  model.Matches,
		Goals:    model.Goals,
		Assists:  model.Assists,
		YCards:   model.YCards,
		RCards:   model.RCards,
		Saves:    model.Saves,
		Position: model.Position,
		MVPCount: model.MVPCount,
	}
}

// Global legacy functions for backward compatibility
// These maintain the same signatures as the old api/mapper functions

var globalPlayerMapper = NewPlayerMapper()

func ToPlayerShort(model *model.Player) *dto.PlayerShort {
	return globalPlayerMapper.ModelToShortDTO(model)
}

func CreatePlayerRequestToDomain(dto *dto.CreatePlayerRequest) *domain.Player {
	return globalPlayerMapper.DTOToDomain(dto)
}

func UpdatePlayerRequestToDomain(target *domain.Player, dto *dto.UpdatePlayerRequest) {
	if dto == nil || target == nil {
		return
	}

	if dto.NickName != nil {
		target.NickName = *dto.NickName
	}
	if dto.Height != nil {
		target.Height = *dto.Height
	}
	if dto.Country != nil {
		target.Country = *dto.Country
	}
	if dto.Country2 != nil {
		target.Country2 = *dto.Country2
	}
	if dto.Foot != nil {
		target.Foot = *dto.Foot
	}
	if dto.Age != nil {
		target.Age = *dto.Age
	}
	if dto.SquadNumber != nil {
		target.SquadNumber = *dto.SquadNumber
	}
	if dto.Rating != nil {
		target.Rating = *dto.Rating
	}
	if dto.Position != nil {
		target.Position = *dto.Position
	}
	if dto.Injured != nil {
		target.Injured = *dto.Injured
	}
	if dto.CareerSummary != nil {
		target.CareerSummary = *dto.CareerSummary
	}
	if dto.UserID != nil {
		target.UserID = dto.UserID
	}
}

func PlayerDomainToShort(entity *domain.Player) *dto.PlayerShort {
	return globalPlayerMapper.DomainToShortDTO(entity)
}
