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

type SeasonHandler struct {
	SeasonService service.SeasonService
}

func NewSeasonHandler(seasonService service.SeasonService) *SeasonHandler {
	return &SeasonHandler{
		SeasonService: seasonService,
	}
}

func (h *SeasonHandler) CreateSeason(c *gin.Context) {
	var season model.Seasons
	if err := c.ShouldBindJSON(&season); err != nil {
		helper.HandleError(c, helper.NewAppError(http.StatusBadRequest, constants.ErrInvalidInput, err.Error()), true)
		return
	}

	ctx := c.Request.Context()
	err := h.SeasonService.CreateSeason(ctx, &season)
	if err != nil {
		helper.HandleError(c, helper.NewAppError(http.StatusInternalServerError, "Failed to create season", err.Error()), true)
		return
	}

	helper.HandleSuccess(c, http.StatusCreated, season, "Season created successfully")
}

func (h *SeasonHandler) GetSeasonByID(c *gin.Context) {
	id, err := strconv.ParseUint(c.Param("id"), 10, 8)
	if err != nil {
		helper.HandleError(c, helper.NewAppError(http.StatusBadRequest, constants.ErrInvalidSeasonID, err.Error()), true)
		return
	}

	ctx := c.Request.Context()
	season, err := h.SeasonService.GetSeasonByID(ctx, uint8(id))
	if err != nil {
		helper.HandleError(c, helper.NewAppError(http.StatusNotFound, constants.ErrSeasonNotFound.Error(), ""), false)
		return
	}

	helper.HandleSuccess(c, http.StatusOK, season, "Season retrieved successfully")
}

func (h *SeasonHandler) ListSeasons(c *gin.Context) {
	pageStr := c.DefaultQuery("page", "1")
	page, err := strconv.ParseUint(pageStr, 10, 64)
	if err != nil {
		helper.HandleError(c, helper.NewAppError(http.StatusBadRequest, "Invalid page number", err.Error()), true)
		return
	}

	ctx := c.Request.Context()
	seasons, err := h.SeasonService.ListSeasons(ctx, page)
	if err != nil {
		helper.HandleError(c, helper.NewAppError(http.StatusInternalServerError, "Failed to list seasons", err.Error()), true)
		return
	}

	helper.HandleSuccess(c, http.StatusOK, seasons, "Seasons listed successfully")
}

func (h *SeasonHandler) UpdateSeason(c *gin.Context) {
	id, err := strconv.ParseUint(c.Param("id"), 10, 8)
	if err != nil {
		helper.HandleError(c, helper.NewAppError(http.StatusBadRequest, constants.ErrInvalidSeasonID, err.Error()), true)
		return
	}

	var season model.Seasons
	if err := c.ShouldBindJSON(&season); err != nil {
		helper.HandleError(c, helper.NewAppError(http.StatusBadRequest, constants.ErrInvalidInput, err.Error()), true)
		return
	}
	season.ID = uint8(id)

	ctx := c.Request.Context()
	err = h.SeasonService.UpdateSeason(ctx, &season)
	if err != nil {
		helper.HandleError(c, helper.NewAppError(http.StatusInternalServerError, "Failed to update season", err.Error()), true)
		return
	}

	helper.HandleSuccess(c, http.StatusOK, season, "Season updated successfully")
}

func (h *SeasonHandler) DeleteSeason(c *gin.Context) {
	id, err := strconv.ParseUint(c.Param("id"), 10, 8)
	if err != nil {
		helper.HandleError(c, helper.NewAppError(http.StatusBadRequest, constants.ErrInvalidSeasonID, err.Error()), true)
		return
	}

	ctx := c.Request.Context()
	err = h.SeasonService.DeleteSeason(ctx, uint8(id))
	if err != nil {
		helper.HandleError(c, helper.NewAppError(http.StatusInternalServerError, "Failed to delete season", err.Error()), true)
		return
	}

	helper.HandleSuccess(c, http.StatusOK, nil, "Season deleted successfully")
}
