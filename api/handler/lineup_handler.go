package handler

import (
	"context"
	"errors"
	"net/http"
	"strconv"
	"time"

	"github.com/EdwinRincon/browersfc-api/api/constants"
	"github.com/EdwinRincon/browersfc-api/api/dto"
	"github.com/EdwinRincon/browersfc-api/api/mapper"
	"github.com/EdwinRincon/browersfc-api/api/model"
	"github.com/EdwinRincon/browersfc-api/api/service"
	"github.com/EdwinRincon/browersfc-api/helper"
	"github.com/gin-gonic/gin"
)

type LineupHandler struct {
	LineupService service.LineupService
	PlayerService service.PlayerService
	MatchService  service.MatchService
}

func NewLineupHandler(lineupService service.LineupService, playerService service.PlayerService, matchService service.MatchService) *LineupHandler {
	return &LineupHandler{
		LineupService: lineupService,
		PlayerService: playerService,
		MatchService:  matchService,
	}
}

// respondServiceError processes errors from service layer operations, providing appropriate HTTP responses.
// It differentiates between "not found" errors and other internal errors, translating them
// into consistent API responses with the correct HTTP status codes.
// Parameters:
//   - c: Gin context for the HTTP response
//   - err: The error returned from the service layer
//   - resourceName: Name of the resource being accessed (e.g., "player", "lineup")
func (h *LineupHandler) respondServiceError(c *gin.Context, err error, resourceName string) {
	if errors.Is(err, constants.ErrRecordNotFound) {
		helper.WriteErrorResponse(c, helper.NewNotFoundError(resourceName))
		return
	}
	helper.WriteErrorResponse(c, helper.NewInternalServerError(err))
}

// validateReferences ensures foreign key references in the lineup update request are valid.
// It checks that any provided player or match IDs exist in the database.
// Returns an appropriate AppError if validation fails.
func (h *LineupHandler) validateReferences(ctx context.Context, dto *dto.UpdateLineupRequest) *helper.AppError {
	// Validate player ID if provided
	if dto.PlayerID != nil {
		_, err := h.PlayerService.GetPlayerByID(ctx, *dto.PlayerID)
		if err != nil {
			if errors.Is(err, constants.ErrRecordNotFound) {
				return helper.NewNotFoundError("player")
			}
			return helper.NewInternalServerError(err)
		}
	}

	// Validate match ID if provided
	if dto.MatchID != nil {
		_, err := h.MatchService.GetMatchByID(ctx, *dto.MatchID)
		if err != nil {
			if errors.Is(err, constants.ErrRecordNotFound) {
				return helper.NewNotFoundError("match")
			}
			return helper.NewInternalServerError(err)
		}
	}

	return nil
}

// CreateLineup godoc
// @Summary      Create a new lineup entry
// @Description  Creates a new lineup entry for a player in a match
// @Tags         lineups
// @ID           createLineup
// @Accept       json
// @Produce      json
// @Param        lineup  body      dto.CreateLineupRequest  true  "Lineup data"
// @Success      201     {object}  dto.LineupResponse  "Lineup created successfully"
// @Failure      400     {object}  helper.AppError "Invalid input"
// @Failure      404     {object}  helper.AppError "Player or Match not found"
// @Failure      500     {object}  helper.AppError "Internal server error"
// @Router       /admin/lineups [post]
// @Security     ApiKeyAuth
func (h *LineupHandler) CreateLineup(c *gin.Context) {
	var createRequest dto.CreateLineupRequest
	if err := c.ShouldBindJSON(&createRequest); err != nil {
		helper.WriteErrorResponse(c, helper.BuildValidationErrorFromBinding(err, "body", "Invalid lineup data"))
		return
	}

	// Wrap context with timeout for DB/service calls
	ctx, cancel := context.WithTimeout(c.Request.Context(), 5*time.Second)
	defer cancel()

	// Verify that the player exists
	_, err := h.PlayerService.GetPlayerByID(ctx, createRequest.PlayerID)
	if err != nil {
		if errors.Is(err, constants.ErrRecordNotFound) {
			helper.WriteErrorResponse(c, helper.NewNotFoundError("player"))
		} else {
			helper.WriteErrorResponse(c, helper.NewInternalServerError(err))
		}
		return
	}

	// Verify that the match exists
	_, err = h.MatchService.GetMatchByID(ctx, createRequest.MatchID)
	if err != nil {
		if errors.Is(err, constants.ErrRecordNotFound) {
			helper.WriteErrorResponse(c, helper.NewNotFoundError("match"))
		} else {
			helper.WriteErrorResponse(c, helper.NewInternalServerError(err))
		}
		return
	}

	// Create the lineup
	lineup := mapper.ToLineup(&createRequest)
	createdLineup, err := h.LineupService.CreateLineup(ctx, lineup)
	if err != nil {
		helper.WriteErrorResponse(c, helper.NewInternalServerError(err))
		return
	}

	// Get the created lineup with all relationships loaded
	fullLineup, err := h.LineupService.GetLineupByID(ctx, createdLineup.ID)
	if err != nil {
		helper.WriteErrorResponse(c, helper.NewInternalServerError(err))
		return
	}

	// Map to response DTO
	response := mapper.ToLineupResponse(fullLineup)
	helper.WriteSuccessResponse(c, http.StatusCreated, response, constants.MsgLineupCreated)
}

// GetLineupByID godoc
// @Summary      Get a lineup by ID
// @Description  Retrieves a lineup by its ID with related player and match information
// @Tags         lineups
// @ID           getLineupByID
// @Param        id   path      int  true  "Lineup ID"
// @Success      200  {object}  dto.LineupResponse "Lineup retrieved successfully"
// @Failure      400  {object}  helper.AppError "Invalid input"
// @Failure      404  {object}  helper.AppError "Lineup not found"
// @Failure      500  {object}  helper.AppError "Internal server error"
// @Router       /lineups/{id} [get]
// @Security     ApiKeyAuth
func (h *LineupHandler) GetLineupByID(c *gin.Context) {
	id, err := strconv.ParseUint(c.Param("id"), 10, 64)
	if err != nil {
		helper.WriteErrorResponse(c, helper.NewBadRequestError("id", constants.MsgInvalidLineupID))
		return
	}

	// Wrap context with timeout for DB/service calls
	ctx, cancel := context.WithTimeout(c.Request.Context(), 5*time.Second)
	defer cancel()

	lineup, err := h.LineupService.GetLineupByID(ctx, id)
	if err != nil {
		h.respondServiceError(c, err, "lineup")
		return
	}

	response := mapper.ToLineupResponse(lineup)
	helper.WriteSuccessResponse(c, http.StatusOK, response, constants.MsgLineupRetrieved)
}

// GetLineupsByMatchID godoc
// @Summary      Get lineups by match ID
// @Description  Retrieves all lineups for a specific match, organized by starting XI and substitutes
// @Tags         lineups
// @ID           getLineupsByMatchID
// @Param        id  path      int  true  "Match ID"
// @Success      200  {object}  dto.MatchLineupResponse "Lineups retrieved successfully"
// @Failure      400  {object}  helper.AppError "Invalid input"
// @Failure      404  {object}  helper.AppError "Match not found"
// @Failure      500  {object}  helper.AppError "Internal server error"
// @Router       /matches/{id}/lineups [get]
// @Security     ApiKeyAuth
func (h *LineupHandler) GetLineupsByMatchID(c *gin.Context) {
	matchID, err := strconv.ParseUint(c.Param("id"), 10, 64)
	if err != nil {
		helper.WriteErrorResponse(c, helper.NewBadRequestError("id", constants.MsgInvalidMatchID))
		return
	}

	// Wrap context with timeout for DB/service calls
	ctx, cancel := context.WithTimeout(c.Request.Context(), 5*time.Second)
	defer cancel()

	matchLineups, err := h.LineupService.GetMatchLineups(ctx, matchID)
	if err != nil {
		h.respondServiceError(c, err, "match")
		return
	}

	helper.WriteSuccessResponse(c, http.StatusOK, matchLineups, constants.MsgLineupsRetrieved)
}

// GetLineupsByPlayerID godoc
// @Summary      Get lineups by player ID
// @Description  Retrieves all lineups for a specific player
// @Tags         lineups
// @ID           getLineupsByPlayerID
// @Param        id  path      int  true  "Player ID"
// @Success      200  {array}   dto.LineupResponse "Lineups retrieved successfully"
// @Failure      400  {object}  helper.AppError "Invalid input"
// @Failure      404  {object}  helper.AppError "Player not found"
// @Failure      500  {object}  helper.AppError "Internal server error"
// @Router       /players/{id}/lineups [get]
// @Security     ApiKeyAuth
func (h *LineupHandler) GetLineupsByPlayerID(c *gin.Context) {
	playerID, err := strconv.ParseUint(c.Param("id"), 10, 64)
	if err != nil {
		helper.WriteErrorResponse(c, helper.NewBadRequestError("id", constants.MsgInvalidPlayerID))
		return
	}

	// Check if the player exists
	ctx, cancel := context.WithTimeout(c.Request.Context(), 5*time.Second)
	defer cancel()

	_, err = h.PlayerService.GetPlayerByID(ctx, playerID)
	if err != nil {
		h.respondServiceError(c, err, "player")
		return
	}

	lineups, err := h.LineupService.GetLineupsByPlayerID(ctx, playerID)
	if err != nil {
		h.respondServiceError(c, err, "lineup")
		return
	}

	response := mapper.ToLineupResponseList(lineups)
	helper.WriteSuccessResponse(c, http.StatusOK, response, constants.MsgLineupsRetrieved)
}

// GetPaginatedLineups godoc
// @Summary      Get paginated lineups
// @Description  Retrieves a paginated list of lineups with sorting and ordering
// @Tags         lineups
// @ID           getPaginatedLineups
// @Param        page      query     int     false  "Page number (0-based)"
// @Param        pageSize  query     int     false  "Number of items per page (default 10)"
// @Param        sort      query     string  false  "Sort field (e.g., position, created_at)"
// @Param        order     query     string  false  "Sort order (asc/desc)"
// @Success      200       {object}  map[string]interface{} "Lineups retrieved successfully"
// @Failure      400       {object}  helper.AppError "Invalid input"
// @Failure      500       {object}  helper.AppError "Internal server error"
// @Router       /lineups [get]
// @Security     ApiKeyAuth
func (h *LineupHandler) GetPaginatedLineups(c *gin.Context) {
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
	if err := helper.ValidateSort(model.Lineup{}, sort); err != nil {
		helper.WriteErrorResponse(c, helper.NewBadRequestError("sort", err.Error()))
		return
	}

	// Wrap context with timeout for DB/service calls
	ctx, cancel := context.WithTimeout(c.Request.Context(), 5*time.Second)
	defer cancel()

	lineups, total, err := h.LineupService.GetPaginatedLineups(ctx, sort, order, page, pageSize)
	if err != nil {
		helper.WriteErrorResponse(c, helper.NewInternalServerError(err))
		return
	}

	response := helper.PaginatedResponse{
		Items:      mapper.ToLineupResponseList(lineups),
		TotalCount: total,
	}

	helper.WriteSuccessResponse(c, http.StatusOK, response, constants.MsgLineupsRetrieved)
}

// UpdateLineup godoc
// @Summary      Update an existing lineup
// @Description  Updates an existing lineup's information by ID
// @Tags         lineups
// @ID           updateLineup
// @Accept       json
// @Produce      json
// @Param        id     path      int                   true  "Lineup ID"
// @Param        lineup body      dto.UpdateLineupRequest true  "Updated lineup data"
// @Success      200    {object}  dto.LineupResponse   "Lineup updated successfully"
// @Failure      400    {object}  helper.AppError "Invalid input"
// @Failure      404    {object}  helper.AppError "Lineup not found"
// @Failure      500    {object}  helper.AppError "Internal server error"
// @Router       /admin/lineups/{id} [put]
// @Security     ApiKeyAuth
func (h *LineupHandler) UpdateLineup(c *gin.Context) {
	id, err := strconv.ParseUint(c.Param("id"), 10, 64)
	if err != nil {
		helper.WriteErrorResponse(c, helper.NewBadRequestError("id", constants.MsgInvalidLineupID))
		return
	}

	var lineupUpdateDTO dto.UpdateLineupRequest
	if err := c.ShouldBindJSON(&lineupUpdateDTO); err != nil {
		helper.WriteErrorResponse(c, helper.BuildValidationErrorFromBinding(err, "body", "Invalid lineup data"))
		return
	}

	// Wrap context with timeout for DB/service calls
	ctx, cancel := context.WithTimeout(c.Request.Context(), 5*time.Second)
	defer cancel()

	if err := h.validateReferences(ctx, &lineupUpdateDTO); err != nil {
		helper.WriteErrorResponse(c, err)
		return
	}

	updatedLineup, err := h.LineupService.UpdateLineup(ctx, id, &lineupUpdateDTO)
	if err != nil {
		h.respondServiceError(c, err, "lineup")
		return
	}

	response := mapper.ToLineupResponse(updatedLineup)
	helper.WriteSuccessResponse(c, http.StatusOK, response, constants.MsgLineupUpdated)
}

// DeleteLineup godoc
// @Summary      Delete a lineup
// @Description  Deletes a lineup by its ID
// @Tags         lineups
// @ID           deleteLineup
// @Param        id   path      int  true  "Lineup ID"
// @Success      204  "No Content"
// @Failure      400  {object}  helper.AppError "Invalid input"
// @Failure      404  {object}  helper.AppError "Lineup not found"
// @Failure      500  {object}  helper.AppError "Internal server error"
// @Router       /admin/lineups/{id} [delete]
// @Security     ApiKeyAuth
func (h *LineupHandler) DeleteLineup(c *gin.Context) {
	id, err := strconv.ParseUint(c.Param("id"), 10, 64)
	if err != nil {
		helper.WriteErrorResponse(c, helper.NewBadRequestError("id", constants.MsgInvalidLineupID))
		return
	}

	// Wrap context with timeout for DB/service calls
	ctx, cancel := context.WithTimeout(c.Request.Context(), 5*time.Second)
	defer cancel()

	err = h.LineupService.DeleteLineup(ctx, id)
	if err != nil {
		h.respondServiceError(c, err, "lineup")
		return
	}

	c.Status(http.StatusNoContent)
}
