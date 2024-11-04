package handler

import (
	"errors"
	"net/http"
	"strconv"

	"github.com/EdwinRincon/browersfc-api/api/constants"
	"github.com/EdwinRincon/browersfc-api/api/model"
	"github.com/EdwinRincon/browersfc-api/api/repository"
	"github.com/EdwinRincon/browersfc-api/api/service"
	"github.com/EdwinRincon/browersfc-api/helper"
	"github.com/gin-gonic/gin"
)

type TeamHandler struct {
	TeamService service.TeamService
}

func NewTeamHandler(teamService service.TeamService) *TeamHandler {
	return &TeamHandler{
		TeamService: teamService,
	}
}

func (h *TeamHandler) CreateTeam(c *gin.Context) {
	var team model.Teams
	if err := c.ShouldBindJSON(&team); err != nil {
		helper.HandleError(c, helper.NewAppError(http.StatusBadRequest, constants.ErrInvalidInput, err.Error()), true)
		return
	}

	ctx := c.Request.Context()
	err := h.TeamService.CreateTeam(ctx, &team)
	if err != nil {
		helper.HandleError(c, helper.NewAppError(http.StatusInternalServerError, "Failed to create team", err.Error()), true)
		return
	}

	helper.HandleSuccess(c, http.StatusCreated, team, "Team created successfully")
}

func (h *TeamHandler) GetTeamByID(c *gin.Context) {
	id, err := strconv.ParseUint(c.Param("id"), 10, 64)
	if err != nil {
		helper.HandleError(c, helper.NewAppError(http.StatusBadRequest, constants.ErrInvalidTeamID, err.Error()), true)
		return
	}

	ctx := c.Request.Context()
	team, err := h.TeamService.GetTeamByID(ctx, id)
	if err != nil {
		helper.HandleError(c, helper.NewAppError(http.StatusNotFound, constants.ErrTeamNotFound.Error(), ""), false)
		return
	}

	helper.HandleSuccess(c, http.StatusOK, team, "Team retrieved successfully")
}

func (h *TeamHandler) ListTeams(c *gin.Context) {
	pageStr := c.DefaultQuery("page", "1")
	page, err := strconv.ParseUint(pageStr, 10, 64)
	if err != nil {
		helper.HandleError(c, helper.NewAppError(http.StatusBadRequest, "Invalid page number", err.Error()), true)
		return
	}

	ctx := c.Request.Context()
	teams, err := h.TeamService.ListTeams(ctx, page)
	if err != nil {
		helper.HandleError(c, helper.NewAppError(http.StatusInternalServerError, "Failed to list teams", err.Error()), true)
		return
	}

	helper.HandleSuccess(c, http.StatusOK, teams, "Teams listed successfully")
}

func (h *TeamHandler) UpdateTeam(c *gin.Context) {
	id, err := strconv.ParseUint(c.Param("id"), 10, 64)
	if err != nil {
		helper.HandleError(c, helper.NewAppError(http.StatusBadRequest, constants.ErrInvalidTeamID, err.Error()), true)
		return
	}

	var team model.Teams
	if err := c.ShouldBindJSON(&team); err != nil {
		helper.HandleError(c, helper.NewAppError(http.StatusBadRequest, constants.ErrInvalidInput, err.Error()), true)
		return
	}
	team.ID = id

	ctx := c.Request.Context()
	err = h.TeamService.UpdateTeam(ctx, &team)
	if err != nil {
		if errors.Is(err, repository.ErrTeamNotFound) {
			helper.HandleError(c, helper.NewAppError(http.StatusNotFound, constants.ErrTeamNotFound.Error(), ""), true)
			return
		}
		helper.HandleError(c, helper.NewAppError(http.StatusInternalServerError, "Failed to update team", err.Error()), true)
		return
	}

	helper.HandleSuccess(c, http.StatusOK, team, "Team updated successfully")
}

func (h *TeamHandler) DeleteTeam(c *gin.Context) {
	id, err := strconv.ParseUint(c.Param("id"), 10, 64)
	if err != nil {
		helper.HandleError(c, helper.NewAppError(http.StatusBadRequest, constants.ErrInvalidTeamID, err.Error()), true)
		return
	}

	ctx := c.Request.Context()
	err = h.TeamService.DeleteTeam(ctx, id)
	if err != nil {
		if errors.Is(err, repository.ErrTeamNotFound) {
			helper.HandleError(c, helper.NewAppError(http.StatusNotFound, constants.ErrTeamNotFound.Error(), ""), true)
			return
		}
		helper.HandleError(c, helper.NewAppError(http.StatusInternalServerError, "Failed to delete team", err.Error()), true)
		return
	}

	helper.HandleSuccess(c, http.StatusOK, nil, "Team deleted successfully")
}
