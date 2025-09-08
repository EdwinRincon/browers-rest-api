package mapper

import (
	"github.com/EdwinRincon/browersfc-api/api/dto"
	"github.com/EdwinRincon/browersfc-api/api/model"
	"github.com/EdwinRincon/browersfc-api/domain"
)

type PlayerTeamMapper struct{}

func NewPlayerTeamMapper() *PlayerTeamMapper {
	return &PlayerTeamMapper{}
}

// DTO to Domain Conversions (HTTP layer)
func (m *PlayerTeamMapper) DTOToDomain(dto *dto.CreatePlayerTeamRequest) *domain.PlayerTeam {
	if dto == nil {
		return nil
	}

	return &domain.PlayerTeam{
		PlayerID:  dto.PlayerID,
		TeamID:    dto.TeamID,
		SeasonID:  dto.SeasonID,
		StartDate: dto.StartDate,
		EndDate:   dto.EndDate,
	}
}

func (m *PlayerTeamMapper) UpdateDTOToDomain(dto *dto.UpdatePlayerTeamRequest) *domain.PlayerTeam {
	if dto == nil {
		return nil
	}

	playerTeam := &domain.PlayerTeam{}

	if dto.StartDate != nil {
		playerTeam.StartDate = *dto.StartDate
	}
	if dto.EndDate != nil {
		playerTeam.EndDate = dto.EndDate
	}

	return playerTeam
}

func (m *PlayerTeamMapper) DomainToDTO(entity *domain.PlayerTeam) *dto.PlayerTeamResponse {
	if entity == nil {
		return nil
	}

	response := &dto.PlayerTeamResponse{
		ID:        entity.ID,
		PlayerID:  entity.PlayerID,
		TeamID:    entity.TeamID,
		SeasonID:  entity.SeasonID,
		StartDate: entity.StartDate,
		EndDate:   entity.EndDate,
		CreatedAt: entity.CreatedAt,
		UpdatedAt: entity.UpdatedAt,
	}

	// Handle Player relationship
	if entity.Player != nil {
		response.Player = dto.PlayerShort{
			ID:       entity.Player.ID,
			NickName: entity.Player.NickName,
			Position: entity.Player.Position,
		}
	}

	// Handle Team relationship
	if entity.Team != nil {
		response.Team = dto.TeamShort{
			ID:        entity.Team.ID,
			FullName:  entity.Team.FullName,
			ShortName: entity.Team.ShortName,
		}
	}

	// Handle Season relationship
	if entity.Season != nil {
		response.Season = dto.SeasonShort{
			ID:   entity.Season.ID,
			Year: entity.Season.Year,
		}
	}

	return response
}

func (m *PlayerTeamMapper) DomainListToDTO(entities []domain.PlayerTeam) []dto.PlayerTeamResponse {
	if entities == nil {
		return nil
	}

	result := make([]dto.PlayerTeamResponse, len(entities))
	for i, entity := range entities {
		response := m.DomainToDTO(&entity)
		if response != nil {
			result[i] = *response
		}
	}
	return result
}

// Domain to Model Conversions (Infrastructure layer)
func (m *PlayerTeamMapper) DomainToModel(entity *domain.PlayerTeam) *model.PlayerTeam {
	if entity == nil {
		return nil
	}

	modelEntity := &model.PlayerTeam{
		ID:        entity.ID,
		PlayerID:  entity.PlayerID,
		TeamID:    entity.TeamID,
		SeasonID:  entity.SeasonID,
		StartDate: entity.StartDate,
		EndDate:   entity.EndDate,
		CreatedAt: entity.CreatedAt,
		UpdatedAt: entity.UpdatedAt,
	}

	// Convert related entities if present
	if entity.Player != nil {
		playerMapper := NewPlayerMapper()
		modelEntity.Player = playerMapper.DomainToModel(entity.Player)
	}

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

func (m *PlayerTeamMapper) ModelToDomain(model *model.PlayerTeam) *domain.PlayerTeam {
	if model == nil {
		return nil
	}

	domainEntity := &domain.PlayerTeam{
		ID:        model.ID,
		PlayerID:  model.PlayerID,
		TeamID:    model.TeamID,
		SeasonID:  model.SeasonID,
		StartDate: model.StartDate,
		EndDate:   model.EndDate,
		CreatedAt: model.CreatedAt,
		UpdatedAt: model.UpdatedAt,
	}

	// Convert related entities if present
	if model.Player != nil {
		playerMapper := NewPlayerMapper()
		domainEntity.Player = playerMapper.ModelToDomain(model.Player)
	}

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

func (m *PlayerTeamMapper) ModelListToDomain(models []model.PlayerTeam) []domain.PlayerTeam {
	if models == nil {
		return nil
	}

	result := make([]domain.PlayerTeam, len(models))
	for i, model := range models {
		domain := m.ModelToDomain(&model)
		if domain != nil {
			result[i] = *domain
		}
	}
	return result
}

func (m *PlayerTeamMapper) DomainToShortDTO(entity *domain.PlayerTeam) *dto.PlayerTeamShort {
	if entity == nil {
		return nil
	}

	return &dto.PlayerTeamShort{
		ID:        entity.ID,
		PlayerID:  entity.PlayerID,
		TeamID:    entity.TeamID,
		SeasonID:  entity.SeasonID,
		StartDate: entity.StartDate,
		EndDate:   entity.EndDate,
	}
}
