package handler

import (
	"net/http"
	"strconv"

	"github.com/EdwinRincon/browersfc-api/api/constants"
	"github.com/EdwinRincon/browersfc-api/api/model"
	"github.com/EdwinRincon/browersfc-api/api/service"
	"github.com/EdwinRincon/browersfc-api/helper"
	"github.com/gin-gonic/gin"
)

type LineupHandler struct {
	LineupService service.LineupService
}

func NewLineupHandler(lineupService service.LineupService) *LineupHandler {
	return &LineupHandler{LineupService: lineupService}
}

func (h *LineupHandler) CreateLineup(c *gin.Context) {
	var lineup model.Lineup
	if err := c.ShouldBindJSON(&lineup); err != nil {
		helper.HandleValidationError(c, err)
		return
	}

	ctx := c.Request.Context()
	err := h.LineupService.CreateLineup(ctx, &lineup)
	if err != nil {
		helper.HandleGormError(c, err)
		return
	}

	helper.HandleSuccess(c, http.StatusCreated, lineup, "Lineup created successfully")
}

func (h *LineupHandler) GetLineupByID(c *gin.Context) {
	lineupID := c.Param("id")
	id, err := strconv.ParseUint(lineupID, 10, 8)
	if err != nil {
		helper.HandleError(c, helper.NewAppError(http.StatusBadRequest, constants.ErrInvalidInput, err.Error()), true)
		return
	}

	ctx := c.Request.Context()
	lineup, err := h.LineupService.GetLineupByID(ctx, id)
	if err != nil {
		helper.HandleGormError(c, err)
		return
	}

	helper.HandleSuccess(c, http.StatusOK, lineup, "Lineup found successfully")
}

func (h *LineupHandler) ListLineups(c *gin.Context) {
	pageStr := c.DefaultQuery("page", "1")
	page, err := strconv.ParseUint(pageStr, 10, 64)
	if err != nil {
		helper.HandleError(c, helper.NewAppError(http.StatusBadRequest, "Invalid page number", err.Error()), true)
		return
	}

	ctx := c.Request.Context()
	lineups, err := h.LineupService.ListLineups(ctx, page)
	if err != nil {
		helper.HandleError(c, helper.NewAppError(http.StatusInternalServerError, "Failed to list lineups", err.Error()), true)
		return
	}

	helper.HandleSuccess(c, http.StatusOK, lineups, "Lineup listed successfully")
}

func (h *LineupHandler) UpdateLineup(c *gin.Context) {
	lineupID := c.Param("id")
	id, err := strconv.ParseUint(lineupID, 10, 8)
	if err != nil {
		helper.HandleError(c, helper.NewAppError(http.StatusBadRequest, constants.ErrInvalidInput, err.Error()), true)
		return
	}

	var lineup model.Lineup
	if err := c.ShouldBindJSON(&lineup); err != nil {
		helper.HandleValidationError(c, err)
		return
	}

	ctx := c.Request.Context()
	lineup.ID = uint64(id)
	err = h.LineupService.UpdateLineup(ctx, &lineup)
	if err != nil {
		helper.HandleGormError(c, err)
		return
	}

	helper.HandleSuccess(c, http.StatusOK, lineup, "Lineup updated successfully")
}

func (h *LineupHandler) DeleteLineup(c *gin.Context) {
	lineupID := c.Param("id")
	id, err := strconv.ParseUint(lineupID, 10, 8)
	if err != nil {
		helper.HandleError(c, helper.NewAppError(http.StatusBadRequest, constants.ErrInvalidInput, err.Error()), true)
		return
	}

	ctx := c.Request.Context()
	err = h.LineupService.DeleteLineup(ctx, uint64(id))
	if err != nil {
		helper.HandleGormError(c, err)
		return
	}

	helper.HandleSuccess(c, http.StatusOK, nil, "Lineup deleted successfully")
}

func (h *LineupHandler) GetLineupsByMatch(c *gin.Context) {
	matchID, err := strconv.ParseUint(c.Param("matchID"), 10, 64)
	if err != nil {
		helper.HandleError(c, helper.NewAppError(http.StatusBadRequest, "Invalid match ID", err.Error()), true)
		return
	}

	ctx := c.Request.Context()
	lineups, err := h.LineupService.GetLineupsByMatch(ctx, matchID)
	if err != nil {
		helper.HandleError(c, helper.NewAppError(http.StatusInternalServerError, "Failed to retrieve lineups by match", err.Error()), false)
		return
	}

	helper.HandleSuccess(c, http.StatusOK, lineups, "Lineup by match retrieved successfully")
}
