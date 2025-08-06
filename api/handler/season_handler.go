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

type SeasonHandler struct {
	SeasonService service.SeasonService
}

func NewSeasonHandler(seasonService service.SeasonService) *SeasonHandler {
	return &SeasonHandler{
		SeasonService: seasonService,
	}
}

// CreateSeason godoc
// @Summary      Create a new season
// @Description  Creates a new season with the provided data
// @Tags         seasons
// @ID           createSeason
// @Accept       json
// @Produce      json
// @Param        season  body      model.Season  true  "Season data"
// @Success      201     {object}  model.Season "Created"
// @Failure      400     {object}  helper.AppError "Invalid input"
// @Failure      409     {object}  helper.AppError "Conflict (e.g., season already exists)"
// @Failure      500     {object}  helper.AppError "Internal server error"
// @Router       /seasons [post]
// @Security     ApiKeyAuth
func (h *SeasonHandler) CreateSeason(c *gin.Context) {
	var season model.Season
	if err := c.ShouldBindJSON(&season); err != nil {
		helper.RespondWithError(c, helper.BadRequest("body", "Invalid season data"))
		return
	}

	ctx := c.Request.Context()
	err := h.SeasonService.CreateSeason(ctx, &season)
	if err != nil {
		if errors.Is(err, gorm.ErrDuplicatedKey) {
			helper.RespondWithError(c, helper.Conflict("season", "A season with these details already exists"))
		} else {
			helper.RespondWithError(c, helper.InternalError(err))
		}
		return
	}

	helper.HandleSuccess(c, http.StatusCreated, season, "Season created successfully")
}

// GetSeasonByID godoc
// @Summary      Get a season by ID
// @Description  Returns the details of a season by its ID
// @Tags         seasons
// @ID           getSeasonByID
// @Param        id   path      int  true  "Season ID"
// @Success      200  {object}  model.Season "Success"
// @Failure      400  {object}  helper.AppError "Invalid input"
// @Failure      404  {object}  helper.AppError "Season not found"
// @Failure      500  {object}  helper.AppError "Internal server error"
// @Router       /seasons/{id} [get]
// @Security     ApiKeyAuth
func (h *SeasonHandler) GetSeasonByID(c *gin.Context) {
	seasonID := c.Param("id")
	id, err := strconv.ParseUint(seasonID, 10, 32)
	if err != nil {
		helper.RespondWithError(c, helper.BadRequest("id", "Invalid season ID"))
		return
	}

	ctx := c.Request.Context()
	season, err := h.SeasonService.GetSeasonByID(ctx, uint(id))
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			helper.RespondWithError(c, helper.NotFound("season"))
		} else {
			helper.RespondWithError(c, helper.InternalError(err))
		}
		return
	}

	helper.HandleSuccess(c, http.StatusOK, season, "Season found successfully")
}

// GetAllSeasons godoc
// @Summary      List seasons
// @Description  Retrieves a paginated list of seasons
// @Tags         seasons
// @ID           listSeasons
// @Param        page  query     int  false  "Page number"
// @Success      200   {array}   model.Season "Seasons listed successfully"
// @Failure      400   {object}  helper.AppError "Invalid input"
// @Failure      500   {object}  helper.AppError "Internal server error"
// @Router       /seasons [get]
// @Security     ApiKeyAuth
func (h *SeasonHandler) GetAllSeasons(c *gin.Context) {
	pageStr := c.DefaultQuery("page", "1")
	page, err := strconv.ParseUint(pageStr, 10, 32)
	if err != nil || page < 1 {
		helper.RespondWithError(c, helper.BadRequest("page", "Invalid page number"))

		return
	}

	ctx := c.Request.Context()
	seasons, err := h.SeasonService.GetAllSeasons(ctx, page)
	if err != nil {
		helper.RespondWithError(c, helper.InternalError(err))
		return
	}

	helper.HandleSuccess(c, http.StatusOK, seasons, "Season listed successfully")
}

// UpdateSeason godoc
// @Summary      Update an existing season
// @Description  Updates the details of an existing season by ID
// @Tags         seasons
// @ID           updateSeason
// @Accept       json
// @Produce      json
// @Param        id      path      int         true  "Season ID"
// @Param        season  body      model.Season true  "Updated season data"
// @Success      200     {object}  model.Season "Updated"
// @Failure      400     {object}  helper.AppError "Invalid input"
// @Failure      404     {object}  helper.AppError "Season not found"
// @Failure      500     {object}  helper.AppError "Internal server error"
// @Router       /seasons/{id} [put]
// @Security     ApiKeyAuth
func (h *SeasonHandler) UpdateSeason(c *gin.Context) {
	seasonID := c.Param("id")
	id, err := strconv.ParseUint(seasonID, 10, 32)
	if err != nil {
		helper.RespondWithError(c, helper.BadRequest("id", "Invalid season ID"))
		return
	}

	var season model.Season
	if err := c.ShouldBindJSON(&season); err != nil {
		helper.RespondWithError(c, helper.BadRequest("body", "Invalid season data"))
		return
	}

	ctx := c.Request.Context()
	season.ID = uint(id)
	err = h.SeasonService.UpdateSeason(ctx, &season)
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			helper.RespondWithError(c, helper.NotFound("season"))
		} else {
			helper.RespondWithError(c, helper.InternalError(err))
		}
		return
	}

	helper.HandleSuccess(c, http.StatusOK, season, "Season updated successfully")
}

// DeleteSeason godoc
// @Summary      Delete a season
// @Description  Deletes a season by its ID
// @Tags         seasons
// @ID           deleteSeason
// @Param        id   path      int  true  "Season ID"
// @Success      204 "No Content"
// @Failure      400  {object}  helper.AppError "Invalid input"
// @Failure      404  {object}  helper.AppError "Season not found"
// @Failure      500  {object}  helper.AppError "Internal server error"
// @Router       /seasons/{id} [delete]
// @Security     ApiKeyAuth
func (h *SeasonHandler) DeleteSeason(c *gin.Context) {
	seasonID := c.Param("id")
	id, err := strconv.ParseUint(seasonID, 10, 32)
	if err != nil {
		helper.RespondWithError(c, helper.BadRequest("id", "Invalid season ID"))
		return
	}

	ctx := c.Request.Context()
	err = h.SeasonService.DeleteSeason(ctx, uint(id))
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			helper.RespondWithError(c, helper.NotFound("season"))
		} else {
			helper.RespondWithError(c, helper.InternalError(err))
		}
		return
	}

	c.Status(http.StatusNoContent)
}
