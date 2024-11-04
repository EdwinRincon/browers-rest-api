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

type TeamStatsHandler struct {
	TeamStatsService service.TeamStatsService
}

func NewTeamStatsHandler(teamStatsService service.TeamStatsService) *TeamStatsHandler {
	return &TeamStatsHandler{
		TeamStatsService: teamStatsService,
	}
}

func (h *TeamStatsHandler) CreateTeamStats(c *gin.Context) {
	var teamStats model.TeamsStats
	if err := c.ShouldBindJSON(&teamStats); err != nil {
		helper.HandleError(c, helper.NewAppError(http.StatusBadRequest, constants.ErrInvalidInput, err.Error()), true)
		return
	}

	ctx := c.Request.Context()
	err := h.TeamStatsService.CreateTeamStats(ctx, &teamStats)
	if err != nil {
		helper.HandleError(c, helper.NewAppError(http.StatusInternalServerError, "Failed to create team stats", err.Error()), true)
		return
	}

	helper.HandleSuccess(c, http.StatusCreated, teamStats, "Team stats created successfully")
}

func (h *TeamStatsHandler) GetTeamStatsByID(c *gin.Context) {
	id, err := strconv.ParseUint(c.Param("id"), 10, 64)
	if err != nil {
		helper.HandleError(c, helper.NewAppError(http.StatusBadRequest, "Invalid ID", err.Error()), true)
		return
	}

	ctx := c.Request.Context()
	teamStats, err := h.TeamStatsService.GetTeamStatsByID(ctx, id)
	if err != nil {
		helper.HandleError(c, helper.NewAppError(http.StatusInternalServerError, "Failed to retrieve team stats", err.Error()), false)
		return
	}

	helper.HandleSuccess(c, http.StatusOK, teamStats, "Team stats retrieved successfully")
}

func (h *TeamStatsHandler) ListTeamStats(c *gin.Context) {
	pageStr := c.DefaultQuery("page", "1")
	page, err := strconv.ParseUint(pageStr, 10, 64)
	if err != nil {
		helper.HandleError(c, helper.NewAppError(http.StatusBadRequest, "Invalid page number", err.Error()), true)
		return
	}

	ctx := c.Request.Context()
	teamStats, err := h.TeamStatsService.ListTeamStats(ctx, page)
	if err != nil {
		helper.HandleError(c, helper.NewAppError(http.StatusInternalServerError, "Failed to list team stats", err.Error()), true)
		return
	}

	helper.HandleSuccess(c, http.StatusOK, teamStats, "Team stats listed successfully")
}

func (h *TeamStatsHandler) UpdateTeamStats(c *gin.Context) {
	var teamStats model.TeamsStats
	if err := c.ShouldBindJSON(&teamStats); err != nil {
		helper.HandleError(c, helper.NewAppError(http.StatusBadRequest, constants.ErrInvalidInput, err.Error()), true)
		return
	}

	ctx := c.Request.Context()
	err := h.TeamStatsService.UpdateTeamStats(ctx, &teamStats)
	if err != nil {
		helper.HandleError(c, helper.NewAppError(http.StatusInternalServerError, "Failed to update team stats", err.Error()), true)
		return
	}

	helper.HandleSuccess(c, http.StatusOK, teamStats, "Team stats updated successfully")
}

func (h *TeamStatsHandler) DeleteTeamStats(c *gin.Context) {
	id, err := strconv.ParseUint(c.Param("id"), 10, 64)
	if err != nil {
		helper.HandleError(c, helper.NewAppError(http.StatusBadRequest, "Invalid ID", err.Error()), true)
		return
	}

	ctx := c.Request.Context()
	err = h.TeamStatsService.DeleteTeamStats(ctx, id)
	if err != nil {
		helper.HandleError(c, helper.NewAppError(http.StatusInternalServerError, "Failed to delete team stats", err.Error()), true)
		return
	}

	helper.HandleSuccess(c, http.StatusOK, nil, "Team stats deleted successfully")
}
