package handler

import (
	"net/http"
	"strconv"

	"github.com/EdwinRincon/browersfc-api/adapter/mapper"
	"github.com/EdwinRincon/browersfc-api/api/constants"
	"github.com/EdwinRincon/browersfc-api/api/dto"
	"github.com/EdwinRincon/browersfc-api/helper"
	domainservice "github.com/EdwinRincon/browersfc-api/internal/domain/service"
	"github.com/gin-gonic/gin"
)

type LineupHandler struct {
	LineupDomainService *domainservice.LineupDomainService
	LineupMapper        *mapper.LineupMapper
}

func NewLineupHandler(lineupDomainService *domainservice.LineupDomainService) *LineupHandler {
	return &LineupHandler{
		LineupDomainService: lineupDomainService,
		LineupMapper:        mapper.NewLineupMapper(),
	}
}

// @Summary Create a new lineup
// @Description Add a new lineup for a match
// @Tags lineups
// @Accept json
// @Produce json
// @Param lineup body dto.CreateLineupRequest true "Lineup data"
// @Success 201 {object} dto.LineupResponse
// @Failure 400 {object} helper.AppError "Invalid input"
// @Failure 500 {object} helper.AppError "Internal server error"
// @Security BearerAuth
// @Router /lineups [post]
func (h *LineupHandler) CreateLineup(c *gin.Context) {
	var createRequest dto.CreateLineupRequest
	if err := c.ShouldBindJSON(&createRequest); err != nil {
		helper.WriteErrorResponse(c, helper.NewBadRequestError("body", "Invalid lineup data"))
		return
	}

	lineupEntity := h.LineupMapper.DTOToDomain(&createRequest)

	if err := h.LineupDomainService.CreateLineup(c.Request.Context(), lineupEntity); err != nil {
		if err == constants.ErrInvalidData {
			helper.WriteErrorResponse(c, helper.NewBadRequestError("data", "Invalid lineup data"))
			return
		}
		if err == constants.ErrPlayerNotFound {
			helper.WriteErrorResponse(c, helper.NewNotFoundError("player"))
			return
		}
		if err == constants.ErrMatchNotFound {
			helper.WriteErrorResponse(c, helper.NewNotFoundError("match"))
			return
		}
		helper.WriteErrorResponse(c, helper.NewInternalServerError(err))
		return
	}

	response := h.LineupMapper.DomainToDTO(lineupEntity)
	helper.WriteSuccessResponse(c, http.StatusCreated, response, "Lineup created successfully")
}

// @Summary Get lineup by ID
// @Description Get a lineup by its ID
// @Tags lineups
// @Accept json
// @Produce json
// @Param id path int true "Lineup ID"
// @Success 200 {object} dto.LineupResponse
// @Failure 400 {object} helper.AppError "Invalid lineup ID"
// @Failure 404 {object} helper.AppError "Lineup not found"
// @Failure 500 {object} helper.AppError "Internal server error"
// @Security BearerAuth
// @Router /lineups/{id} [get]
func (h *LineupHandler) GetLineupByID(c *gin.Context) {
	idStr := c.Param("id")
	id, err := strconv.ParseUint(idStr, 10, 64)
	if err != nil {
		helper.WriteErrorResponse(c, helper.NewBadRequestError("id", "Invalid lineup ID"))
		return
	}

	lineup, err := h.LineupDomainService.GetLineupByID(c.Request.Context(), id)
	if err != nil {
		if err == constants.ErrLineupNotFound {
			helper.WriteErrorResponse(c, helper.NewNotFoundError("lineup"))
			return
		}
		helper.WriteErrorResponse(c, helper.NewInternalServerError(err))
		return
	}

	response := h.LineupMapper.DomainToDTO(lineup)
	helper.WriteSuccessResponse(c, http.StatusOK, response, "Lineup retrieved successfully")
}

// @Summary Update lineup
// @Description Update an existing lineup
// @Tags lineups
// @Accept json
// @Produce json
// @Param id path int true "Lineup ID"
// @Param lineup body dto.UpdateLineupRequest true "Lineup data"
// @Success 200 {object} dto.LineupResponse
// @Failure 400 {object} helper.AppError "Invalid input"
// @Failure 404 {object} helper.AppError "Lineup not found"
// @Failure 500 {object} helper.AppError "Internal server error"
// @Security BearerAuth
// @Router /lineups/{id} [put]
func (h *LineupHandler) UpdateLineup(c *gin.Context) {
	idStr := c.Param("id")
	id, err := strconv.ParseUint(idStr, 10, 64)
	if err != nil {
		helper.WriteErrorResponse(c, helper.NewBadRequestError("id", "Invalid lineup ID"))
		return
	}

	var updateRequest dto.UpdateLineupRequest
	if err := c.ShouldBindJSON(&updateRequest); err != nil {
		helper.WriteErrorResponse(c, helper.NewBadRequestError("body", "Invalid request format"))
		return
	}

	lineupEntity := h.LineupMapper.UpdateDTOToDomain(&updateRequest)

	if err := h.LineupDomainService.UpdateLineup(c.Request.Context(), id, lineupEntity); err != nil {
		if err == constants.ErrLineupNotFound {
			helper.WriteErrorResponse(c, helper.NewNotFoundError("lineup"))
			return
		}
		if err == constants.ErrInvalidData {
			helper.WriteErrorResponse(c, helper.NewBadRequestError("data", "Invalid lineup data"))
			return
		}
		helper.WriteErrorResponse(c, helper.NewInternalServerError(err))
		return
	}

	// Get updated lineup
	updatedLineup, err := h.LineupDomainService.GetLineupByID(c.Request.Context(), id)
	if err != nil {
		helper.WriteErrorResponse(c, helper.NewInternalServerError(err))
		return
	}

	response := h.LineupMapper.DomainToDTO(updatedLineup)
	helper.WriteSuccessResponse(c, http.StatusOK, response, "Lineup updated successfully")
}

// @Summary Delete lineup
// @Description Delete a lineup by ID
// @Tags lineups
// @Accept json
// @Produce json
// @Param id path int true "Lineup ID"
// @Success 200 {object} helper.MessageResponse
// @Failure 400 {object} helper.AppError "Invalid lineup ID"
// @Failure 404 {object} helper.AppError "Lineup not found"
// @Failure 500 {object} helper.AppError "Internal server error"
// @Security BearerAuth
// @Router /lineups/{id} [delete]
func (h *LineupHandler) DeleteLineup(c *gin.Context) {
	idStr := c.Param("id")
	id, err := strconv.ParseUint(idStr, 10, 64)
	if err != nil {
		helper.WriteErrorResponse(c, helper.NewBadRequestError("id", "Invalid lineup ID"))
		return
	}

	if err := h.LineupDomainService.DeleteLineup(c.Request.Context(), id); err != nil {
		if err == constants.ErrLineupNotFound {
			helper.WriteErrorResponse(c, helper.NewNotFoundError("lineup"))
			return
		}
		helper.WriteErrorResponse(c, helper.NewInternalServerError(err))
		return
	}

	helper.WriteSuccessResponse(c, http.StatusOK, nil, "Lineup deleted successfully")
}

// @Summary Get paginated lineups
// @Description Retrieve lineups with pagination, optional sorting, and ordering.
// @Tags lineups
// @Accept json
// @Produce json
// @Param page query int false "Page number" default(0)
// @Param pageSize query int false "Page size" default(10)
// @Param sort query string false "Sort field"
// @Param order query string false "Sort order" Enums(asc, desc) default(asc)
// @Success 200 {object} helper.AppSuccess{data=helper.PaginatedResponse{items=[]dto.LineupResponse, totalCount=int}}
// @Failure 400 {object} helper.AppError "Invalid input"
// @Failure 500 {object} helper.AppError "Internal server error"
// @Security BearerAuth
// @Router /lineups [get]
func (h *LineupHandler) GetPaginatedLineups(c *gin.Context) {
	page, _ := strconv.Atoi(c.DefaultQuery("page", "0"))
	pageSize, _ := strconv.Atoi(c.DefaultQuery("pageSize", "10"))
	sort := c.DefaultQuery("sort", "")
	order := c.DefaultQuery("order", "asc")

	lineups, total, err := h.LineupDomainService.GetPaginatedLineups(c.Request.Context(), sort, order, page, pageSize)
	if err != nil {
		if err == constants.ErrInvalidPaginationParams {
			helper.WriteErrorResponse(c, helper.NewBadRequestError("pagination", "Invalid pagination parameters"))
			return
		}
		helper.WriteErrorResponse(c, helper.NewInternalServerError(err))
		return
	}

	response := helper.PaginatedResponse{
		Items:      h.LineupMapper.DomainListToDTO(lineups),
		TotalCount: total,
	}

	helper.WriteSuccessResponse(c, http.StatusOK, response, "Lineups retrieved successfully")
}

// @Summary Get lineups by match ID
// @Description Get all lineups for a specific match
// @Tags lineups
// @Accept json
// @Produce json
// @Param matchId path int true "Match ID"
// @Success 200 {array} dto.LineupResponse
// @Failure 400 {object} helper.AppError "Invalid match ID"
// @Failure 404 {object} helper.AppError "Lineup not found"
// @Failure 500 {object} helper.AppError "Internal server error"
// @Security BearerAuth
// @Router /lineups/match/{matchId} [get]
func (h *LineupHandler) GetLineupsByMatchID(c *gin.Context) {
	matchIDStr := c.Param("matchId")
	matchID, err := strconv.ParseUint(matchIDStr, 10, 64)
	if err != nil {
		helper.WriteErrorResponse(c, helper.NewBadRequestError("id", "Invalid match ID"))
		return
	}

	lineups, err := h.LineupDomainService.GetLineupsByMatchID(c.Request.Context(), matchID)
	if err != nil {
		if err == constants.ErrMatchNotFound {
			helper.WriteErrorResponse(c, helper.NewNotFoundError("match"))
			return
		}
		helper.WriteErrorResponse(c, helper.NewInternalServerError(err))
		return
	}

	response := h.LineupMapper.DomainListToDTO(lineups)
	helper.WriteSuccessResponse(c, http.StatusOK, response, "Lineups retrieved successfully")
}

// @Summary Get starting lineups by match ID
// @Description Get all starting lineups for a specific match
// @Tags lineups
// @Accept json
// @Produce json
// @Param matchId path int true "Match ID"
// @Success 200 {array} dto.LineupResponse
// @Failure 400 {object} helper.AppError "Invalid match ID"
// @Failure 404 {object} helper.AppError "Lineup not found"
// @Failure 500 {object} helper.AppError "Internal server error"
// @Security BearerAuth
// @Router /lineups/match/{matchId}/starting [get]
func (h *LineupHandler) GetStartingLineupsByMatchID(c *gin.Context) {
	matchIDStr := c.Param("matchId")
	matchID, err := strconv.ParseUint(matchIDStr, 10, 64)
	if err != nil {
		helper.WriteErrorResponse(c, helper.NewBadRequestError("id", "Invalid match ID"))
		return
	}

	lineups, err := h.LineupDomainService.GetStartingLineupsByMatchID(c.Request.Context(), matchID)
	if err != nil {
		if err == constants.ErrMatchNotFound {
			helper.WriteErrorResponse(c, helper.NewNotFoundError("match"))
			return
		}
		helper.WriteErrorResponse(c, helper.NewInternalServerError(err))
		return
	}

	response := h.LineupMapper.DomainListToDTO(lineups)
	helper.WriteSuccessResponse(c, http.StatusOK, response, "Starting lineups retrieved successfully")
}

// @Summary Get substitute lineups by match ID
// @Description Get all substitute lineups for a specific match
// @Tags lineups
// @Accept json
// @Produce json
// @Param matchId path int true "Match ID"
// @Success 200 {array} dto.LineupResponse
// @Failure 400 {object} helper.AppError "Invalid match ID"
// @Failure 404 {object} helper.AppError "Lineup not found"
// @Failure 500 {object} helper.AppError "Internal server error"
// @Security BearerAuth
// @Router /lineups/match/{matchId}/substitutes [get]
func (h *LineupHandler) GetSubstitutesLineupsByMatchID(c *gin.Context) {
	matchIDStr := c.Param("matchId")
	matchID, err := strconv.ParseUint(matchIDStr, 10, 64)
	if err != nil {
		helper.WriteErrorResponse(c, helper.NewBadRequestError("id", "Invalid match ID"))
		return
	}

	lineups, err := h.LineupDomainService.GetSubstitutesLineupsByMatchID(c.Request.Context(), matchID)
	if err != nil {
		if err == constants.ErrMatchNotFound {
			helper.WriteErrorResponse(c, helper.NewNotFoundError("match"))
			return
		}
		helper.WriteErrorResponse(c, helper.NewInternalServerError(err))
		return
	}

	response := h.LineupMapper.DomainListToDTO(lineups)
	helper.WriteSuccessResponse(c, http.StatusOK, response, "Substitute lineups retrieved successfully")
}

// @Summary Get lineups by player ID
// @Description Get all lineups for a specific player
// @Tags lineups
// @Accept json
// @Produce json
// @Param playerId path int true "Player ID"
// @Success 200 {array} dto.LineupResponse
// @Failure 400 {object} helper.AppError "Invalid player ID"
// @Failure 404 {object} helper.AppError "Lineup not found"
// @Failure 500 {object} helper.AppError "Internal server error"
// @Security BearerAuth
// @Router /lineups/player/{playerId} [get]
func (h *LineupHandler) GetLineupsByPlayerID(c *gin.Context) {
	playerIDStr := c.Param("playerId")
	playerID, err := strconv.ParseUint(playerIDStr, 10, 64)
	if err != nil {
		helper.WriteErrorResponse(c, helper.NewBadRequestError("id", "Invalid player ID"))
		return
	}

	lineups, err := h.LineupDomainService.GetLineupsByPlayerID(c.Request.Context(), playerID)
	if err != nil {
		if err == constants.ErrPlayerNotFound {
			helper.WriteErrorResponse(c, helper.NewNotFoundError("player"))
			return
		}
		helper.WriteErrorResponse(c, helper.NewInternalServerError(err))
		return
	}

	response := h.LineupMapper.DomainListToDTO(lineups)
	helper.WriteSuccessResponse(c, http.StatusOK, response, "Lineups retrieved successfully")
}
