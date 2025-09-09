package handler

import (
	"context"
	"errors"
	"net/http"
	"strconv"
	"time"

	httpMapper "github.com/EdwinRincon/browersfc-api/adapter/http"
	"github.com/EdwinRincon/browersfc-api/api/constants"
	"github.com/EdwinRincon/browersfc-api/api/dto"
	"github.com/EdwinRincon/browersfc-api/helper"
	domainservice "github.com/EdwinRincon/browersfc-api/internal/domain/service"
	"github.com/EdwinRincon/browersfc-api/internal/infrastructure/persistence/model"
	"github.com/gin-gonic/gin"
)

// Domain-specific constants for player teams
const (
	msgPlayerTeamRelationship             = "player-team relationship"
	msgPlayerTeamRelationshipRetrievedOK  = "Player-team relationship retrieved successfully"
	msgPlayerTeamRelationshipCreatedOK    = "Player-team relationship created successfully"
	msgPlayerTeamRelationshipUpdatedOK    = "Player-team relationship updated successfully"
	msgPlayerTeamRelationshipsRetrievedOK = "Player-team relationships retrieved successfully"
)

type PlayerTeamHandler struct {
	PlayerTeamDomainService *domainservice.PlayerTeamDomainService
	PlayerTeamMapper        *httpMapper.PlayerTeamHTTPMapper
}

func NewPlayerTeamHandler(playerTeamDomainService *domainservice.PlayerTeamDomainService) *PlayerTeamHandler {
	return &PlayerTeamHandler{
		PlayerTeamDomainService: playerTeamDomainService,
		PlayerTeamMapper:        httpMapper.NewPlayerTeamHTTPMapper(),
	}
}

// CreatePlayerTeam godoc
// @Summary Create a new player-team relationship
// @Tags playerTeams
// @ID createPlayerTeam
// @Accept json
// @Produce json
// @Param playerTeam body dto.CreatePlayerTeamRequest true "Player-Team relationship data"
// @Success 201 {object} dto.PlayerTeamResponse "Player-Team relationship created successfully"
// @Failure 400 {object} helper.AppError "Invalid input"
// @Failure 409 {object} helper.AppError "Conflict (e.g., date overlap)"
// @Failure 500 {object} helper.AppError "Internal server error"
// @Router /admin/player-teams [post]
// @Security BearerAuth
func (h *PlayerTeamHandler) CreatePlayerTeam(c *gin.Context) {
	var createRequest dto.CreatePlayerTeamRequest
	if err := c.ShouldBindJSON(&createRequest); err != nil {
		helper.WriteErrorResponse(c, helper.BuildValidationErrorFromBinding(err, "body", "Invalid player-team data"))
		return
	}

	// Wrap context with timeout for DB/service calls
	ctx, cancel := context.WithTimeout(c.Request.Context(), 5*time.Second)
	defer cancel()

	// Map DTO to domain
	playerTeam := h.PlayerTeamMapper.DTOToDomain(&createRequest)

	createdPlayerTeam, err := h.PlayerTeamDomainService.CreatePlayerTeam(ctx, playerTeam)
	if err != nil {
		switch {
		case errors.Is(err, constants.ErrPlayerNotFound):
			helper.WriteErrorResponse(c, helper.NewNotFoundError("player"))
			return
		case errors.Is(err, constants.ErrTeamNotFound):
			helper.WriteErrorResponse(c, helper.NewNotFoundError("team"))
			return
		case errors.Is(err, constants.ErrSeasonNotFound):
			helper.WriteErrorResponse(c, helper.NewNotFoundError("season"))
			return
		case errors.Is(err, constants.ErrOverlappingDates):
			helper.WriteErrorResponse(c, helper.NewConflictError("date_range", "Date range overlaps with an existing record"))
			return
		default:
			helper.WriteErrorResponse(c, helper.NewInternalServerError(err))
			return
		}
	}

	playerTeamResponse := h.PlayerTeamMapper.DomainToDTO(createdPlayerTeam)
	helper.WriteSuccessResponse(c, http.StatusCreated, playerTeamResponse, msgPlayerTeamRelationshipCreatedOK)
}

// GetPlayerTeamByID godoc
// @Summary Get player-team relationship by ID
// @Tags playerTeams
// @ID getPlayerTeamByID
// @Param id path int true "PlayerTeam ID"
// @Success 200 {object} dto.PlayerTeamResponse "Player-team relationship retrieved successfully"
// @Failure 400 {object} helper.AppError "Invalid input"
// @Failure 404 {object} helper.AppError "Not found"
// @Failure 500 {object} helper.AppError "Internal server error"
// @Router /player-teams/{id} [get]
// @Security BearerAuth
func (h *PlayerTeamHandler) GetPlayerTeamByID(c *gin.Context) {
	id, err := strconv.ParseUint(c.Param("id"), 10, 64)
	if err != nil {
		helper.WriteErrorResponse(c, helper.NewBadRequestError("id", constants.MsgInvalidIDSimple))
		return
	}

	// Wrap context with timeout for DB/service calls
	ctx, cancel := context.WithTimeout(c.Request.Context(), 5*time.Second)
	defer cancel()

	playerTeam, err := h.PlayerTeamDomainService.GetPlayerTeamByID(ctx, id)
	if err != nil {
		if errors.Is(err, constants.ErrRecordNotFound) {
			helper.WriteErrorResponse(c, helper.NewNotFoundError(msgPlayerTeamRelationship))
		} else {
			helper.WriteErrorResponse(c, helper.NewInternalServerError(err))
		}
		return
	}

	playerTeamResponse := h.PlayerTeamMapper.DomainToDTO(playerTeam)

	helper.WriteSuccessResponse(c, http.StatusOK, playerTeamResponse, msgPlayerTeamRelationshipRetrievedOK)
}

// GetPlayerTeamsByPlayerID godoc
// @Summary Get player-team relationships by player ID
// @Tags playerTeams
// @ID get-player-teams-by-player-id
// @Param id path int true "Player ID"
// @Produce json
// @Success 200 {array} dto.PlayerTeamResponse "Player-team relationships retrieved successfully"
// @Failure 400 {object} helper.AppError "Invalid input"
// @Failure 404 {object} helper.AppError "Player not found"
// @Failure 500 {object} helper.AppError "Internal server error"
// @Router /players/{id}/teams [get]
// @Security BearerAuth
func (h *PlayerTeamHandler) GetPlayerTeamsByPlayerID(c *gin.Context) {
	playerID, err := strconv.ParseUint(c.Param("id"), 10, 64)
	if err != nil {
		helper.WriteErrorResponse(c, helper.NewBadRequestError("id", "Invalid player ID"))
		return
	}

	// Wrap context with timeout for DB/service calls
	ctx, cancel := context.WithTimeout(c.Request.Context(), 5*time.Second)
	defer cancel()

	playerTeams, err := h.PlayerTeamDomainService.GetPlayerTeamsByPlayerID(ctx, playerID)
	if err != nil {
		if errors.Is(err, constants.ErrPlayerNotFound) {
			helper.WriteErrorResponse(c, helper.NewNotFoundError("player"))
			return
		}
		helper.WriteErrorResponse(c, helper.NewInternalServerError(err))
		return
	}

	playerTeamResponses := h.PlayerTeamMapper.DomainListToDTO(playerTeams)
	helper.WriteSuccessResponse(c, http.StatusOK, playerTeamResponses, msgPlayerTeamRelationshipRetrievedOK)
}

// GetPlayerTeamsByTeamID godoc
// @Summary Get player-team relationships by team ID
// @Tags playerTeams
// @ID getPlayerTeamsByTeamID
// @Param id path int true "Team ID"
// @Produce json
// @Success 200 {array} dto.PlayerTeamResponse "Player-team relationships retrieved successfully"
// @Failure 400 {object} helper.AppError "Invalid input"
// @Failure 404 {object} helper.AppError "Team not found"
// @Failure 500 {object} helper.AppError "Internal server error"
// @Router /teams/{id}/players [get]
// @Security BearerAuth
func (h *PlayerTeamHandler) GetPlayerTeamsByTeamID(c *gin.Context) {
	teamID, err := strconv.ParseUint(c.Param("id"), 10, 64)
	if err != nil {
		helper.WriteErrorResponse(c, helper.NewBadRequestError("id", "Invalid team ID"))
		return
	}

	// Wrap context with timeout for DB/service calls
	ctx, cancel := context.WithTimeout(c.Request.Context(), 5*time.Second)
	defer cancel()

	playerTeams, err := h.PlayerTeamDomainService.GetPlayerTeamsByTeamID(ctx, teamID)
	if err != nil {
		if errors.Is(err, constants.ErrTeamNotFound) {
			helper.WriteErrorResponse(c, helper.NewNotFoundError("team"))
			return
		}
		helper.WriteErrorResponse(c, helper.NewInternalServerError(err))
		return
	}

	playerTeamResponses := h.PlayerTeamMapper.DomainListToDTO(playerTeams)
	helper.WriteSuccessResponse(c, http.StatusOK, playerTeamResponses, msgPlayerTeamRelationshipRetrievedOK)
}

// GetPlayerTeamsBySeasonID godoc
// @Summary Get player-team relationships by season ID
// @Tags playerTeams
// @ID getPlayerTeamsBySeasonID
// @Param id path int true "Season ID"
// @Produce json
// @Success 200 {array} dto.PlayerTeamResponse "Player-team relationships retrieved successfully"
// @Failure 400 {object} helper.AppError "Invalid input"
// @Failure 404 {object} helper.AppError "Season not found"
// @Failure 500 {object} helper.AppError "Internal server error"
// @Router /seasons/{id}/player-teams [get]
// @Security BearerAuth
func (h *PlayerTeamHandler) GetPlayerTeamsBySeasonID(c *gin.Context) {
	seasonID, err := strconv.ParseUint(c.Param("id"), 10, 64)
	if err != nil {
		helper.WriteErrorResponse(c, helper.NewBadRequestError("id", "Invalid season ID"))
		return
	}

	// Wrap context with timeout for DB/service calls
	ctx, cancel := context.WithTimeout(c.Request.Context(), 5*time.Second)
	defer cancel()

	playerTeams, err := h.PlayerTeamDomainService.GetPlayerTeamsBySeasonID(ctx, seasonID)
	if err != nil {
		if errors.Is(err, constants.ErrSeasonNotFound) {
			helper.WriteErrorResponse(c, helper.NewNotFoundError("season"))
			return
		}
		helper.WriteErrorResponse(c, helper.NewInternalServerError(err))
		return
	}

	playerTeamResponses := h.PlayerTeamMapper.DomainListToDTO(playerTeams)
	helper.WriteSuccessResponse(c, http.StatusOK, playerTeamResponses, msgPlayerTeamRelationshipRetrievedOK)
}

// GetPaginatedPlayerTeams godoc
// @Summary Get paginated player-team relationships
// @Tags playerTeams
// @Accept json
// @Produce json
// @Param page query int false "Page number" default(0)
// @Param pageSize query int false "Page size" default(10)
// @Param sort query string false "Sort field"
// @Param order query string false "Sort order" Enums(asc, desc) default(desc)
// @Success 200 {object} helper.AppSuccess{data=helper.PaginatedResponse{items=[]dto.PlayerTeamResponse, totalCount=int}}
// @Failure 400 {object} helper.AppError "Invalid input"
// @Failure 500 {object} helper.AppError "Internal server error"
// @Router /player-teams [get]
// @Security BearerAuth
func (h *PlayerTeamHandler) GetPaginatedPlayerTeams(c *gin.Context) {
	sort := c.DefaultQuery("sort", "created_at")
	order := c.DefaultQuery("order", "desc")
	page, _ := strconv.Atoi(c.DefaultQuery("page", "0"))
	pageSize, _ := strconv.Atoi(c.DefaultQuery("pageSize", "10"))

	if page < 0 {
		page = 0
	}
	if pageSize < 1 || pageSize > 200 {
		pageSize = 10
	}
	if order != "asc" && order != "desc" {
		order = "asc"
	}

	// Validate sort field
	if err := helper.ValidateSort(model.PlayerTeam{}, sort); err != nil {
		helper.WriteErrorResponse(c, helper.NewBadRequestError("sort", err.Error()))
		return
	}

	// Wrap context with timeout for DB/service calls
	ctx, cancel := context.WithTimeout(c.Request.Context(), 5*time.Second)
	defer cancel()

	playerTeams, total, err := h.PlayerTeamDomainService.GetPaginatedPlayerTeams(ctx, sort, order, page, pageSize)
	if err != nil {
		helper.WriteErrorResponse(c, helper.NewInternalServerError(err))
		return
	}

	playerTeamResponses := h.PlayerTeamMapper.DomainListToDTO(playerTeams)

	response := helper.PaginatedResponse{
		Items:      playerTeamResponses,
		TotalCount: total,
	}

	helper.WriteSuccessResponse(c, http.StatusOK, response, msgPlayerTeamRelationshipRetrievedOK)
}

// UpdatePlayerTeam godoc
// @Summary Update a player-team relationship
// @Tags playerTeams
// @ID updatePlayerTeam
// @Accept json
// @Produce json
// @Param id path int true "PlayerTeam ID"
// @Param playerTeam body dto.UpdatePlayerTeamRequest true "Updated player-team data"
// @Success 200 {object} dto.PlayerTeamResponse "Player-team relationship updated successfully"
// @Failure 400 {object} helper.AppError "Invalid input"
// @Failure 404 {object} helper.AppError "Not found"
// @Failure 409 {object} helper.AppError "Conflict (e.g., date overlap)"
// @Failure 500 {object} helper.AppError "Internal server error"
// @Router /admin/player-teams/{id} [put]
// @Security BearerAuth
func (h *PlayerTeamHandler) UpdatePlayerTeam(c *gin.Context) {
	id, err := strconv.ParseUint(c.Param("id"), 10, 64)
	if err != nil {
		helper.WriteErrorResponse(c, helper.NewBadRequestError("id", constants.MsgInvalidIDSimple))
		return
	}

	var updateRequest dto.UpdatePlayerTeamRequest
	if err := c.ShouldBindJSON(&updateRequest); err != nil {
		helper.WriteErrorResponse(c, helper.BuildValidationErrorFromBinding(err, "body", "Invalid player-team data"))
		return
	}

	// Wrap context with timeout for DB/service calls
	ctx, cancel := context.WithTimeout(c.Request.Context(), 5*time.Second)
	defer cancel()

	// Get existing player team
	existingPlayerTeam, err := h.PlayerTeamDomainService.GetPlayerTeamByID(ctx, id)
	if err != nil {
		if errors.Is(err, constants.ErrRecordNotFound) {
			helper.WriteErrorResponse(c, helper.NewNotFoundError(msgPlayerTeamRelationship))
		} else {
			helper.WriteErrorResponse(c, helper.NewInternalServerError(err))
		}
		return
	}

	// Update only the provided fields
	if updateRequest.StartDate != nil {
		existingPlayerTeam.StartDate = *updateRequest.StartDate
	}
	if updateRequest.EndDate != nil {
		existingPlayerTeam.EndDate = updateRequest.EndDate
	}

	playerTeam, err := h.PlayerTeamDomainService.UpdatePlayerTeam(ctx, id, existingPlayerTeam)
	if err != nil {
		switch {
		case errors.Is(err, constants.ErrRecordNotFound):
			helper.WriteErrorResponse(c, helper.NewNotFoundError(msgPlayerTeamRelationship))
			return
		case errors.Is(err, constants.ErrOverlappingDates):
			helper.WriteErrorResponse(c, helper.NewConflictError("date_range", "Date range overlaps with an existing record"))
			return
		default:
			helper.WriteErrorResponse(c, helper.NewInternalServerError(err))
			return
		}
	}

	playerTeamResponse := h.PlayerTeamMapper.DomainToDTO(playerTeam)
	helper.WriteSuccessResponse(c, http.StatusOK, playerTeamResponse, msgPlayerTeamRelationshipUpdatedOK)
}

// DeletePlayerTeam godoc
// @Summary Delete a player-team relationship
// @Tags playerTeams
// @ID deletePlayerTeam
// @Param id path int true "PlayerTeam ID"
// @Success 204 "No Content"
// @Failure 400 {object} helper.AppError "Invalid input"
// @Failure 404 {object} helper.AppError "Not found"
// @Failure 500 {object} helper.AppError "Internal server error"
// @Router /admin/player-teams/{id} [delete]
// @Security BearerAuth
func (h *PlayerTeamHandler) DeletePlayerTeam(c *gin.Context) {
	id, err := strconv.ParseUint(c.Param("id"), 10, 64)
	if err != nil {
		helper.WriteErrorResponse(c, helper.NewBadRequestError("id", constants.MsgInvalidIDSimple))
		return
	}

	// Wrap context with timeout for DB/service calls
	ctx, cancel := context.WithTimeout(c.Request.Context(), 5*time.Second)
	defer cancel()

	err = h.PlayerTeamDomainService.DeletePlayerTeam(ctx, id)
	if err != nil {
		if errors.Is(err, constants.ErrRecordNotFound) {
			helper.WriteErrorResponse(c, helper.NewNotFoundError(msgPlayerTeamRelationship))
		} else {
			helper.WriteErrorResponse(c, helper.NewInternalServerError(err))
		}
		return
	}

	c.Status(http.StatusNoContent)
}
