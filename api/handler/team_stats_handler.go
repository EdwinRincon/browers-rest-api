package handler

import (
	"context"
	"errors"
	"net/http"
	"strconv"
	"time"

	"github.com/EdwinRincon/browersfc-api/adapter/mapper"
	"github.com/EdwinRincon/browersfc-api/api/constants"
	"github.com/EdwinRincon/browersfc-api/api/dto"
	"github.com/EdwinRincon/browersfc-api/internal/infrastructure/persistence/model"
	"github.com/EdwinRincon/browersfc-api/helper"
	domainservice "github.com/EdwinRincon/browersfc-api/internal/domain/service"
	"github.com/gin-gonic/gin"
)

type TeamStatsHandler struct {
	TeamStatsDomainService *domainservice.TeamStatsDomainService
	TeamStatsMapper        *mapper.TeamStatsMapper
}

func NewTeamStatsHandler(teamStatsDomainService *domainservice.TeamStatsDomainService) *TeamStatsHandler {
	return &TeamStatsHandler{
		TeamStatsDomainService: teamStatsDomainService,
		TeamStatsMapper:        mapper.NewTeamStatsMapper(),
	}
}

// CreateTeamStats godoc
// @Summary      Create new team stats
// @Description  Creates new team statistics for a season
// @Tags         team-stats
// @ID           createTeamStats
// @Accept       json
// @Produce      json
// @Param        teamStats  body      dto.CreateTeamStatsRequest  true  "Team stats data"
// @Success      201        {object}  dto.TeamStatsShort "Created"
// @Failure      400        {object}  helper.AppError "Invalid input"
// @Failure      404        {object}  helper.AppError "Team or season not found"
// @Failure      409        {object}  helper.AppError "Conflict (e.g., stats already exist for this team/season)"
// @Failure      500        {object}  helper.AppError "Internal server error"
// @Router       /admin/team-stats [post]
// @Security     BearerAuth
func (h *TeamStatsHandler) CreateTeamStats(c *gin.Context) {
	var createRequest dto.CreateTeamStatsRequest
	if err := c.ShouldBindJSON(&createRequest); err != nil {
		helper.WriteErrorResponse(c, helper.BuildValidationErrorFromBinding(err, "body", "Invalid team stats data"))
		return
	}

	// Wrap context with timeout for DB/service calls
	ctx, cancel := context.WithTimeout(c.Request.Context(), 5*time.Second)
	defer cancel()

	// Map DTO to domain entity
	teamStats := h.TeamStatsMapper.CreateRequestToDomain(&createRequest)

	err := h.TeamStatsDomainService.CreateTeamStats(ctx, teamStats)
	if err != nil {
		switch {
		case errors.Is(err, constants.ErrRecordAlreadyExists):
			helper.WriteErrorResponse(c, helper.NewConflictError("team_stats", "Team stats already exist for this team and season"))
			return
		case errors.Is(err, constants.ErrTeamNotFound):
			helper.WriteErrorResponse(c, helper.NewNotFoundError("team"))
			return
		case errors.Is(err, constants.ErrSeasonNotFound):
			helper.WriteErrorResponse(c, helper.NewNotFoundError("season"))
			return
		default:
			helper.WriteErrorResponse(c, helper.NewInternalServerError(err))
			return
		}
	}

	// Return the created team stats as a short response
	createdTeamStatsShort := h.TeamStatsMapper.DomainToShortDTO(teamStats)
	helper.WriteSuccessResponse(c, http.StatusCreated, createdTeamStatsShort, "Team stats created successfully")
}

// GetTeamStatsByID godoc
// @Summary      Get team stats by ID
// @Description  Returns team statistics by their ID
// @Tags         team-stats
// @ID           getTeamStatsByID
// @Param        id   path      int  true  "Team Stats ID"
// @Success      200  {object}  dto.TeamStatsResponse "Success"
// @Failure      400  {object}  helper.AppError "Invalid input"
// @Failure      404  {object}  helper.AppError "Team stats not found"
// @Failure      500  {object}  helper.AppError "Internal server error"
// @Router       /team-stats/{id} [get]
// @Security     BearerAuth
func (h *TeamStatsHandler) GetTeamStatsByID(c *gin.Context) {
	idStr := c.Param("id")
	id, err := strconv.ParseUint(idStr, 10, 64)
	if err != nil {
		helper.WriteErrorResponse(c, helper.NewBadRequestError("id", "Invalid team stats ID"))
		return
	}

	// Wrap context with timeout for DB/service calls
	ctx, cancel := context.WithTimeout(c.Request.Context(), 5*time.Second)
	defer cancel()

	teamStats, err := h.TeamStatsDomainService.GetTeamStatsByID(ctx, id)
	if err != nil {
		if errors.Is(err, constants.ErrRecordNotFound) {
			helper.WriteErrorResponse(c, helper.NewNotFoundError("team_stats"))
		} else {
			helper.WriteErrorResponse(c, helper.NewInternalServerError(err))
		}
		return
	}

	response := h.TeamStatsMapper.DomainToResponse(teamStats)
	helper.WriteSuccessResponse(c, http.StatusOK, response, "Team stats retrieved successfully")
}

// GetTeamStatsBySeasonID godoc
// @Summary      Get team stats by season
// @Description  Returns all team statistics for a specific season, ordered by rank
// @Tags         team-stats
// @ID           getTeamStatsBySeasonID
// @Param        id   path      int  true  "Season ID"
// @Success      200        {object}  []dto.TeamStatsResponse "Success"
// @Failure      400        {object}  helper.AppError "Invalid input"
// @Failure      404        {object}  helper.AppError "Season not found"
// @Failure      500        {object}  helper.AppError "Internal server error"
// @Router       /seasons/{id}/team-stats [get]
// @Security     BearerAuth
func (h *TeamStatsHandler) GetTeamStatsBySeasonID(c *gin.Context) {
	seasonIDStr := c.Param("id")
	seasonID, err := strconv.ParseUint(seasonIDStr, 10, 64)
	if err != nil {
		helper.WriteErrorResponse(c, helper.NewBadRequestError("id", "Invalid season ID"))
		return
	}

	// Wrap context with timeout for DB/service calls
	ctx, cancel := context.WithTimeout(c.Request.Context(), 5*time.Second)
	defer cancel()

	teamStats, err := h.TeamStatsDomainService.GetTeamStatsBySeasonID(ctx, seasonID)
	if err != nil {
		if errors.Is(err, constants.ErrSeasonNotFound) {
			helper.WriteErrorResponse(c, helper.NewNotFoundError("season"))
		} else {
			helper.WriteErrorResponse(c, helper.NewInternalServerError(err))
		}
		return
	}

	response := h.TeamStatsMapper.DomainListToResponse(teamStats)
	helper.WriteSuccessResponse(c, http.StatusOK, response, "Team stats for season retrieved successfully")
}

// GetTeamStatsByTeamID godoc
// @Summary      Get team stats by team
// @Description  Returns all statistics for a specific team across all seasons
// @Tags         team-stats
// @ID           getTeamStatsByTeamID
// @Param        id   path      int  true  "Team ID"
// @Success      200      {object}  []dto.TeamStatsResponse "Success"
// @Failure      400      {object}  helper.AppError "Invalid input"
// @Failure      404      {object}  helper.AppError "Team not found"
// @Failure      500      {object}  helper.AppError "Internal server error"
// @Router       /teams/{id}/stats [get]
// @Security     BearerAuth
func (h *TeamStatsHandler) GetTeamStatsByTeamID(c *gin.Context) {
	teamIDStr := c.Param("id")
	teamID, err := strconv.ParseUint(teamIDStr, 10, 64)
	if err != nil {
		helper.WriteErrorResponse(c, helper.NewBadRequestError("id", "Invalid team ID"))
		return
	}

	// Wrap context with timeout for DB/service calls
	ctx, cancel := context.WithTimeout(c.Request.Context(), 5*time.Second)
	defer cancel()

	teamStats, err := h.TeamStatsDomainService.GetTeamStatsByTeamID(ctx, teamID)
	if err != nil {
		if errors.Is(err, constants.ErrTeamNotFound) {
			helper.WriteErrorResponse(c, helper.NewNotFoundError("team"))
		} else {
			helper.WriteErrorResponse(c, helper.NewInternalServerError(err))
		}
		return
	}

	response := h.TeamStatsMapper.DomainListToResponse(teamStats)
	helper.WriteSuccessResponse(c, http.StatusOK, response, "Team stats for team retrieved successfully")
}

// GetPaginatedTeamStats godoc
// @Summary      Get paginated team stats
// @Description  Returns a paginated list of team statistics
// @Tags         team-stats
// @Summary Get paginated team stats
// @Description Retrieve team stats with pagination, optional sorting, and ordering.
// @Tags team-stats
// @Accept json
// @Produce json
// @Param page query int false "Page number" default(0)
// @Param pageSize query int false "Page size" default(10)
// @Param sort query string false "Sort field"
// @Param order query string false "Sort order" Enums(asc, desc) default(asc)
// @Success 200 {object} helper.AppSuccess{data=helper.PaginatedResponse{items=[]dto.TeamStatsResponse, totalCount=int}}
// @Failure 400 {object} helper.AppError "Invalid input"
// @Failure 500 {object} helper.AppError "Internal server error"
// @Router /team-stats [get]
// @Security BearerAuth
func (h *TeamStatsHandler) GetPaginatedTeamStats(c *gin.Context) {
	sort := c.DefaultQuery("sort", "rank")
	order := c.DefaultQuery("order", "asc")
	page, _ := strconv.Atoi(c.DefaultQuery("page", "0"))
	pageSize, _ := strconv.Atoi(c.DefaultQuery("pageSize", "10"))

	if page < 0 {
		page = 0
	}
	if pageSize < 1 || pageSize > 100 {
		pageSize = 10
	}

	// Validate sort field
	if err := helper.ValidateSort(model.TeamStat{}, sort); err != nil {
		helper.WriteErrorResponse(c, helper.NewBadRequestError("sort", err.Error()))
		return
	}

	// Wrap context with timeout for DB/service calls
	ctx, cancel := context.WithTimeout(c.Request.Context(), 5*time.Second)
	defer cancel()

	teamStats, total, err := h.TeamStatsDomainService.GetPaginatedTeamStats(ctx, sort, order, page, pageSize)
	if err != nil {
		helper.WriteErrorResponse(c, helper.NewInternalServerError(err))
		return
	}

	response := helper.PaginatedResponse{
		Items:      h.TeamStatsMapper.DomainListToResponse(teamStats),
		TotalCount: total,
	}

	helper.WriteSuccessResponse(c, http.StatusOK, response, "Team stats retrieved successfully")
}

// UpdateTeamStats godoc
// @Summary      Update team stats
// @Description  Updates team statistics with the provided data
// @Tags         team-stats
// @ID           updateTeamStats
// @Accept       json
// @Produce      json
// @Param        id         path      int                        true  "Team Stats ID"
// @Param        teamStats  body      dto.UpdateTeamStatsRequest  true  "Updated team stats data"
// @Success      200        {object}  dto.TeamStatsResponse "Success"
// @Failure      400        {object}  helper.AppError "Invalid input"
// @Failure      404        {object}  helper.AppError "Team stats, team, or season not found"
// @Failure      409        {object}  helper.AppError "Conflict (e.g., duplicate team/season combination)"
// @Failure      500        {object}  helper.AppError "Internal server error"
// @Router       /admin/team-stats/{id} [put]
// @Security     BearerAuth
func (h *TeamStatsHandler) UpdateTeamStats(c *gin.Context) {
	idStr := c.Param("id")
	id, err := strconv.ParseUint(idStr, 10, 64)
	if err != nil {
		helper.WriteErrorResponse(c, helper.NewBadRequestError("id", "Invalid team stats ID"))
		return
	}

	var updateRequest dto.UpdateTeamStatsRequest
	if err = c.ShouldBindJSON(&updateRequest); err != nil {
		helper.WriteErrorResponse(c, helper.BuildValidationErrorFromBinding(err, "body", "Invalid team stats data"))
		return
	}

	// Wrap context with timeout for DB/service calls
	ctx, cancel := context.WithTimeout(c.Request.Context(), 5*time.Second)
	defer cancel()

	// Get the current team stats to apply partial updates
	currentTeamStats, err := h.TeamStatsDomainService.GetTeamStatsByID(ctx, id)
	if err != nil {
		if errors.Is(err, constants.ErrRecordNotFound) {
			helper.WriteErrorResponse(c, helper.NewNotFoundError("team_stats"))
		} else {
			helper.WriteErrorResponse(c, helper.NewInternalServerError(err))
		}
		return
	}

	// Apply the updates to the domain entity
	updatedTeamStats := h.TeamStatsMapper.UpdateRequestToDomain(currentTeamStats, &updateRequest)

	updatedResult, err := h.TeamStatsDomainService.UpdateTeamStats(ctx, id, updatedTeamStats)
	if err != nil {
		switch {
		case errors.Is(err, constants.ErrRecordNotFound):
			helper.WriteErrorResponse(c, helper.NewNotFoundError("team_stats"))
			return
		case errors.Is(err, constants.ErrRecordAlreadyExists):
			helper.WriteErrorResponse(c, helper.NewConflictError("team_stats", "Team stats already exist for this team and season combination"))
			return
		case errors.Is(err, constants.ErrTeamNotFound):
			helper.WriteErrorResponse(c, helper.NewNotFoundError("team"))
			return
		case errors.Is(err, constants.ErrSeasonNotFound):
			helper.WriteErrorResponse(c, helper.NewNotFoundError("season"))
			return
		default:
			helper.WriteErrorResponse(c, helper.NewInternalServerError(err))
			return
		}
	}

	response := h.TeamStatsMapper.DomainToResponse(updatedResult)
	helper.WriteSuccessResponse(c, http.StatusOK, response, "Team stats updated successfully")
}

// DeleteTeamStats godoc
// @Summary      Delete team stats
// @Description  Deletes team statistics by their ID
// @Tags         team-stats
// @ID           deleteTeamStats
// @Param        id   path      int  true  "Team Stats ID"
// @Success      204 "No Content"
// @Failure      400  {object}  helper.AppError "Invalid input"
// @Failure      500  {object}  helper.AppError "Internal server error"
// @Router       /admin/team-stats/{id} [delete]
// @Security     BearerAuth
func (h *TeamStatsHandler) DeleteTeamStats(c *gin.Context) {
	idStr := c.Param("id")
	id, err := strconv.ParseUint(idStr, 10, 64)
	if err != nil {
		helper.WriteErrorResponse(c, helper.NewBadRequestError("id", "Invalid team stats ID"))
		return
	}

	// Wrap context with timeout for DB/service calls
	ctx, cancel := context.WithTimeout(c.Request.Context(), 5*time.Second)
	defer cancel()

	err = h.TeamStatsDomainService.DeleteTeamStats(ctx, id)
	if err != nil && !errors.Is(err, constants.ErrRecordNotFound) {
		helper.WriteErrorResponse(c, helper.NewInternalServerError(err))
		return
	}

	c.Status(http.StatusNoContent)
}
