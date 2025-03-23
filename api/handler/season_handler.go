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
	var season model.Season
	if err := c.ShouldBindJSON(&season); err != nil {
		helper.HandleValidationError(c, err)
		return
	}

	ctx := c.Request.Context()
	err := h.SeasonService.CreateSeason(ctx, &season)
	if err != nil {
		helper.HandleGormError(c, err)
		return
	}

	helper.HandleSuccess(c, http.StatusCreated, season, "Season created successfully")
}

func (h *SeasonHandler) GetSeasonByID(c *gin.Context) {
	seasonID := c.Param("id")
	id, err := strconv.ParseUint(seasonID, 10, 8)
	if err != nil {
		helper.HandleError(c, helper.NewAppError(http.StatusBadRequest, constants.ErrInvalidInput, err.Error()), true)
		return
	}

	ctx := c.Request.Context()
	season, err := h.SeasonService.GetSeasonByID(ctx, uint8(id))
	if err != nil {
		helper.HandleGormError(c, err)
		return
	}

	helper.HandleSuccess(c, http.StatusOK, season, "Season found successfully")
}

func (h *SeasonHandler) GetAllSeasons(c *gin.Context) {
	pageStr := c.DefaultQuery("page", "1")
	page, err := strconv.ParseUint(pageStr, 10, 8)
	if err != nil {
		helper.HandleError(c, helper.NewAppError(http.StatusBadRequest, "Invalid page number", err.Error()), true)
		return
	}

	ctx := c.Request.Context()
	seasons, err := h.SeasonService.GetAllSeasons(ctx, page)
	if err != nil {
		helper.HandleError(c, helper.NewAppError(http.StatusInternalServerError, "Failed to list seasons", err.Error()), true)
		return
	}

	helper.HandleSuccess(c, http.StatusOK, seasons, "Season listed successfully")
}

func (h *SeasonHandler) UpdateSeason(c *gin.Context) {
	seasonID := c.Param("id")
	id, err := strconv.ParseUint(seasonID, 10, 8)
	if err != nil {
		helper.HandleError(c, helper.NewAppError(http.StatusBadRequest, constants.ErrInvalidInput, err.Error()), true)
		return
	}

	var season model.Season
	if err := c.ShouldBindJSON(&season); err != nil {
		helper.HandleValidationError(c, err)
		return
	}

	ctx := c.Request.Context()
	season.ID = uint(id)
	err = h.SeasonService.UpdateSeason(ctx, &season)
	if err != nil {
		helper.HandleGormError(c, err)
		return
	}

	helper.HandleSuccess(c, http.StatusOK, season, "Season updated successfully")
}

func (h *SeasonHandler) DeleteSeason(c *gin.Context) {
	seasonID := c.Param("id")
	id, err := strconv.ParseUint(seasonID, 10, 8)
	if err != nil {
		helper.HandleError(c, helper.NewAppError(http.StatusBadRequest, constants.ErrInvalidInput, err.Error()), true)
		return
	}

	ctx := c.Request.Context()
	err = h.SeasonService.DeleteSeason(ctx, uint8(id))
	if err != nil {
		helper.HandleGormError(c, err)
		return
	}

	helper.HandleSuccess(c, http.StatusOK, nil, "Season deleted successfully")
}
