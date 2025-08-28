package handler

import (
	"context"
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

const errInvalidSeasonID = "Invalid season ID"

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
// @Param        season  body      dto.CreateSeasonRequest  true  "Season data"
// @Success      201     {object}  dto.SeasonResponse "Created"
// @Failure      400     {object}  helper.AppError "Invalid input"
// @Failure      409     {object}  helper.AppError "Conflict (e.g., season already exists)"
// @Failure      500     {object}  helper.AppError "Internal server error"
// @Router       /admin/seasons [post]
// @Security     ApiKeyAuth
func (h *SeasonHandler) CreateSeason(c *gin.Context) {
	var createRequest dto.CreateSeasonRequest
	if err := c.ShouldBindJSON(&createRequest); err != nil {
		helper.RespondWithError(c, helper.BadRequest("body", "Invalid season data"))
		return
	}

	// Validate additional business rules
	if createRequest.EndDate.Before(createRequest.StartDate) {
		helper.RespondWithError(c, helper.BadRequest("end_date", "End date must be after start date"))
		return
	}

	// Wrap context with timeout for DB/service calls
	ctx, cancel := context.WithTimeout(c.Request.Context(), 5*time.Second)
	defer cancel()

	seasonResponse, err := h.SeasonService.CreateSeason(ctx, &createRequest)
	if err != nil {
		if err == constants.ErrRecordAlreadyExists {
			helper.RespondWithError(c, helper.Conflict("season", "A season with this year already exists"))
		} else {
			helper.RespondWithError(c, helper.InternalError(err))
		}
		return
	}

	helper.HandleSuccess(c, http.StatusCreated, seasonResponse, "Season created successfully")
}

// GetSeasonByID godoc
// @Summary      Get a season by ID
// @Description  Returns the details of a season by its ID
// @Tags         seasons
// @ID           getSeasonByID
// @Param        id   path      int  true  "Season ID"
// @Success      200  {object}  dto.SeasonResponse "Success"
// @Failure      400  {object}  helper.AppError "Invalid input"
// @Failure      404  {object}  helper.AppError "Season not found"
// @Failure      500  {object}  helper.AppError "Internal server error"
// @Router       /seasons/{id} [get]
// @Security     ApiKeyAuth
func (h *SeasonHandler) GetSeasonByID(c *gin.Context) {
	seasonID := c.Param("id")
	id, err := strconv.ParseUint(seasonID, 10, 64)
	if err != nil {
		helper.RespondWithError(c, helper.BadRequest("id", errInvalidSeasonID))
		return
	}

	// Wrap context with timeout for DB/service calls
	ctx, cancel := context.WithTimeout(c.Request.Context(), 5*time.Second)
	defer cancel()

	season, err := h.SeasonService.GetSeasonByID(ctx, id)
	if err != nil {
		if err == constants.ErrRecordNotFound {
			helper.RespondWithError(c, helper.NotFound("season"))
		} else {
			helper.RespondWithError(c, helper.InternalError(err))
		}
		return
	}

	seasonResponse := mapper.ToSeasonResponse(season)
	helper.HandleSuccess(c, http.StatusOK, seasonResponse, "Season found successfully")
}

// GetCurrentSeason godoc
// @Summary      Get the current season
// @Description  Returns the details of the currently active season
// @Tags         seasons
// @ID           getCurrentSeason
// @Success      200  {object}  dto.SeasonResponse "Success"
// @Failure      404  {object}  helper.AppError "No current season found"
// @Failure      500  {object}  helper.AppError "Internal server error"
// @Router       /seasons/current [get]
// @Security     ApiKeyAuth
func (h *SeasonHandler) GetCurrentSeason(c *gin.Context) {
	// Wrap context with timeout for DB/service calls
	ctx, cancel := context.WithTimeout(c.Request.Context(), 5*time.Second)
	defer cancel()

	season, err := h.SeasonService.GetCurrentSeason(ctx)
	if err != nil {
		if err == constants.ErrRecordNotFound {
			helper.RespondWithError(c, helper.NotFound("current season"))
		} else {
			helper.RespondWithError(c, helper.InternalError(err))
		}
		return
	}

	seasonResponse := mapper.ToSeasonResponse(season)
	helper.HandleSuccess(c, http.StatusOK, seasonResponse, "Current season retrieved successfully")
}

// GetPaginatedSeasons godoc
// @Summary      List seasons with pagination and sorting
// @Description  Retrieves a paginated list of seasons with sorting options
// @Tags         seasons
// @ID           listSeasons
// @Param        page      query     int     false  "Page number (0-based, default: 0)"
// @Param        pageSize  query     int     false  "Items per page (default: 10)"
// @Param        sort      query     string  false  "Sort field (e.g., year, created_at)"
// @Param        order     query     string  false  "Sort order (asc/desc)"
// @Success      200       {object}  helper.PaginatedResponse "Seasons listed successfully"
// @Failure      400       {object}  helper.AppError "Invalid input"
// @Failure      500       {object}  helper.AppError "Internal server error"
// @Router       /seasons [get]
// @Security     ApiKeyAuth
func (h *SeasonHandler) GetPaginatedSeasons(c *gin.Context) {
	sort := c.DefaultQuery("sort", "year")
	order := c.DefaultQuery("order", "desc")
	page, _ := strconv.Atoi(c.DefaultQuery("page", "0"))
	pageSize, _ := strconv.Atoi(c.DefaultQuery("pageSize", "10"))

	// Validate pagination parameters
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
	if err := helper.ValidateSort(model.Season{}, sort); err != nil {
		helper.RespondWithError(c, helper.BadRequest("sort", err.Error()))
		return
	}

	// Wrap context with timeout for DB/service calls
	ctx, cancel := context.WithTimeout(c.Request.Context(), 5*time.Second)
	defer cancel()

	seasons, total, err := h.SeasonService.GetPaginatedSeasons(ctx, sort, order, page, pageSize)
	if err != nil {
		helper.RespondWithError(c, helper.InternalError(err))
		return
	}

	seasonResponses := mapper.ToSeasonResponseList(seasons)

	response := helper.PaginatedResponse{
		Items:      seasonResponses,
		TotalCount: total,
	}

	helper.HandleSuccess(c, http.StatusOK, response, "Seasons retrieved successfully")
}

// UpdateSeason godoc
// @Summary      Update an existing season
// @Description  Updates the details of an existing season by ID
// @Tags         seasons
// @ID           updateSeason
// @Accept       json
// @Produce      json
// @Param        id      path      int                   true  "Season ID"
// @Param        season  body      dto.UpdateSeasonRequest true  "Updated season data"
// @Success      200     {object}  dto.SeasonResponse "Updated"
// @Failure      400     {object}  helper.AppError "Invalid input"
// @Failure      404     {object}  helper.AppError "Season not found"
// @Failure      409     {object}  helper.AppError "Conflict (e.g., year already exists)"
// @Failure      500     {object}  helper.AppError "Internal server error"
// @Router       /admin/seasons/{id} [put]
// @Security     ApiKeyAuth
func (h *SeasonHandler) UpdateSeason(c *gin.Context) {
	seasonID := c.Param("id")
	id, err := strconv.ParseUint(seasonID, 10, 64)
	if err != nil {
		helper.RespondWithError(c, helper.BadRequest("id", errInvalidSeasonID))
		return
	}

	var updateRequest dto.UpdateSeasonRequest
	if err = c.ShouldBindJSON(&updateRequest); err != nil {
		helper.RespondWithError(c, helper.ProcessValidationError(err, "body", constants.MsgInvalidSeasonData))
		return
	}

	// Additional validations for date consistency
	if updateRequest.StartDate != nil && updateRequest.EndDate != nil {
		if updateRequest.EndDate.Before(*updateRequest.StartDate) {
			helper.RespondWithError(c, helper.BadRequest("end_date", "End date must be after start date"))
			return
		}
	}

	// Wrap context with timeout for DB/service calls
	ctx, cancel := context.WithTimeout(c.Request.Context(), 5*time.Second)
	defer cancel()

	seasonResponse, err := h.SeasonService.UpdateSeason(ctx, id, &updateRequest)
	if err != nil {
		switch err {
		case constants.ErrRecordNotFound:
			helper.RespondWithError(c, helper.NotFound("season"))
		case constants.ErrRecordAlreadyExists:
			helper.RespondWithError(c, helper.Conflict("year", "A season with this year already exists"))
		default:
			helper.RespondWithError(c, helper.InternalError(err))
		}
		return
	}

	helper.HandleSuccess(c, http.StatusOK, seasonResponse, "Season updated successfully")
}

// SetCurrentSeason godoc
// @Summary      Set a season as the current active season
// @Description  Sets the specified season as the current active season
// @Tags         seasons
// @ID           setCurrentSeason
// @Param        id   path      int  true  "Season ID"
// @Success      200  {object}  map[string]interface{} "Success message"
// @Failure      400  {object}  helper.AppError "Invalid input"
// @Failure      404  {object}  helper.AppError "Season not found"
// @Failure      500  {object}  helper.AppError "Internal server error"
// @Router       /admin/seasons/{id}/set-current [put]
// @Security     ApiKeyAuth
func (h *SeasonHandler) SetCurrentSeason(c *gin.Context) {
	seasonID := c.Param("id")
	id, err := strconv.ParseUint(seasonID, 10, 64)
	if err != nil {
		helper.RespondWithError(c, helper.BadRequest("id", errInvalidSeasonID))
		return
	}

	// Wrap context with timeout for DB/service calls
	ctx, cancel := context.WithTimeout(c.Request.Context(), 5*time.Second)
	defer cancel()

	err = h.SeasonService.SetCurrentSeason(ctx, id)
	if err != nil {
		if err == constants.ErrRecordNotFound {
			helper.RespondWithError(c, helper.NotFound("season"))
		} else {
			helper.RespondWithError(c, helper.InternalError(err))
		}
		return
	}

	response := map[string]interface{}{
		"message": "Season set as current successfully",
		"id":      id,
	}

	helper.HandleSuccess(c, http.StatusOK, response, "Season set as current successfully")
}

// DeleteSeason godoc
// @Summary      Delete a season
// @Description  Deletes a season by its ID
// @Tags         seasons
// @ID           deleteSeason
// @Param        id   path      int  true  "Season ID"
// @Success      204  "No Content"
// @Failure      400  {object}  helper.AppError "Invalid input"
// @Failure      404  {object}  helper.AppError "Season not found"
// @Failure      500  {object}  helper.AppError "Internal server error"
// @Router       /admin/seasons/{id} [delete]
// @Security     ApiKeyAuth
func (h *SeasonHandler) DeleteSeason(c *gin.Context) {
	seasonID := c.Param("id")
	id, err := strconv.ParseUint(seasonID, 10, 64)
	if err != nil {
		helper.RespondWithError(c, helper.BadRequest("id", errInvalidSeasonID))
		return
	}

	// Wrap context with timeout for DB/service calls
	ctx, cancel := context.WithTimeout(c.Request.Context(), 5*time.Second)
	defer cancel()

	err = h.SeasonService.DeleteSeason(ctx, id)
	if err != nil {
		if err == constants.ErrRecordNotFound {
			helper.RespondWithError(c, helper.NotFound("season"))
		} else {
			helper.RespondWithError(c, helper.InternalError(err))
		}
		return
	}

	c.Status(http.StatusNoContent)
}
