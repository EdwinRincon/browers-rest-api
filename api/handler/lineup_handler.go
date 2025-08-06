package handler

import (
	"errors"
	"net/http"
	"strconv"

	"github.com/EdwinRincon/browersfc-api/api/model"
	"github.com/EdwinRincon/browersfc-api/api/service"
	"github.com/EdwinRincon/browersfc-api/helper"
	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
)

type LineupHandler struct {
	LineupService service.LineupService
}

func NewLineupHandler(lineupService service.LineupService) *LineupHandler {
	return &LineupHandler{LineupService: lineupService}
}

// CreateLineup godoc
// @Summary      Create a new lineup
// @Description  Creates a new lineup with the provided data
// @Tags         lineups
// @ID           createLineup
// @Accept       json
// @Produce      json
// @Param        lineup  body      model.Lineup  true  "Lineup data"
// @Success      201     {object}  model.Lineup "Created"
// @Failure      400     {object}  helper.AppError "Invalid input"
// @Failure      409     {object}  helper.AppError "Conflict (e.g., lineup already exists)"
// @Failure      500     {object}  helper.AppError "Internal server error"
// @Router       /lineups [post]
// @Security     ApiKeyAuth
func (h *LineupHandler) CreateLineup(c *gin.Context) {
	var lineup model.Lineup
	if err := c.ShouldBindJSON(&lineup); err != nil {
		helper.RespondWithError(c, helper.BadRequest("body", "Invalid lineup data"))
		return
	}

	ctx := c.Request.Context()
	err := h.LineupService.CreateLineup(ctx, &lineup)
	if err != nil {
		if errors.Is(err, gorm.ErrDuplicatedKey) {
			helper.RespondWithError(c, helper.Conflict("lineup", "A lineup with these details already exists"))
		} else {
			helper.RespondWithError(c, helper.InternalError(err))
		}
		return
	}

	helper.HandleSuccess(c, http.StatusCreated, lineup, "Lineup created successfully")
}

// GetLineupByID godoc
// @Summary      Get a lineup by ID
// @Description  Returns the details of a lineup by its ID
// @Tags         lineups
// @ID           getLineupByID
// @Param        id   path      int  true  "Lineup ID"
// @Success      200  {object}  model.Lineup "Success"
// @Failure      400  {object}  helper.AppError "Invalid input"
// @Failure      404  {object}  helper.AppError "Lineup not found"
// @Failure      500  {object}  helper.AppError "Internal server error"
// @Router       /lineups/{id} [get]
// @Security     ApiKeyAuth
func (h *LineupHandler) GetLineupByID(c *gin.Context) {
	lineupID := c.Param("id")
	id, err := strconv.ParseUint(lineupID, 10, 64)
	if err != nil {
		helper.RespondWithError(c, helper.BadRequest("id", "Invalid lineup ID"))
		return
	}

	ctx := c.Request.Context()
	lineup, err := h.LineupService.GetLineupByID(ctx, id)
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			helper.RespondWithError(c, helper.NotFound("lineup"))
		} else {
			helper.RespondWithError(c, helper.InternalError(err))
		}
		return
	}

	helper.HandleSuccess(c, http.StatusOK, lineup, "Lineup found successfully")
}

// ListLineups godoc
// @Summary      List lineups
// @Description  Retrieves a paginated list of lineups
// @Tags         lineups
// @ID           listLineups
// @Param        page  query     int  false  "Page number"
// @Success      200   {array}   model.Lineup "Lineups listed successfully"
// @Failure      400   {object}  helper.AppError "Invalid input"
// @Failure      500   {object}  helper.AppError "Internal server error"
// @Router       /lineups [get]
// @Security     ApiKeyAuth
func (h *LineupHandler) ListLineups(c *gin.Context) {
	pageStr := c.DefaultQuery("page", "1")
	page, err := strconv.ParseUint(pageStr, 10, 64)
	if err != nil {
		helper.RespondWithError(c, helper.BadRequest("page", "Invalid page number"))
		return
	}

	ctx := c.Request.Context()
	lineups, err := h.LineupService.ListLineups(ctx, page)
	if err != nil {
		helper.RespondWithError(c, helper.InternalError(err))
		return
	}

	helper.HandleSuccess(c, http.StatusOK, lineups, "Lineup listed successfully")
}

// UpdateLineup godoc
// @Summary      Update an existing lineup
// @Description  Updates the details of an existing lineup by ID
// @Tags         lineups
// @ID           updateLineup
// @Accept       json
// @Produce      json
// @Param        id      path      int         true  "Lineup ID"
// @Param        lineup  body      model.Lineup true  "Updated lineup data"
// @Success      200     {object}  model.Lineup "Updated"
// @Failure      400     {object}  helper.AppError "Invalid input"
// @Failure      404     {object}  helper.AppError "Lineup not found"
// @Failure      500     {object}  helper.AppError "Internal server error"
// @Router       /lineups/{id} [put]
// @Security     ApiKeyAuth
func (h *LineupHandler) UpdateLineup(c *gin.Context) {
	lineupID := c.Param("id")
	id, err := strconv.ParseUint(lineupID, 10, 64)
	if err != nil {
		helper.RespondWithError(c, helper.BadRequest("id", "Invalid lineup ID"))
		return
	}

	var lineup model.Lineup
	if err := c.ShouldBindJSON(&lineup); err != nil {
		helper.RespondWithError(c, helper.BadRequest("body", "Invalid lineup data"))
		return
	}

	ctx := c.Request.Context()
	lineup.ID = uint64(id)
	err = h.LineupService.UpdateLineup(ctx, &lineup)
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			helper.RespondWithError(c, helper.NotFound("lineup"))
		} else if errors.Is(err, gorm.ErrDuplicatedKey) {
			helper.RespondWithError(c, helper.Conflict("lineup", "A lineup with these details already exists"))
		} else {
			helper.RespondWithError(c, helper.InternalError(err))
		}
		return
	}

	helper.HandleSuccess(c, http.StatusOK, lineup, "Lineup updated successfully")
}

// DeleteLineup godoc
// @Summary      Delete a lineup
// @Description  Deletes a lineup by its ID
// @Tags         lineups
// @ID           deleteLineup
// @Param        id   path      int  true  "Lineup ID"
// @Success      204  {object}  nil "No Content"
// @Failure      400  {object}  helper.AppError "Invalid input"
// @Failure      404  {object}  helper.AppError "Lineup not found"
// @Failure      500  {object}  helper.AppError "Internal server error"
// @Router       /lineups/{id} [delete]
// @Security     ApiKeyAuth
func (h *LineupHandler) DeleteLineup(c *gin.Context) {
	lineupID := c.Param("id")
	id, err := strconv.ParseUint(lineupID, 10, 64)
	if err != nil {
		helper.RespondWithError(c, helper.BadRequest("id", "Invalid lineup ID"))
		return
	}

	ctx := c.Request.Context()
	err = h.LineupService.DeleteLineup(ctx, uint64(id))
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			helper.RespondWithError(c, helper.NotFound("lineup"))
		} else {
			helper.RespondWithError(c, helper.InternalError(err))
		}
		return
	}

	c.Status(http.StatusNoContent)
}

// GetLineupsByMatch godoc
// @Summary      Get lineups by match ID
// @Description  Retrieves all lineups for a given match ID
// @Tags         lineups
// @ID           getLineupsByMatch
// @Param        matchID  path      int  true  "Match ID"
// @Success      200      {array}   model.Lineup "Lineups retrieved successfully"
// @Failure      400      {object}  helper.AppError "Invalid input"
// @Failure      500      {object}  helper.AppError "Internal server error"
// @Router       /lineups/match/{matchID} [get]
// @Security     ApiKeyAuth
func (h *LineupHandler) GetLineupsByMatch(c *gin.Context) {
	matchID, err := strconv.ParseUint(c.Param("matchID"), 10, 64)
	if err != nil {
		helper.RespondWithError(c, helper.BadRequest("matchID", "Invalid match ID"))
		return
	}

	ctx := c.Request.Context()
	lineups, err := h.LineupService.GetLineupsByMatch(ctx, matchID)
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			helper.RespondWithError(c, helper.NotFound("lineup"))
		} else {
			helper.RespondWithError(c, helper.InternalError(err))
		}
		return
	}

	helper.HandleSuccess(c, http.StatusOK, lineups, "Lineup by match retrieved successfully")
}
