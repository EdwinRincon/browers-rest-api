package mapper

import (
	"time"

	"github.com/EdwinRincon/browersfc-api/api/dto"
	"github.com/EdwinRincon/browersfc-api/api/model"
	"github.com/EdwinRincon/browersfc-api/domain"
)

type TeamMapper struct{}

func NewTeamMapper() *TeamMapper {
	return &TeamMapper{}
}

// DTO to Domain Conversions (HTTP layer)
func (m *TeamMapper) DTOToDomain(dto *dto.CreateTeamRequest) *domain.Team {
	if dto == nil {
		return nil
	}

	return &domain.Team{
		FullName:    dto.FullName,
		ShortName:   dto.ShortName,
		Color:       dto.Color,
		Color2:      dto.Color2,
		Shield:      dto.Shield,
		NextMatchID: dto.NextMatchID,
		CreatedAt:   time.Now(),
		UpdatedAt:   time.Now(),
	}
}

func (m *TeamMapper) UpdateDTOToDomain(dto *dto.UpdateTeamRequest) *domain.Team {
	if dto == nil {
		return nil
	}

	team := &domain.Team{}

	if dto.FullName != nil {
		team.FullName = *dto.FullName
	}
	if dto.ShortName != nil {
		team.ShortName = *dto.ShortName
	}
	if dto.Color != nil {
		team.Color = *dto.Color
	}
	if dto.Color2 != nil {
		team.Color2 = *dto.Color2
	}
	if dto.Shield != nil {
		team.Shield = *dto.Shield
	}
	if dto.NextMatchID != nil {
		team.NextMatchID = dto.NextMatchID
	}

	return team
}

func (m *TeamMapper) DomainToDTO(entity *domain.Team) *dto.TeamResponse {
	if entity == nil {
		return nil
	}

	return &dto.TeamResponse{
		ID:          entity.ID,
		FullName:    entity.FullName,
		ShortName:   entity.ShortName,
		Color:       entity.Color,
		Color2:      entity.Color2,
		Shield:      entity.Shield,
		NextMatchID: entity.NextMatchID,
		CreatedAt:   entity.CreatedAt,
		UpdatedAt:   entity.UpdatedAt,
	}
}

func (m *TeamMapper) DomainListToDTO(entities []domain.Team) []dto.TeamResponse {
	if entities == nil {
		return nil
	}

	responses := make([]dto.TeamResponse, len(entities))
	for i, entity := range entities {
		response := m.DomainToDTO(&entity)
		if response != nil {
			responses[i] = *response
		}
	}

	return responses
}

// Domain to Model Conversions (Infrastructure layer)
func (m *TeamMapper) DomainToModel(entity *domain.Team) *model.Team {
	if entity == nil {
		return nil
	}

	return &model.Team{
		ID:          entity.ID,
		FullName:    entity.FullName,
		ShortName:   entity.ShortName,
		Color:       entity.Color,
		Color2:      entity.Color2,
		Shield:      entity.Shield,
		NextMatchID: entity.NextMatchID,
		CreatedAt:   entity.CreatedAt,
		UpdatedAt:   entity.UpdatedAt,
	}
}

func (m *TeamMapper) ModelToDomain(model *model.Team) *domain.Team {
	if model == nil {
		return nil
	}

	return &domain.Team{
		ID:          model.ID,
		FullName:    model.FullName,
		ShortName:   model.ShortName,
		Color:       model.Color,
		Color2:      model.Color2,
		Shield:      model.Shield,
		NextMatchID: model.NextMatchID,
		CreatedAt:   model.CreatedAt,
		UpdatedAt:   model.UpdatedAt,
	}
}

func (m *TeamMapper) ModelListToDomain(models []model.Team) []domain.Team {
	if models == nil {
		return nil
	}

	domains := make([]domain.Team, len(models))
	for i, model := range models {
		domain := m.ModelToDomain(&model)
		if domain != nil {
			domains[i] = *domain
		}
	}

	return domains
}

// Legacy Support (For backward compatibility)
func (m *TeamMapper) ModelToShortDTO(model *model.Team) *dto.TeamShort {
	if model == nil {
		return nil
	}

	return &dto.TeamShort{
		ID:        model.ID,
		FullName:  model.FullName,
		ShortName: model.ShortName,
	}
}

func (m *TeamMapper) DomainToShortDTO(entity *domain.Team) *dto.TeamShort {
	if entity == nil {
		return nil
	}

	return &dto.TeamShort{
		ID:        entity.ID,
		FullName:  entity.FullName,
		ShortName: entity.ShortName,
	}
}
