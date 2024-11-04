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

type MatchHandler struct {
	MatchService service.MatchService
}

func NewMatchHandler(matchService service.MatchService) *MatchHandler {
	return &MatchHandler{
		MatchService: matchService,
	}
}

func (h *MatchHandler) CreateMatch(c *gin.Context) {
	var match model.Matches
	if err := c.ShouldBindJSON(&match); err != nil {
		helper.HandleError(c, helper.NewAppError(http.StatusBadRequest, constants.ErrInvalidInput, err.Error()), true)
		return
	}

	ctx := c.Request.Context()
	err := h.MatchService.CreateMatch(ctx, &match)
	if err != nil {
		helper.HandleError(c, helper.NewAppError(http.StatusInternalServerError, "Failed to create match", err.Error()), true)
		return
	}

	helper.HandleSuccess(c, http.StatusCreated, match, "Match created successfully")
}

func (h *MatchHandler) GetMatchByID(c *gin.Context) {
	id, err := strconv.ParseUint(c.Param("id"), 10, 64)
	if err != nil {
		helper.HandleError(c, helper.NewAppError(http.StatusBadRequest, constants.ErrInvalidMatchID, err.Error()), true)
		return
	}

	ctx := c.Request.Context()
	match, err := h.MatchService.GetMatchByID(ctx, id)
	if err != nil {
		helper.HandleError(c, helper.NewAppError(http.StatusNotFound, constants.ErrMatchNotFound.Error(), ""), false)
		return
	}

	helper.HandleSuccess(c, http.StatusOK, match, "Match retrieved successfully")
}

func (h *MatchHandler) ListMatches(c *gin.Context) {
	page, err := strconv.ParseUint(c.DefaultQuery("page", "1"), 10, 64)
	if err != nil {
		helper.HandleError(c, helper.NewAppError(http.StatusBadRequest, "Invalid page number", err.Error()), true)
		return
	}

	pageSize, err := strconv.ParseUint(c.DefaultQuery("pageSize", "10"), 10, 64)
	if err != nil {
		helper.HandleError(c, helper.NewAppError(http.StatusBadRequest, "Invalid page size", err.Error()), true)
		return
	}

	ctx := c.Request.Context()
	matches, err := h.MatchService.ListMatches(ctx, page, pageSize)
	if err != nil {
		helper.HandleError(c, helper.NewAppError(http.StatusInternalServerError, "Failed to list matches", err.Error()), true)
		return
	}

	helper.HandleSuccess(c, http.StatusOK, matches, "Matches listed successfully")
}

func (h *MatchHandler) UpdateMatch(c *gin.Context) {
	id, err := strconv.ParseUint(c.Param("id"), 10, 64)
	if err != nil {
		helper.HandleError(c, helper.NewAppError(http.StatusBadRequest, constants.ErrInvalidMatchID, err.Error()), true)
		return
	}

	var match model.Matches
	if err := c.ShouldBindJSON(&match); err != nil {
		helper.HandleError(c, helper.NewAppError(http.StatusBadRequest, constants.ErrInvalidInput, err.Error()), true)
		return
	}
	match.ID = id

	ctx := c.Request.Context()
	err = h.MatchService.UpdateMatch(ctx, &match)
	if err != nil {
		helper.HandleError(c, helper.NewAppError(http.StatusInternalServerError, "Failed to update match", err.Error()), true)
		return
	}

	helper.HandleSuccess(c, http.StatusOK, match, "Match updated successfully")
}

func (h *MatchHandler) DeleteMatch(c *gin.Context) {
	id, err := strconv.ParseUint(c.Param("id"), 10, 64)
	if err != nil {
		helper.HandleError(c, helper.NewAppError(http.StatusBadRequest, constants.ErrInvalidMatchID, err.Error()), true)
		return
	}

	ctx := c.Request.Context()
	err = h.MatchService.DeleteMatch(ctx, id)
	if err != nil {
		helper.HandleError(c, helper.NewAppError(http.StatusInternalServerError, "Failed to delete match", err.Error()), true)
		return
	}

	helper.HandleSuccess(c, http.StatusOK, nil, "Match deleted successfully")
}