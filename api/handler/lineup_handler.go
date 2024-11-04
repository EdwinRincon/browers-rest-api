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
	return &LineupHandler{
		LineupService: lineupService,
	}
}

func (h *LineupHandler) CreateLineup(c *gin.Context) {
	var lineup model.Lineups
	if err := c.ShouldBindJSON(&lineup); err != nil {
		helper.HandleError(c, helper.NewAppError(http.StatusBadRequest, constants.ErrInvalidInput, err.Error()), true)
		return
	}

	ctx := c.Request.Context()
	err := h.LineupService.CreateLineup(ctx, &lineup)
	if err != nil {
		helper.HandleError(c, helper.NewAppError(http.StatusInternalServerError, "Failed to create lineup", err.Error()), true)
		return
	}

	helper.HandleSuccess(c, http.StatusCreated, lineup, "Lineup created successfully")
}

func (h *LineupHandler) GetLineupByID(c *gin.Context) {
	id, err := strconv.ParseUint(c.Param("id"), 10, 64)
	if err != nil {
		helper.HandleError(c, helper.NewAppError(http.StatusBadRequest, "Invalid ID", err.Error()), true)
		return
	}

	ctx := c.Request.Context()
	lineup, err := h.LineupService.GetLineupByID(ctx, id)
	if err != nil {
		helper.HandleError(c, helper.NewAppError(http.StatusInternalServerError, "Failed to retrieve lineup", err.Error()), false)
		return
	}

	helper.HandleSuccess(c, http.StatusOK, lineup, "Lineup retrieved successfully")
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

	helper.HandleSuccess(c, http.StatusOK, lineups, "Lineups listed successfully")
}

func (h *LineupHandler) UpdateLineup(c *gin.Context) {
	var lineup model.Lineups
	if err := c.ShouldBindJSON(&lineup); err != nil {
		helper.HandleError(c, helper.NewAppError(http.StatusBadRequest, constants.ErrInvalidInput, err.Error()), true)
		return
	}

	ctx := c.Request.Context()
	err := h.LineupService.UpdateLineup(ctx, &lineup)
	if err != nil {
		helper.HandleError(c, helper.NewAppError(http.StatusInternalServerError, "Failed to update lineup", err.Error()), true)
		return
	}

	helper.HandleSuccess(c, http.StatusOK, lineup, "Lineup updated successfully")
}

func (h *LineupHandler) DeleteLineup(c *gin.Context) {
	id, err := strconv.ParseUint(c.Param("id"), 10, 64)
	if err != nil {
		helper.HandleError(c, helper.NewAppError(http.StatusBadRequest, "Invalid ID", err.Error()), true)
		return
	}

	ctx := c.Request.Context()
	err = h.LineupService.DeleteLineup(ctx, id)
	if err != nil {
		helper.HandleError(c, helper.NewAppError(http.StatusInternalServerError, "Failed to delete lineup", err.Error()), true)
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

	helper.HandleSuccess(c, http.StatusOK, lineups, "Lineups by match retrieved successfully")
}
